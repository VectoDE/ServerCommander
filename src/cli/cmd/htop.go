package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/pterm/pterm"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func HtopCommand() {
	if err := keyboard.Open(); err != nil {
		pterm.Error.Println("Error initializing keyboard monitoring:", err)
		return
	}
	defer keyboard.Close()

	cpuInfo, _ := cpu.Info()
	cpuName := "Unknown"
	if len(cpuInfo) > 0 {
		cpuName = cpuInfo[0].ModelName
	}

	hostInfo, _ := host.Info()

	for {
		fmt.Print("\033[H\033[2J")

		s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString("Resource Monitor")).Srender()
		pterm.Println(s)

		cpuPercent, _ := cpu.Percent(500*time.Millisecond, false)
		memStats, _ := mem.VirtualMemory()
		diskStats := getDiskUsage()
		gpuName, gpuUsage := getGPUInfo()
		ramModules := getMemoryModules()

		data := [][]string{
			{"Component", "Name", "Usage"},
			{"CPU", cpuName, fmt.Sprintf("%.2f%%", cpuPercent[0])},
			{"RAM Total", "Memory", fmt.Sprintf("%.2f%% (%.2f GB / %.2f GB)", memStats.UsedPercent, float64(memStats.Used)/1e9, float64(memStats.Total)/1e9)},
		}
		for _, ram := range ramModules {
			data = append(data, []string{"RAM Module", ram.Name, fmt.Sprintf("%.2f GB", ram.Capacity)})
		}
		for _, d := range diskStats {
			data = append(data, []string{"Storage", d.Name, fmt.Sprintf("%.2f%% (%.2f GB / %.2f GB)", d.UsedPercent, d.Used, d.Total)})
		}
		data = append(data, []string{"GPU", gpuName, fmt.Sprintf("%.2f%%", gpuUsage)})

		pterm.DefaultTable.WithHasHeader(true).WithData(data).Render()

		pterm.Info.Println("Operating System:", hostInfo.Platform, hostInfo.PlatformVersion)
		pterm.Info.Println("Press F10 to exit.")

		select {
		case <-time.After(1 * time.Second):
		default:
			if char, key, err := keyboard.GetKey(); err == nil {
				if key == keyboard.KeyF10 || char == 0 {
					pterm.Warning.Println("Exiting monitoring...")
					return
				}
			}
		}
	}
}

type DiskUsage struct {
	Name        string
	Total       float64
	Used        float64
	UsedPercent float64
}

func getDiskUsage() []DiskUsage {
	partitions, _ := disk.Partitions(false)
	var diskUsages []DiskUsage
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err == nil {
			diskUsages = append(diskUsages, DiskUsage{
				Name:        partition.Mountpoint,
				Total:       float64(usage.Total) / 1e9,
				Used:        float64(usage.Used) / 1e9,
				UsedPercent: usage.UsedPercent,
			})
		}
	}
	return diskUsages
}

type RAMModule struct {
	Name     string
	Capacity float64
}

func getMemoryModules() []RAMModule {
	var modules []RAMModule
	if runtime.GOOS == "windows" {
		out, err := exec.Command("wmic", "memorychip", "get", "Capacity,Manufacturer").Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines[1:] {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					capacity, _ := strconv.ParseFloat(fields[0], 64)
					modules = append(modules, RAMModule{Name: fields[1], Capacity: capacity / 1e9})
				}
			}
		}
	}
	return modules
}

func getGPUInfo() (string, float64) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsGPUInfo()
	case "linux":
		return getLinuxGPUInfo()
	case "darwin":
		return getMacOSGPUInfo()
	default:
		return "GPU not recognized", 0.0
	}
}

func getWindowsGPUInfo() (string, float64) {
	out, err := exec.Command("wmic", "path", "win32_videocontroller", "get", "name").Output()
	if err != nil {
		return "Unknown GPU", 0.0
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 1 {
		return strings.TrimSpace(lines[1]), 0.0
	}
	return "Unknown GPU", 0.0
}

func getLinuxGPUInfo() (string, float64) {
	out, err := exec.Command("sh", "-c", "lspci | grep -i 'vga\\|3d\\|2d'").Output()
	if err != nil {
		return "Unknown GPU", 0.0
	}
	gpuName := strings.TrimSpace(string(out))

	nvidiaOut, err := exec.Command("sh", "-c", "nvidia-smi --query-gpu=utilization.gpu --format=csv,noheader,nounits").Output()
	if err == nil {
		usage, _ := strconv.ParseFloat(strings.TrimSpace(string(nvidiaOut)), 64)
		return gpuName, usage
	}

	return gpuName, 0.0
}

func getMacOSGPUInfo() (string, float64) {
	out, err := exec.Command("sh", "-c", "system_profiler SPDisplaysDataType | grep 'Chipset Model'").Output()
	if err != nil {
		return "Unknown GPU", 0.0
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		gpuName := strings.TrimSpace(strings.Split(lines[0], ":")[1])
		return gpuName, 0.0
	}
	return "Unknown GPU", 0.0
}
