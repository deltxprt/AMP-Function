package data

import (
	"context"
	"database/sql"
	"errors"
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
	DB *sql.DB
}

func (m InstancesModel) Create(instance InstanceStatus) error {

	query := `
	INSERT INTO instances (instance_id, name, module, running, suspended, cpu_usage, cpu_max, cpu_percent, memory_usage, memory_max, memory_percent, users_active, users_max, users_percent)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id, name, running`

	args := []any{instance.InstanceID, instance.FriendlyName, instance.Module, instance.Running, instance.Suspended, instance.Metrics.CPUUsage.RawValue, instance.Metrics.CPUUsage.MaxValue, instance.Metrics.CPUUsage.Percent, instance.Metrics.MemoryUsage.RawValue, instance.Metrics.MemoryUsage.MaxValue, instance.Metrics.MemoryUsage.Percent, instance.Metrics.ActiveUsers.RawValue, instance.Metrics.ActiveUsers.MaxValue, instance.Metrics.ActiveUsers.Percent}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&instance.InstanceID, &instance.FriendlyName, &instance.Running)
}

func (m InstancesModel) Get(instanceName string) (*InstanceStatus, error) {

	query := `SELECT name, module, running, cpu_usage, cpu_max, cpu_percent, memory_usage, memory_max, memory_percent, users_active, users_max, users_percent FROM instances WHERE name = $1`

	var instance InstanceStatus

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, instanceName).Scan(
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
		&instance.Metrics.ActiveUsers.Percent,
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
	SET running = $1, suspended = $2, cpu_usage = $3, cpu_usage_max = $4, cpu_usage_percent = $5, memory_usage = $6, memory_usage_max = $7, memory_usage_percent = $8, active_users = $9, active_users_max = $10, active_users_percent = $11
	WHERE name = $12
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
		instances.Metrics.ActiveUsers.Percent,
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
	query := `SELECT name, module, running, cpu_usage, cpu_max, cpu_percent, memory_usage, memory_max, memory_percent, users_active, users_max, users_percent FROM instances`

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
			&instance.Metrics.ActiveUsers.Percent,
		)

		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}

	return instances, nil
}
