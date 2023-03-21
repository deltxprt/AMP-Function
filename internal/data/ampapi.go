package data

type Instances struct {
	Result []struct {
		AvailableInstances []struct {
			InstanceID   string `json:"InstanceID"`
			InstanceName string `json:"InstanceName"`
			FriendlyName string `json:"FriendlyName"`
		} `json:"AvailableInstances"`
	} `json:"result"`
}

type Status struct {
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
