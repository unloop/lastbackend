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

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

type IngressView struct{}

func (nv *IngressView) New(obj *models.Ingress) *Ingress {
	n := Ingress{}
	n.Meta = nv.ToIngressMeta(obj.Meta)
	n.Status = nv.ToIngressStatus(obj.Status)
	return &n
}

func (nv *IngressView) ToIngressMeta(meta models.IngressMeta) IngressMeta {
	m := IngressMeta{}
	m.Name = meta.Name
	m.Description = meta.Description
	m.Created = meta.Created
	m.Updated = meta.Updated
	return m
}

func (nv *IngressView) ToIngressStatus(status models.IngressStatus) IngressStatus {
	return IngressStatus{
		Ready: status.Ready,
	}
}

func (obj *Ingress) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *IngressView) NewList(obj *models.IngressList) *IngressList {
	if obj == nil {
		return nil
	}
	ingresses := make(IngressList, 0)
	for _, v := range obj.Items {
		nn := nv.New(v)
		ingresses[nn.Meta.Name] = nn
	}

	return &ingresses
}

func (obj *IngressList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *IngressView) NewManifest(obj *models.IngressManifest) *IngressManifest {

	manifest := IngressManifest{
		Endpoints: make(map[string]*models.EndpointManifest, 0),
		Routes:    make(map[string]*models.RouteManifest, 0),
		Subnets:   make(map[string]*models.SubnetManifest, 0),
	}

	if obj == nil {
		return nil
	}

	manifest.Meta.Initial = obj.Meta.Initial
	manifest.Resolvers = obj.Resolvers

	manifest.Endpoints = obj.Endpoints
	manifest.Routes = obj.Routes
	manifest.Subnets = obj.Network

	return &manifest
}

func (obj *IngressManifest) Decode() *models.IngressManifest {

	manifest := models.IngressManifest{
		Routes:    make(map[string]*models.RouteManifest, 0),
		Endpoints: make(map[string]*models.EndpointManifest, 0),
		Network:   make(map[string]*models.SubnetManifest, 0),
	}

	manifest.Meta.Initial = obj.Meta.Initial
	manifest.Resolvers = make(map[string]*models.ResolverManifest, 0)

	for i, r := range obj.Resolvers {
		manifest.Resolvers[i] = r
	}

	for i, r := range obj.Routes {
		manifest.Routes[i] = r
	}

	for i, e := range obj.Endpoints {
		manifest.Endpoints[i] = e
	}

	for i, e := range obj.Subnets {
		manifest.Network[i] = e
	}

	return &manifest
}

func (obj *IngressManifest) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
