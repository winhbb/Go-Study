package task

import (
	"time"
)

/**
 * author hebingbin
 * 用于控制Executor的执行，设置超时时间等
 */
type ChanFuture struct {
	resultChan chan Result
	cancel     chan struct{}
	done       bool
	canceled   bool
}

func NewChanFuture() *ChanFuture {
	return &ChanFuture{
		resultChan: make(chan Result),
		cancel:     make(chan struct{}),
	}
}

/**
 * goroutine被设计为匿名worker，不存在id之类可以唯一标识协程的属性，因此也不存在可以从协程外部终止协程的方法
 * 只能通过channel、context等方式通知协程，让协程主动结束，但是如果协程在接收通知前，自身已死循环或发生I/O阻塞等现象，
 * 可能会导致协程泄漏，所以使用协程时一定要避免这些情况。
 * 目前Cancel的作用，实际上仅仅是解除Get方法的等待，任务协程有可能依旧在执行。
 */
func (f *ChanFuture) Cancel() {
	f.Close()
}

func (f *ChanFuture) IsCancelled() bool {
	return f.canceled
}

func (f *ChanFuture) IsDone() bool {
	return f.done
}

/**
 * 阻塞直到任务执行结束
 */
func (f *ChanFuture) Get() Result {
	return <-f.resultChan
}

/**
 * 等待指定的到达指定时间强制终止任务或者任务在指定时间内结束
 */
func (f *ChanFuture) WaitAndGet(d time.Duration) Result {
	var result Result
	//定时器
	timer := time.NewTimer(d)
	select {
	case <-timer.C:
		//超时发送终止请求
		f.cancel <- struct{}{}
		return nil
	case result = <-f.resultChan:
		//清理timer
		timer.Stop()
	}
	return result
}

func (f *ChanFuture) Close() {
	select {
	case _, ok := <-f.resultChan:
		if ok {
			close(f.resultChan)
		}
	default:
		close(f.resultChan)
	}
	select {
	case _, ok := <-f.cancel:
		if ok {
			close(f.cancel)
		}
	default:
		close(f.cancel)
	}
}
