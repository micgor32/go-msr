package msr

import (
	"fmt"
	"os"
)

type MSRDev struct {
	dev *os.File
}

func MSR(cpu int, fn func(MSRDev) error) error {
	perCpuDev := fmt.Sprintf(defaultFmtStr, cpu)
	f, err := os.OpenFile(perCpuDev, os.O_RDWR, 777)
	if err != nil {
		return err
	}
	defer f.Close()

	return fn(MSRDev{dev: f})
}

func MSRWithLocation(cpu int, fmtString string, fn func(MSRDev) error) error {
	cpuDir := fmt.Sprintf(fmtString, cpu)
	f, err := os.OpenFile(cpuDir, os.O_RDWR, 777)
	if err != nil {
		return err
	}
	defer f.Close()

	return fn(MSRDev{dev: f})
}
