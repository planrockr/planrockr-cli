package pkg

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	jiraApi "github.com/andygrunwald/go-jira"
	"gopkg.in/cheggaaa/pb.v1"
	"net/http"

	"github.com/planrockr/planrockr-cli/config"
)

type JiraImporter struct {
	user   string
	pass   string
	url    string
	client *jiraApi.Client
}

type hook struct {
	EventName string                 `json:"eventName"`
	EventData map[string]interface{} `json:"eventData"`
	Data      IssueHook              `json:"data"`
}

// IssueHook represent the hook sended by Jira webhook.
type IssueHook struct {
	Timestamp    int64                  `json:"timestamp"`
	WebhookEvent string                 `json:"webhookEvent"`
	TypeName     string                 `json:"issue_event_type_name,omitempty"`
	User         *jiraApi.User          `json:"user,omitempty"`
	Issue        jiraApi.Issue          `json:"issue,omitempty"`
	Changelog    *changelog             `json:"changelog,omitempty"`
	Comment      *jiraApi.Comment       `json:"comment,omitempty"`
	Worklog      *jiraApi.WorklogRecord `json:"worklog,omitempty"`
}

type changelog struct {
	ID    string          `json:"id"`
	Items []changelogItem `json:"items"`
}

type changelogItem struct {
	Field      string      `json:"field"`
	Fieldtype  string      `json:"fieldtype"`
	From       interface{} `json:"from"`
	FromString string      `json:"fromString"`
	ToString   string      `json:"toString"`
	To         string      `json:"to"`
}

type ProjectData struct {
	Project struct {
		Id int
	}
}

type BoardData struct {
	Id int
}

var projectData ProjectData
var configData config.Config

