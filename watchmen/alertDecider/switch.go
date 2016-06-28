package alertDecider

import (
    "github.com/stianeikeland/go-rpio"
    "log"
    "time"
)

type Switch struct {
    switchGpioPin *rpio.Pin
    State bool
    lastState bool
}

func InitSwitch(switchPin uint8) *Switch {
    rpio.Open()

    spin := rpio.Pin(switchPin)
    spin.Input()
    spin.PullUp()

    s := &Switch{
        switchGpioPin: &spin,
        State: false,
    }
    start(s)
    return s
}

func start(s *Switch) {
    go func() {
        for {
            s.updateState(s.getState())
            time.Sleep(time.Second / 2)
        }
    }()
}

func (s *Switch) getState() bool {
    if state := s.switchGpioPin.Read(); state == 0 {
        return false
    }
    return true
}

func (s *Switch) updateState(state bool) {
    s.State = s.getState()
    if s.lastState != s.State {
        log.Printf("[switch] changed state to: %v\n", s.State)
    }
    s.lastState = s.State
}