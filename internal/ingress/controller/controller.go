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

package controller

import (
	"context"
	"sync"
	"time"

	"github.com/lastbackend/lastbackend/internal/ingress/envs"
	"github.com/lastbackend/lastbackend/internal/ingress/runtime"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logPrefix = "ingress:controller"
	logLevel  = 3
)

type Controller struct {
	ctx     context.Context
	runtime *runtime.Runtime
	cache   struct {
		lock      sync.RWMutex
		resources models.IngressStatus
		routes    map[string]*models.RouteStatus
	}
}

func New(r *runtime.Runtime) *Controller {
	var c = new(Controller)
	c.ctx = context.Background()
	c.runtime = r
	c.cache.routes = make(map[string]*models.RouteStatus, 0)
	return c
}

func (c *Controller) Connect(ctx context.Context) error {

	log.Debugf("%s:connect:> connect init", logPrefix)

	opts := v1.Request().Ingress().IngressConnectOptions()
	opts.Info = envs.Get().GetState().Ingress().Info
	opts.Status = envs.Get().GetState().Ingress().Status

	var net = envs.Get().GetNet()
	if net != nil {
		opts.Network = *net.Info(ctx)
	}

	for {
		err := envs.Get().GetClient().Connect(ctx, opts)
		if err == nil {
			log.Debugf("%s connected", logPrefix)
			return nil
		}

		log.Errorf("%s connect err: %s", logPrefix, err.Error())
		time.Sleep(3 * time.Second)
	}

	return nil
}

func (c *Controller) Sync() {

	log.Debugf("Start ingress sync")

	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {
		opts := new(request.IngressStatusOptions)
		opts.Routes = make(map[string]*models.RouteStatus, 0)

		c.cache.lock.Lock()
		var i = 0
		for r, status := range c.cache.routes {
			i++
			if i > 10 {
				break
			}

			if status != nil {
				opts.Routes[r] = status
			} else {
				delete(c.cache.routes, r)
			}
		}

		spec, err := envs.Get().GetClient().SetStatus(c.ctx, opts)
		if err != nil {
			log.Errorf("ingress:exporter:dispatch err: %s", err.Error())
		}

		for r := range opts.Routes {
			delete(c.cache.routes, r)
		}

		c.cache.lock.Unlock()
		if spec != nil {
			c.runtime.Sync(spec.Decode())
		} else {
			log.Debug("received spec is nil, skip apply changes")
		}
	}

}

func (c *Controller) Subscribe() {
	var (
		routes = make(chan string)
		done   = make(chan bool)
	)

	go func() {
		log.Debugf("%s subscribe state", logPrefix)

		for {
			select {
			case r := <-routes:
				log.Debugf("%s route changed: %s", logPrefix, r)
				c.cache.lock.Lock()
				st := envs.Get().GetState().Routes().GetRouteStatus(r)
				if st == nil {
					c.cache.routes[r] = &models.RouteStatus{State: models.StateDestroyed}
				} else {
					c.cache.routes[r] = st
				}
				c.cache.lock.Unlock()
				break
			}
		}

	}()

	go envs.Get().GetState().Routes().Watch(routes, done)
	<-done
}
