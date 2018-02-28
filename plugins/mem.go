package plugins

import (
	"runtime"
	"strconv"
)

var mems runtime.MemStats

func Mem() string {
	runtime.ReadMemStats(&mems)

	var mem string
	if mems.Alloc > 1000000 {
		mem = strconv.Itoa(int(mems.Alloc/1000000)) + "MB"
	} else {
		mem = strconv.Itoa(int(mems.Alloc/1000)) + "KB"
	}
	return mem
}
