//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package network

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/network/state"
)

const logLevel = 3

func (n *Network) Subnets() *state.SubnetState {
	return n.state.Subnets()
}

func (n *Network) Info(ctx context.Context) *types.NetworkState {
	return n.cni.Info(ctx)
}

func (n *Network) SubnetRestore(ctx context.Context) error {

	sn, err := n.cni.Subnets(ctx)
	if err != nil {
		log.Errorf("Can-not get subnets from CNI err: %v", err)
	}

	for cidr, s := range sn {
		n.state.Subnets().SetSubnet(cidr, s)
	}

	return nil
}

func (n *Network) SubnetManage(ctx context.Context, cidr string, sn *types.SubnetManifest) error {

	subnets := n.state.Subnets().GetSubnets()
	if state, ok := subnets[cidr]; ok {

		if sn.State == types.StateDestroy {
			n.cni.Destroy(ctx, &state)
			n.state.Subnets().DelSubnet(cidr)
			return nil
		}

		// TODO: check if network manifest changes
		// if changes then update routes and interfaces
		return nil
	}

	if sn.State == types.StateDestroy {
		return nil
	}

	state, err := n.cni.Create(ctx, sn)
	if err != nil {
		log.Errorf("Can not create network subnet: %s", err.Error())
		return err
	}

	n.state.Subnets().AddSubnet(cidr, state)
	return nil
}

func (n *Network) SubnetDestroy(ctx context.Context, cidr string) error {

	sn := n.state.Subnets().GetSubnet(cidr)

	if err := n.cni.Destroy(ctx, sn); err != nil {
		log.Errorf("Can not destroy network subnet: %s", err.Error())
		return err
	}

	n.state.Subnets().DelSubnet(cidr)
	return nil
}