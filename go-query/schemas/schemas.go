package schemas

import "github.com/dominikbraun/graph"

type GuacIDArtifact struct {
	Name string
	IDs  []GuacID
}

type GuacID struct {
	Ecosystem string   `json:"ecosystem,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
	Name      string   `json:"name,omitempty"`
	Version   string   `json:"version,omitempty"`
	Arch      string   `json:"arch,omitempty"`
	Other     []string `json:"other,omitempty"`
	SubPath   string   `json:"subpath,omitempty"`
	PkgRel    string   `json:"pkgrel,omitempty"`
	Edition   string   `json:"edition,omitempty"`
}

type Purl struct {
	Scheme    string `json:"scheme,omitempty"`
	Type      string `json:"type,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Version   string `json:"version,omitempty"`
	QualArch  string `json:"qual_arch,omitempty"`
	QualX     string `json:"qual_x,omitempty"`
	SubPath   string `json:"subpath,omitempty"`
}

type CPE struct {
	TargetSW  string   `json:"target_sw,omitempty"`
	Vendor    string   `json:"vendor,omitempty"`
	Product   string   `json:"product,omitempty"`
	Version   string   `json:"version,omitempty"`
	TargetHW  string   `json:"target_hw,omitempty"`
	Update    string   `json:"update,omitempty"`
	Edition   string   `json:"edition,omitempty"`
	Language  string   `json:"language,omitempty"`
	SWEdition string   `json:"sw_edition,omitempty"`
	Other     []string `json:"other,omitempty"`
}

type Community struct {
	CommunityID string
	Size        int
	GraphSubset *graph.Graph[string, *GuacIDNode]
}

type GuacIDNode struct {
	NodeID     string
	NodeType   GuacIDNodeType
	NodeWeight float32
}

type GuacIDNodeType int

const (
	NodeTypeSoft GuacIDNodeType = iota
	NodeTypeHard
)
