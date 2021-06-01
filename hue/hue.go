package hue

import (
	"fmt"
	"github.com/amimof/huego"
	"strings"
)

type HSB struct {
	Hue         uint16
	Saturation  uint8
	Brightness  uint8
	Temperature uint16
}

var (
	Red    = HSB{65535, 254, 254, 0}
	Orange = HSB{6800, 254, 254, 0}
	Yellow = HSB{10900, 254, 254, 0}
	Green  = HSB{24400, 254, 254, 0}
	Cyan   = HSB{41600, 254, 254, 0}
	Blue   = HSB{46000, 254, 254, 0}
	Purple = HSB{48500, 254, 254, 0}
	Pink   = HSB{59800, 254, 254, 0}
	White  = HSB{0, 0, 254, 1}
	Mood   = HSB{0, 0, 254, 550}
)

func Discover(name string) (*huego.Bridge, string, error) {
	bridge, err := huego.Discover()
	if err != nil {

		return nil, "", err
	}
	user, err := bridge.CreateUser(name)
	if err != nil {
		return nil, "", err
	}
	bridge = bridge.Login(user)
	return bridge, user, nil
}

func Connect(host, appID string, ids []int) (*huego.Bridge, []*huego.Light, error) {
	var lights []*huego.Light
	bridge := huego.New(host, appID)
	for _, id := range ids {
		light, err := bridge.GetLight(id)
		if err != nil {
			return nil, nil, err
		}
		lights = append(lights, light)
	}
	return bridge, lights, nil
}

func GetColor(color string) HSB {
	color = strings.ToLower(color)
	switch color {
	case "red":
		return Red
	case "orange":
		return Orange
	case "yellow":
		return Yellow
	case "green":
		return Green
	case "cyan":
		return Cyan
	case "blue":
		return Blue
	case "purple":
		return Purple
	case "pink":
		return Pink
	case "white":
		return White
	case "mood":
		return Mood
	default:
		return White
	}
}

func Switch(light *huego.Light) error {
	if light.State.On {
		return light.Off()
	}
	return light.On()
}

func SetColor(light *huego.Light, hsb HSB) error {
	if !light.State.On {
		err := light.On()
		if err != nil {
			return err
		}
	}
	err := light.Hue(hsb.Hue)
	if err != nil {
		return err
	}
	err = light.Sat(hsb.Saturation)
	if err != nil {
		return err
	}
	err = light.Bri(hsb.Brightness)
	if err != nil {
		return err
	}
	if hsb.Temperature == 0 {
		return nil
	}
	return light.Ct(hsb.Temperature)
}

func SetHue(light *huego.Light, hue uint16) error {
	if !light.State.On {
		err := light.On()
		if err != nil {
			return err
		}
	}
	return light.Hue(hue)
}

func SetSaturation(light *huego.Light, saturation uint8) error {
	if !light.State.On {
		err := light.On()
		if err != nil {
			return err
		}
	}
	return light.Sat(saturation)
}

func SetBrightness(light *huego.Light, brightness uint8) error {
	if !light.State.On {
		err := light.On()
		if err != nil {
			return err
		}
	}
	return light.Bri(brightness)
}

func SetTemperature(light *huego.Light, temperature uint16) error {
	if !light.State.On {
		err := light.On()
		if err != nil {
			return err
		}
	}
	return light.Ct(temperature)
}

func PrintInfo(light *huego.Light) {
	fmt.Printf("Light %d:\n", light.ID)
	fmt.Printf("  Hue: %d\n", light.State.Hue)
	fmt.Printf("  Sat: %d\n", light.State.Sat)
	fmt.Printf("  Bri: %d\n", light.State.Bri)
	fmt.Printf("  Ct:  %d\n", light.State.Ct)
}
