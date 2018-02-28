package plugins

import (
	"os"
	"strconv"

	"github.com/struCoder/pidusage"
)

func Stats() string {
	sysInfo, _ := pidusage.GetStat(os.Getpid())
	var mem string
	if sysInfo.Memory > 1000000 {
		mem = strconv.Itoa(int(sysInfo.Memory/1000000)) + "MB"
	} else {
		mem = strconv.Itoa(int(sysInfo.Memory/1000)) + "KB"
	}
	cpu := strconv.Itoa(int(sysInfo.CPU)) + "%"
	s := mem + "\n" + cpu
	return s
}
