package gtio

// Key : /dev/key
import (
	"os"
	"syscall"
)

// Key key
type Key struct {
	f *os.File
}

// OpenKey open key
func OpenKey(name string) (*Key, error) {
	f, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return &Key{f}, nil
}

// Read read key  ioctl(fd,0,0)
func (sf Key) Read() (int, error) {
	result, _, err := syscall.Syscall(syscall.SYS_IOCTL, sf.f.Fd(), 0, 0)
	return int(result), err
}

// Close close
func (sf Key) Close() error {
	return sf.f.Close()
}
