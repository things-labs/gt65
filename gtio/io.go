// Package gtio gt65xx网关
package gtio

// SYS led 	  : /dev/led
// SYS2 led   : /dev/led2
// USB power  : /dev/usbpwr
// buzzer 	  : /dev/buzzer

import (
	"os"
)

// IO µÖ«ÚÇÜio
type IO struct {
	f *os.File
}

// OpenIO open io
func OpenIO(name string) (*IO, error) {
	f, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return &IO{f}, nil
}

// On io on
func (sf IO) On() error {
	_, err := sf.f.Write([]byte{'1'})
	return err
}

// Off io off│
func (sf IO) Off() error {
	_, err := sf.f.Write([]byte{'0'})
	return err
}

// Close io close
func (sf IO) Close() error {
	return sf.f.Close()
}
