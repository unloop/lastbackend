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

package state

import (
	"context"
	"github.com/lastbackend/lastbackend/tools/logger"
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

type NetworkState struct {
	lock    sync.RWMutex
	subnets map[string]models.NetworkState
}

func (n *NetworkState) GetSubnets() map[string]models.NetworkState {
	return n.subnets
}

func (n *NetworkState) AddSubnet(cidr string, sn *models.NetworkState) {
	log := logger.WithContext(context.Background())
	log.Debugf("Stage: NetworkState: add subnet: %s", cidr)
	n.SetSubnet(cidr, sn)
}

func (n *NetworkState) SetSubnet(cidr string, sn *models.NetworkState) {
	log := logger.WithContext(context.Background())
	log.Debugf("Stage: NetworkState: set subnet: %s", cidr)
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, ok := n.subnets[cidr]; ok {
		delete(n.subnets, cidr)
	}

	n.subnets[cidr] = *sn
}

func (n *NetworkState) GetSubnet(cidr string) *models.NetworkState {
	log := logger.WithContext(context.Background())
	log.Debugf("Stage: NetworkState: get subnet: %s", cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	s, ok := n.subnets[cidr]
	if !ok {
		return nil
	}
	return &s
}

func (n *NetworkState) DelSubnet(cidr string) {
	log := logger.WithContext(context.Background())
	log.Debugf("Stage: NetworkState: del subnet: %s", cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.subnets[cidr]; ok {
		delete(n.subnets, cidr)
	}
}
