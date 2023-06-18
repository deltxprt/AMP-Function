package main

import (
	"amp-management-api/internal/data"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func (app *application) updateInstancesHandler() {
	url := app.config.AMP.Url + "/API/ADSModule/GetInstances"
	app.ampLogin()
	sessionId := app.config.AMP.SessionId
	var listInstances data.InstancesData

	dataBody := map[string]string{"SESSIONID": sessionId}
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(dataBody)

	if err != nil {
		app.logger.PrintError(err, map[string]string{"detail": "error encoding dataBody"})
	}

	request, err := http.NewRequest("POST", url, buffer)
	request.Header.Set("accept", "application/json; charset=UTF-8")

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.logger.PrintError(err, map[string]string{"detail": "error doing request"})
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			app.logger.PrintError(err, map[string]string{"detail": "error closing response body"})
		}
	}(response.Body)

	body, _ := io.ReadAll(response.Body)

	err = json.Unmarshal([]byte(body), &listInstances)
	if err != nil {
		app.logger.PrintError(err, map[string]string{"detail": "error unmarshalling response body"})
	}

	for i := 0; i < len(listInstances.Result); i++ {

		hostInstances := listInstances.Result[i].AvailableInstances

		for j := 0; j < len(hostInstances); j++ {

			var instanceInformation data.InstanceStatus
			instanceInformation = hostInstances[j].InstanceStatus

			if instanceInformation.InstanceName == "ADS01" {
				continue
			}
			_, err := app.dbmodels.Instance.Get(instanceInformation.FriendlyName)
			if err != nil {

				err := app.dbmodels.Instance.Create(instanceInformation)
				if err != nil {
					log.Printf("Error creating instance with error: %s", err)
				}
				//app.logger.PrintInfo("validating ttl presence on "+instanceData.FriendlyName, nil)
				ttlIsPresent, err := app.rdbmodels.Instance.GetTTL(instanceInformation.InstanceName)
				if instanceInformation.Metrics.ActiveUsers.RawValue > 0 && ttlIsPresent == false {
					//app.logger.PrintInfo("Setting TTL for instance"+instanceData.FriendlyName, nil)
					expireTime, err := time.ParseDuration("2h")
					instanceTTL := data.InstanceTTL{
						InstanceName:         instanceInformation.InstanceName,
						InstanceFriendlyName: instanceInformation.FriendlyName,
						TTL:                  expireTime,
					}
					err = app.rdbmodels.Instance.SetTTL(instanceTTL)

					if err != nil {
						app.logger.PrintError(err, nil)
					}
				} else if instanceInformation.Metrics.ActiveUsers.RawValue > 0 && ttlIsPresent == true {
					//app.logger.PrintInfo("Updating TTL for instance"+instanceData.FriendlyName, nil)
					err := app.rdbmodels.Instance.UpdateTTL(instanceInformation.InstanceID, 120)
					if err != nil {
						app.logger.PrintError(err, nil)
					}
				} else if instanceInformation.Metrics.ActiveUsers.RawValue == 0 && ttlIsPresent == false {
					// stop the instance
					app.logger.PrintInfo("[system]initiated a stop on "+instanceInformation.FriendlyName, nil)
					stopUrl := app.config.AMP.Url + "/API/ADSModule/StopInstance"
					_, err := app.dbmodels.Instance.InstanceAction(stopUrl, instanceInformation.InstanceName, sessionId)
					if err != nil {
						app.logger.PrintError(err, map[string]string{"detail": "error stopping instance"})
					}
				}
			} else {
				//app.logger.PrintInfo("Updating instance"+instanceData.FriendlyName, nil)
				err := app.dbmodels.Instance.Update(instanceInformation)
				if err != nil {
					log.Println("Error updating instance")
					log.Println(err)
				}

				//app.logger.PrintInfo("validating ttl presence on "+instanceData.FriendlyName, nil)
				ttlIsPresent, err := app.rdbmodels.Instance.GetTTL(instanceInformation.InstanceName)

				if instanceInformation.Metrics.ActiveUsers.RawValue > 0 && ttlIsPresent == false {
					//app.logger.PrintInfo("Setting TTL for instance"+instanceData.FriendlyName, nil)
					expireTime, err := time.ParseDuration("2h")
					instanceTTL := data.InstanceTTL{
						InstanceName:         instanceInformation.InstanceName,
						InstanceFriendlyName: instanceInformation.FriendlyName,
						TTL:                  expireTime,
					}
					err = app.rdbmodels.Instance.SetTTL(instanceTTL)

					if err != nil {
						app.logger.PrintError(err, nil)
					}
				} else if instanceInformation.Metrics.ActiveUsers.RawValue > 0 && ttlIsPresent == true {
					//app.logger.PrintInfo("Updating TTL for instance"+instanceData.FriendlyName, nil)
					expireTime, err := time.ParseDuration("2h")
					if err != nil {
						app.logger.PrintError(err, nil)
					}
					err = app.rdbmodels.Instance.UpdateTTL(instanceInformation.InstanceID, expireTime)
					if err != nil {
						app.logger.PrintError(err, nil)
					}
				} else if instanceInformation.Metrics.ActiveUsers.RawValue == 0 && ttlIsPresent == false && instanceInformation.Running == true {
					// stop the instance
					stopUrl := app.config.AMP.Url + "/API/ADSModule/StopInstance"
					//app.logger.PrintInfo(fmt.Sprintf("[SYSTEM] name: %s | players: %d | ttl: %t | running: %t", instanceInformation.InstanceName, instanceInformation.Metrics.ActiveUsers.RawValue, ttlIsPresent, instanceInformation.Running), nil)
					app.logger.PrintInfo("[system]initiated a stop on "+instanceInformation.FriendlyName, nil)
					_, err := app.dbmodels.Instance.InstanceAction(stopUrl, instanceInformation.InstanceName, sessionId)
					if err != nil {
						app.logger.PrintError(err, map[string]string{"detail": "error stopping instance"})
					}
				}
			}

			//log.Println(instanceID, instanceFriendlyName, instanceName)
		}
	}
}

