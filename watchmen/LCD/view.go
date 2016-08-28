package LCD

import (
    "watchmen/alertDecider"
    "watchmen/horn"
    "sync"
    "time"
	"fmt"
)

type lcdView struct {
    lcdWrapper      *LcdWrapper
    stateUnits      *StateUnits
    wg              *sync.WaitGroup
    lastStateTime   *lastStateTime
    histogram       *histogram
}

type StateUnits struct {
    Switch_     *alertDecider.Switch
    Sensor      *alertDecider.Sensors
    Bluetooth   *alertDecider.Bluetooth
    Horn        *horn.Horn
}

type lastStateTime struct {
    switch_     time.Time
    horn        time.Time
}

func NewLCD(iicBussAddr, devAddr byte, height, width int, units *StateUnits) *lcdView {
    lcd := newLcdWrapper(&iicBussAddr, &devAddr, &height, &width)
    var wg sync.WaitGroup
    
    return &lcdView {
        lcdWrapper: lcd,
        stateUnits: units,
        wg:         &wg,
        lastStateTime: &lastStateTime{},
        histogram: NewHistogram(lcd, 4, 2),
    }
}

func (l *lcdView) Display() {
    var switchLastState, bluetoothLastState bool
    go func() {
        l.showStateTimes()
        for i := 0;; i++ {

            // switch
            if l.stateUnits.Switch_.State != switchLastState {
                l.lastStateTime.switch_ = time.Now()
                l.changeIconState(&l.stateUnits.Switch_.State, &StateIcon)
                switchLastState = l.stateUnits.Switch_.State
                l.showStateTimes()
            }

            // bluetooth
            if l.stateUnits.Bluetooth.AnyDevAlive != bluetoothLastState {
                l.changeIconState(&l.stateUnits.Bluetooth.AnyDevAlive, &BluetoothIcon)
                bluetoothLastState = l.stateUnits.Bluetooth.AnyDevAlive
            }

            // motion sensor
            if l.stateUnits.Sensor.AnySensorRaised() {
                l.blinkIcon(&MotionIcon)
                l.histogram.dataCollector.fill(1)
            }

            // horn
            if l.stateUnits.Horn.IsRunning {
                l.lastStateTime.horn = time.Now()
                l.blinkIcon(&HornIcon)
                l.showStateTimes()
            }

            // draw histogram
            if (i % 7200) == 0 {
            //if (i % 40) == 0 {
                l.histogram.plot()
            }

            time.Sleep(time.Second / 10)
        }
    }()

}

func (l *lcdView) blinkIcon(ico *Icon) {
    l.wg.Wait()
    l.wg.Add(1)
    go func() {
        defer l.wg.Done()
        l.lcdWrapper.BlinkIcon(ico, 2, (1 / 2))
    }()
}

func (l *lcdView) changeIconState(state *bool, ico *Icon) {
    l.wg.Wait()
    l.wg.Add(1)
    go func() {
        defer l.wg.Done()
        if *state {
            l.lcdWrapper.DrawIcon(ico)
        } else {
            l.lcdWrapper.ClearIcon(ico)
        }
    }()
}

func (l *lcdView) showStateTimes() {
    l.wg.Wait()
    l.wg.Add(1)

    go func() {
        defer l.wg.Done()
        time := l.formatTime(&l.lastStateTime.switch_) 
        l.lcdWrapper.DisplayString("on/off: " + time , 384)

        time = l.formatTime(&l.lastStateTime.horn) 
        l.lcdWrapper.DisplayString("alarm: " + time, 512)
    }()
}

func (l *lcdView) formatTime(t *time.Time) string {
    if t.IsZero() {
        return "n/a"
    } else {
        return fmt.Sprintf("%02d:%02d %02d.%02d",
            t.Hour(), t.Minute(), t.Month(), t.Day())
    }
}