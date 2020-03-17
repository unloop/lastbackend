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
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logTaskPrefix = "node:runtime:task"
)

func (r Runtime) taskExecute(ctx context.Context, pod string, task types.SpecRuntimeTask, m types.ContainerManifest, ps *types.PodStatus) error {

	status := ps.AddTask(task.Name)
	status.SetStarted()

	r.state.Pods().SetPod(pod, ps)
	log.V(logLevel).Debugf("%s task %s start", logTaskPrefix, task.Name)

	m.Name = ""
	m.Labels[types.ContainerTypeRuntime] = types.ContainerTypeRuntimeTask

	var (
		c   types.PodContainer
		err error
	)

	m.RestartPolicy.Policy = "no"
	status.AddTaskCommandContainer(&c)
	r.state.Pods().SetPod(pod, ps)

	//========================================================================================
	// create container ======================================================================
	//========================================================================================

	c.ID, err = r.cri.Create(ctx, &m)
	if err != nil {
		switch err {
		case context.Canceled:
			log.Errorf("%s stop creating container: %s", logTaskPrefix, err.Error())
		}

		log.Errorf("%s can-not create container: %s", logTaskPrefix, err)
		status.SetExited(true, err.Error())
		r.state.Pods().SetPod(pod, ps)
		return err

	}

	c.State.Created = types.PodContainerStateCreated{
		Created: time.Now().UTC(),
	}

	//========================================================================================
	// start container =======================================================================
	//========================================================================================
	log.V(logLevel).Debugf("%s container created: %s", logTaskPrefix, c.ID)

	if err := r.cri.Start(ctx, c.ID); err != nil {

		log.Errorf("%s can-not start container: %s", logTaskPrefix, err)
		switch err {
		case context.Canceled:
			log.Errorf("%s stop starting container err: %s", logTaskPrefix, err.Error())
		}

		log.Errorf("%s can-not start container: %s", logTaskPrefix, err)
		c.State.Error = types.PodContainerStateError{
			Error:   true,
			Message: err.Error(),
			Exit: types.PodContainerStateExit{
				Timestamp: time.Now().UTC(),
			},
		}
		status.SetExited(true, err.Error())
		r.state.Pods().SetPod(pod, ps)

		return r.taskCommandFinish(ctx, &c)
	}

	c.Ready = true
	c.State.Started = types.PodContainerStateStarted{
		Started:   true,
		Timestamp: time.Now().UTC(),
	}

	//========================================================================================
	// wait container ========================================================================
	//========================================================================================

	//go func() {
	//	req, err := r.cri.Logs(ctx, c.ID, true, true, true)
	//	if err != nil {
	//		log.Errorf("%s error get logs stream %s", logPodPrefix, err)
	//		return
	//	}
	//
	//	io.Copy(os.Stdout, req)
	//}()

	log.V(logLevel).Debugf("%s container wait: %s", logTaskPrefix, c.ID)
	if err := r.cri.Wait(ctx, c.ID); err != nil {
		log.Errorf("%s error: %s", logTaskPrefix, err.Error())
		c.State.Error = types.PodContainerStateError{
			Error:   true,
			Message: err.Error(),
			Exit: types.PodContainerStateExit{
				Timestamp: time.Now().UTC(),
			},
		}
		status.SetExited(true, err.Error())
		r.state.Pods().SetPod(pod, ps)
		return r.taskCommandFinish(ctx, &c)
	}

	info, err := r.cri.Inspect(ctx, c.ID)
	if err != nil {
		log.Errorf("%s error: %s", logTaskPrefix, err.Error())
		c.State.Error = types.PodContainerStateError{
			Error:   true,
			Message: err.Error(),
			Exit: types.PodContainerStateExit{
				Timestamp: time.Now().UTC(),
			},
		}
		status.SetExited(true, err.Error())
		r.state.Pods().SetPod(pod, ps)
		return r.taskCommandFinish(ctx, &c)
	}

	if err := r.containerInspect(context.Background(), &c); err != nil {
		log.Errorf("%s inspect container after create: err %s", logServicePrefix, err.Error())
		return err
	}

	c.Ready = true
	c.State.Stopped = types.PodContainerStateStopped{
		Stopped: true,
		Exit: types.PodContainerStateExit{
			Code:      info.ExitCode,
			Timestamp: time.Now().UTC(),
		},
	}

	if info.ExitCode != 0 {
		status.SetExited(true, info.Error)
		r.state.Pods().SetPod(pod, ps)
		return r.taskCommandFinish(ctx, &c)
	}

	if err := r.taskCommandFinish(ctx, &c); err != nil {
		log.Errorf("%s task %s cleanup failed: %s", logTaskPrefix, task.Name, err.Error())
	}

	status.SetExited(false, types.EmptyString)
	r.state.Pods().SetPod(pod, ps)
	return nil
}

func (r Runtime) taskCommandFinish(ctx context.Context, c *types.PodContainer) error {

	log.V(logLevel).Debugf("%s container remove: %s", logTaskPrefix, c.ID)
	if err := r.cri.Remove(ctx, c.ID, true, true); err != nil {
		log.Errorf("%s error: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}
