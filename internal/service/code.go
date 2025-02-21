package service

import (
	"Project/internal/respository"
	"Project/internal/service/ShortMessage"
	"context"
	"fmt"
	"math/rand"
)

var ErrCodeVerifyToManyTimes = respository.ErrCodeVerifyTooManyTimes

type CodeService interface {
	SendCode(ctx context.Context, biz, phone string) error
	VerifyCode(ctx context.Context, biz, phone, input string) error
}
type CodeServiceImpl struct {
	codeRepo respository.CodeCodeRepository
	smsSvc   ShortMessage.Service
}

func NewCodeService(codeRepo respository.CodeCodeRepository, smsSvc ShortMessage.Service) CodeService {
	return &CodeServiceImpl{
		codeRepo: codeRepo,
		smsSvc:   smsSvc,
	}
}

func (codeSvc *CodeServiceImpl) SendCode(ctx context.Context, biz, phone string) error {
	//生成验证码
	code := codeSvc.generateCode()
	//保存
	fmt.Println(code)

	err := codeSvc.codeRepo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	//发送出去
	err = codeSvc.smsSvc.Sends(ctx, "tplID", []string{code}, phone)
	return err
}

func (codeSvc *CodeServiceImpl) VerifyCode(ctx context.Context, biz, phone, input string) error {
	return codeSvc.codeRepo.Verify(ctx, biz, phone, input)
}

func (codeSvc *CodeServiceImpl) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
