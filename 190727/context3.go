//创建类型安全
package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessReq("jane", "abc123")
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
	return c.Value(ctxToken).(string)
}

func ProcessReq(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxToken, authToken)
	Handle(ctx)
}

func Handle(ctx context.Context) {
	fmt.Printf(
		"handle response for %v (ahth: %v)",
		UserID(ctx),
		AuthToken(ctx),
	)
}
