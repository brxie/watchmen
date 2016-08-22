package alertDecider

import (
    "os/exec"
    "log"
    //"os"
)

type Bluetooth struct {
    DevicesMAC []string
    AnyDevAlive bool
    lastState bool
}

func InitBluetooth(devices []string) *Bluetooth {
    return &Bluetooth {
        DevicesMAC: devices,
        AnyDevAlive: true,
        lastState: true,
    }
}

func (b *Bluetooth) Start() {
    go func() {
        deadCnt := 0
        devIdx := 0
        for {
            if devIdx >= len(b.DevicesMAC) {
                devIdx = 0
            }
            /*
             l2ping parameters
              -s - size
              -c - count
            */
            cmd := exec.Command("l2ping", "-s 1", "-c 1", b.DevicesMAC[devIdx])
            //cmd.Stderr = os.Stderr
            exitErr := cmd.Run()

            if exitErr == nil {
                deadCnt = 0
                b.updateState(true)
                continue
            } else {
                devIdx++
                deadCnt++
            }

            // clear AnyDevAlive flag if no device alive
            if deadCnt >= len(b.DevicesMAC) {
                b.updateState(false)
                deadCnt = 0
            }
        }
    }()
}

func (b *Bluetooth) updateState(state bool) {
    b.AnyDevAlive = state
    if b.lastState != b.AnyDevAlive {
        log.Printf("[bluetooth] any device alive: %v\n", b.AnyDevAlive)
    }
    b.lastState = state
}