package msr

import (
	"fmt"
	"os"
)

const defaultFmtStr = "/dev/cpu/%d/msr"

type MSRDev struct {
	dev *os.File
}

func MSR(cpu int, fn func(MSRDev) error) error {
	perCpuDev := fmt.Sprintf(defaultFmtStr, cpu)
	f, err := os.OpenFile(perCpuDev, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck // We don't need to check the error in this context, but CI will complain

	return fn(MSRDev{dev: f})
}

func MSRWithLocation(cpu int, fmtString string, fn func(MSRDev) error) error {
	cpuDir := fmt.Sprintf(fmtString, cpu)
	f, err := os.OpenFile(cpuDir, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck

	return fn(MSRDev{dev: f})
}
