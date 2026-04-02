package context

import (
	"context"
	"fmt"
	"time"
)

// Basic demonstrates basic context usage
func Basic() {
	// 创建根 context
	ctx := context.Background()

	// 带取消的 context
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(2 * time.Second)
		cancel() // 2秒后取消
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Cancelled:", ctx.Err()) // context canceled
	case <-time.After(3 * time.Second):
		fmt.Println("Timeout")
	}
}
