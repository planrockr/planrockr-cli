package pkg

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	jiraApi "github.com/andygrunwald/go-jira"

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
	projects, err := GetProjects(i)
	if err != nil {
		log.Fatalf("[IMPORTER] Failed to get The project list: %v", err)
		return err
	}
	projects = selectProject(projects)

	var wg sync.WaitGroup
	//file for debug
	producer, err := os.Create("/tmp/dat2")
	// Start processing the importers.
	fmt.Println("Importing...")
	for _, proj := range projects {
		pID, err := strconv.Atoi(proj.ID)
		if err != nil {
			log.Errorf("[IMPORTER] Failed to convert projectID to integer: %v", err)
			continue
		}
		pName := proj.Name
		pKey := proj.Key
		err = createProject(pName)
		if err != nil {
			fmt.Println(err)
		}
		boardId := proj.ID + "_" + proj.Key
		err = createBoard(boardId, "3")
		if err != nil {
			fmt.Println(err)
		}
		wg.Add(1)
		wLog := log.StandardLogger().WriterLevel(log.ErrorLevel)
		defer wLog.Close()
		go func() {
			Process(i, pID, pKey, producer, wLog)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Finished importing projects")

	err = createHook(i.url, i.user, i.pass)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Webhook created")

	return nil
}

func createHook(host string, user string, password string) error {
	// Disable HTTP/2
	http.DefaultClient.Transport = &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	type Payload struct {
		Name                string   `json:"name"`
		URL                 string   `json:"url"`
		Events              []string `json:"events"`
		ExcludeIssueDetails bool     `json:"excludeIssueDetails"`
	}

	data := Payload{
		Name:                "Planrockr",
		URL:                 "https://app.planrockr.com/hook/jira/" + strconv.Itoa(configData.Auth.Id) + "/${project.id}/${project.key}",
		Events:              []string{"jira:issue_created", "jira:issue_updated", "worklog_created", "worklog_updated", "worklog_deleted", "comment_created", "comment_updated", "comment_deleted", "project_deleted", "project_updated", "jira:issue_deleted", "project_created", "jira:worklog_updated"},
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

func createProject(name string) error {
	body := strings.NewReader("parameters%5Bname%5D=" + name)
	req, err := http.NewRequest("POST", configData.BaseUrl+"/rpc/v1/project/create", body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "planrockr-cli")
	req.Header.Set("Authorization-Coderockr", configData.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized {
			panic("You must login")
		}
		return errors.New("Error creating project")
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &projectData)
	if err != nil {
		return errors.New("Error parsing project data")
	}

	return nil
}

func createBoard(boardId string, boardType string) error {
	q, err := url.ParseQuery("parameters[projectId]=" + strconv.Itoa(projectData.Project.Id) + "&parameters[boardId]=" + boardId + "&parameters[boardType]=" + boardType)
	if err != nil {
		return err
	}
	body := strings.NewReader(q.Encode())
	req, err := http.NewRequest("POST", configData.BaseUrl+"/rpc/v1/project/addBoard", body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "planrockr-cli")
	req.Header.Set("Authorization-Coderockr", configData.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusCreated {
		fmt.Println(resp.StatusCode)
		return errors.New("Error creating board")
	}

	return nil
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
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		err = errors.New("Failed to get the list of projects. Resp: " + string(body))
	}
	return *projects, err
}

func selectProject(list jiraApi.ProjectList) jiraApi.ProjectList {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Projects available:")
	for i, op := range list {
		fmt.Printf("\t%d - %s(%s)\n", i, op.Name, op.Key)
	}
	fmt.Print("Select options(space to separate): ")
	str, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	selected := strings.Split(string(str), " ")
	res := make(jiraApi.ProjectList, 0, len(list))
	for _, s := range selected {
		i, err := strconv.Atoi(s)
		if err != nil || i >= len(list) || i < 0 {
			log.Fatal("Option invalid", err)
		}
		res = append(res, list[i])
	}
	return res
}

func Process(i JiraImporter, jiraProjectId int, jiraProjectKey string, w io.WriteCloser, wErr io.Writer) {
	op := jiraApi.SearchOptions{
		StartAt:    0,
		MaxResults: 50,
	}

	for {
		issues, resp, err := i.client.Issue.Search(fmt.Sprintf("project=%d", jiraProjectId), &op)
		if err != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode != 429 {
				wErr.Write([]byte(errors.New("Failed to get issues. Resp:" + string(body)).Error()))
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
			wErr.Write([]byte(errors.New("Failed to get issues by Rate limit. Slepping for %d ms. Resp: " + string(body)).Error()))
			time.Sleep(t)
			continue
		}

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
				wErr.Write([]byte(errors.New("Failed to Marshal the hook data").Error()))
				continue
			}
			err = enqueue(string(j))
			if err != nil {
				fmt.Println(err)
			}
		}
		op.StartAt = op.StartAt + op.MaxResults
		if len(issues) < op.MaxResults {
			return
		}
	}
}
