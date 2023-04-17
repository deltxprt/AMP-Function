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
	url := app.config.AMP.Url + "/API/ADSModule/GetInstances"
	sessionId := app.config.AMP.SessionId
	var listInstances data.InstancesData

	dataBody := map[string]string{"SESSIONID": sessionId}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(dataBody)

	request, err := http.NewRequest("POST", url, buffer)
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

	err = json.Unmarshal([]byte(body), &listInstances)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(listInstances.Result); i++ {
		hostInstances := listInstances.Result[i].AvailableInstances
		for j := 0; j < len(hostInstances); j++ {
			var instanceInformation data.InstanceStatus
			instanceInformation = hostInstances[j].InstanceStatus
			_, err := app.dbmodels.Instance.Get(instanceInformation.FriendlyName)
			if err != nil {
				err := app.dbmodels.Instance.Create(instanceInformation)
				if err != nil {
					log.Printf("Error creating instance with error: %s", err)
				}
			} else {
				err := app.dbmodels.Instance.Update(instanceInformation)
				if err != nil {
					log.Println(err)
				}
			}
			//log.Println(instanceID, instanceFriendlyName, instanceName)
		}
	}
}

func (app *application) listInstancesHandler(w http.ResponseWriter, r *http.Request) {
	app.ampLogin()
	go app.updateInstancesHandler()

	instances, err := app.dbmodels.Instance.GetAll()
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

	instances, err := app.dbmodels.Instance.Get(instance)
	if instances.FriendlyName == "" || err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Content": instances, "StatusCode": 0}, nil)

}
