package watchmen

import (
    "gopkg.in/yaml.v2"
    "fmt"
    "io/ioutil"
)

type Config struct {
    Switch      SwitchConfig
    Bluetooth   BluetoothConfig
    Horn        HornConfig
    Camera      CameraConfig
    Uploader    UploaderConfig
    Sensors     SensorsConfig
    Notifier    NotifierConfig
}

type SwitchConfig struct {
    Pin uint8
}

type BluetoothConfig struct {
    Devices []string
}

type HornConfig struct {
    Pin int
    Duration int64
}

type CameraConfig struct {
    Device string
    ImagesDir string `yaml:"images_dir"`
    Quality uint8
    Resolution string
}

type SensorsConfig struct {
    Pins []uint8
}

type UploaderConfig struct {
    Ftp FtpConfig
}

type FtpConfig struct {
    IP string
    Port uint16
    User string
    Password string
}

type NotifierConfig struct {
    Mail MailConfig
}

type MailConfig struct {
    User string
    Password string
    Host string
    Port int
    From string
    Recipients []string
}

func GetConfig(fileName string) *Config {
    config := new(Config)

    data := readFile(fileName)
    err := yaml.Unmarshal(*data, config)
    if err != nil {
            panic(fmt.Sprintf("error: %v", err))
    }

    return config
}

func readFile(fileName string) *[]byte {
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
    	panic(fmt.Sprintf("error: %v", err))
    }
    return &data
}