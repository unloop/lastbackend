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

package job

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/master/envs"
	"github.com/lastbackend/lastbackend/internal/master/state/cluster"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)
}

func testJobObserver(t *testing.T, name, werr string, wjs *JobState, js *JobState, job *models.Job) {
	var (
		ctx = context.Background()
		stg = envs.Get().GetStorage()
		err error
	)

	err = stg.Del(ctx, stg.Collection().Task(), "")
	if !assert.NoError(t, err) {
		return
	}

	err = stg.Del(ctx, stg.Collection().Pod(), "")
	if !assert.NoError(t, err) {
		return
	}

	t.Run(name, func(t *testing.T) {

		err := jobObserve(js, job)
		if werr != models.EmptyString {

			if assert.NoError(t, err, "error should be presented") {
				return
			}

			if !assert.Equal(t, werr, err.Error(), "err message different") {
				return
			}

			return
		}

		if wjs.job == nil {
			if assert.Nil(t, js.job, "job should be nil") {
				return
			}
		}

		if err := compareJobStateProperties(wjs, js); assert.NoError(t, err) {
			return
		}

	})

}

func TestHandleJobStateCreated(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			job      *models.Job
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle created job state"}

		job := getJobAsset(models.StateCreated, models.EmptyString)
		js := getJobStateAsset(job)

		s.args.jobState = js
		s.args.job = job

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting

		return s
	}())

	for _, tt := range tests {
		testJobObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.job)
	}
}

func TestHandleJobStateRunning(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			job      *models.Job
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle running job state"}

		job := getJobAsset(models.StateRunning, models.EmptyString)
		js := getJobStateAsset(job)

		s.args.jobState = js
		s.args.job = job

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateRunning

		return s
	}())

	for _, tt := range tests {
		testJobObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.job)
	}
}

func TestHandleJobStatePaused(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			job      *models.Job
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle paused job state"}

		job := getJobAsset(models.StatePaused, models.EmptyString)
		js := getJobStateAsset(job)

		s.args.jobState = js
		s.args.job = job

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StatePaused

		return s
	}())

	for _, tt := range tests {
		testJobObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.job)
	}
}

func TestHandleJobStateError(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			job      *models.Job
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle error job state"}

		job := getJobAsset(models.StateError, models.EmptyString)
		js := getJobStateAsset(job)

		s.args.jobState = js
		s.args.job = job

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateError

		return s
	}())

	for _, tt := range tests {
		testJobObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.job)
	}
}

func TestHandleJobStateDestroy(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			job      *models.Job
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroy job state without tasks"}

		job := getJobAsset(models.StateDestroy, models.EmptyString)
		js := getJobStateAsset(job)

		s.args.jobState = js
		s.args.job = job

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateDestroyed

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroy job state with tasks"}

		job := getJobAsset(models.StateDestroy, models.EmptyString)
		js := getJobStateAsset(job)
		task1 := getTaskAsset(job, models.StateDestroyed, models.EmptyString)
		task2 := getTaskAsset(job, models.StateQueued, models.EmptyString)

		s.args.jobState = js
		s.args.job = job
		s.args.jobState.task.list[task1.SelfLink().String()] = task1
		s.args.jobState.task.list[task2.SelfLink().String()] = task2

		wt1 := getTaskCopy(task1)
		wt2 := getTaskCopy(task2)
		wt2.Spec.State.Destroy = true
		wt2.Status.State = models.StateDestroyed

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.task.list[wt1.SelfLink().String()] = wt1
		s.want.jobState.task.list[wt2.SelfLink().String()] = wt2

		return s
	}())

	for _, tt := range tests {
		testJobObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.job)
	}
}

func TestHandleJobStateDestroyed(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			job      *models.Job
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroyed job state without tasks"}

		job := getJobAsset(models.StateDestroyed, models.EmptyString)
		js := getJobStateAsset(job)

		s.args.jobState = js
		s.args.job = job

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job = nil

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroyed job state with tasks"}

		job := getJobAsset(models.StateDestroyed, models.EmptyString)
		js := getJobStateAsset(job)
		task1 := getTaskAsset(job, models.StateCreated, models.EmptyString)
		task2 := getTaskAsset(job, models.StateDestroyed, models.EmptyString)

		s.args.jobState = js
		s.args.job = job
		s.args.jobState.task.list[task1.SelfLink().String()] = task1
		s.args.jobState.task.list[task2.SelfLink().String()] = task2

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateDestroy

		return s
	}())

	for _, tt := range tests {
		testJobObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.job)
	}
}

