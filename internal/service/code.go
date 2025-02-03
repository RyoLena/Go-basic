package service

import (
	"Project/webBook_git/internal/respository"
	"Project/webBook_git/internal/service/ShortMessage"
	"context"
	"fmt"
	"math/rand"
)

var ErrCodeVerifyToManyTimes = respository.ErrCodeVerifyTooManyTimes

type CodeService struct {
	codeRepo *respository.CodeRepo
	smsSvc   ShortMessage.Service
}

func NewCodeService(codeRepo *respository.CodeRepo, smsSvc ShortMessage.Service) *CodeService {
	return &CodeService{
		codeRepo: codeRepo,
		smsSvc:   smsSvc,
	}
}

func (codeSvc *CodeService) SendCode(ctx context.Context, biz, phone string) error {
	//生成验证码
	code := codeSvc.generateCode()
	//保存
	fmt.Println(code)

	err := codeSvc.codeRepo.Store(ctx, biz, phone, code)
	if err != nil {
		fmt.Println("走这里出错了嘛")
		return err
	}
	//发送出去
	err = codeSvc.smsSvc.Sends(ctx, "tplID", []string{code}, phone)
	return err
}

func (codeSvc *CodeService) VerifyCode(ctx context.Context, biz, phone, input string) error {
	return codeSvc.codeRepo.Verify(ctx, biz, phone, input)
}

func (codeSvc *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
