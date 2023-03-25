package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type sessionIDStruct struct {
	sessionId string `json:"SESSIONID"`
}

func (app *application) ampLogin() {
	loginUrl := app.config.AMP.Url + "/API/Core/Login"

	data := map[string]string{
		"username":   app.config.AMP.Username,
		"password":   app.config.AMP.Password,
		"token":      app.config.AMP.Token,
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

	body, _ := io.ReadAll(response.Body)
	var sessionId map[string]interface{}
	json.Unmarshal(body, &sessionId)

	sessionIdString := sessionId["sessionID"].(string)

	app.config.AMP.SessionId = sessionIdString
}
