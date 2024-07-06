//go:build tinygo

package keyboard

import (
	"machine"
	"time"
)

type MatrixKeyboard struct {
	State    []State
	Keys     [][]Keycode
	options  Options
	callback Callback

	Col           []machine.Pin
	Row           []machine.Pin
	cycleCounter  []uint8
	debounce      uint8
	sleepDuration time.Duration
}

func (d *Device) AddMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][]Keycode, opt ...Option) *MatrixKeyboard {
	col := len(colPins)
	row := len(rowPins)
	state := make([]State, row*col)
	cycleCnt := make([]uint8, len(state))

	for c := range colPins {
		colPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}
	for r := range rowPins {
		rowPins[r].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	o := Options{}
	for _, f := range opt {
		f(&o)
	}

	keydef := make([][]Keycode, LayerCount)
	for l := 0; l < len(keydef); l++ {
		keydef[l] = make([]Keycode, len(state))
	}
	for l := 0; l < len(keys); l++ {
		for kc := 0; kc < len(keys[l]); kc++ {
			keydef[l][kc] = keys[l][kc]
		}
	}

	k := &MatrixKeyboard{
		Col:           colPins,
		Row:           rowPins,
		State:         state,
		Keys:          keydef,
		options:       o,
		callback:      func(layer, index int, state State) {},
		cycleCounter:  cycleCnt,
		debounce:      8,
		sleepDuration: 0,
	}
	if o.MatrixScanPeriod != 0 {
		// Calculate the sleep duration before and after each key read, based on the total scan period for the entire key matrix.
		k.sleepDuration = o.MatrixScanPeriod / time.Duration(len(k.Col)*len(k.Row))
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *MatrixKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *MatrixKeyboard) Callback(layer, index int, state State) {
	if d.callback != nil {
		d.callback(layer, index, state)
	}
}

func (d *MatrixKeyboard) UpdateKeyState(idx int, current bool) {
	switch d.State[idx] {
	case None:
		if current {
			if d.cycleCounter[idx] >= d.debounce {
				d.State[idx] = NoneToPress
				d.cycleCounter[idx] = 0
			} else {
				d.cycleCounter[idx]++
			}
		} else {
			d.cycleCounter[idx] = 0
		}
	case NoneToPress:
		d.State[idx] = Press
	case Press:
		if current {
			d.cycleCounter[idx] = 0
		} else {
			if d.cycleCounter[idx] >= d.debounce {
				d.State[idx] = PressToRelease
				d.cycleCounter[idx] = 0
			} else {
				d.cycleCounter[idx]++
			}
		}
	case PressToRelease:
		d.State[idx] = None
	}
}

func (d *MatrixKeyboard) readPin(out, in []machine.Pin, i, j int) bool {
	out[i].Configure(machine.PinConfig{Mode: machine.PinOutput})
	out[i].High()
	time.Sleep(d.sleepDuration) // 読み取り前に待機
	return in[j].Get()
}

func (d *MatrixKeyboard) Get() []State {
	current := false
	if !d.options.InvertDiode {
		for c := range d.Col {
			for r := range d.Row {
				current = d.readPin(d.Col, d.Row, c, r)
				d.UpdateKeyState(r*len(d.Col)+c, current)
			}
			d.Col[c].Low()
			d.Col[c].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
		}
	} else {
		for r := range d.Row {
			for c := range d.Col {
				current = d.readPin(d.Col, d.Row, c, r)
				d.UpdateKeyState(r*len(d.Col)+c, current)
			}
			d.Row[r].Low()
			d.Row[r].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
		}
	}
	return d.State
}

func (d *MatrixKeyboard) Key(layer, index int) Keycode {
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *MatrixKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *MatrixKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *MatrixKeyboard) Init() error {
	return nil
}
