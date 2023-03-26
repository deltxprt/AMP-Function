package main

import (
	"amp-management-api/internal/data"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (app *application) updateInstancesHandler() {
	listInstances := app.config.AMP.Url + "/API/ADSModule/GetInstances"
	sessionId := app.config.AMP.SessionId
	var list_Instances data.InstancesData

	dataBody := map[string]string{"SESSIONID": sessionId}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(dataBody)

	request, err := http.NewRequest("POST", listInstances, buffer)
	request.Header.Set("accept", "application/json; charset=UTF-8")

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	err = json.Unmarshal([]byte(body), &list_Instances)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(list_Instances.Result); i++ {
		hostInstances := list_Instances.Result[i].AvailableInstances
		for j := 0; j < len(hostInstances); j++ {
			var instanceInformation data.InstanceStatus
			instanceInformation = hostInstances[j].InstanceStatus
			//log.Println(instanceID, instanceFriendlyName, instanceName)
			err := app.models.Instance.Update(instanceInformation)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (app *application) listInstancesHandler(w http.ResponseWriter, r *http.Request) {
	app.ampLogin()
	go app.updateInstancesHandler()

	instances, err := app.models.Instance.GetAll()
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Content": instances, "StatusCode": 0}, nil)
}

func (app *application) ListInstanceHandler(w http.ResponseWriter, r *http.Request) {

	app.ampLogin()
	go app.updateInstancesHandler()

	instance, err := app.readInstanceParam(r)
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	instances, err := app.models.Instance.Get(instance)
	if instances.FriendlyName == "" || err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Content": instances, "StatusCode": 0}, nil)

}
