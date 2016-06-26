package alertDecider

import (
	"github.com/stianeikeland/go-rpio"
)

type Switch struct {
    switchGpioPin *rpio.Pin
}

func InitSwitch(switchPin uint8) *Switch {
	rpio.Open()
	
	spin := rpio.Pin(switchPin)
	spin.Input()
	spin.PullUp()
	
	return &Switch{
		switchGpioPin: &spin,
	}
}

func (a *Switch) GetState() (active byte) {
	active = byte(a.switchGpioPin.Read())
	return
}
