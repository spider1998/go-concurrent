/*将一个或多个完成的 channel 合并到 完成
channel 中，该 channel 在任何组件 channel 关闭时关闭*/
package main

import (
	`fmt`
	`time`
)

func main()  {
	var or func(channels ...<-chan interface{})<-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}


	sig := func(after time.Duration) <-chan interface{}{
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v",time.Since(start))

}

/*这种模式在系统中的模块交汇处非常有用。在这些交汇处，调用堆
中应该有复数种的用来取消 goroutine 的决策树 使用 or 函数，可以简
地将它们组合在 起并将其传递给堆枝。*/

