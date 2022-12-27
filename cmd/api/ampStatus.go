package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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
	message := Response{
		Content: StatusInstance,
	}
	jsonResponse, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	return string(jsonResponse)
}
