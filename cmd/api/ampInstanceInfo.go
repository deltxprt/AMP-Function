package main

import (
	"ampstatus-azfunction/internal/data"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func listInstances(url string, sessionId string) *data.Instances {
	listInstances := url + "/API/ADSModule/GetInstances"
	var list_Instances data.Instances

	data := map[string]string{"SESSIONID": sessionId}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	request, err := http.NewRequest("POST", listInstances, bytes.NewBuffer(jsonData))
	request.Header.Set("accept", "application/json; charset=UTF-8")

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		log.Fatal(error)
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	err = json.Unmarshal([]byte(body), &list_Instances)
	if err != nil {
		log.Fatal(err)
	}
	return &list_Instances
}

func statusInstances(url string, sessionId string, instanceID data.Instances) *[]data.Status {
	statPerInstance := url + "/API/ADSModule/GetInstance"
	var allinstancesStatus []data.Status
	for _, instance := range instanceID.Result[0].AvailableInstances {
		var getStatus data.Status
		data := map[string]string{"SESSIONID": sessionId, "InstanceId": instance.InstanceID}
		json_Data, err := json.Marshal(data)

		if err != nil {
			log.Fatal(err)
		}

		request, err := http.NewRequest("POST", statPerInstance, bytes.NewBuffer(json_Data))
		request.Header.Set("accept", "application/json; charset=UTF-8")

		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{}
		response, error := client.Do(request)

		if error != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		err = json.Unmarshal([]byte(body), &getStatus)
		if err != nil {
			log.Fatal(err)
		}
		allinstancesStatus = append(allinstancesStatus, getStatus)

	}
	return &allinstancesStatus
}
