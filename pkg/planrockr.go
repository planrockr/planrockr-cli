package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/planrockr/planrockr-cli/config"
)

var configData config.Config

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
