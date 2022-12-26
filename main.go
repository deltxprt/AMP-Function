package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	//	"fmt"
	"log"
	"net/http"
	"os"
)

type sessionIDStruct struct {
	sessionId string `json:"SESSIONID"`
}

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
	Content *[]Status
}

func ampLogin(url, user, pass string) *sessionIDStruct {
	loginUrl := url + "/API/Core/Login"
	data := map[string]string{
		"username":   user,
		"password":   pass,
		"token":      "",
		"rememberMe": "false",
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(data)

	request, err := http.NewRequest("POST", loginUrl, payloadBuf)
	request.Header.Set("accept", "application/json; charset=UTF-8")

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		fmt.Println(error)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	var sessionId map[string]interface{}
	json.Unmarshal(body, &sessionId)

	sessionIdString := sessionId["sessionID"].(string)

	return &sessionIDStruct{sessionId: sessionIdString}
}

func listInstances(url string, sessionId string) *Instances {
	listInstances := url + "/API/ADSModule/GetInstances"

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

	body, _ := ioutil.ReadAll(response.Body)

	var list_Instances Instances
	err = json.Unmarshal([]byte(body), &list_Instances)
	if err != nil {
		log.Fatal(err)
	}
	return &list_Instances
}

func statusInstances(url string, sessionId string, instanceID Instances) *[]Status {
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
			log.Fatal(err)
		}

		client := &http.Client{}
		response, error := client.Do(request)

		if error != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)
		var getStatus Status
		err = json.Unmarshal([]byte(body), &getStatus)
		if err != nil {
			log.Fatal(err)
		}
		allinstancesStatus = append(allinstancesStatus, getStatus)

	}
	return &allinstancesStatus
}

func ampStatus() string {
	ampUrl := os.Getenv("AMPUrl")
	ampUser := os.Getenv("AMPUser")
	ampPass := os.Getenv("AMPPass")
	if ampUrl == "" || ampUser == "" || ampPass == "" {
		fmt.Println("Please set the environment variables")
	}
	sessionIdToken := ampLogin(ampUrl, ampUser, ampPass)
	allInstances := listInstances(ampUrl, sessionIdToken.sessionId)
	StatusInstance := statusInstances(ampUrl, sessionIdToken.sessionId, *allInstances)
	// return Response{
	// 	Content: StatusInstance,
	// }
	jsonResponse, err := json.Marshal(StatusInstance)
	if err != nil {
		fmt.Println(err)
	}
	return string(jsonResponse)
}

func ampInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	result := ampStatus()
	fmt.Fprint(w, result)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "OK")
}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/api/AMPStatus", ampInfoHandler)
	http.HandleFunc("/api/HealthCheck", healthCheckHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
