package main

import (
	"ampstatus-azfunction/internal/data"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Response struct {
	Content *[]data.Status
}

type config struct {
	port int
	// env  string
	// db   struct {
	// 	dsn          string
	// 	maxOpenConns int
	// 	maxIdleConns int
	// 	maxIdleTime  string
	// }
}

type application struct {
	config config
	logger *log.Logger
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
	message := Response{
		Content: StatusInstance,
	}
	jsonResponse, err := json.Marshal(message)
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
	var cfg config

	listenAddr := 8080
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal(err)
		}
		listenAddr = port
	}
	flag.IntVar(&cfg.port, "port", listenAddr, "API server port")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting server on %s", srv.Addr)

	srv.ListenAndServe()
	// logger.Fatal(err)
}
