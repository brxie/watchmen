package horn

import (
    "github.com/stianeikeland/go-rpio"
    "time"
    "log"
)

type Horn struct {
    gpioPin *rpio.Pin
    duration int64
    asycWorker *asyncWorker
    IsRunning bool
}

type asyncWorker struct {
    workers chan byte
    workRequest chan bool
}

func NewHorn(gpio int, duration int64) *Horn {
    rpio.Open()
    pin := rpio.Pin(gpio)
    pin.Output()
    pin.Low()

    return &Horn{
        gpioPin: &pin,
        duration: duration,
        asycWorker: &asyncWorker{
            workers: make(chan byte, 10),
            workRequest: make(chan bool, 10),
        },
        IsRunning: false,
    }
}

func (h *Horn) StartAsync() {
    log.Printf("[horn] activating horn")
    h.setState(true)

    go func() {
        // disable horn on finish
        defer h.gpioPin.Low()
        defer h.setState(false)

        // when trying to create new horn, fist kill the current one if is running
        if len(h.asycWorker.workers) != 0 {
            h.killCurrWorker()
        }

        h.asycWorker.workers <- 1
        startTime := time.Now().Unix()
        h.gpioPin.High()
        for time.Now().Unix() < startTime + h.duration {
            h.asycWorker.workRequest <- true
            if v := <-h.asycWorker.workRequest; v {
                time.Sleep(time.Millisecond * 120)
            } else {
                <-h.asycWorker.workRequest
                <-h.asycWorker.workers
                return
            }
        }
        <-h.asycWorker.workers
    }()
    
}

func (h *Horn) setState (state bool) {
    h.IsRunning = state
}

func (h *Horn) killCurrWorker(){
    h.ForceStop()
    for len(h.asycWorker.workers) != 0 {
        time.Sleep(time.Millisecond * 100)
    }
}

func (h *Horn) ForceStop() {
    h.asycWorker.workRequest <- false
}