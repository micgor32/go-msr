package msr

import (
	"encoding/binary"
	"fmt"
)

func (d MSRDev) Read(msr int64) (uint64, error) {
	regBuf := make([]byte, 8)

	rc, err := d.dev.ReadAt(regBuf, msr)

	if err != nil {
		return 0, err
	}

	if rc != 8 {
		return 0, fmt.Errorf("Read wrong count of bytes: %d", rc)
	}

	// x86 processors are little-endian, put out of principle of not
	// assuming things, let's use native-endian (since it is exposed on
	// std library since Go 1.21)
	msrOut := binary.NativeEndian.Uint64(regBuf)

	return msrOut, nil
}

func ReadMSR(cpu int, msr int64) (uint64, error) {
	var msrD uint64

	err := MSR(cpu, func(dev MSRDev) error {
		var readErr error
		msrD, readErr = dev.Read(msr)
		return readErr
	})

	if err != nil {
		return 0, err
	}

	return msrD, nil
}
