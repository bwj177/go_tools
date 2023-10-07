package string

import "math/rand"

type stringUtil struct {
}

var StringUtil = new(stringUtil)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyz123456789"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// RandomStr
//
//	@Description: 随机生成字符串
//	@receiver su
//	@param strLen 字符串长度
//	@return string
func (su *stringUtil) RandomStr(strLen int) string {
	b := make([]byte, strLen)

	for i, cache, remain := strLen-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
