package watchmen
import (
	"time"
	"watchmen/alertDecider"
	"watchmen/horn"
	"watchmen/camera"
	"watchmen/uploader"
	"log"
)

type watcher struct {
    horn *horn.Horn
    camera *camera.Camera
    uploader *uploader.Uploader
    alertDecider *alertDecider.AlertDecider
}


func Run() {
	log.Println("Watchmen start") 
	
    // init alert decider
    switch_ := alertDecider.InitSwitch(23)
    sensors  := alertDecider.InitSensorsGpio([]uint8{24})
    btDevices := []string{"D8:9E:3F:DD:7A:DA", "CC:07:E4:2C:E9:C5"}
    bluetooth := alertDecider.InitBluetooth(btDevices)
    bluetooth.Start()

    decider := &alertDecider.AlertDecider {
        Sensors: sensors,
        Switch: switch_,
        Bluetooth: bluetooth,
    }
    
    // init horn
	horn := horn.NewHorn(25)
    
    // init camera
    cam := camera.NewCamera("/dev/video0", "/var/watchmen/DCIM", 75)
    
    // init uploader
	upld := uploader.NewUploader("5.1.1.1", 21, "*****", "*****")
    upld.ScanAndSend.ScanPath = "/var/watchmen/DCIM"
    
    w := &watcher {
        horn: horn,
        camera: cam,
        uploader: upld,
        alertDecider: decider,
    }
    
    startWatch(w)
}

func startWatch(w *watcher) {
    // scan directory with captured images and upload thought ftp
    w.uploader.PeriodicalScanAndSend()
    for {
        if w.alertDecider.ShouldBeLaunched() {
            w.horn.StartAsync(15)
            w.camera.CaptureAsync()
        } 
        
        // force deactivate alarm if it is requested
        if w.horn.IsHornRunning() {
            if w.alertDecider.ShouldBeStopped() {
                w.horn.ForceStop()
            }
        }

        
        time.Sleep(time.Second)
    }
}