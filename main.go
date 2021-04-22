package main

import (
	"io/ioutil"
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
		cpuTemp = "CPU TEMP: " + string(
			strings.Split(
				sensors.Chips["coretemp-isa-0000"]["Core 0"], " ",
			)[0],
		)
	}

	return cpuTemp, nil
}

func main() {
	var cpuFrequency string
	for {
		frequency, err := getCPUFrequency()
		if err != nil {
			cpuFrequency = "CPU FREQUENCY: error"
			log.Error(cpuFrequency)
			log.Fatal(err)
		} else {
			cpuFrequency = "CPU FREQUENCY: " + frequency + " Mhz"
		}

		cpuTemperature, err := getCpuTemperature()
		if err != nil {
			log.Fatal(err)
		}

		log.Info(cpuTemperature + " | " + cpuFrequency)

		time.Sleep(sleepTime * time.Second)
	}
}
