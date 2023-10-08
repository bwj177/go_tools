package common

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
	"context"
	"sync"
	"errors"
)

// MaxInt64 返回数组中最大的
func MaxInt64(i int64, nums ...int64) int64 {
	max := i
	for _, num := range nums {
		if max < num {
			max = num
		}
	}
	return max
}

// MinInt64 返回数组中最小的
func MinInt64(i int64, nums ...int64) int64 {
	min := i
	for _, num := range nums {
		if min > num {
			min = num
		}
	}
	return min
}

func TypeName(data interface{}) string {
	if data == nil {
		return ""
	}
	t := reflect.TypeOf(data)
	for t.Kind() == reflect.Ptr { // 解引用嵌套指针
		t = t.Elem()
	}

	return t.String()
}

// use before function like defer CostTime()() get function cost time
func CostTime(funcName string) func() {
	now := time.Now()
	return func() {
		costTime := time.Now().Sub(now)
		logrus.Infof("%v total cost time :%v", funcName, costTime)
	}
}




// -----------------------------------wairGrop decoration--------------

func NewWaitGroup() *WaitGroup {
	return new(WaitGroup)
}

func WithContext(ctx context.Context) (*WaitGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &WaitGroup{cancel: cancel}, ctx
}

type WaitGroup struct {
	cancel  func()   // context中取消信号
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

func (g *WaitGroup) Run(fn func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		if err := fn(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func (g *WaitGroup) RunWithRecover(fn func() error) {
	g.wg.Add(1)
	GoWithRecover(func() { // handler
		defer g.wg.Done()
		if err := fn(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}, func(r interface{}) { // recoverHandler
		g.errOnce.Do(func() {
			g.err = errors.New("panic")
			if g.cancel != nil {
				g.cancel()
			}
		})
	})
}

//
// Wait
//  @Description: 在原有wait方法中装饰了cancel的执行
//  @receiver g
//  @return error
//
func (g *WaitGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}


//
// GoWithRecover
//  @Description: 携带捕获panic的函数执行
//  @param handler 执行函数
//  @param recoverHandler 捕获panic后的处理函数
//
func GoWithRecover(handler func(), recoverHandler func(r interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Error(r)
				if recoverHandler != nil {
					go func() {
						defer func() {
							if p := recover(); p != nil {
								logrus.Error(p)
							}
						}()
						recoverHandler(r)
					}()
				}
			}
		}()
		handler()
	}()
}
