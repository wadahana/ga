package core

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/wadahana/memu/log"
)

type WebSetting struct {
	Port int `yaml:"port"`
}

type IceServer struct {
	Url  string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RtcSetting struct {
	IceServers   []IceServer `yaml:"iceServers"`
	FrameRate    int         `yaml:"frameRate"`
	Bitrate      int         `yaml:"bitrate"`
}
type MEmuSetting struct {
	Path  string `yaml:"path"`
}
type AppSetting struct {
	Web     WebSetting     `yaml:"web"`
	Rtc     RtcSetting     `yaml:"rtc"`
	MEmu    MEmuSetting    `yaml:"memu"`
}

var Setting *AppSetting = nil

func ParseYamlFile(filepath string) error {

	if Setting == nil {
		Setting = new(AppSetting)
	}

	yamlFile, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Errorf("yamlFile.Get err #%v ", err)
		return err;
	}

	log.Infof("yamlFile: \n%s\r\n", string(yamlFile))

	err = yaml.Unmarshal(yamlFile, Setting)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}

	return nil
}
