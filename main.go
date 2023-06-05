package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/reconquest/pkg/log"
	"github.com/ssimunic/gosensors"
)

const (
	sleepTime = 1
)

func main() {
	var cpuFrequency string
	for {
		frequency, err := getCPUFrequency()
		if err != nil {
			log.Fatal(err)
		} else {
			cpuFrequency = "i9 F: " + frequency + " Mhz"
		}

		cpuTemperature, err := getCpuTemperature()
		if err != nil {
			log.Fatal(err)
		}

		gpuTemperature, err := getGpuTemperature()
		if err != nil {
			log.Fatal(err)
		}

		gpuFrequency, err := getGpuFrequency()
		if err != nil {
			log.Fatal(err)
		}

		log.Info("--------------------------------------"+
			"\n"+cpuTemperature+"    | "+cpuFrequency+"\n\n",
			gpuTemperature+" | "+gpuFrequency,
		)

		time.Sleep(sleepTime * time.Second)
	}
}

func getCPUFrequency() (string, error) {
	contents, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	var total []float64
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if strings.Contains(line, "cpu MHz") {
			frequency := strings.Split(line, ": ")
			intFrequency, err := strconv.ParseFloat(frequency[1], 64)
			if err != nil {
				return "", err
			}

			total = append(total, intFrequency)
		}
	}

	var totalInt []int

	for _, item := range total {
		convertedItem := int(item)
		totalInt = append(totalInt, convertedItem)
	}

	sort.Slice(totalInt, func(i, j int) bool {
		return totalInt[i] > totalInt[j]
	})

	return strconv.Itoa(totalInt[0]), nil
}

func getCpuTemperature() (string, error) {
	sensors, err := gosensors.NewFromSystem()
	var cpuTemp string
	if err != nil {
		return "", err
	} else {
		cpuTemp = "i9 T: " + string(
			strings.Split(
				sensors.Chips["coretemp-isa-0000"]["Core 0"], " ",
			)[0],
		)
	}

	return cpuTemp, nil
}

func getGpuTemperature() (string, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("RTX-2080 T: %s'C", strings.TrimSpace(string(output))), nil
}

func getGpuFrequency() (string, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=clocks.gr", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace("RTX-2080 F: " + string(output)), nil
}
