package usersWire

import (
	"Project/internal/service/ShortMessage"
	"Project/internal/service/ShortMessage/fakerSMS"
)

func InitFakerSMS() ShortMessage.Service {
	// InitSMSService()
	return fakerSMS.NewService()
}
