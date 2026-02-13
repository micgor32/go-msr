// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// !race
package msr_test

import (
	"os"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

// Tests both sample implementations of MSR R/W.
// TODO: add note about extending
func TestRDMSRQEMU(t *testing.T) {
	// Out of principle again, building this
	// library on other arch's makes no sense
	// anyways.
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	var ovmf string
	var bzImage string

	if ovmf = os.Getenv("VMTEST_OVMF"); len(ovmf) == 0 {
		t.Skipf("VMTEST_OVMF not set!!")
	}

	if _, err := os.Stat(ovmf); err != nil {
		t.Skipf("OVMF.fd image is not found: %s\n", ovmf)
	}

	if bzImage = os.Getenv("VMTEST_KERNEL"); len(bzImage) == 0 {
		t.Skipf("VMTEST_KERNEL not set!!")
	}

	if _, err := os.Stat(bzImage); err != nil {
		t.Skipf("Linux kernel image is not found: %s\n", bzImage)
	}

	vm := scriptvm.Start(t, "rdmsr", "",
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/micgor32/go-msr/examples",
			),
			uimage.WithUinitCommand("/bbin/examples"),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(5*time.Minute),
			qemu.ArbitraryArgs("-machine", "q35"),
			qemu.ArbitraryArgs("-m", "4096"),
			// arb. nr of cpus, actually we would need
			// only 2 to test what we whant
			qemu.ArbitraryArgs("-smp", "4"),
			qemu.ArbitraryArgs("-bios", ovmf),
			qemu.ArbitraryArgs("-kernel", bzImage),
		),
	)

	if _, err := vm.Console.Expect(expect.String("placeholder for now (i.e. will fail :D)")); err != nil {
		t.Errorf("VM output did not match expectations: %v", err)
	}

	if err := vm.Kill(); err != nil {
		t.Errorf("Wait for VM process to be killed: %v", err)
	}

	vm.Wait()
}
