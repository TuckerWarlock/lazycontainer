package commands

import (
	"fmt"
	"time"
)

// Container represents an Apple container
type Container struct {
	Status        string            `json:"status"`
	StartedDate   float64           `json:"startedDate"`
	Configuration ContainerConfig   `json:"configuration"`
	Networks      []NetworkAttachment `json:"networks"`
}

// GetID returns the container ID
func (c *Container) GetID() string {
	return c.Configuration.ID
}

// GetName returns the container name (same as ID for Apple containers)
func (c *Container) GetName() string {
	return c.Configuration.ID
}

// GetImage returns the image reference
func (c *Container) GetImage() string {
	return c.Configuration.Image.Reference
}

// GetStatus returns the container status
func (c *Container) GetStatus() string {
	return c.Status
}

// IsRunning returns true if the container is running
func (c *Container) IsRunning() bool {
	return c.Status == "running"
}

// GetStartedAt returns the container start time
func (c *Container) GetStartedAt() time.Time {
	// Apple uses CFAbsoluteTime (seconds since Jan 1, 2001)
	cfEpoch := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	return cfEpoch.Add(time.Duration(c.StartedDate) * time.Second)
}

// GetUptime returns a human-readable uptime string
func (c *Container) GetUptime() string {
	if !c.IsRunning() {
		return "-"
	}
	uptime := time.Since(c.GetStartedAt())
	if uptime < time.Minute {
		return fmt.Sprintf("%ds", int(uptime.Seconds()))
	} else if uptime < time.Hour {
		return fmt.Sprintf("%dm", int(uptime.Minutes()))
	} else if uptime < 24*time.Hour {
		return fmt.Sprintf("%dh", int(uptime.Hours()))
	}
	return fmt.Sprintf("%dd", int(uptime.Hours()/24))
}

// GetCPUs returns allocated CPU count
func (c *Container) GetCPUs() int {
	return c.Configuration.Resources.CPUs
}

// GetMemory returns allocated memory in bytes
func (c *Container) GetMemory() int64 {
	return c.Configuration.Resources.MemoryInBytes
}

// GetMemoryHuman returns human-readable memory string
func (c *Container) GetMemoryHuman() string {
	mem := c.Configuration.Resources.MemoryInBytes
	if mem >= 1<<30 {
		return fmt.Sprintf("%.1fG", float64(mem)/(1<<30))
	} else if mem >= 1<<20 {
		return fmt.Sprintf("%.0fM", float64(mem)/(1<<20))
	}
	return fmt.Sprintf("%dB", mem)
}

// GetPorts returns published ports as a string
func (c *Container) GetPorts() string {
	if len(c.Configuration.PublishedPorts) == 0 {
		return ""
	}
	ports := ""
	for i, p := range c.Configuration.PublishedPorts {
		if i > 0 {
			ports += ", "
		}
		ports += fmt.Sprintf("%d:%d", p.HostPort, p.ContainerPort)
	}
	return ports
}

type ContainerConfig struct {
	ID              string            `json:"id"`
	Image           ImageRef          `json:"image"`
	InitProcess     InitProcess       `json:"initProcess"`
	Resources       Resources         `json:"resources"`
	Mounts          []Mount           `json:"mounts"`
	Networks        []NetworkConfig   `json:"networks"`
	PublishedPorts  []PortMapping     `json:"publishedPorts"`
	Labels          map[string]string `json:"labels"`
	RuntimeHandler  string            `json:"runtimeHandler"`
	Rosetta         bool              `json:"rosetta"`
	SSH             bool              `json:"ssh"`
	Virtualization  bool              `json:"virtualization"`
	ReadOnly        bool              `json:"readOnly"`
	Platform        Platform          `json:"platform"`
	DNS             DNSConfig         `json:"dns"`
}

type ImageRef struct {
	Reference  string         `json:"reference"`
	Descriptor ImageDescriptor `json:"descriptor"`
}

type ImageDescriptor struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

type InitProcess struct {
	Executable          string            `json:"executable"`
	Arguments           []string          `json:"arguments"`
	Environment         []string          `json:"environment"`
	WorkingDirectory    string            `json:"workingDirectory"`
	User                UserID            `json:"user"`
	Terminal            bool              `json:"terminal"`
	SupplementalGroups  []int             `json:"supplementalGroups"`
	Rlimits             []Rlimit          `json:"rlimits"`
}

type UserID struct {
	ID struct {
		UID int `json:"uid"`
		GID int `json:"gid"`
	} `json:"id"`
}

