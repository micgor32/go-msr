// This pkg is intended to be run on QEMU by running msr_test.go
// from the root dir of go-msr. It can be, however, used

package main

import (
	"fmt"
	"github.com/micgor32/go-msr"
	"runtime"
)

const (
	MsrSMBase             int64 = 0x9e
	MsrMTRRCap            int64 = 0xfe
	MsrSMRRPhysBase       int64 = 0x1F2
	MsrSMRRPhysMask       int64 = 0x1F3
	MsrFeatureControl     int64 = 0x3A
	MsrPlatformID         int64 = 0x17
	MsrIA32DebugInterface int64 = 0xC80
	TimeStampCounter      int64 = 0x10
	MsrFsbFreq            int64 = 0x000000cd
	MsrPlatformInfo       int64 = 0x000000ce
	Ia32Efer              int64 = 0xC0000080
)

type msrs struct {
	name string
	msr  int64
}

// Yes, both are doing reads per cpu, so the name of the function might be
// bit misleading. The idea is, that we first just use MSR and feed it with
// closure that is doing essentially the same logic msr.ReadMSR().
func sessionPerCpu(cpu int, testName string, msrAddr int64) (uint64, error) {
	var msrData uint64
	err := msr.MSR(cpu, func(dev msr.MSRDev) error {
		var readErr error
		msrData, readErr = dev.Read(msrAddr)
		return readErr
	})
	if err != nil {
		return 0xff, fmt.Errorf("msr.Read aborted with: %v", err)
	}

	fmt.Printf("%s, core %d, 0x%x\n", testName, cpu, msrData)
	return msrData, nil
}

func singleRead(cpu int, testName string, msrAddr int64) (uint64, error) {
	res, err := msr.ReadMSR(cpu, msrAddr)
	if err != nil {
		return 0xff, fmt.Errorf("msrReadMSR aborted with: %v", err)
	}

	fmt.Printf("%s, core %d, 0x%x\n", testName, cpu, res)
	return res, nil
}

// todo: also check for 0xff, if thats the case also fail
func eq(e []uint64) bool {
	if len(e) <= 1 {
		return true
	}

	f := e[0]
	for _, v := range e[1:] {
		if v != f {
			return false
		}
	}
	return true
}

func main() {
	tests := []msrs{
		{"SMBASE", MsrSMBase},
		{"MTRR_CAP", MsrMTRRCap},
		{"SMRR_PHYS_BASE", MsrSMRRPhysBase},
		{"SMRR_PHYS_MASK", MsrSMRRPhysMask},
		{"FEATURE_CONTROL", MsrFeatureControl},
		{"PLATFORM_ID", MsrPlatformID},
		{"IA32_DEBUG_INTERFACE", MsrIA32DebugInterface},
		{"TSC", TimeStampCounter},
		{"FSB_FREQ", MsrFsbFreq},
		{"PLATFORM_INFO", MsrPlatformInfo},
		{"IA32_EFER", Ia32Efer},
	}

	noCpus := runtime.NumCPU()

	// nested loop i know, but its still O(n * m), where
	// n = 11 (range tests, so constant basically), m = noCpus
	// so still linear :D
	for _, test := range tests {
		var sessionRet []uint64
		var singleRet []uint64

		for i := 0; i <= (noCpus - 1); i++ {
			retSes, err := sessionPerCpu(i, test.name, test.msr)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			sessionRet = append(sessionRet, retSes)

			retSin, err := singleRead(i, test.name, test.msr)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			singleRet = append(singleRet, retSin)
		}

		if eq(sessionRet) {
			fmt.Printf("session %s consistent accross all cpus\n", test.name)
		} else {
			fmt.Printf("session %s not consistent accross all cpus\n", test.name)
		}

		if eq(singleRet) {
			fmt.Printf("single %s consistent accross all cpus\n", test.name)
		} else {
			fmt.Printf("single %s not consistent accross all cpus\n", test.name)
		}
	}
}
