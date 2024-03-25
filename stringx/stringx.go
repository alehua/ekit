package stringx

import "unsafe"

func UnsafeToBytes(val string) []byte {
	sh := (*[2]uintptr)(unsafe.Pointer(&val))
	bh := [3]uintptr{sh[0], sh[1], sh[1]}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// UnsafeToString 非安全 []byte 转 string 他必须遵守上述规则
func UnsafeToString(val []byte) string {
	bh := (*[3]uintptr)(unsafe.Pointer(&val))
	sh := [2]uintptr{bh[0], bh[1]}
	return *(*string)(unsafe.Pointer(&sh))
}
