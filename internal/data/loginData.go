package data

type Login struct {
	username   string `json:"username"`
	password   string `json:"password"`
	token      string `json:"token"`
	rememberme bool   `json:"rememberMe"`
}

type Token struct {
	sessionid string `json:"SESSIONID"`
}
