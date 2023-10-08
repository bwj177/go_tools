package common

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_TypeName(t *testing.T) {
	convey.Convey("基本用例集合", t, func() {
		tt := []struct {
			name string
			kind interface{}
			want string
		}{
			{"用例1", []int{0}, "[]int"},
			{"用例2", "5555", "string"},
			{"用例3", true, "bool"},
		}
		for _, tc := range tt {
			convey.Convey(tc.name, func() {
				got := TypeName(tc.kind)
				convey.So(got, convey.ShouldResemble, tc.want)
			})
		}
	})
}

func TestCostTime(t *testing.T) {
	defer CostTime("测试costTIme by sleep")()
	time.Sleep(5 * time.Second)
}