type Rlimit struct {
	Type string `json:"type"`
	Hard int64  `json:"hard"`
	Soft int64  `json:"soft"`
}

type Resources struct {
	CPUs          int   `json:"cpus"`
	MemoryInBytes int64 `json:"memoryInBytes"`
}

type Mount struct {
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Type        MountType `json:"type"`
	Options     []string  `json:"options"`
}

type MountType struct {
	Volume *VolumeMount `json:"volume,omitempty"`
	Bind   *BindMount   `json:"bind,omitempty"`
}

type VolumeMount struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

type BindMount struct {
	Source   string `json:"source"`
	ReadOnly bool   `json:"readOnly"`
}

type NetworkConfig struct {
	Network string        `json:"network"`
	Options NetworkOptions `json:"options"`
}

type NetworkOptions struct {
	Hostname string `json:"hostname"`
}

type NetworkAttachment struct {
	Network string `json:"network"`
	IP      string `json:"ip,omitempty"`
}

type PortMapping struct {
	HostPort      int    `json:"hostPort"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol"`
	HostIP        string `json:"hostIP,omitempty"`
}

type Platform struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
}

type DNSConfig struct {
	Nameservers   []string `json:"nameservers"`
	SearchDomains []string `json:"searchDomains"`
	Options       []string `json:"options"`
}

// Image represents an Apple container image
type Image struct {
	Reference string  `json:"reference"`
	Digest    string  `json:"digest"`
	Size      int64   `json:"size"`
	CreatedAt float64 `json:"createdAt"`
}

// GetSizeHuman returns human-readable size
func (i *Image) GetSizeHuman() string {
	if i.Size >= 1<<30 {
		return fmt.Sprintf("%.1fG", float64(i.Size)/(1<<30))
	} else if i.Size >= 1<<20 {
		return fmt.Sprintf("%.1fM", float64(i.Size)/(1<<20))
	} else if i.Size >= 1<<10 {
		return fmt.Sprintf("%.1fK", float64(i.Size)/(1<<10))
	}
	return fmt.Sprintf("%dB", i.Size)
}

// GetDigestShort returns truncated digest
func (i *Image) GetDigestShort() string {
	if len(i.Digest) > 19 {
		return i.Digest[:19]
	}
	return i.Digest
}

// Volume represents an Apple container volume
type Volume struct {
	Name        string            `json:"name"`
	Source      string            `json:"source"`
	Format      string            `json:"format"`
	SizeInBytes int64             `json:"sizeInBytes"`
	CreatedAt   float64           `json:"createdAt"`
	Labels      map[string]string `json:"labels"`
	Driver      string            `json:"driver"`
	Options     map[string]string `json:"options"`
}

// GetSizeHuman returns human-readable size
func (v *Volume) GetSizeHuman() string {
	if v.SizeInBytes >= 1<<40 {
		return fmt.Sprintf("%.1fT", float64(v.SizeInBytes)/(1<<40))
	} else if v.SizeInBytes >= 1<<30 {
		return fmt.Sprintf("%.1fG", float64(v.SizeInBytes)/(1<<30))
	} else if v.SizeInBytes >= 1<<20 {
		return fmt.Sprintf("%.1fM", float64(v.SizeInBytes)/(1<<20))
	}
	return fmt.Sprintf("%dB", v.SizeInBytes)
}

// Network represents an Apple container network
type Network struct {
	ID     string        `json:"id"`
	State  string        `json:"state"`
	Config NetworkDetail `json:"config"`
	Status NetworkStatus `json:"status"`
}

type NetworkDetail struct {
	ID           string            `json:"id"`
	Mode         string            `json:"mode"`
	CreationDate float64           `json:"creationDate"`
	Labels       map[string]string `json:"labels"`
}

type NetworkStatus struct {
	IPv4Subnet  string `json:"ipv4Subnet"`
	IPv4Gateway string `json:"ipv4Gateway"`
	IPv6Subnet  string `json:"ipv6Subnet"`
}

// ContainerStats represents container resource usage
type ContainerStats struct {
	ContainerID string
	CPUPercent  float64
	MemoryUsage int64
	MemoryLimit int64
	NetworkRx   int64
	NetworkTx   int64
}

// GetMemoryPercent returns memory usage percentage
func (s *ContainerStats) GetMemoryPercent() float64 {
	if s.MemoryLimit == 0 {
		return 0
	}
	return float64(s.MemoryUsage) / float64(s.MemoryLimit) * 100
}
