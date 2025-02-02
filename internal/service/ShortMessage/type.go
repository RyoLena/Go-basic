package ShortMessage

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, args []string, number ...string) error
}
