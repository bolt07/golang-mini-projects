package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port = ":8080"

type systemInfo struct {
	HostName            string  `json:"hostname"`
	TotalMemory         uint64  `json:"total_memory"`
	FreeMemory          uint64  `json:"free_memory"`
	MemoryUsePercentage float64 `json:"memory_used_percentage"`
	Architecture        string  `json:"architecture"`
	OperatingSytem      string  `json:"os"`
	NumberOfCPUCores    int     `json:"number_of_cpu_cores"`
	CPUUsedPercentage   float64 `json:"cpu_used_percentage"`
}

func main() {
	http.HandleFunc("/", getSysInfo)

	fmt.Printf("Server is running on %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getSysInfo(w http.ResponseWriter, r *http.Request) {

	totalMem, freeMem, usedMemPercent := getMemInfo()
	hostName, arch, os := getHostInfo()
	cpuCores, cpuUsage := getCpuInfo()

	info := systemInfo{
		HostName:            hostName,
		TotalMemory:         totalMem,
		FreeMemory:          freeMem,
		MemoryUsePercentage: usedMemPercent,
		Architecture:        arch,
		OperatingSytem:      os,
		NumberOfCPUCores:    cpuCores,
		CPUUsedPercentage:   cpuUsage,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Failed to encode system info", http.StatusInternalServerError)
		log.Printf("JSON encode error: %v", err)
	}
}
