package main

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func getMemInfo() (uint64, uint64, float64) {

	v, _ := mem.VirtualMemory()

	totalMemory := v.Total / 1000000      //in MB
	freeMemory := v.Free / 1000000        // in MB
	usedMemoryPercentage := v.UsedPercent // in %

	return totalMemory, freeMemory, usedMemoryPercentage
}

func getHostInfo() (string, string, string) {

	architecture, _ := host.KernelArch()
	hostInfo, _ := host.Info()
	hostName := hostInfo.Hostname
	operationSystem := hostInfo.OS
	return hostName, architecture, operationSystem
}

func getCpuInfo() (int, float64) {

	cpuNumCores, _ := cpu.Counts(true)

	cpuPercentage, _ := cpu.Percent(time.Second, false)

	return cpuNumCores, cpuPercentage[0]

}
