package failover

import (
	"Project/internal/service/ShortMessage"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

type FailOverService struct {
	svcs []ShortMessage.Service

	idx uint64
}

func NewFailOverService(svc ...ShortMessage.Service) *FailOverService {
	return &FailOverService{
		svcs: svc,
	}
}

func (f *FailOverService) Sends(ctx context.Context, tpl string, args []string, number ...string) error {
	for _, svc := range f.svcs {
		err := svc.Sends(ctx, tpl, args, number...)
		if err == nil {
			return nil
		}

		log.Println("切换服务商")
	}
	return errors.New("百分之百网络断了")
}

func (f *FailOverService) SendsV1(ctx context.Context, tpl string, args []string, number ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Sends(ctx, tpl, args, number...)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return err
		}
		//其他情况会走在这里需要打印日志
	}
	return errors.New("发送失败哦，所有产商都已经尝试过了")
}
