package main

import (
	"amp-management-api/internal/data"
	"amp-management-api/internal/jsonlog"
	"amp-management-api/internal/vcs"
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

// Config
// DSN example: "postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable"
type Config struct {
	Port            int    `yaml:"port"`
	Env             string `yaml:"env"`
	RefreshInterval string `yaml:"RefreshInterval"`
	//Redis           struct {
	//	Address         string `yaml:"address"`
	//	Password        string `yaml:"password"`
	//	Database        int    `yaml:"database"`
	//	MaxIdleConns    int    `yaml:"maxIdleConns"`
	//	ConnMaxIdleTime int    `yaml:"connMaxIdleTime"`
	//} `yaml:"redis"`
	Postgres struct {
		Dsn          string `yaml:"dsn"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxIdleTime  string `yaml:"maxIdleTime"`
	} `yaml:"postgres"`
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
	//	rdbmodels data.RDBModels
	dbmodels data.DBModels
	wg       sync.WaitGroup
}

func main() {
	var cfg Config
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	displayVersion := flag.Bool("version", false, "Display version and exit")
	configFilePath := flag.String("config", "./examples/config.yaml", "Path to the configuration file")

	flag.Parse()

	// Load the configuration settings from the config.yml file.

	cfg, err := getConfig(*configFilePath)

	if err != nil {
		logger.PrintFatal(err, map[string]string{"details": "Failed to load the configuration file."})
	}

	logger.PrintInfo(cfg.Env, nil)

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	//cfg.Database.Database = 0
	//cfg.Redis.MaxIdleConns = 10
	//cfg.Redis.ConnMaxIdleTime = 5

	// dropping the auto closing feature, might add it later
	//rdb, err := openRedis(cfg)

	//if err != nil {
	//	logger.PrintFatal(err, nil)
	//}

	//defer rdb.Close()

	db, err := openDB(cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	expvar.NewString("version").Set(version)

	app := &application{
		config: cfg,
		logger: logger,
		//rdbmodels: data.NewModels(rdb),
		dbmodels: data.NewDBModels(db),
	}

	go updateInstance(app)

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}

// dropping the auto closing feature, might add it later
//func openRedis(cfg Config) (*redis.Client, error) {
//	rdb := redis.NewClient(&redis.Options{
//		Addr:         cfg.Redis.Address,
//		Password:     cfg.Redis.Password,
//		DB:           cfg.Redis.Database,
//		MaxIdleConns: cfg.Redis.MaxIdleConns,
//	})
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	err := rdb.Ping(ctx).Err()
//	if err != nil {
//		return nil, err
//	}
//
//	return rdb, nil
//}

func openDB(cfg Config) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg.Postgres.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.Postgres.MaxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func updateInstance(app *application) {
	RefreshDuration, err := time.ParseDuration(app.config.RefreshInterval)
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}

	for _ = range time.Tick(RefreshDuration) {
		app.updateInstancesHandler()
	}
}

func getConfig(path string) (Config, error) {
	var cfg Config
	var err error
	_, err = os.Stat(path)
	if err != nil {
		log.Println("Config file not found, validating environment variables")
	} else {
		configFile, err := os.ReadFile(path)
		if err != nil {
			log.Println("file not found")
		}

		err = yaml.Unmarshal(configFile, &cfg)
		if err != nil {
			log.Println(err)
		}
		return cfg, nil
	}

	log.Println("searching environment variables")

	cfg.Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		cfg.Port = 8080
	}

	cfg.AMP.Url = os.Getenv("AMP_URL")
	if cfg.AMP.Url == "" {
		return cfg, err
	}
	cfg.AMP.Username = os.Getenv("AMP_USERNAME")
	if cfg.AMP.Username == "" {
		return cfg, err
	}
	cfg.AMP.Password = os.Getenv("AMP_PASSWORD")
	if cfg.AMP.Password == "" {
		return cfg, err
	}
	cfg.AMP.Token = os.Getenv("AMP_TOKEN")

	cfg.AMP.RememberMe, err = strconv.ParseBool(os.Getenv("AMP_REMEMBER_ME"))
	if err != nil {
		cfg.AMP.RememberMe = false
	}

	cfg.Postgres.Dsn = os.Getenv("POSTGRES_DSN")

	cfg.Postgres.MaxOpenConns, err = strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS"))
	if err != nil {
		cfg.Postgres.MaxOpenConns = 5
	}

	cfg.Postgres.MaxIdleConns, err = strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS"))
	if err != nil {
		cfg.Postgres.MaxIdleConns = 5
	}

	cfg.Postgres.MaxIdleTime = os.Getenv("POSTGRES_MAX_IDLE_TIME")
	cfg.RefreshInterval = os.Getenv("REFRESH_INTERVAL")

	return cfg, nil
}
