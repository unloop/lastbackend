//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package views

import (
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

// Node - default node structure
// swagger:model views_node
type Node struct {
	Meta   NodeMeta   `json:"meta"`
	Status NodeStatus `json:"status"`
	Spec   NodeSpec   `json:"spec"`
}

// NodeList - node map list
// swagger:model views_node_list
type NodeList map[string]*Node

// NodeMeta - node metadata structure
// swagger:model views_node_meta
type NodeMeta struct {
	NodeInfo
	Name     string            `json:"name"`
	Labels   map[string]string `json:"labels"`
	SelfLink string            `json:"self_link"`
	Created  time.Time         `json:"created"`
	Updated  time.Time         `json:"updated"`
}

// NodeInfo - node info struct
// swagger:model views_node_info
type NodeInfo struct {
	Hostname     string `json:"hostname"`
	OSName       string `json:"os_name"`
	OSType       string `json:"os_type"`
	Architecture string `json:"architecture"`
	IP           struct {
		External string `json:"external"`
		Internal string `json:"internal"`
	} `json:"ip"`
	CIDR    string `json:"cidr"`
	Version string `json:"version"`
}

type NodeStatusState struct {
	Ready bool                     `json:"ready"`
	CRI   NodeStatusInterfaceState `json:"cri"`
	CNI   NodeStatusInterfaceState `json:"cni"`
	CPI   NodeStatusInterfaceState `json:"cpi"`
	CSI   NodeStatusInterfaceState `json:"csi"`
}

type NodeStatusInterfaceState struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	State   string `json:"state"`
	Message string `json:"message"`
}

// NodeStatus - node state struct
// swagger:model views_node_status
type NodeStatus struct {
	State     NodeStatusState `json:"state"`
	Online    bool            `json:"online"`
	Capacity  NodeResources   `json:"capacity"`
	Allocated NodeResources   `json:"allocated"`
}

// swagger:ignore
// swagger:model types_node_spec
type NodeSpec struct {
	Security NodeSecurity `json:"security"`
}

type NodeSecurity struct {
	TLS bool `json:"tls"`
}

// NodeResources - node resources structure
// swagger:model views_node_resources
type NodeResources struct {
	Containers int    `json:"containers"`
	Pods       int    `json:"pods"`
	Memory     string `json:"ram"`
	CPU        string `json:"cpu"`
	Storage    string `json:"storage"`
}

// swagger:model views_node_spec
type NodeManifest struct {
	Meta      NodeManifestMeta                    `json:"meta"`
	Discovery map[string]*models.ResolverManifest `json:"discovery,omitempty"`
	Exporter  *models.ExporterManifest            `json:"exporter,omitempty"`
	Configs   map[string]*models.ConfigManifest   `json:"configs,omitempty"`
	Secrets   map[string]*models.SecretManifest   `json:"secrets,omitempty"`
	Network   map[string]*models.SubnetManifest   `json:"network,omitempty"`
	Pods      map[string]*models.PodManifest      `json:"pods,omitempty"`
	Volumes   map[string]*models.VolumeManifest   `json:"volumes,omitempty"`
	Endpoints map[string]*models.EndpointManifest `json:"endpoints,omitempty"`
}

type NodeManifestMeta struct {
	Initial bool `json:"initial"`
}
