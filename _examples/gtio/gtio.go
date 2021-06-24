package main

import (
	"time"

	"github.com/things-labs/gt65/gtio"
)

func main() {
	led, err := gtio.OpenIO("/dev/led")
	if err != nil {
		panic(err)

	}

	for {
		led.On()
		time.Sleep(time.Second * 1)
		led.Off()
		time.Sleep(time.Second * 1)
	}
}
