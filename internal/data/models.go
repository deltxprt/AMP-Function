package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// dropping the auto closing feature, might add it later
//type RDBModels struct {
//	Instance InstanceTTLModel
//}

type DBModels struct {
	Instance InstancesModel
}

// dropping the auto closing feature, might add it later
//func NewModels(rdb *redis.Client) RDBModels {
//	return RDBModels{
//		Instance: InstanceTTLModel{DB: rdb},
//	}
//}

func NewDBModels(db *sql.DB) DBModels {
	return DBModels{
		Instance: InstancesModel{DB: db},
	}
}
