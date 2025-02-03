package ShortMessage

import "context"

type Service interface {
	Sends(ctx context.Context, tpl string, args []string, number ...string) error
}
