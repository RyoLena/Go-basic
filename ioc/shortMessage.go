package ioc

import (
	"Project/internal/service/ShortMessage"
	"Project/internal/service/ShortMessage/fakerSMS"
)

func InitFakeSMS() ShortMessage.Service {
	// InitSMSService()
	return fakerSMS.NewService()
}
