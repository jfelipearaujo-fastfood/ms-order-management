package tests

import "context"

type ctxKeyType string

const ctxKey ctxKeyType = "tests"

func enrichContext[T any](ctx context.Context, data *T) context.Context {
	if ctx == nil {
      return context.Background()
  }
  return context.WithValue(ctx, ctxKey, data)
}

func fromContext[T any](ctx context.Context) *T {
  if ctx == nil {
    return new(T)
  }
  data, ok := ctx.Value(ctxKey).(*T)
  if !ok {
    return new(T)
  }
  return data
}
