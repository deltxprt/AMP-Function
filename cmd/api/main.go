package main

import (
	"amp-management-api/internal/data"
	"amp-management-api/internal/jsonlog"
	"amp-management-api/internal/vcs"
	"context"
	"expvar"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"time"
)

var (
	version = vcs.Version()
)

type Config struct {
	Port     int    `yaml:"port"`
	Env      string `yaml:"env"`
	Database struct {
		Address         string `yaml:"address"`
		Password        string `yaml:"password"`
		Database        int    `yaml:"database"`
		MaxIdleConns    int    `yaml:"maxIdleConns"`
		ConnMaxIdleTime int    `yaml:"connMaxIdleTime"`
	} `yaml:"database"`
	AMP struct {
		Url        string `yaml:"url"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Token      string `yaml:"token"`
		RememberMe bool   `yaml:"rememberMe"`
		SessionId  string `yaml:"sessionId"`
	} `yaml:"amp"`
}

type application struct {
	config Config
	logger *jsonlog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	var cfg Config
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	// Load the configuration settings from the config.yml file.
	configFile, err := os.ReadFile("/config.yaml")
	if err != nil {
		logger.PrintError(err, nil)
	}

	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		logger.PrintError(err, nil)
	}

	logger.PrintInfo(cfg.Env, nil)

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	//cfg.Database.Database = 0
	cfg.Database.MaxIdleConns = 10
	cfg.Database.ConnMaxIdleTime = 5

	rdb, err := openDB(cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer rdb.Close()

	expvar.NewString("version").Set(version)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(rdb),
	}
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Database.Address,
		Password:     cfg.Database.Password,
		DB:           cfg.Database.Database,
		MaxIdleConns: cfg.Database.MaxIdleConns,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
