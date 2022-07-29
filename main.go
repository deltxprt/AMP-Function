package main

import (
	"bytes"
	"encoding/json"

	//	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Instances struct {
	Result []struct {
		AvailableInstances []struct {
			InstanceID       string `json:"InstanceID"`
			InstanceName     string `json:"InstanceName"`
			FriendlyName     string `json:"FriendlyName"`
			Module           string `json:"Module"`
			InstalledVersion struct {
				Major         int `json:"Major"`
				Minor         int `json:"Minor"`
				Build         int `json:"Build"`
				Revision      int `json:"Revision"`
				MajorRevision int `json:"MajorRevision"`
				MinorRevision int `json:"MinorRevision"`
			} `json:"InstalledVersion"`
			Running   bool `json:"Running"`
			Suspended bool `json:"Suspended"`
		} `json:"AvailableInstances"`
	} `json:"result"`
}

type Status struct {
	InstanceID   string `json:"InstanceID"`
	FriendlyName string `json:"FriendlyName"`
	Module       string `json:"Module"`
	Running      bool   `json:"Running"`
	Suspended    bool   `json:"Suspended"`
	Metrics      struct {
		CPUUsage struct {
			RawValue int    `json:"RawValue"`
			MaxValue int    `json:"MaxValue"`
			Percent  int    `json:"Percent"`
			Units    string `json:"Units"`
		} `json:"CPU Usage"`
		MemoryUsage struct {
			RawValue int    `json:"RawValue"`
			MaxValue int    `json:"MaxValue"`
			Percent  int    `json:"Percent"`
			Units    string `json:"Units"`
		} `json:"Memory Usage"`
		ActiveUsers struct {
			RawValue int    `json:"RawValue"`
			MaxValue int    `json:"MaxValue"`
			Percent  int    `json:"Percent"`
			Units    string `json:"Units"`
		} `json:"Active Users"`
	} `json:"Metrics"`
}

type Response struct {
	Body *[]Status
}

func ampLogin(url, user, pass string) string {
	loginUrl := url + "/API/Core/Login"

	data := map[string]string{"username": user, "password": pass, "token": "", "rememberMe": "false"}
	json_Data, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(json_Data))
	request.Header.Set("accept", "application/json; charset=UTF-8")

	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	var sessionId map[string]interface{}
	json.Unmarshal(body, &sessionId)

	return sessionId["sessionID"].(string)
}

func listInstances(url, sessionId string) *Instances {
	listInstances := url + "/API/ADSModule/GetInstances"

	data := map[string]string{"SESSIONID": sessionId}
	json_Data, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", listInstances, bytes.NewBuffer(json_Data))
	request.Header.Set("accept", "application/json; charset=UTF-8")

	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		panic(error)
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	var list_Instances Instances
	err = json.Unmarshal([]byte(body), &list_Instances)
	if err != nil {
		panic(err)
	}
	return &list_Instances
}

func statusInstances(url, sessionId string, instanceID Instances) *[]Status {
	statPerInstance := url + "/API/ADSModule/GetInstance"
	var allinstancesStatus []Status
	for _, instance := range instanceID.Result[0].AvailableInstances {
		data := map[string]string{"SESSIONID": sessionId, "InstanceId": instance.InstanceID}
		json_Data, err := json.Marshal(data)

		if err != nil {
			log.Fatal(err)
		}

		request, err := http.NewRequest("POST", statPerInstance, bytes.NewBuffer(json_Data))
		request.Header.Set("accept", "application/json; charset=UTF-8")

		if err != nil {
			panic(err)
		}

		client := &http.Client{}
		response, error := client.Do(request)

		if error != nil {
			panic(error)
		}

		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)
		var getStatus Status
		err = json.Unmarshal([]byte(body), &getStatus)
		if err != nil {
			panic(err)
		}
		allinstancesStatus = append(allinstancesStatus, getStatus)

	}
	return &allinstancesStatus
}

func Main() *Response {
	ampUrl := os.Getenv("AMPUrl")
	ampUser := os.Getenv("AMPUser")
	ampPass := os.Getenv("AMPPass")
	sessionId := ampLogin(ampUrl, ampUser, ampPass)
	allInstances := listInstances(ampUrl, sessionId)
	StatusInstance := statusInstances(ampUrl, sessionId, *allInstances)

	return &Response{
		Body: StatusInstance,
	}
}
