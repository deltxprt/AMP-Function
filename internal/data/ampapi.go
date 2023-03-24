package data

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type InstancesData struct {
	Result []struct {
		ID           uint8  `json:"Id"`
		InstanceID   string `json:"InstanceId"`
		FriendlyName string `json:"FriendlyName"`
		Platform     struct {
			CPUInfo struct {
				Sockets      uint8  `json:"Sockets"`
				Cores        uint8  `json:"Cores"`
				Threads      uint8  `json:"Threads"`
				Vendor       string `json:"Vendor"`
				ModelName    string `json:"ModelName"`
				TotalCores   uint8  `json:"TotalCores"`
				TotalThreads uint8  `json:"TotalThreads"`
			} `json:"CPUInfo"`
		} `json:"Platform"`
		State              uint8                      `json:"State"`
		LastUpdated        string                     `json:"LastUpdated"`
		AvailableInstances []struct{ InstanceStatus } `json:"AvailableInstances"`
		AvailableIPs       []string                   `json:"AvailableIPs"`
	} `json:"result"`
}

type InstanceStatus struct {
	InstanceID   string `json:"InstanceID"`
	FriendlyName string `json:"FriendlyName"`
	Module       string `json:"Module"`
	Running      bool   `json:"Running"`
	Suspended    bool   `json:"Suspended"`
	Metrics      struct {
		CPUUsage struct {
			RawValue uint8  `json:"RawValue"`
			MaxValue uint8  `json:"MaxValue"`
			Percent  uint8  `json:"Percent"`
			Units    string `json:"Units"`
		} `json:"CPU Usage"`
		MemoryUsage struct {
			RawValue uint16 `json:"RawValue"`
			MaxValue uint16 `json:"MaxValue"`
			Percent  uint8  `json:"Percent"`
			Units    string `json:"Units"`
		} `json:"Memory Usage"`
		ActiveUsers struct {
			RawValue uint8 `json:"RawValue"`
			MaxValue uint8 `json:"MaxValue"`
			Percent  uint8 `json:"Percent"`
		} `json:"Active Users"`
	} `json:"Metrics"`
}

type InstancesModel struct {
	DB *redis.Client
}

func (i InstancesModel) Update(instances InstanceStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//	instanceInfo := map[string]string{
	//		"InstanceID":   instances.InstanceID,
	//		"InstanceName": instances.InstanceName,
	//	}
	//	log.Print(instanceInfo)
	jsonFormat, err := json.Marshal(instances)
	if err != nil {
		return err
	}
	return i.DB.Set(ctx, instances.FriendlyName, jsonFormat, 0).Err()
}

func (i InstancesModel) GetAll() ([]InstanceStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	iter := i.DB.Scan(ctx, 0, "*", 0).Iterator()

	var instances []InstanceStatus
	for iter.Next(ctx) {
		var instance InstanceStatus
		jsonResult := i.DB.Get(ctx, iter.Val()).Val()
		err := json.Unmarshal([]byte(jsonResult), &instance)

		if err != nil {
			return nil, ErrRecordNotFound
		}
		//log.Print(instance.InstanceID)
		instances = append(instances, instance)
	}

	return instances, nil
}

func (i InstancesModel) Get(instanceName string) (InstanceStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var instance InstanceStatus

	instanceData := i.DB.Get(ctx, instanceName).Val()
	err := json.Unmarshal([]byte(instanceData), &instance)
	if err != nil {
		return instance, err
	}

	return instance, nil
}