func TestJobTaskProvision(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful job state waiting without tasks in queue"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)

		s.args.jobState = js

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful job state running without tasks in queue"}

		job := getJobAsset(models.StateRunning, models.EmptyString)
		js := getJobStateAsset(job)

		s.args.jobState = js

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful job state waiting with tasks in queue"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, models.StateQueued, models.EmptyString)

		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.queue[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = models.StateProvision

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateRunning
		s.want.jobState.task.queue[wt.SelfLink().String()] = wt

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful job state running with tasks in queue and active"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task1 := getTaskAsset(job, models.StatePaused, models.EmptyString)
		task2 := getTaskAsset(job, models.StateQueued, models.EmptyString)
		task3 := getTaskAsset(job, models.StateQueued, models.EmptyString)

		s.args.jobState = js
		s.args.jobState.task.list[task1.SelfLink().String()] = task1
		s.args.jobState.task.list[task2.SelfLink().String()] = task2
		s.args.jobState.task.list[task2.SelfLink().String()] = task3
		s.args.jobState.task.active[task1.SelfLink().String()] = task1

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful job state running with tasks in queue and active and limit 2"}

		job := getJobAsset(models.StateRunning, models.EmptyString)
		js := getJobStateAsset(job)
		task1 := getTaskAsset(job, models.StateProvision, models.EmptyString)
		task2 := getTaskAsset(job, models.StateQueued, models.EmptyString)
		task3 := getTaskAsset(job, models.StateQueued, models.EmptyString)

		s.args.jobState = js
		s.args.jobState.job.Spec.Concurrency.Limit = 2
		s.args.jobState.task.list[task1.SelfLink().String()] = task1
		s.args.jobState.task.list[task2.SelfLink().String()] = task2
		s.args.jobState.task.list[task2.SelfLink().String()] = task3
		s.args.jobState.task.active[task1.SelfLink().String()] = task1
		s.args.jobState.task.queue[task2.SelfLink().String()] = task2
		s.args.jobState.task.queue[task3.SelfLink().String()] = task3

		wt1 := getTaskCopy(task1)
		wt1.Status.State = models.StateProvision
		wt2 := getTaskCopy(task2)
		wt2.Status.State = models.StateProvision
		wt3 := getTaskCopy(task3)
		wt3.Status.State = models.StateQueued

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateRunning
		s.want.jobState.task.active[wt1.SelfLink().String()] = wt1
		s.want.jobState.task.queue[wt2.SelfLink().String()] = wt2
		s.want.jobState.task.queue[wt3.SelfLink().String()] = wt3

		return s
	}())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := jobTaskProvision(tt.args.jobState)
			if !assert.NoError(t, err) {
				return
			}

			err = compareJobStateProperties(tt.want.jobState, tt.args.jobState)
			if !assert.NoError(t, err) {
				return
			}
		})
	}
}

func getJobAsset(state, message string) *models.Job {
	j := new(models.Job)

	j.Meta.Namespace = "test"
	j.Meta.Name = "job"
	j.Meta.SelfLink = *models.NewJobSelfLink(j.Meta.Namespace, j.Meta.Name)

	j.Status.State = state
	j.Status.Message = message

	j.Spec.Enabled = true
	j.Spec.Concurrency.Limit = 1

	return j
}

func getJobStateAsset(job *models.Job) *JobState {

	n := new(models.Node)

	n.Meta.Name = "node"
	n.Meta.Hostname = "node.local"
	n.Status.Capacity = models.NodeResources{
		Containers: 10,
		Pods:       10,
		RAM:        1000,
		CPU:        1,
		Storage:    1000,
	}
	n.Meta.SelfLink = *models.NewNodeSelfLink(n.Meta.Hostname)

	cs := cluster.NewClusterState()
	cs.SetNode(n)

	s := NewJobState(cs, job)

	return s
}

func getJobStateCopy(js *JobState) *JobState {

	j := *js.job

	njs := NewJobState(js.cluster, &j)

	njs.task.list = make(map[string]*models.Task, 0)
	for k, t := range js.task.list {
		njs.task.list[k] = &(*t)
	}

	njs.task.active = make(map[string]*models.Task, 0)
	for k, t := range js.task.active {
		task := *t
		njs.task.active[k] = &task
	}

	njs.task.queue = make(map[string]*models.Task, 0)
	for k, t := range js.task.queue {
		task := *t
		njs.task.queue[k] = &task
	}

	njs.task.finished = make([]*models.Task, 0)
	for _, t := range js.task.finished {
		task := *t
		njs.task.finished = append(njs.task.finished, &task)
	}

	njs.pod.list = make(map[string]*models.Pod, 0)
	for k, p := range js.pod.list {
		pod := *p
		njs.pod.list[k] = &pod
	}

	return njs
}

func compareJobStateProperties(old *JobState, new *JobState) error {

	if old.job != nil {
		if old.job.Status.State != new.job.Status.State {
			return fmt.Errorf("job status state is different %s != %s", old.job.Status.State, new.job.Status.State)
		}
		if old.job.Status.Message != new.job.Status.Message {
			return fmt.Errorf("job status message is different %s != %s", old.job.Status.Message, new.job.Status.Message)
		}

	}

	if len(old.task.list) != len(new.task.list) {
		return errors.New("list tasks count is different")
	}

	for k, v := range new.task.list {
		if _, ok := old.task.list[k]; !ok {
			return errors.New("list tasks is different")
		}

		if err := compareTaskProperties(old.task.list[k], v); err != nil {
			return err
		}
	}

	// check queue tasks count
	if len(old.task.queue) != len(new.task.queue) {
		return errors.New("queue tasks count is different")
	}

	for k, v := range new.task.queue {
		if _, ok := old.task.queue[k]; !ok {
			return errors.New("queue tasks is different")
		}

		if err := compareTaskProperties(old.task.queue[k], v); err != nil {
			return err
		}
	}

	// check active tasks count
	if len(old.task.active) != len(new.task.active) {
		return errors.New("active tasks count is different")
	}

	for k, v := range new.task.active {
		if _, ok := old.task.active[k]; !ok {
			return errors.New("active tasks is different")
		}

		if err := compareTaskProperties(old.task.active[k], v); err != nil {
			return err
		}
	}

	// check finished tasks count
	if len(old.task.finished) != len(new.task.finished) {
		return errors.New("finished tasks count is different")
	}

	finished := make(map[string]*models.Task, 0)
	for _, v := range old.task.finished {
		finished[v.SelfLink().String()] = v
	}

	for _, v := range new.task.finished {
		if _, ok := finished[v.SelfLink().String()]; !ok {
			return errors.New("finished tasks is different")
		}

		if err := compareTaskProperties(finished[v.SelfLink().String()], v); err != nil {
			return err
		}
	}

	// check pods count
	if len(old.pod.list) != len(new.pod.list) {
		return errors.New("pods count is different")
	}

	return nil
}
