package system

import (
	"fmt"
	"runtime"
	"syscall"

	"github.com/shirou/gopsutil/mem"
)

func CheckFreeRam() (int64, error) {
	var mem syscall.Sysinfo_t
	err := syscall.Sysinfo(&mem)
	if err != nil {
		return 0, err
	}
	totalRam := int64(mem.Totalram)
	freeRam := int64(mem.Freeram)
	bufferRam := int64(mem.Bufferram)
	usedRam := totalRam - freeRam - bufferRam
	usagePercentage := (float64(usedRam) / float64(totalRam)) * 100
	return int64(usagePercentage), nil
}

func CheckFreeRamSecond() (int64, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	total := int64(mem.TotalAlloc)
	usage := int64(mem.Sys)
	fmt.Println(total, usage)
	percent := (float64(total-usage) / float64(total)) * 100
	return int64(percent), nil
}

func CheckFreeRamThird() (int64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	total := vmStat.Total / 1024 / 1024
	available := vmStat.Available / 1024 / 1024
	percent := (float64(total-available) / float64(total)) * 100
	return int64(percent), nil
}
