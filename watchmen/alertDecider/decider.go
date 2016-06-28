package alertDecider

import (
    "log"
)

type AlertDecider struct {
    Sensors *Sensors
    Switch *Switch
    Bluetooth *Bluetooth
}

func (a *AlertDecider) ShouldBeLaunched() bool {
    // check if alarm is turned on
    if a.Switch.State == false {
        return false
    }

    // check sensors
    raisedSensors := a.Sensors.GetRaisedSensors()
    if len(raisedSensors) == 0 {
        return false
    }
    log.Printf("[decider] sensor(s) active: %v\n", raisedSensors)

    // check alive bluetooth devices
    if a.Bluetooth.anyDevAlive == true {
        return false
    }
    log.Println("[decider] no bluetooth device active")

    return true
}

func (a *AlertDecider) ShouldBeStopped() bool {

    // 0 means alarm is turned off
    if a.Switch.State == false {
        return true
    }

    if a.Bluetooth.anyDevAlive == true {
        return true
    }

    return false
}
