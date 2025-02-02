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

	for {
		cpuInfo, _ := cpu.Info()
		cpuPercent, _ := cpu.Percent(time.Second, false)
		memStats, _ := mem.VirtualMemory()
		diskStats, _ := disk.Usage("/")
		hostInfo, _ := host.Info()

		gpuName, gpuUsage := getGPUInfo()

		pterm.Println("\033[H\033[2J")

		s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString("Resource Monitor")).Srender()
		pterm.DefaultCenter.Println(s)

		cpuName := "Unknown"
		if len(cpuInfo) > 0 {
			cpuName = cpuInfo[0].ModelName
		}

		pterm.DefaultTable.WithHasHeader(true).WithData([][]string{
			{"Component", "Name", "Usage (%)"},
			{"CPU", cpuName, fmt.Sprintf("%.2f%%", cpuPercent[0])},
			{"RAM", "Memory", fmt.Sprintf("%.2f%%", memStats.UsedPercent)},
			{"Storage", "Disk", fmt.Sprintf("%.2f%%", diskStats.UsedPercent)},
			{"GPU", gpuName, fmt.Sprintf("%.2f%%", gpuUsage)},
		}).Render()

		pterm.DefaultBarChart.WithHorizontal().WithBars([]pterm.Bar{
			{Label: "CPU", Value: int(cpuPercent[0]), Style: pterm.NewStyle(pterm.FgGreen)},
			{Label: "RAM", Value: int(memStats.UsedPercent), Style: pterm.NewStyle(pterm.FgBlue)},
			{Label: "Storage", Value: int(diskStats.UsedPercent), Style: pterm.NewStyle(pterm.FgCyan)},
			{Label: "GPU", Value: int(gpuUsage), Style: pterm.NewStyle(pterm.FgMagenta)},
		}).Render()

		pterm.Info.Println("Operating System:", hostInfo.Platform, hostInfo.PlatformVersion)
		pterm.Info.Println("Press F10 to exit.")

		time.Sleep(2 * time.Second)

		if char, key, err := keyboard.GetKey(); err == nil {
			if key == keyboard.KeyF10 || char == 0 {
				pterm.Warning.Println("Exiting monitoring...")
				break
			}
		}
	}
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
