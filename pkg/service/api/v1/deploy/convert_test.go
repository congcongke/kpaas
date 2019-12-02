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

package deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/common"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
)

func TestConvertModelAuthenticationTypeToAPIAuthenticationType(t *testing.T) {

	assert.Equal(t, api.AuthenticationTypePassword, convertModelAuthenticationTypeToAPIAuthenticationType(wizard.AuthenticationTypePassword))
	assert.Equal(t, api.AuthenticationTypePrivateKey, convertModelAuthenticationTypeToAPIAuthenticationType(wizard.AuthenticationTypePrivateKey))
	assert.Equal(t, api.AuthenticationType("unknown(OtherType)"), convertModelAuthenticationTypeToAPIAuthenticationType("OtherType"))
}

func TestConvertModelLabelToAPILabel(t *testing.T) {

	assert.EqualValues(t, api.Label{
		Key:   "key",
		Value: "value",
	}, convertModelLabelToAPILabel(&wizard.Label{
		Key:   "key",
		Value: "value",
	}))
}

func TestConvertModelAnnotationToAPIAnnotation(t *testing.T) {

	assert.EqualValues(t,
		api.Annotation{
			Key:   "key",
			Value: "value",
		}, convertModelAnnotationToAPIAnnotation(&wizard.Annotation{
			Key:   "key",
			Value: "value",
		}),
	)
}

func TestConvertModelTaintToAPITaint(t *testing.T) {

	assert.EqualValues(t,
		api.Taint{
			Key:    "key",
			Value:  "value",
			Effect: api.TaintEffectNoSchedule,
		},
		convertModelTaintToAPITaint(&wizard.Taint{
			Key:    "key",
			Value:  "value",
			Effect: wizard.TaintEffectNoSchedule,
		}),
	)
}

func TestConvertModelTaintEffectToAPITaintEffect(t *testing.T) {

	assert.Equal(t, api.TaintEffectNoSchedule, convertModelTaintEffectToAPITaintEffect(wizard.TaintEffectNoSchedule))
	assert.Equal(t, api.TaintEffectNoExecute, convertModelTaintEffectToAPITaintEffect(wizard.TaintEffectNoExecute))
	assert.Equal(t, api.TaintEffectPreferNoSchedule, convertModelTaintEffectToAPITaintEffect(wizard.TaintEffectPreferNoSchedule))
	assert.Equal(t, api.TaintEffect("unknown(OtherType)"), convertModelTaintEffectToAPITaintEffect("OtherType"))
}

func TestConvertModelCheckResultToAPICheckResult(t *testing.T) {

	assert.Equal(t, api.CheckResultNotRunning, convertModelCheckResultToAPICheckResult(wizard.CheckResultNotRunning))
	assert.Equal(t, api.CheckResultChecking, convertModelCheckResultToAPICheckResult(wizard.CheckResultChecking))
	assert.Equal(t, api.CheckResultPassed, convertModelCheckResultToAPICheckResult(wizard.CheckResultPassed))
	assert.Equal(t, api.CheckResultFailed, convertModelCheckResultToAPICheckResult(wizard.CheckResultFailed))
	assert.Equal(t, api.CheckResult("unknown(OtherType)"), convertModelCheckResultToAPICheckResult("OtherType"))
}

func TestConvertModelErrorToAPIError(t *testing.T) {

	assert.EqualValues(t,
		&api.Error{
			Reason:     "reason",
			Detail:     "detail",
			FixMethods: "fixMethods",
			LogId:      1234,
		},
		convertModelErrorToAPIError(&common.FailureDetail{
			Reason:     "reason",
			Detail:     "detail",
			FixMethods: "fixMethods",
			LogId:      1234,
		}),
	)
}

