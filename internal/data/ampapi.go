package data

type Instances struct {
	Result []struct {
		AvailableInstances []struct {
			InstanceID       string `json:"InstanceID"`
			InstanceName     string `json:"InstanceName"`
			FriendlyName     string `json:"FriendlyName"`
			Module           string `json:"Module"`
			InstalledVersion struct {
				Major         int `json:"Major"`
				Minor         int `json:"Minor"`
				Build         int `json:"Build"`
				Revision      int `json:"Revision"`
				MajorRevision int `json:"MajorRevision"`
				MinorRevision int `json:"MinorRevision"`
			} `json:"InstalledVersion"`
			Running   bool `json:"Running"`
			Suspended bool `json:"Suspended"`
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
			RawValue int    `json:"RawValue"`
			MaxValue int    `json:"MaxValue"`
			Percent  int    `json:"Percent"`
			Units    string `json:"Units"`
		} `json:"CPU Usage"`
		MemoryUsage struct {
			RawValue int    `json:"RawValue"`
			MaxValue int    `json:"MaxValue"`
			Percent  int    `json:"Percent"`
			Units    string `json:"Units"`
		} `json:"Memory Usage"`
		ActiveUsers struct {
			RawValue int    `json:"RawValue"`
			MaxValue int    `json:"MaxValue"`
			Percent  int    `json:"Percent"`
			Units    string `json:"Units"`
		} `json:"Active Users"`
	} `json:"Metrics"`
}
