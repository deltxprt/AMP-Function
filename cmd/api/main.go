package main

import (
	"ampstatus-azfunction/internal/data"
	"ampstatus-azfunction/internal/jsonlog"
	"context"
	"flag"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type config struct {
	port int
	db   struct {
		address         string
		password        string
		DB              int
		MaxIdleConns    int
		ConnMaxIdleTime int
	}
	amp struct {
		url        string
		username   string
		password   string
		token      string
		rememberMe bool
		sessionId  string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if os.Getenv("AMPURL") == "" || os.Getenv("AMPUser") == "" || os.Getenv("AMPPass") == "" {
		log.Fatal("Please set the environment variables")
	}

	cfg.amp.url = os.Getenv("AMPURL")
	cfg.amp.username = os.Getenv("AMPUser")
	cfg.amp.password = os.Getenv("AMPPass")

	cfg.db.address = os.Getenv("REDISADDR")
	cfg.db.password = os.Getenv("REDISPW")
	cfg.db.DB = 0
	cfg.db.MaxIdleConns = 10
	cfg.db.ConnMaxIdleTime = 5

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

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	log.Println("before openDB")
	rdb, err := openDB(cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer rdb.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(rdb),
	}
	log.Println("before serve")
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.db.address,
		//		Password: cfg.db.password,
		DB: cfg.db.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