func TestConvertModelDeployClusterStatusToAPIDeployClusterStatus(t *testing.T) {

	assert.Equal(t, api.DeployClusterStatusNotRunning, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusNotRunning))
	assert.Equal(t, api.DeployClusterStatusRunning, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusRunning))
	assert.Equal(t, api.DeployClusterStatusSuccessful, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusSuccessful))
	assert.Equal(t, api.DeployClusterStatusFailed, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusFailed))
	assert.Equal(t, api.DeployClusterStatusWorkedButHaveError, convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.DeployClusterStatusWorkedButHaveError))
	assert.Equal(t, api.DeployClusterStatus("unknown(OtherType)"), convertModelDeployClusterStatusToAPIDeployClusterStatus("OtherType"))
}

func TestConvertModelMachineRoleToAPIMachineRole(t *testing.T) {

	assert.Equal(t, api.MachineRoleMaster, convertModelMachineRoleToAPIMachineRole(wizard.MachineRoleMaster))
	assert.Equal(t, api.MachineRoleWorker, convertModelMachineRoleToAPIMachineRole(wizard.MachineRoleWorker))
	assert.Equal(t, api.MachineRoleEtcd, convertModelMachineRoleToAPIMachineRole(wizard.MachineRoleEtcd))
	assert.Equal(t, api.MachineRole("unknown(OtherType)"), convertModelMachineRoleToAPIMachineRole("OtherType"))
}

func TestConvertModelDeployStatusToAPiDeployStatus(t *testing.T) {

	assert.Equal(t, api.DeployStatusPending, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusPending))
	assert.Equal(t, api.DeployStatusDeploying, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusDeploying))
	assert.Equal(t, api.DeployStatusCompleted, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusCompleted))
	assert.Equal(t, api.DeployStatusFailed, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusFailed))
	assert.Equal(t, api.DeployStatusAborted, convertModelDeployStatusToAPIDeployStatus(wizard.DeployStatusAborted))
	assert.Equal(t, api.DeployStatus("unknown(OtherType)"), convertModelDeployStatusToAPIDeployStatus("OtherType"))
}

func TestConvertAPIMachineRoleToModelMachineRole(t *testing.T) {

	assert.Equal(t, wizard.MachineRoleMaster, convertAPIMachineRoleToModelMachineRole(api.MachineRoleMaster))
	assert.Equal(t, wizard.MachineRoleWorker, convertAPIMachineRoleToModelMachineRole(api.MachineRoleWorker))
	assert.Equal(t, wizard.MachineRoleEtcd, convertAPIMachineRoleToModelMachineRole(api.MachineRoleEtcd))
	assert.Equal(t, wizard.MachineRole("unknown(OtherType)"), convertAPIMachineRoleToModelMachineRole("OtherType"))
}

func TestConvertAPITaintEffectToModelTaintEffect(t *testing.T) {

	assert.Equal(t, wizard.TaintEffectNoSchedule, convertAPITaintEffectToModelTaintEffect(api.TaintEffectNoSchedule))
	assert.Equal(t, wizard.TaintEffectNoExecute, convertAPITaintEffectToModelTaintEffect(api.TaintEffectNoExecute))
	assert.Equal(t, wizard.TaintEffectPreferNoSchedule, convertAPITaintEffectToModelTaintEffect(api.TaintEffectPreferNoSchedule))
	assert.Equal(t, wizard.TaintEffect("unknown(OtherType)"), convertAPITaintEffectToModelTaintEffect("OtherType"))
}

func TestConvertAPIAuthenticationTypeToModelAuthenticationType(t *testing.T) {

	assert.Equal(t, wizard.AuthenticationTypePassword, convertAPIAuthenticationTypeToModelAuthenticationType(api.AuthenticationTypePassword))
	assert.Equal(t, wizard.AuthenticationTypePrivateKey, convertAPIAuthenticationTypeToModelAuthenticationType(api.AuthenticationTypePrivateKey))
	assert.Equal(t, wizard.AuthenticationType("unknown(OtherType)"), convertAPIAuthenticationTypeToModelAuthenticationType("OtherType"))
}