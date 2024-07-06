package main

import (
	"context"
	"machine"

	keyboard "github.com/xcd0/tinygo-keyboard"
	"github.com/xcd0/tinygo-keyboard/keycodes/jp"
)

func main() {
	d := keyboard.New()

	gpioPins := []machine.Pin{
		machine.D0,
		machine.D3,
	}

	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	d.AddGpioKeyboard(gpioPins, [][]keyboard.Keycode{
		{
			jp.KeyA,
			jp.KeyB,
		},
	})

	loadKeyboardDef()
	d.Loop(context.Background())
}
