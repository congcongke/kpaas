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

package init

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

func CheckHaproxyParameter(ipAddresses ...string) error {
	logger := logrus.WithFields(logrus.Fields{
		"error_reason": operation.ErrPara,
	})

	if len(ipAddresses) == 0 {
		logger.Errorf("%v", operation.ErrParaEmpty)
		return fmt.Errorf("%v", operation.ErrParaEmpty)
	}

	for _, ip := range ipAddresses {
		if ok := operation.CheckIPValid(ip); ok {
			continue
		}

		logrus.WithFields(logrus.Fields{
			"error_reason": operation.ErrPara,
		}).Errorf("%v", operation.ErrInvalid)
		return fmt.Errorf("%v", operation.ErrInvalid)
	}

	return nil
}