func JiraImport(host string, user string, password string) error {
	err := config.Init()
	if err != nil {
		return errors.New("Error reading config file")
	}

	configData = config.Get()
	if configData.Auth.Token == "" {
		return errors.New("Missing token")
	}

	i := JiraImporter{
		user: user,
		pass: password,
		url:  host,
	}
	c, err := jiraApi.NewClient(nil, i.url)
	if err != nil {
		return errors.New("Error creating Jira Client")
	}
	i.client = c
	i.client.Authentication.SetBasicAuth(i.user, i.pass)
	var wg sync.WaitGroup
	jql, err := getJql()
	if err != nil {
		return errors.New("Error reading JQL")
	}
	if jql != "\n" && jql != "" {
		err = processWithJql(jql, i, &wg)
		return err
	}

	projects, err := GetProjects(i)
	if err != nil {
		log.Fatalf("[IMPORTER] Failed to get The project list: %v", err)
		return err
	}
	projects = selectProject(projects)

	// Start processing the importers.
	fmt.Println("Importing...")
	var projectId int
	var boardId int
	for _, proj := range projects {
		pID, err := strconv.Atoi(proj.ID)
		jql = fmt.Sprintf("project=%d", pID)
		if err != nil {
			log.Errorf("[IMPORTER] Failed to convert projectID to integer: %v", err)
			continue
		}
		pName := proj.Name
		pKey := proj.Key
		projectId, err = createProject(pName)
		if err != nil {
			fmt.Println(err)
		}
		boardId, err = createBoard(proj.ID+"_"+proj.Key, "3")
		if err != nil {
			fmt.Println(err)
			return err
		}
		wg.Add(1)
		go func() {
			Process(i, pID, pKey, jql)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Finished importing projects")

	err = createHook(i.url, jql, projectId, boardId, i.user, i.pass)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Webhook created")

	return nil
}

func processWithJql(jql string, i JiraImporter, wg *sync.WaitGroup) error {
	var projectId int
	var boardId int
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter project name: ")
	input, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	pName := string(input)
	projectId, err = createProject(pName)
	if err != nil {
		return err
	}
	boardId, err = createBoard(pName, "3")
	if err != nil {
		return err
	}
	Process(i, 0, "0", jql)

	err = createHook(i.url, jql, projectId, boardId, i.user, i.pass)
	if err != nil {
		return err
	}
	fmt.Println("Webhook created")

	return nil
}

func createHook(host string, jql string, projectId int, boardId int, user string, password string) error {
	// Disable HTTP/2
	http.DefaultClient.Transport = &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	type Filters struct {
		IssueRelatedEventsSection string `json:"issue-related-events-section"`
	}
	type Payload struct {
		Name                string   `json:"name"`
		URL                 string   `json:"url"`
		Events              []string `json:"events"`
		JqlFilter           string   `json:"jqlFilter"`
		Filters             Filters  `json:"filters"`
		ExcludeIssueDetails bool     `json:"excludeIssueDetails"`
	}

	if jql == "\n" {
		jql = ""
	}
	data := Payload{
		Name:                "Planrockr",
		URL:                 "https://app.planrockr.com/hook/jira/${project.id}/${project.key}/" + strconv.Itoa(configData.Auth.Id) + "/" + strconv.Itoa(projectId) + "/" + strconv.Itoa(boardId),
		Events:              []string{"jira:issue_created", "jira:issue_updated", "worklog_created", "worklog_updated", "worklog_deleted", "comment_created", "comment_updated", "comment_deleted", "project_deleted", "project_updated", "jira:issue_deleted", "project_created", "jira:worklog_updated"},
		JqlFilter:           jql,
		Filters:             Filters{IssueRelatedEventsSection: jql},
		ExcludeIssueDetails: false,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", host+"/rest/webhooks/1.0/webhook", body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func createProject(name string) (int, error) {
	body := strings.NewReader("parameters%5Bname%5D=" + name)
	req, err := http.NewRequest("POST", configData.BaseUrl+"/rpc/v1/project/create", body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "planrockr-cli")
	req.Header.Set("Authorization-Coderockr", configData.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized {
			panic("You must login")
		}
		return 0, errors.New("Error creating project")
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &projectData)
	if err != nil {
		return 0, errors.New("Error parsing project data")
	}

	return projectData.Project.Id, nil
}

func createBoard(boardId string, boardType string) (int, error) {
	q, err := url.ParseQuery("parameters[projectId]=" + strconv.Itoa(projectData.Project.Id) + "&parameters[boardId]=" + boardId + "&parameters[boardType]=" + boardType)
	if err != nil {
		return 0, err
	}
	body := strings.NewReader(q.Encode())
	req, err := http.NewRequest("POST", configData.BaseUrl+"/rpc/v1/project/addBoard", body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "planrockr-cli")
	req.Header.Set("Authorization-Coderockr", configData.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusCreated {
		fmt.Println(resp.StatusCode)
		return 0, errors.New("Error creating board")
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	b := BoardData{}
	err = json.Unmarshal(buf, &b)
	if err != nil {
		return 0, errors.New("Error parsing board data")
	}

	return b.Id, nil
}

func enqueue(toImport string) error {
	var jsonStr = []byte(toImport)
	body := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("POST", configData.BaseUrl+"/importer/"+strconv.Itoa(projectData.Project.Id), body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "planrockr-cli")
	req.Header.Set("Authorization-Coderockr", configData.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return errors.New("Error sending to queue")
	}

	return nil
}

func GetProjects(i JiraImporter) (jiraApi.ProjectList, error) {
	projects, resp, err := i.client.Project.GetList()
	if err != nil || resp == nil || projects == nil {
		err = errors.New("Failed to get the list of projects. Jira's response: " + string(resp.Status))
		return nil, err
	}
	return *projects, err
}

func getJql() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("JQL query(enter to select a project): ")
	jql, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	return string(jql), nil
}

func selectProject(list jiraApi.ProjectList) jiraApi.ProjectList {
	reader := bufio.NewReader(os.Stdin)
	res := make(jiraApi.ProjectList, 0, len(list))

	fmt.Println("Projects available:")
	for i, op := range list {
		fmt.Printf("\t%d - %s(%s)\n", i, op.Name, op.Key)
	}
	fmt.Print("Select option: ")
	selected, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	i, err := strconv.Atoi(string(selected))
	if err != nil {
		log.Fatal("Option invalid", err)
	}
	res = append(res, list[i])

	return res
}

func Process(i JiraImporter, jiraProjectId int, jiraProjectKey string, jql string) {
	op := jiraApi.SearchOptions{
		StartAt:    0,
		MaxResults: 50,
	}

	for {
		searchString := fmt.Sprintf("project=%d", jiraProjectId)
		if jql != "\n" && jql != "" {
			searchString = jql
		}
		issues, resp, err := i.client.Issue.Search(searchString, &op)
		if err != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode != 429 {
				fmt.Println("Failed to get issues. Resp:" + string(body))
				continue
			}
			h := resp.Header.Get("X-Ratelimit-Reset")
			t := 30 * time.Second
			if len(h) > 0 {
				reset, err := strconv.Atoi(h)
				if err == nil {
					t1 := time.Unix(int64(reset), 0)
					t = t1.Sub(t1)
				}
			}
			fmt.Println("Failed to get issues by Rate limit. Slepping for %d ms. Resp: " + string(body))
			time.Sleep(t)
			continue
		}

		bar := pb.StartNew(len(issues))
		for _, issue := range issues {
			hook := hook{
				EventName: "jira:issue_imported",
				EventData: map[string]interface{}{"project_id": jiraProjectId, "project_key": jiraProjectKey},
				Data: IssueHook{
					Issue:        issue,
					Timestamp:    time.Now().Unix() * 1000,
					WebhookEvent: "jira:issue_imported",
					TypeName:     "issue_imported",
				},
			}
			j, err := json.Marshal(hook)
			if err != nil {
				fmt.Println("Failed to Marshal the hook data")
				continue
			}
			err = enqueue(string(j))
			if err != nil {
				fmt.Println(err)
			}
			bar.Increment()
		}
		bar.FinishPrint("Issues imported")
		op.StartAt = op.StartAt + op.MaxResults
		if len(issues) < op.MaxResults {
			return
		}
	}
}
