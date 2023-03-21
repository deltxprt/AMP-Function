package main

import (
	"ampstatus-azfunction/internal/data"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
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
}

func main() {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
