package main

import (
	"fmt"
	"github.com/micgor32/go-msr"
)

const (
	MsrSMBase             int64 = 0x9e
	MsrMTRRCap            int64 = 0xfe
	MsrSMRRPhysBase       int64 = 0x1F2
	MsrSMRRPhysMask       int64 = 0x1F3
	MsrFeatureControl     int64 = 0x3A
	MsrPlatformID         int64 = 0x17
	MsrIA32DebugInterface int64 = 0xC80
)

func main() {
	t, _ := msr.ReadMSR(0, MsrSMBase)

	fmt.Printf("0x%x\n", t)
}
