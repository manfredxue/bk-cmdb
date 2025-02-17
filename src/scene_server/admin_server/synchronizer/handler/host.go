/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handler

import (
	"context"
	"fmt"

	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/synchronizer/meta"
	"configcenter/src/scene_server/admin_server/synchronizer/utils"
)

// HandleHostSync do sync host of one business
func (ih *IAMHandler) HandleHostSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)

	// step1 get instances by business from core service
	bizID := businessSimplify.BKAppIDField
	hosts, err := ih.authManager.CollectHostByBusinessID(context.Background(), *header, bizID)
	if err != nil {
		blog.Errorf("get host by business %d failed, err: %+v", businessSimplify.BKAppIDField, err)
		return err
	}
	resources, err := ih.authManager.MakeResourcesByHosts(context.Background(), *header, authmeta.EmptyAction, hosts...)
	if err != nil {
		blog.Errorf("make host resources failed, bizID: %d, err: %+v", businessSimplify.BKAppIDField, err)
		return err
	}

	// step2 get host by business from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.HostInstance,
		},
		BusinessID: bizID,
	}

	taskName := fmt.Sprintf("sync host for business: %d", businessSimplify.BKAppIDField)
	iamIDPrefix := ""
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
