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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/ingress/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logRuntimePrefix = "ingress:runtime"
	logLevel         = 3
)

type Runtime struct {
	ctx     context.Context
	spec    chan *models.IngressManifest
	process *Process
	config  *conf
	iface   string
}

func New(iface string, cfg *conf) *Runtime {
	r := new(Runtime)
	r.ctx = context.Background()
	r.spec = make(chan *models.IngressManifest)
	r.process = new(Process)
	r.iface = iface
	r.config = cfg
	return r
}

// Restore node runtime state
func (r *Runtime) Restore() {
	log.Debugf("%s:restore:> restore init", logRuntimePrefix)
	var network = envs.Get().GetNet()

	if network != nil {
		if err := envs.Get().GetNet().EndpointRestore(r.ctx); err != nil {
			log.Errorf("%s:> can not restore endpoins: %s", logRuntimePrefix, err.Error())
		}

		if err := envs.Get().GetNet().SubnetRestore(r.ctx); err != nil {
			log.Errorf("%s:> can not restore network: %s", logRuntimePrefix, err.Error())
		}

		if err := envs.Get().GetNet().ResolverManage(r.ctx); err != nil {
			log.Errorf("%s:> can not manage resolver:%s", logRuntimePrefix, err.Error())
		}
	}

	if err := r.config.Check(); err != nil {
		log.Errorf("can no sync config: %s", err.Error())
		return
	}

	if err := r.process.manage(); err != nil {
		log.Errorf("can not manage haproxy process: %s", err.Error())
		return
	}
}

// Sync node runtime with new spec
func (r *Runtime) Sync(spec *models.IngressManifest) error {
	log.Debugf("%s:sync:> sync runtime state", logRuntimePrefix)
	r.spec <- spec
	return nil
}

func (r *Runtime) Loop() {

	var network = envs.Get().GetNet()

	log.Debugf("%s:loop:> start runtime loop", logRuntimePrefix)

	go func(ctx context.Context) {
		for {
			select {
			case spec := <-r.spec:

				log.Debugf("%s:loop:> provision new spec", logRuntimePrefix)

				if spec.Meta.Initial && network != nil {

					log.Debugf("%s> clean up endpoints", logRuntimePrefix)
					endpoints := envs.Get().GetNet().Endpoints().GetEndpoints()
					for e := range endpoints {

						log.Debugf("check endpoint: %s", e)

						if e == envs.Get().GetNet().GetResolverEndpointKey() {
							continue
						}

						if _, ok := spec.Endpoints[e]; !ok {
							network.EndpointDestroy(context.Background(), e, endpoints[e])
						}
					}

					log.Debugf("%s> clean up networks", logRuntimePrefix)
					nets := network.Subnets().GetSubnets()

					for cidr := range nets {
						if _, ok := spec.Network[cidr]; !ok {
							network.SubnetDestroy(ctx, cidr)
						}
					}

				}

				if len(spec.Resolvers) != 0 {
					log.Debugf("%s>set cluster dns ips: %#v", logRuntimePrefix, spec.Resolvers)

					resolvers := make(map[string]uint16, 0)
					for key, res := range spec.Resolvers {
						resolvers[res.IP] = res.Port

						envs.Get().SetResolvers(resolvers)

						if network != nil {
							network.Resolvers().SetResolver(key, res)
						}

					}
				}

				log.Debugf("%s> provision init", logRuntimePrefix)

				if network != nil {
					log.Debugf("%s> provision endpoints", logRuntimePrefix)
					for e, spec := range spec.Endpoints {
						log.Debugf("endpoint: %v", e)
						if err := network.EndpointManage(ctx, e, spec); err != nil {
							log.Errorf("Upstream [%s] manage err: %s", e, err.Error())
						}
					}

					log.Debugf("%s> provision networks", logRuntimePrefix)
					for cidr, n := range spec.Network {
						log.Debugf("network: %v", n)
						if err := network.SubnetManage(ctx, cidr, n); err != nil {
							log.Errorf("Subnet [%s] create err: %s", n.CIDR, err.Error())
						}
					}
				}

				log.Debugf("%s> provision routes", logRuntimePrefix)
				var upd = make(map[string]bool, 0)
				for e, spec := range spec.Routes {
					if err := r.RouteManage(e, spec); err != nil {
						log.Errorf("Route [%s] manage err: %s", e, err.Error())
						continue
					}
					upd[e] = true
				}

				if len(upd) > 0 {

					if err := r.process.reload(); err != nil {
						log.Errorf("reload process err: %s", err.Error())
					} else {

						for r := range upd {
							st := envs.Get().GetState().Routes().GetRouteStatus(r)
							if st == nil {
								continue
							}

							if st.State == models.StateProvision {
								st.State = models.StateReady
								envs.Get().GetState().Routes().SetRouteStatus(r, st)
							}
						}
					}
				}

			}
		}
	}(r.ctx)
}
