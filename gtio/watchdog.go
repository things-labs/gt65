package gtio

// watchdog: /dev/watchdog
import (
	"os"
)

// defined timeout
const (
	WdtTimeout1s  = 1
	WdtTimeout3s  = 3
	WdtTimeout14s = 14 // default timeout
)

// Watchdog watch dog
type Watchdog struct {
	f *os.File
}

// OpenWatchDog open dog
func OpenWatchDog(name string) (*Watchdog, error) {
	f, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return &Watchdog{f}, nil
}

// Feed feed dog
func (sf Watchdog) Feed() error {
	_, err := sf.f.Write([]byte{'1'})
	return err
}

// SetTimeout set timeout  WDIOC_SETTIMEOUT 1,3,14þºÆ
func (sf Watchdog) SetTimeout(_ int) {
	panic("unknown WDIOC_SETTIMEOUT")
	// syscall.Ioctl(sf.f,)
}

// Close close dog
func (sf Watchdog) Close() error {
	_, err := sf.f.Write([]byte{'V'})
	if err != nil {
		return err
	}
	return sf.f.Close()
}
