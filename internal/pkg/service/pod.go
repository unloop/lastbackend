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

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logPodPrefix = "distribution:pod"
)

type Pod struct {
	context context.Context
	storage storage.IStorage
}

func (p *Pod) Runtime() (*models.System, error) {

	log.Debugf("%s:get:> get pod runtime info", logPodPrefix)
	runtime, err := p.storage.Info(p.context, p.storage.Collection().Pod(), "")
	if err != nil {
		log.Errorf("%s:get:> get runtime info error: %s", logPodPrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil
}

// Get pod info from storage
func (p *Pod) Get(selflink string) (*models.Pod, error) {
	log.Debugf("%s:get:> get by name %s", logPodPrefix, selflink)

	pod := new(models.Pod)

	err := p.storage.Get(p.context, p.storage.Collection().Pod(),
		selflink, pod, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.Warnf("%s:get:> `%s` not found", logPodPrefix, selflink)
			return nil, nil
		}

		log.Debugf("%s:get:> get Pod `%s` err: %v", logPodPrefix, selflink, err)
		return nil, err
	}

	return pod, nil
}

// Create new pod
func (p *Pod) Put(pod *models.Pod) (*models.Pod, error) {

	if err := p.storage.Put(p.context, p.storage.Collection().Pod(),
		pod.SelfLink().String(), pod, nil); err != nil {
		log.Errorf("%s:create:> insert pod err %v", logPodPrefix, err)
		return nil, err
	}

	return pod, nil
}

// ListByNamespace returns pod list in selected namespace
func (p *Pod) ListByNamespace(namespace string) (*models.PodList, error) {
	log.Debugf("%s:listbynamespace:> get pod list by namespace %s", logPodPrefix, namespace)

	list := models.NewPodList()
	filter := p.storage.Filter().Pod().ByNamespace(namespace)

	err := p.storage.List(p.context, p.storage.Collection().Pod(), filter, list, nil)
	if err != nil {
		log.Debugf("%s:listbynamespace:> get pod list by deployment id `%s` err: %v", logPodPrefix, namespace, err)
		return nil, err
	}

	return list, nil
}

// ListByService returns pod list in selected service
func (p *Pod) ListByService(namespace, service string) (*models.PodList, error) {
	log.Debugf("%s:listbyservice:> get pod list by service id %s/%s", logPodPrefix, namespace, service)

	list := models.NewPodList()
	filter := p.storage.Filter().Pod().ByService(namespace, service)

	err := p.storage.List(p.context, p.storage.Collection().Pod(), filter, list, nil)
	if err != nil {
		log.Debugf("%s:listbyservice:> get pod list by service id `%s` err: %v", logPodPrefix, namespace, service, err)
		return nil, err
	}

	return list, nil
}

// ListByDeployment returns pod list in selected deployment
func (p *Pod) ListByDeployment(namespace, service, deployment string) (*models.PodList, error) {
	log.Debugf("%s:listbydeployment:> get pod list by id %s/%s/%s", logPodPrefix, namespace, service, deployment)

	list := models.NewPodList()
	filter := p.storage.Filter().Pod().ByDeployment(namespace, service, deployment)

	err := p.storage.List(p.context, p.storage.Collection().Pod(), filter, list, nil)
	if err != nil {
		log.Debugf("%s:listbydeployment:> get pod list by deployment id `%s/%s/%s` err: %v",
			logPodPrefix, namespace, service, deployment, err)
		return nil, err
	}

	return list, nil
}

// ListByJob returns pod list in selected job
func (p *Pod) ListByJob(namespace, job string) (*models.PodList, error) {
	log.Debugf("%s:listbyjob:> get pod list by id %s/%s", logPodPrefix, namespace, job)

	list := models.NewPodList()
	filter := p.storage.Filter().Pod().ByJob(namespace, job)

	err := p.storage.List(p.context, p.storage.Collection().Pod(), filter, list, nil)
	if err != nil {
		log.Debugf("%s:listbyjob:> get pod list by deployment id `%s/%s` err: %v",
			logPodPrefix, namespace, job, err)
		return nil, err
	}

	return list, nil
}

// ListByTask returns pod list in selected task
func (p *Pod) ListByTask(namespace, job, task string) (*models.PodList, error) {
	log.Debugf("%s:listbytask:> get pod list by id %s/%s/%s", logPodPrefix, namespace, job, task)

	list := models.NewPodList()
	filter := p.storage.Filter().Pod().ByTask(namespace, job, task)

	err := p.storage.List(p.context, p.storage.Collection().Pod(), filter, list, nil)
	if err != nil {
		log.Debugf("%s:listbytask:> get pod list by deployment id `%s/%s/%s` err: %v",
			logPodPrefix, namespace, job, task, err)
		return nil, err
	}

	return list, nil
}

// SetNode - set node info to pod
func (p *Pod) SetNode(pod *models.Pod, node *models.Node) error {
	log.Debugf("%s:setnode:> set node for pod: %s", logPodPrefix, pod.Meta.Name)

	pod.Meta.Node = node.Meta.Name

	if err := p.storage.Set(p.context, p.storage.Collection().Pod(),
		pod.SelfLink().String(), pod, nil); err != nil {
		log.Errorf("%s:setnode:> pod set node err: %v", logPodPrefix, err)
		return err
	}

	return nil
}

// SetStatus - set state for pod
func (p *Pod) Update(pod *models.Pod) error {

	log.Debugf("%s:update:> update pod: %s", logPodPrefix, pod.Meta.Name)

	if err := p.storage.Set(p.context, p.storage.Collection().Pod(),
		pod.SelfLink().String(),
		pod, nil); err != nil {
		log.Errorf("%s:update:> pod update err: %v", logPodPrefix, err)
		return err
	}

	return nil
}

// Destroy pod
func (p *Pod) Destroy(pod *models.Pod) error {

	pod.Spec.State.Destroy = true

	if err := p.storage.Set(p.context, p.storage.Collection().Pod(),
		pod.SelfLink().String(), pod, nil); err != nil {
		log.Errorf("%s:destroy:> mark pod for destroy error: %v", logPodPrefix, err)
		return err
	}
	return nil
}

// Remove pod from storage
func (p *Pod) Remove(pod *models.Pod) error {
	if err := p.storage.Del(p.context, p.storage.Collection().Pod(),
		pod.SelfLink().String()); err != nil {
		log.Errorf("%s:remove:> mark pod for destroy error: %v", logPodPrefix, err)
		return err
	}
	return nil
}

func (p *Pod) Watch(ch chan models.PodEvent, rev *int64) error {
	log.Debugf("%s:watch:> watch pod, from revision %d", logPodPrefix, *rev)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-p.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := models.PodEvent{}
				res.Action = e.Action
				res.Name = e.Name

				obj := new(models.Pod)

				if err := json.Unmarshal(e.Data.([]byte), obj); err != nil {
					log.Errorf("%s:watch:> parse json", logPodPrefix)
					continue
				}

				res.Data = obj

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := p.storage.Watch(p.context, p.storage.Collection().Pod(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func (p *Pod) ManifestMap(node string) (*models.PodManifestMap, error) {
	log.Debugf("%s:PodManifestMap:> ", logPodPrefix)

	var (
		mf = models.NewPodManifestMap()
	)

	if err := p.storage.Map(p.context, p.storage.Collection().Manifest().Pod(node), models.EmptyString, mf, nil); err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s:PodManifestMap:> err: %s", logPodPrefix, err.Error())
			return nil, err
		}

		return nil, nil
	}

	return mf, nil
}

func (p *Pod) ManifestGet(node, pod string) (*models.PodManifest, error) {
	log.Debugf("%s:PodManifestGet:> ", logPodPrefix)

	var (
		mf = new(models.PodManifest)
	)

	if err := p.storage.Get(p.context, p.storage.Collection().Manifest().Pod(node), pod, &mf, nil); err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return mf, nil
}

func (p *Pod) ManifestAdd(node, pod string, manifest *models.PodManifest) error {
	log.Debugf("%s:PodManifestAdd:> ", logPodPrefix)

	if err := p.storage.Put(p.context, p.storage.Collection().Manifest().Pod(node), pod, manifest, nil); err != nil {
		log.Errorf("%s:PodManifestAdd:> err :%s", logPodPrefix, err.Error())
		return err
	}

	return nil
}

func (p *Pod) ManifestSet(node, pod string, manifest *models.PodManifest) error {
	log.Debugf("%s:PodManifestSet:> ", logPodPrefix)

	if err := p.storage.Set(p.context, p.storage.Collection().Manifest().Pod(node), pod, manifest, nil); err != nil {
		log.Errorf("%s:PodManifestSet:> err :%s", logPodPrefix, err.Error())
		return err
	}

	return nil
}

func (p *Pod) ManifestDel(node, pod string) error {
	log.Debugf("%s:PodManifestDel:> %s on node %s", logPodPrefix, pod, node)

	if err := p.storage.Del(p.context, p.storage.Collection().Manifest().Pod(node), pod); err != nil {
		log.Errorf("%s:PodManifestDel:> err :%s", logPodPrefix, err.Error())
		return err
	}

	return nil
}

func (p *Pod) ManifestWatch(node string, ch chan models.PodManifestEvent, rev *int64) error {

	log.Debugf("%s:watch:> watch pod manifest ", logPodPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	var f, c string

	if node != models.EmptyString {
		f = fmt.Sprintf(`\b.+\/%s\/%s\/(.+)\b`, node, storage.PodKind)
		c = p.storage.Collection().Manifest().Pod(node)
	} else {
		f = fmt.Sprintf(`\b.+\/(.+)\/%s\/(.+)\b`, storage.PodKind)
		c = p.storage.Collection().Manifest().Node()
	}

	r, err := regexp.Compile(f)
	if err != nil {
		log.Errorf("%s:> filter compile err: %v", logPodPrefix, err.Error())
		return err
	}

	go func() {
		for {
			select {
			case <-p.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				keys := r.FindStringSubmatch(e.Storage.Key)
				if len(keys) == 0 {
					continue
				}

				res := models.PodManifestEvent{}
				res.Action = e.Action
				res.Name = e.Name
				res.SelfLink = e.SelfLink
				if node != models.EmptyString {
					res.Node = node
				} else {
					res.Node = keys[1]
				}

				manifest := new(models.PodManifest)

				if err := json.Unmarshal(e.Data.([]byte), manifest); err != nil {
					log.Errorf("%s:> parse data err: %v", logPodPrefix, err)
					continue
				}

				res.Data = manifest

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := p.storage.Watch(p.context, c, watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewPodModel(ctx context.Context, stg storage.IStorage) *Pod {
	return &Pod{ctx, stg}
}