func (app *application) listInstancesHandler(w http.ResponseWriter, r *http.Request) {

	instances, err := app.dbmodels.Instance.GetAll()
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Content": instances, "StatusCode": 0}, nil)
}

func (app *application) ListInstanceHandler(w http.ResponseWriter, r *http.Request) {

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

// a function that will start the instance based on the friendly name

func (app *application) actionInstanceHandler(w http.ResponseWriter, r *http.Request) {

	action, err := app.readActionParam(r)
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	type instanceJson struct {
		InstanceName string `json:"InstanceName"`
	}
	var instance instanceJson

	err = app.readJSON(w, r, &instance)

	if err != nil {
		app.logger.PrintError(err, map[string]string{"detail": "error reading instance name in json"})
	}

	instanceName := instance.InstanceName

	instanceInformation, err := app.dbmodels.Instance.Get(instanceName)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	expireTime, err := time.ParseDuration("2h")

	switch action {
	case "Start":
		err := app.rdbmodels.Instance.SetTTL(data.InstanceTTL{
			InstanceName:         instanceInformation.InstanceName,
			InstanceFriendlyName: instanceInformation.FriendlyName,
			TTL:                  expireTime,
		})
		if err != nil {
			app.logger.PrintError(err, map[string]string{"detail": "error setting ttl on instance " + instanceInformation.FriendlyName})
		}
	case "Stop":
		err := app.rdbmodels.Instance.DeleteTTL(instanceInformation.InstanceName)
		if err != nil {
			app.logger.PrintError(err, map[string]string{"detail": "error deleting ttl on instance " + instanceInformation.FriendlyName})
		}
	default:
		app.logger.PrintError(err, map[string]string{"detail": "error reading instance name in json"})
	}

	app.ampLogin()
	sessionId := app.config.AMP.SessionId

	url := app.config.AMP.Url + "/API/ADSModule/" + action + "Instance"
	app.logger.PrintInfo("[user]initiated a "+action+" on "+instanceName, nil)

	actionInstanceResponse, err := app.dbmodels.Instance.InstanceAction(url, instanceInformation.InstanceName, sessionId)

	if err != nil {
		app.logger.PrintError(err, nil)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"InstanceName": instanceName, "action": action, "Status": actionInstanceResponse.Result.Status}, nil)

}
