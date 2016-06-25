package uploader

import (
    "github.com/dutchcoders/goftp"
    "strconv"
    "log"
    "os"
    "path"
    "path/filepath"
    "time"
)

type Uploader struct {
    host string
    port uint16
    user string
    passwd string
    destPath string
    client *goftp.FTP
    ScanAndSend *ScanAndSend
}

type ScanAndSend struct {
    BackoffTime time.Duration
    ScanPath string
}

func NewUploader(host string, port uint16, user string, passwd string) *Uploader {
    upld := &Uploader {
        host: host,
        port: port,
        user: user,
        passwd: passwd,
        destPath: "/watchmen",
        client: nil,
        ScanAndSend: &ScanAndSend {
            BackoffTime: 1,
            ScanPath: "/var/watchmen/DCIM",
        },
    }
    upld.connect()
    return upld
}


func (u *Uploader) connect() { 
    address := u.host + ":" + strconv.Itoa(int(u.port))
    for {
        client, err := goftp.Connect(address)
        if err != nil {
            log.Println(err)
            continue
        }
        err = client.Login(u.user, u.passwd)
        if err != nil {
            log.Println(err)
            continue
        }
        u.client = client
        break
    }
}


func (u *Uploader) PeriodicalScanAndSend() {
    go func() {
        for {
            //log.Printf("Scan and sending %v\n", u.ScanAndSend.ScanPath)
            filepath.Walk(u.ScanAndSend.ScanPath, func(path string, f os.FileInfo, err error) error {
                if f.IsDir() == false {
                    fullPath := u.ScanAndSend.ScanPath + "/" + f.Name()
                    log.Printf("Uploading: %v\n", fullPath)
                    err := u.send(&fullPath)
                    if err != nil {
                        log.Printf("%v\n", err)
                        // reinitialize connection when uploading failed
                        u.connect()
                    } else {
                        removeFile(&fullPath)
                    }
                }
                return nil
            })         
            time.Sleep(time.Second * u.ScanAndSend.BackoffTime)
        }
    }()
    
}

func removeFile(fullPath *string) {
    err := os.Remove(*fullPath)
    log.Println("Removing:", *fullPath)
    if err != nil {
        log.Printf("%v\n", err)
    }
}

func (u *Uploader) send(fileName *string) error {
    file, err := os.Open(*fileName)
    if err != nil {
        return err
    }
    
    if err := u.client.Cwd(u.destPath); err != nil {
        return err
    }

    if err := u.client.Stor(path.Base(*fileName), file); err != nil {
        return err
    }
    
    return nil
}