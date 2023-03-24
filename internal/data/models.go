package data

import (
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Instance InstancesModel
}

func NewModels(rdb *redis.Client) Models {
	return Models{
		Instance: InstancesModel{DB: rdb},
	}
}
