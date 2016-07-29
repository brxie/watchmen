package watchmen
import (
    "time"
    "watchmen/alertDecider"
    "watchmen/horn"
    "watchmen/camera"
    "watchmen/uploader"
    "watchmen/notifier"
    "log"
)

type watcher struct {
    horn *horn.Horn
    camera *camera.Camera
    uploader *uploader.Uploader
    alertDecider *alertDecider.AlertDecider
    notifier *notifier.Notifier
}

func Run() {
    log.Println("Watchmen start")
    config := GetConfig("config.yml")

    // init alert decider
    switch_ := alertDecider.InitSwitch(config.Switch.Pin)
    sensors  := alertDecider.InitSensorsGpio(config.Sensors.Pins)

    bluetooth := alertDecider.InitBluetooth(config.Bluetooth.Devices)
    bluetooth.Start()

    decider := &alertDecider.AlertDecider {
        Sensors: sensors,
        Switch: switch_,
        Bluetooth: bluetooth,
    }

    // init horn
    horn := horn.NewHorn(config.Horn.Pin, config.Horn.Duration)

    // init camera
    cam := camera.NewCamera(config.Camera.Device,
                            config.Camera.ImagesDir,
                            config.Camera.Quality,
                            config.Camera.Resolution)

    // init uploader
    ftpCfg := config.Uploader.Ftp
    upld := uploader.NewUploader(ftpCfg.IP,
                                 ftpCfg.Port,
                                 ftpCfg.User,
                                 ftpCfg.Password,
                                 config.Camera.ImagesDir)

    // init notifier
    mailCfg := config.Notifier.Mail
    mail := notifier.NewMail(mailCfg.User, mailCfg.Password, mailCfg.Host,
                             mailCfg.Port, mailCfg.From, mailCfg.Recipients)
    notifier := &notifier.Notifier {
        Mail: mail,
    }
                                 
    w := &watcher {
        horn: horn,
        camera: cam,
        uploader: upld,
        alertDecider: decider,
        notifier: notifier,
    }

    startWatch(w)
}

func startWatch(w *watcher) {
    // scan directory with captured images and upload thought ftp
    w.uploader.PeriodicalScanAndSend()
    for {
        if w.alertDecider.ShouldBeLaunched() {
            w.horn.StartAsync()
            w.camera.CaptureAsync()
            w.notifier.Mail.SendAsync(new(string))
        }

        // force deactivate alarm if it is requested
        if w.horn.IsHornRunning() {
            if w.alertDecider.ShouldBeStopped() {
                w.horn.ForceStop()
            }
        }

        time.Sleep(time.Second * 3)
    }
}