package alertDecider

import (
    "github.com/stianeikeland/go-rpio"
)

type Sensors struct {
    gpioPins []uint8
}

func InitSensorsGpio(pins []uint8) *Sensors {
    rpio.Open()
    for _, gpio := range pins {
        pin := rpio.Pin(gpio)
        pin.Input()
        pin.PullDown()
    }

    return &Sensors {
        gpioPins: pins,
    }
}

func (s *Sensors) GetRaisedSensors() []uint8 {
    var rsens []uint8
    for _, gpio := range s.gpioPins {
        pin := rpio.Pin(gpio)
        if state := pin.Read(); uint8(state) == 1 {
            rsens = append(rsens, gpio)
        }
    }
    return rsens
}

func (s *Sensors) AnySensorRaised() (r bool) {
    if rsens := s.GetRaisedSensors(); len(rsens) > 0 {
        r = true
    }
    return
}
