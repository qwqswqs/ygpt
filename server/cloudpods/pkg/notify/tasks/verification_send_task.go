// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tasks

import (
	"context"
	"fmt"
	"strings"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"

	api "yunion.io/x/onecloud/pkg/apis/notify"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/notify/models"
	"yunion.io/x/onecloud/pkg/util/logclient"
)

type VerificationSendTask struct {
	taskman.STask
}

func init() {
	taskman.RegisterTask(VerificationSendTask{})
}

func (self *VerificationSendTask) taskFailed(ctx context.Context, receiver *models.SReceiver, reason string) {
	log.Errorf("fail to send verification: %s", reason)
	logclient.AddActionLogWithContext(ctx, receiver, logclient.ACT_SEND_VERIFICATION, reason, self.UserCred, false)
	self.SetStageFailed(ctx, jsonutils.NewString(reason))
}

func (self *VerificationSendTask) OnInit(ctx context.Context, obj db.IStandaloneModel, body jsonutils.JSONObject) {
	receiver := obj.(*models.SReceiver)
	contactType, _ := self.Params.GetString("contact_type")
	verification, err := models.VerificationManager.Get(receiver.GetId(), contactType)
	if err != nil {
		self.taskFailed(ctx, receiver, fmt.Sprintf("VerificationManager.Get for receiver_id %q and contact_type %q: %s", receiver.GetId(), contactType, err.Error()))
		return
	}
	if err != nil {
		self.taskFailed(ctx, receiver, fmt.Sprintf("fail to get contact(type: %s): %s", contactType, err.Error()))
		return
	}

	emailMsg := &api.SEmailMessage{}
	notifyReceiver := api.SNotifyReceiver{}
	// build message
	var message string
	switch contactType {
	case api.EMAIL:
		info, err := models.TemplateManager.GetCompanyInfo(ctx)
		if err != nil {
			self.taskFailed(ctx, receiver, fmt.Sprintf("fail to get company info: %s", err.Error()))
			return
		}
		data := struct {
			models.SCompanyInfo
			ReceiverName string
			Code         string
		}{
			ReceiverName: receiver.Name,
			Code:         verification.Token,
			SCompanyInfo: info,
		}
		message = jsonutils.Marshal(data).String()
		notifyReceiver = api.SNotifyReceiver{
			Contact:  receiver.Email,
			DomainId: receiver.DomainId,
		}
		emailMsg.To = append(emailMsg.To, receiver.Email)
	case api.MOBILE:
		message = fmt.Sprintf("[\"%s\"]", verification.Token)
		mobileArr := strings.Split(receiver.Mobile, " ")
		mobile := strings.Join(mobileArr, "")
		notifyReceiver = api.SNotifyReceiver{
			Contact:  mobile,
			DomainId: receiver.DomainId,
		}
	default:
		// no way
	}
	tLang, err := receiver.GetTemplateLang(ctx)
	if err != nil {
		self.taskFailed(ctx, receiver, fmt.Sprintf("unable to GetTemplateLang for receiver %q: %v", receiver.Id, err))
	}
	param, err := models.TemplateManager.FillWithTemplate(ctx, tLang, api.SsNotification{
		ContactType: contactType,
		Topic:       "VERIFY",
		Message:     message,
	})
	if err != nil {
		self.taskFailed(ctx, receiver, err.Error())
		return
	}
	param.RemoteTemplateParam = api.SRemoteTemplateParam{Code: verification.Token}
	param.Receivers = notifyReceiver
	if len(emailMsg.Subject) == 0 {
		emailMsg.Subject = param.Title
	}
	if len(emailMsg.Body) == 0 {
		emailMsg.Body = param.Message
	}
	param.EmailMsg = *emailMsg
	driver := models.GetDriver(contactType)
	err = driver.Send(ctx, param)
	// err = models.NotifyService.Send(ctx, contactType, param)
	if err != nil {
		self.taskFailed(ctx, receiver, err.Error())
		return
	}
	logclient.AddActionLogWithContext(ctx, receiver, logclient.ACT_SEND_VERIFICATION, "", self.UserCred, true)
	self.SetStageComplete(ctx, nil)
}
