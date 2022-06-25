package utils

import "context"

type ReqCtxKey string

var reqCtx ReqCtxKey = "req-id"

func CreateReqIDCtx(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, reqCtx, reqID)
}

func GetReqID(ctx context.Context) string {
	id, ok := ctx.Value(reqCtx).(string)
	if !ok {
		return ""
	}
	return id
}
