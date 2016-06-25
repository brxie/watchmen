package horn

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

type Horn struct {
    gpioPin *rpio.Pin
	asycWorker *asyncWorker
}

type asyncWorker struct {
	workers chan byte
 	workRequest chan bool
}

func NewHorn(gpio int) *Horn {
	rpio.Open()
	pin := rpio.Pin(gpio)
	pin.Output()
	pin.Low()

	return &Horn{
		gpioPin: &pin,
		asycWorker: &asyncWorker{
			workers: make(chan byte, 10),
			workRequest: make(chan bool, 10),
		},
	}
}

func (h *Horn) StartAsync(duration int64) {
	go func() {
		defer h.gpioPin.Low()

		// when trying to create new horn, fist kill the current one if is running
		if len(h.asycWorker.workers) != 0 {
			h.killCurrWorker()
		}
		
		h.asycWorker.workers <- 1
		startTime := time.Now().Unix() 
		for time.Now().Unix() < startTime + duration {
			h.asycWorker.workRequest <- true	
			if v := <-h.asycWorker.workRequest; v {
				h.gpioPin.Toggle()
				time.Sleep(time.Millisecond * 120)
			} else {
				<-h.asycWorker.workRequest
				<-h.asycWorker.workers
				return
			}
		}
		<-h.asycWorker.workRequest
		<-h.asycWorker.workers
	}()
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

func (h *Horn) IsHornRunning() bool {
	if len(h.asycWorker.workers) != 0 {
		return true
	}
	return false
}