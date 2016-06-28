package camera

import (
    "os/exec"
    "os"
    "time"
    "log"
    "strconv"
    "fmt"
)

type Camera struct {
    Device string
    ImagesPath string
    Quality uint8
    AsyncWorker *asyncWorker
}

type asyncWorker struct {
    workRequest chan bool
}

func NewCamera(device, imagesPath string, quality uint8) (cam *Camera) {
    // create images dir if not exist
    mkdir(&imagesPath)

    cam = new(Camera)
    cam.Device = device
    cam.ImagesPath = imagesPath
    cam.Quality = quality
    cam.AsyncWorker = &asyncWorker {
        workRequest: make(chan bool, 1),
    }
    return
}

func (c *Camera) Capture() string {
    outBaseName := getUnixTime() + "_000.jpeg"
    log.Printf("[camera] Capturing images in format %v", outBaseName)
    /*
     streamer parameters
        -t times    number of frames or hh:mm:ss
        -r fps      frame rate
        -q          quiet operation
        -j quality  quality for mjpeg or jpeg
        -s size     specify size
    */
    cmd := exec.Command("streamer", "-t 00:00:06", "-r 4", "-o" + outBaseName, "-q",
                        "-j " + strconv.Itoa(int(c.Quality)), "-s " + "800x600")
    cmd.Dir = c.ImagesPath
    cmd.Stderr = os.Stderr
    cmd.Run()
    return outBaseName
}

func (c *Camera) CaptureAsync() {
    go func() {
        c.AsyncWorker.workRequest <- true
        c.Capture()
        <-c.AsyncWorker.workRequest
    }()
}

func getUnixTime() (utime string) {
    utime = strconv.Itoa(int(time.Now().Unix()))
    return
}

func mkdir(path *string) {
    err := os.MkdirAll(*path, 0655)
    if err != nil {
        panic(fmt.Sprintf("%v", err))
    }
}