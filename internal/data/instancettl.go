package data

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type InstanceTTL struct {
	InstanceName         string        `json:"InstanceName"`
	InstanceFriendlyName string        `json:"InstanceFriendlyName"`
	TTL                  time.Duration `json:"ttl"`
}

type InstanceTTLModel struct {
	DB *redis.Client
}

func (m *InstanceTTLModel) SetTTL(instance InstanceTTL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.Set(ctx, instance.InstanceName, instance.InstanceFriendlyName, instance.TTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (m *InstanceTTLModel) GetTTL(instanceName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//var instance InstanceTTL
	_, err := m.DB.Get(ctx, instanceName).Result()

	//err := json.Unmarshal([]byte(status), &instance)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *InstanceTTLModel) UpdateTTL(instanceName string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var instance InstanceTTL
	status := m.DB.Get(ctx, instanceName).Val()

	err := json.Unmarshal([]byte(status), &instance)
	if err != nil {
		return err
	}

	err = m.DB.Expire(ctx, instanceName, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (m *InstanceTTLModel) DeleteTTL(instanceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.Del(ctx, instanceName).Err()

	if err != nil {
		return err
	}

	return nil
}
