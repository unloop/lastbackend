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
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/logger"
	"sync"
)

const logTaskPrefix = "state:task:>"

type TaskState struct {
	lock  sync.RWMutex
	tasks map[string]models.NodeTask
}

func (s *TaskState) AddTask(key string, task *models.NodeTask) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s add cancel func pod: %s", logTaskPrefix, key)
	s.tasks[key] = *task
}

func (s *TaskState) GetTask(key string) *models.NodeTask {
	log := logger.WithContext(context.Background())
	log.Debugf("%s get cancel func pod: %s", logTaskPrefix, key)

	if _, ok := s.tasks[key]; ok {
		t := s.tasks[key]
		return &t
	}

	return nil
}

func (s *TaskState) DelTask(pod *models.Pod) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s del cancel func pod: %s", logTaskPrefix, pod.SelfLink())
	delete(s.tasks, pod.Meta.Name)
}
