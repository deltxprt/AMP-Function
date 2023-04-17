package data

import (
	"database/sql"
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type RDBModels struct {
	Instance InstanceTTLModel
}

type DBModels struct {
	Instance InstancesModel
}

func NewModels(rdb *redis.Client) RDBModels {
	return RDBModels{
		Instance: InstanceTTLModel{DB: rdb},
	}
}

func NewDBModels(db *sql.DB) DBModels {
	return DBModels{
		Instance: InstancesModel{DB: db},
	}
}
