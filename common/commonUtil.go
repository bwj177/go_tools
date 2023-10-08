package common

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"time"
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

// ================! carry context、reconvery waitGroup =================================
func NewWaitGroup() *WaitGroup {
	return new(WaitGroup)
}

func WithContext(ctx context.Context) (*WaitGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &WaitGroup{cancel: cancel}, ctx
}

type WaitGroup struct {
	cancel  func()         //携带cancel信号，当错误发出cancel信号
	wg      sync.WaitGroup // 原生wg
	errOnce sync.Once      // 程序错误后执行的
	err     error          // 接受fn中err
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

// 装饰wait() 携带cancel 的wait方法
func (g *WaitGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// GoWithRecover
//
//	@Description: 携带捕获panic的func执行
//	@param handler ：func执行
//	@param recoverHandler 捕获panic后的 func执行
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
