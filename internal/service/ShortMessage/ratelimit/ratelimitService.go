package ratelimit

import (
	"Project/internal/service/ShortMessage"
	"Project/pkg/ratelimit"
	"context"
	"errors"
	"fmt"
)

var LimitErr = errors.New("已经被限流")

type LimitSMSService struct {
	svc     ShortMessage.Service
	limiter ratelimit.Limiter
}

func NewLimitSMSService(svc ShortMessage.Service, limiter ratelimit.Limiter) *LimitSMSService {
	return &LimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (l *LimitSMSService) Sends(ctx context.Context, tpl string, args []string, number ...string) error {
	//装饰器部分
	//例如：
	limited, err := l.limiter.Limit(ctx, "key")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流异常 %w", err)
	}
	if limited {
		//限流了s
		return LimitErr
	}

	//执行的主体
	err = l.Sends(ctx, tpl, args)

	//这里也可以加上装饰器
	//限流
	//.....
	//
	return err
}
