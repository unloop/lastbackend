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
	"encoding/json"

	"github.com/lastbackend/lastbackend/internal/util/resource"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

type VolumeView struct{}

func (rv *VolumeView) New(obj *models.Volume) *Volume {
	r := Volume{}
	r.Meta = r.ToMeta(obj.Meta)
	r.Spec = r.ToSpec(obj.Spec)
	r.Status = r.ToStatus(obj.Status)
	return &r
}

func (p *Volume) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Volume) ToMeta(obj models.VolumeMeta) VolumeMeta {
	meta := VolumeMeta{}
	meta.Name = obj.Name
	meta.Namespace = obj.Namespace
	meta.SelfLink = obj.SelfLink.String()
	meta.Updated = obj.Updated
	meta.Created = obj.Created

	return meta
}

func (p *Volume) ToSpec(obj models.VolumeSpec) VolumeSpec {
	spec := VolumeSpec{}
	spec.State.Destroy = obj.State.Destroy
	spec.Selector.Node = obj.Selector.Node
	spec.Selector.Labels = obj.Selector.Labels
	spec.HostPath = obj.HostPath
	spec.Type = obj.Type
	spec.AccessMode = obj.AccessMode
	spec.Capacity.Storage = resource.EncodeMemoryResource(obj.Capacity.Storage)
	return spec
}

func (p *Volume) ToStatus(obj models.VolumeStatus) VolumeStatus {
	state := VolumeStatus{
		State:   obj.State,
		Message: obj.Message,
		Status: VolumeState{
			Path:  obj.Status.Path,
			Type:  obj.Status.Type,
			Ready: obj.Status.Ready,
		},
	}
	return state
}

func (rv VolumeView) NewList(obj *models.VolumeList) *VolumeList {
	if obj == nil {
		return nil
	}

	n := make(VolumeList, 0)
	for _, v := range obj.Items {
		n = append(n, rv.New(v))
	}
	return &n
}

func (n *VolumeList) ToJson() ([]byte, error) {
	if n == nil {
		n = &VolumeList{}
	}
	return json.Marshal(n)
}
