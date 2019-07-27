//存储和检索请求范围数据的Context的数据包
package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("jane", "abc123")
}

func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), "userID", userID)
	ctx = context.WithValue(ctx, "authToken", authToken)
	HanleResponse(ctx)
}

func HanleResponse(ctx context.Context) {
	fmt.Printf(
		"handlilng response for %v (%v)",
		ctx.Value("userID"),
		ctx.Value("authToken"),
	)
}
