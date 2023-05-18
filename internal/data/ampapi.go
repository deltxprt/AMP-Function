package data

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
	InstanceID   string `json:"InstanceID,omitempty"`
	InstanceName string `json:"InstanceName"`
	FriendlyName string `json:"FriendlyName"`
	Module       string `json:"Module"`
	Running      bool   `json:"Running"`
	Suspended    bool   `json:"Suspended"`
	Metrics      struct {
		CPUUsage struct {
			RawValue uint8 `json:"RawValue"`
			MaxValue uint8 `json:"MaxValue"`
			Percent  uint8 `json:"Percent"`
		} `json:"CPU Usage"`
		MemoryUsage struct {
			RawValue uint16 `json:"RawValue"`
			MaxValue uint16 `json:"MaxValue"`
			Percent  uint8  `json:"Percent"`
		} `json:"Memory Usage"`
		ActiveUsers struct {
			RawValue uint8 `json:"RawValue"`
			MaxValue uint8 `json:"MaxValue"`
		} `json:"Active Users"`
	} `json:"Metrics"`
}

type TaskActionResult struct {
	Result struct {
		Status bool `json:"Status"`
	} `json:"result"`
}

type InstancesModel struct {
	DB *sql.DB
}

func (m InstancesModel) Create(instance InstanceStatus) error {

	query := `
	INSERT INTO instances (instance_id, instance_name, name, module, running, suspended, cpu_usage, cpu_max, cpu_percent, memory_usage, memory_max, memory_percent, users_active, users_max)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id, name, running`

	args := []any{instance.InstanceID, instance.InstanceName, instance.FriendlyName, instance.Module, instance.Running, instance.Suspended, instance.Metrics.CPUUsage.RawValue, instance.Metrics.CPUUsage.MaxValue, instance.Metrics.CPUUsage.Percent, instance.Metrics.MemoryUsage.RawValue, instance.Metrics.MemoryUsage.MaxValue, instance.Metrics.MemoryUsage.Percent, instance.Metrics.ActiveUsers.RawValue, instance.Metrics.ActiveUsers.MaxValue}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&instance.InstanceID, &instance.FriendlyName, &instance.Running)
}

func (m InstancesModel) Get(instanceName string) (*InstanceStatus, error) {

	query := `SELECT instance_name, name, module, running, cpu_usage, cpu_max, cpu_percent, memory_usage, memory_max, memory_percent, users_active, users_max FROM instances WHERE name = $1`

	var instance InstanceStatus

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, instanceName).Scan(
		&instance.InstanceName,
		&instance.FriendlyName,
		&instance.Module,
		&instance.Running,
		&instance.Metrics.CPUUsage.RawValue,
		&instance.Metrics.CPUUsage.MaxValue,
		&instance.Metrics.CPUUsage.Percent,
		&instance.Metrics.MemoryUsage.RawValue,
		&instance.Metrics.MemoryUsage.MaxValue,
		&instance.Metrics.MemoryUsage.Percent,
		&instance.Metrics.ActiveUsers.RawValue,
		&instance.Metrics.ActiveUsers.MaxValue,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &instance, nil
}

func (m InstancesModel) Update(instances InstanceStatus) error {

	query := `
	UPDATE instances
	SET running = $1, suspended = $2, cpu_usage = $3, cpu_max = $4, cpu_percent = $5, memory_usage = $6, memory_max = $7, memory_percent = $8, users_active = $9, users_max = $10
	WHERE name = $11
	RETURNING id, name, running`
	args := []any{
		instances.Running,
		instances.Suspended,
		instances.Metrics.CPUUsage.RawValue,
		instances.Metrics.CPUUsage.MaxValue,
		instances.Metrics.CPUUsage.Percent,
		instances.Metrics.MemoryUsage.RawValue,
		instances.Metrics.MemoryUsage.MaxValue,
		instances.Metrics.MemoryUsage.Percent,
		instances.Metrics.ActiveUsers.RawValue,
		instances.Metrics.ActiveUsers.MaxValue,
		instances.FriendlyName,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&instances.InstanceID, &instances.FriendlyName, &instances.Running)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m InstancesModel) GetAll() ([]InstanceStatus, error) {
	query := `SELECT name, module, running, cpu_usage, cpu_max, cpu_percent, memory_usage, memory_max, memory_percent, users_active, users_max FROM instances`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var instances []InstanceStatus

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var instance InstanceStatus
		err := rows.Scan(
			&instance.FriendlyName,
			&instance.Module,
			&instance.Running,
			&instance.Metrics.CPUUsage.RawValue,
			&instance.Metrics.CPUUsage.MaxValue,
			&instance.Metrics.CPUUsage.Percent,
			&instance.Metrics.MemoryUsage.RawValue,
			&instance.Metrics.MemoryUsage.MaxValue,
			&instance.Metrics.MemoryUsage.Percent,
			&instance.Metrics.ActiveUsers.RawValue,
			&instance.Metrics.ActiveUsers.MaxValue,
		)

		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

// function that will do api calls to stop an instance
func (m InstancesModel) InstanceAction(url string, instanceName string, token string) (*TaskActionResult, error) {
	var InstanceActionResult TaskActionResult
	dataBody := map[string]string{"SESSIONID": token, "InstanceName": instanceName}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(dataBody)

	request, err := http.NewRequest("POST", url, buffer)
	request.Header.Set("accept", "application/json; charset=UTF-8")

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	err = json.Unmarshal([]byte(body), &InstanceActionResult)
	if err != nil {
		return nil, err
	}
	return &InstanceActionResult, nil
}
