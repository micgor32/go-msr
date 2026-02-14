// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// !race
package msr_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

var (
	msrs = []string{"SMBASE", "MTRR_CAP", "SMRR_PHYS_BASE", "SMRR_PHYS_MASK", "FEATURE_CONTROL", "PLATFORM_ID", "IA32_DEBUG_INTERFACE", "TSC", "FSB_FREQ", "PLATFORM_INFO", "IA32_EFER"}
)

// Tests both sample implementations of MSR R/W.
// TODO: add note about extending
func TestRdmsrInQemu(t *testing.T) {
	// Out of principle again, building this
	// library on other arch's makes no sense
	// anyways.
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	var bzImage string

	if bzImage = os.Getenv("VMTEST_KERNEL"); len(bzImage) == 0 {
		t.Skipf("VMTEST_KERNEL not set!!")
	}

	if _, err := os.Stat(bzImage); err != nil {
		t.Skipf("Linux kernel image is not found: %s\n", bzImage)
	}

	vm := scriptvm.Start(t, "msr", "",
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/micgor32/go-msr/tests/rdmsr",
			),
			uimage.WithUinitCommand("/bbin/rdmsr"),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(5*time.Minute),
			//qemu.ArbitraryArgs("-machine", "q35"),
			qemu.ArbitraryArgs("-m", "4096"),
			qemu.ArbitraryArgs("-cpu", "Skylake-Client"),
			// arb. nr of cpus, actually we would need
			// only 2 to test what we whant
			qemu.ArbitraryArgs("-smp", "4"),
			qemu.WithKernel(bzImage),
			qemu.WithAppendKernel("console=ttyS0,115200", "earlyprintk=serial,ttyS0,115200", "loglevel=8"),
			//qemu.LogSerialByLine(qemu.DefaultPrint("vm", t.Logf)),
		),
	)

	for _, msr := range msrs {
		if _, err := vm.Console.Expect(expect.String(fmt.Sprintf("%s consistent accross all cpus", msr))); err != nil {
			t.Errorf("VM output did not match expectations: %v", err)
		}
	}

	if err := vm.Kill(); err != nil {
		t.Errorf("Wait for VM process to be killed: %v", err)
	}

	vm.Wait()
}
