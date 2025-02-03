package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appID    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appID string, signName string) *Service {
	return &Service{
		client:   client,
		appID:    ekit.ToPtr[string](appID),
		signName: ekit.ToPtr[string](signName),
	}
}

func (s *Service) Sends(ctx context.Context, tpl string, args []string, number ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appID
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tpl)
	req.PhoneNumberSet = s.toStringPtrSlice(number)
	req.TemplateParamSet = s.toStringPtrSlice(number)
	sendSMS, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range sendSMS.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "OK" {
			return fmt.Errorf("发送失败code: %s,原因是：%s",
				*status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string { return &src })
}
