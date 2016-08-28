package LCD

import (
    "time"
)

type Icon struct {
    file string
    position int
}

var (
    StateIcon = Icon {
        file: "img/play-icon.png",
        position: 4,
    }
    MotionIcon = Icon {
        file: "img/motion-icon.png",
        position: 36,
    }
    BluetoothIcon = Icon {
        file: "img/bluetooth-icon.png",
        position: 68,
    }
    HornIcon = Icon {
        file: "img/horn-icon.png",
        position: 97,
    }
)

func (w *LcdWrapper) DrawIcon(icon *Icon) {
    w.WriteImage(icon.file, icon.position)
    w.Dev.Display()
}

func (w *LcdWrapper) ClearIcon(icon *Icon) {
    w.ClearImageArea(icon.file, icon.position)
    w.Dev.Display()
}

func (w *LcdWrapper) BlinkIcon(icon *Icon, ntimes, delay int) {
    for i := 0; i < ntimes; i++ {
        if i % 2 == 0 {
            w.DrawIcon(icon)
            time.Sleep(time.Second * time.Duration(delay))
        } else{
            w.ClearIcon(icon)
        }
        w.Dev.Display()
    }
}

func (w *LcdWrapper) DisplayString(text string, position int) {
    w.WriteText(text, position)
    w.Dev.Display()
}