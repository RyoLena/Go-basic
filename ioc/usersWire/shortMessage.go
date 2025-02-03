package usersWire

import (
	"Project/webBook_git/internal/service/ShortMessage"
	"Project/webBook_git/internal/service/ShortMessage/fakerSMS"
)

func InitFakerSMS() ShortMessage.Service {
	// InitSMSService()
	return fakerSMS.NewService()
}
