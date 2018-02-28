package plugins

import (
	"runtime"
	"strconv"
)

func Mem() string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return strconv.Itoa(int(mem.Alloc/1000)) + "KB"
}
