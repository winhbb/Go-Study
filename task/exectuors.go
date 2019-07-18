package task

import "time"

/**
 * author hebingbin
 * 单任务执行器
 */
type SimpleExecutor struct{}

func NewSimpleExecutor() *SimpleExecutor {
	return &SimpleExecutor{}
}

func (e *SimpleExecutor) Execute(task Runnable) {
	go task()
}

func (e *SimpleExecutor) Submit(task Callable) Future {
	future := NewChanFuture()
	go func() {
		defer future.Close()
		future.resultChan <- task()
	}()
	return future
}

func (e *SimpleExecutor) ScheduleAtFixedDelay(task Callable, d time.Duration) Future {
	future := NewChanFuture()
	go func() {
		//关闭所有使用到的channel，避免阻塞
		defer future.Close()
		for {
			select {
			case <-future.cancel: //超时或者用户手动出发
				future.canceled = true
				return
			default:
				future.resultChan <- task()
			}
			time.Sleep(d)
		}
	}()
	return future
}
