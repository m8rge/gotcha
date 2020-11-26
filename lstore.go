package gotcha

import (
	"context"
	"sync"

	"github.com/modern-go/gls"
)

type lskey string

const glskey lskey = "glskey"

type Tracer func(Context)

func Trace(ctx context.Context, t Tracer, opts ...ContextOpt) {
	gls.WithEmptyGls(func() {
		ctx := NewContext(ctx, opts...)
		gls.Set(glskey, ctx)
		t(ctx)
	})()
}

func GoTrace(ctx context.Context, wg *sync.WaitGroup, t Tracer, opts ...ContextOpt) {
	if wg != nil {
		wg.Add(1)
	}
	go func() {
		if wg != nil {
			defer wg.Done()
		}
		Trace(ctx, t, opts...)
	}()
}

func trackAlloc(bytes, objects int) {
	if v := gls.Get(glskey); v != nil {
		if ctx, ok := v.(Tracker); ok {
			ctx.Add(int64(bytes), int64(objects), 1)
		}
	}
}
