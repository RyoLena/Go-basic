package failover

import (
	"Project/internal/service/ShortMessage"
	"context"
	"sync/atomic"
)

type TimeOutFailOver struct {
	//假如是一个服务商集群
	svcs []ShortMessage.Service

	//下标索引
	idx int32

	//连续超时次数
	cnt int32

	//阈值，单个服务商超时次数大于这个就切换
	threshold int32
}

func (t *TimeOutFailOver) Sends(ctx context.Context, tpl string, args []string, number ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)

	if cnt > t.threshold {
		//切换
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//成功往后移动一位
			atomic.StoreInt32(&t.cnt, 0)
		}
		//出现并发，其他人等待,将切换完的存入idx下标
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Sends(ctx, tpl, args, number...)
	switch err {
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:
		//不知道什么错误

		//在这里也可以选择换下一个，语义则是，
		// -超时错误，可能是偶发的，所以我再试一试
		// -非超时错误，直接下一个

	}
	return err
}
