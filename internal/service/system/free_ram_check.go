package system

import (
	"syscall"
)

func CheckFreeRam() (int64, error) {
	var mem syscall.Sysinfo_t
	err := syscall.Sysinfo(&mem)
	if err != nil {
		return 0, err
	}
	totalRam := int64(mem.Totalram)
	freeRam := int64(mem.Freeram)
	usedRam := totalRam - freeRam
	usagePercentage := (float64(usedRam) / float64(totalRam)) * 100
	return int64(usagePercentage), nil
}
