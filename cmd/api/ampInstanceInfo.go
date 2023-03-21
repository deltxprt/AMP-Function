package main

import (
	"ampstatus-azfunction/internal/data"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"net/http"
)

var ctx = context.Background()

func listInstances(url string, sessionId string) {
	listInstances := url + "/API/ADSModule/GetInstances"
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:32768",
		Password: "redispw",
		DB:       0,
	})
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
	for i := 0; i < len(list_Instances.Result[0].AvailableInstances); i++ {
		instancesInformation := list_Instances.Result[0].AvailableInstances[i]
		instanceID := instancesInformation.InstanceID
		instanceName := instancesInformation.FriendlyName
		rdb.Set(ctx, instanceName, instanceID, 0)
	}
}

func statusInstances(url string, sessionId string) *[]data.Status {
	statPerInstance := url + "/API/ADSModule/GetInstance"
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:32768",
		Password: "redispw",
		DB:       0,
	})
	iter := rdb.Scan(ctx, 0, "*", 0).Iterator()

	var alliancesStatus []data.Status
	for iter.Next(ctx) {
		instanceID, _ := rdb.Get(ctx, iter.Val()).Result()
		var getStatus data.Status
		instanceData := map[string]string{"SESSIONID": sessionId, "InstanceId": instanceID}
		jsonData, err := json.Marshal(instanceData)

		if err != nil {
			log.Fatal(err)
		}

		request, err := http.NewRequest("POST", statPerInstance, bytes.NewBuffer(jsonData))
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
		alliancesStatus = append(alliancesStatus, getStatus)

	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	log.Println("This API path has been called")
	return &alliancesStatus
}
