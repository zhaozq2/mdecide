package ctxdata

import (
	"context"
)

const (
	CtxKeyUserId   = "userId"
	CtxKeyUserName = "userName"
)

func GetUserIdFromCtx(ctx context.Context) int64 {
	v := ctx.Value(CtxKeyUserId)
	if v == nil {
		return 0
	}
	if userId, ok := v.(int64); ok {
		return userId
	}
	return 0
}

func GetUserNameFromCtx(ctx context.Context) string {
	v := ctx.Value(CtxKeyUserName)
	if v == nil {
		return ""
	}
	if userName, ok := v.(string); ok {
		return userName
	}
	return ""
}
