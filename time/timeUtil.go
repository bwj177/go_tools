package time

import (
	"github.com/golang-module/carbon"
	"math"
	"time"
)

type timeUtil struct{}

var TimeUtil = new(timeUtil)

// DaysBetweenByStr 两个日期相差的天数
func (t *timeUtil) DaysBetweenByStr(t1, t2 string) int64 {
	return t.DaysBetween(t.ParseStrToDate(t1), t.ParseStrToDate(t2))
}

// DaysBetween 两个日期相差的天数
func (t *timeUtil) DaysBetween(t1, t2 time.Time) int64 {
	// 原生没有实现
	diff := int64(math.Abs((t1.Sub(t2)).Hours()))
	remainder := diff % 24
	if remainder == 0 {
		return diff / 24
	}
	return diff/24 + 1
}

// 将string格式的date解析为时间戳
func (t *timeUtil) ParseStrToDate(dateStr string) time.Time {
	tt, _ := time.ParseInLocation("2006/01/02", dateStr, time.Local)
	return tt
}

// 获取当前时间戳的date（String格式）
func (t *timeUtil) GetCurrDateStr() string {
	return time.Now().Format("2006/01/02")
}

// 获取今日整点00：00的时间戳
func (t *timeUtil) GetZeroTime() time.Time {
	now := carbon.Now()
	todayMilli := now.StartOfDay().Carbon2Time()
	return todayMilli
}
