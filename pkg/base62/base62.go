package base62

import "strings"

// 避免被人恶意请求，可以将顺序打乱
// const base62Str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var (
	base62Str string
)

// MustInit 要使用base62这个包必须先完成初始化
func MustInit(bs string) {
	if len(bs) == 0 {
		panic("base string is empty")
	}
	base62Str = bs
}

// Base62Encode 十进制数转化为62进制字符串
func Base62Encode(u uint64) string {
	if u == 0 {
		return "0"
	}

	byteStr := make([]byte, 0)
	for u > 0 {
		mod := u % 62
		div := u / 62
		byteStr = append(byteStr, base62Str[mod])
		u = div
	}
	return string(Reverse(byteStr))
}

// Base62Decode 62进制字符串转化为十进制数
func Base62Decode(s string) (res uint64) {
	byteStr := []byte(s)
	newByteStr := Reverse(byteStr)
	for i := 0; i < len(newByteStr); i++ {
		res += uint64(strings.Index(base62Str, string(newByteStr[i]))) * PowInt(62, i)
	}
	return res
}

func Reverse(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func PowInt(a uint64, b int) uint64 {
	res := uint64(1)
	for i := 0; i < b; i++ {
		res *= a
	}
	return res
}
