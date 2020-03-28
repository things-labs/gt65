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

func (sf IO) Write(status bool) error {
	val := byte('0')
	if status {
		val = '1'
	}
	_, err := sf.f.Write([]byte{val})
	return err
}

// On io on
func (sf IO) On() error {
	return sf.Write(true)
}

// Off io off│
func (sf IO) Off() error {
	return sf.Write(false)
}

// Close io close
func (sf IO) Close() error {
	return sf.f.Close()
}
