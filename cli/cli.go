package cli

import (
	"flag"
	"fmt"
	"github.com/MarekSalgovic/hue-go/hue"
	"github.com/amimof/huego"
	"log"
	"os"
)



type CommandLineInterface struct {
	bridge *huego.Bridge
	lights []*huego.Light
	config config
}

func NewCLI() (*CommandLineInterface, error) {
	cli := &CommandLineInterface{}
	err := cli.loadConfig()
	if err != nil {
		return nil, err
	}
	bridge, lights, err := hue.Connect(cli.config.BridgeHost, cli.config.AppID, cli.config.Lights)
	if err != nil {
		return nil, err
	}
	cli.bridge = bridge
	cli.lights = lights
	return cli, nil
}


func (cli *CommandLineInterface) init() error {
	bridge, user, err := hue.Discover("Philips Hue CLI app")
	cli.config.BridgeHost = bridge.Host
	cli.config.AppID = user
	if err != nil {
		log.Printf("error with hue discover: %s\n", err)
		return err
	}
	err = cli.saveConfig()
	if err != nil {
		log.Printf("error with saving config")
		return err
	}
	return nil
}

func (cli *CommandLineInterface) listLights() error {
	lights, err := cli.bridge.GetLights()
	if err != nil {
		return err
	}
	for _, l := range lights {
		fmt.Printf("ID: %d  name: %s\n", l.ID, l.Name)
	}
	return nil
}

func (cli *CommandLineInterface) addLight(id int) error {
	for _, l := range cli.config.Lights {
		if l == id {
			return nil
		}
	}
	cli.config.Lights = append(cli.config.Lights, id)
	return cli.saveConfig()
}

func (cli *CommandLineInterface) removeLight(id int) error {
	for i, l := range cli.config.Lights {
		if l == id {
			cli.config.Lights = append(cli.config.Lights[:i], cli.config.Lights[i+1:]...)
		}
	}
	return cli.saveConfig()
}

func (cli *CommandLineInterface) changeLights(color string, brightness, id int) error {
	for _, l := range cli.lights {
		if id != l.ID && id != 0 {
			continue
		}
		c := hue.GetColor(color)
		if brightness != 0 {
			c.Brightness = uint8(float64(brightness) * (255.0 / 100.0))
		}
		err := hue.SetColor(l, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cli *CommandLineInterface) switchLights() error {
	for _, l := range cli.lights {
		err := hue.Switch(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cli *CommandLineInterface) info() {
	fmt.Printf("Philips Hue CLI App ID:  %s\n", cli.config.AppID)
	fmt.Printf("Philips Hue Bridge host: %s\n", cli.config.BridgeHost)
	for _, l := range cli.lights {
		hue.PrintInfo(l)
	}
}

func (cli *CommandLineInterface) printUsage() {
	fmt.Println()
}

func (cli *CommandLineInterface) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(2)
	}
}

func (cli *CommandLineInterface) Run() {
	cli.validateArgs()

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	changeCmd := flag.NewFlagSet("change", flag.ExitOnError)
	switchCmd := flag.NewFlagSet("switch", flag.ExitOnError)
	lightsCmd := flag.NewFlagSet("lights", flag.ExitOnError)
	addLightCmd := flag.NewFlagSet("add", flag.ExitOnError)
	removeLightCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)

	changeCmdColor := changeCmd.String("color", "", "The color to change the light to.")
	changeCmdBrightness := changeCmd.Int("brightness", 0, "The brightness percentage to change the light to.")
	changeCmdID := changeCmd.Int("id", 0, "The ID of light to be changed.")

	addLightCmdID := addLightCmd.Int("id", 0, "The ID of light to be added to app.")
	removeLightCmdID := removeLightCmd.Int("id", 0, "The ID of light to be removed from app.")

	switch os.Args[1] {
	case "init":
		err := initCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			os.Exit(5)
		}
	case "change":
		err := changeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			os.Exit(5)
		}
	case "switch":
		err := switchCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			os.Exit(5)
		}
	case "lights":
		err := lightsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			os.Exit(5)
		}
	case "add":
		err := addLightCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			os.Exit(5)
		}
	case "remove":
		err := removeLightCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			os.Exit(5)
		}
	case "info":
		err := infoCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			os.Exit(5)
		}
	default:
		log.Println("undefined command")
		os.Exit(2)
	}

	if initCmd.Parsed() {
		err := cli.init()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	if changeCmd.Parsed() {
		err := cli.changeLights(*changeCmdColor, *changeCmdBrightness, *changeCmdID)
		if err != nil {
			log.Printf("hue error %s\n", err.Error())
			os.Exit(1)
		}
	}

	if switchCmd.Parsed() {
		err := cli.switchLights()
		if err != nil {
			log.Printf("hue error %s\n", err.Error())
			os.Exit(1)
		}
	}

	if lightsCmd.Parsed() {
		err := cli.listLights()
		if err != nil {
			log.Printf("hue error %s\n", err.Error())
			os.Exit(1)
		}
	}

	if addLightCmd.Parsed() {
		if *addLightCmdID == 0 {
			addLightCmd.Usage()
			log.Printf("argument error: light id missing\n")
			os.Exit(2)
		}
		err := cli.addLight(*addLightCmdID)
		if err != nil {
			log.Printf("hue error %s\n", err.Error())
			os.Exit(1)
		}
	}

	if removeLightCmd.Parsed() {
		if *removeLightCmdID == 0 {
			removeLightCmd.Usage()
			log.Printf("argument error: light id missing\n")
			os.Exit(2)
		}
		err := cli.removeLight(*removeLightCmdID)
		if err != nil {
			log.Printf("hue error %s\n", err.Error())
			os.Exit(1)
		}
	}

	if infoCmd.Parsed() {
		cli.info()
	}
}
