// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package task

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
)

// nodeCheckProcessor implements the specific logic for the node check task.
type nodeCheckProcessor struct {
}

// Spilt the task into one or more node check actions
func (p *nodeCheckProcessor) SplitTask(t Task) error {
	if err := p.verifyTask(t); err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: t.GetName(),
	})

	logger.Debug("Start to split node check task")

	checkTask := t.(*nodeCheckTask)

	// split task into actions: will create a action for every node, the action type
	// is NodeCheckAction
	actions := make([]action.Action, 0, len(checkTask.nodeConfigs))
	for _, subConfig := range checkTask.nodeConfigs {
		actionCfg := &action.NodeCheckActionConfig{
			NodeCheckConfig: subConfig,
			LogFileBasePath: checkTask.logFilePath,
		}
		act, err := action.NewNodeCheckAction(actionCfg)
		if err != nil {
			return err
		}
		actions = append(actions, act)
	}
	checkTask.actions = actions

	logrus.Debugf("Finish to split node check task: %d actions", len(actions))
	return nil
}

// Verify if the task is valid.
func (p *nodeCheckProcessor) verifyTask(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	nodeCheckTask, ok := t.(*nodeCheckTask)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, t)
	}

	if len(nodeCheckTask.nodeConfigs) == 0 {
		return fmt.Errorf("nodeConfigs is empty")
	}

	return nil
}
