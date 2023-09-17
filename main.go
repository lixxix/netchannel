package main

import (
	"fmt"
	"os"

	"github.com/lixxix/netchannel/pipe"
	"gopkg.in/yaml.v2"
)

type ChannelConfig struct {
	Listen int    `yaml:"listen"`
	Type   string `yaml:"type"`
	Target string `yaml:"target"`
}

func main() {
	buf, err := os.ReadFile("pipe.yaml")
	if err != nil {
		panic(err)
	}

	ConfigData := &ChannelConfig{}
	err = yaml.Unmarshal(buf, ConfigData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("from:%d to %s 类型:%s\n", ConfigData.Listen, ConfigData.Target, ConfigData.Type)
	if ConfigData.Type == "ws" {
		pipe.StartWS(fmt.Sprintf(":%d", ConfigData.Listen), ConfigData.Target)
	} else if ConfigData.Type == "tcp" {
		pipe.StartTCP(fmt.Sprintf(":%d", ConfigData.Listen), ConfigData.Target)
	} else {
		fmt.Println("type只能选择tcp或者ws")
	}

}
