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

	"github.com/lastbackend/lastbackend/internal/master/cache"
	"github.com/lastbackend/lastbackend/internal/master/ipam"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/pkg/system"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logLevel = 7

type Runtime struct {
	ctx      context.Context
	process  *system.Process
	observer *Observer
	cache    *cache.Cache
	storage  storage.IStorage
	ipam     ipam.IPAM
	active   bool
}

func NewRuntime(ctx context.Context, stg storage.IStorage, ipam ipam.IPAM, cache *cache.Cache) *Runtime {
	r := new(Runtime)

	r.ctx = ctx
	r.process = new(system.Process)
	_, err := r.process.Register(ctx, models.KindController, stg)
	if err != nil {
		log.Error(err)
		return nil
	}

	r.observer = NewObserver(stg, ipam)
	r.storage = stg
	r.ipam = ipam
	r.cache = cache

	return r
}

// Loop - runtime main loop watch
func (r *Runtime) Loop() {

	var (
		lead = make(chan bool)
	)

	log.Debug("Controller: Container: Loop")

	go func() {
		for {
			select {
			case <-r.ctx.Done():
				return

			case l := <-lead:
				{

					if l {
						if r.active {
							log.Debug("Container: is already marked as lead -> skip")
							continue
						}

						log.Debug("Container: Mark as lead")
						r.active = true
						r.observer.Loop()

					} else {

						if !r.active {
							log.Debug("Container: is already marked as slave -> skip")
							continue
						}

						log.Debug("Container: Mark as slave")
						r.active = false
						r.observer.Stop()
					}

				}
			}
		}
	}()

	go r.process.HeartBeat(r.ctx, lead)
}
