package task

import "time"

/**
 * author hebingbin
 * 简易的任务执行框架
 */
type Runnable func()

type Callable func() Result

type Cancelable interface {
	Cancel() bool
}

type Executor interface {
	Execute(task Runnable)
	Submit(task Callable) Future
}

type Future interface {
	Cancel()
	IsCancelled() bool
	IsDone() bool
	Get() Result
	WaitAndGet(time.Duration) Result
}

type Result interface {
	GetResult()interface{}
	GetError() error
}

type ResultWithError struct {
	Data interface{}
	Err  error
}

func (result *ResultWithError) GetError() error {
	return result.Err
}

func (result *ResultWithError) GetResult() interface{} {
	return result.Data
}