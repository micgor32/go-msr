package msr

import (
	"encoding/binary"
	"fmt"
)

func (d MSRDev) Write(reg int64, val uint64) error {
	regBuf := make([]byte, 8)

	binary.NativeEndian.PutUint64(regBuf, val)

	count, err := d.dev.WriteAt(regBuf, reg)
	if err != nil {
		return err
	}

	if count != 8 {
		return fmt.Errorf("Write count not a uint64: %d", count)
	}

	return nil
}

func WriteMSR(cpu int, msr int64, val uint64) error {
	err := MSR(cpu, func(dev MSRDev) error {
		err := dev.Write(msr, val)
		return err
	})

	return err
}
