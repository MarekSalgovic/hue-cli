package cli

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
)

const configPath = "/hue/config"
const configFile = "/hue/config/config.data"

type config struct {
	AppID      string `json:"appID"`
	BridgeHost string `json:"bridgeHost"`
	Lights     []int  `json:"lights"`
}

func (cli *CommandLineInterface) saveConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(homeDir+configPath, os.ModePerm)
	if err != nil {
		return err
	}
	file := homeDir + configFile
	if _, err := os.Stat(file); os.IsNotExist(err) {
		_, err := os.Create(file)
		if err != nil {
			return err
		}
	}
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err = encoder.Encode(cli.config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, content.Bytes(), 0644)
	return err
}

func (cli *CommandLineInterface) loadConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	_, err = os.Stat(homeDir + configFile)

	if os.IsNotExist(err) {
		err := cli.saveConfig()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var config config
	fileContent, err := ioutil.ReadFile(homeDir + configFile)
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}
	cli.config = config
	return nil
}

