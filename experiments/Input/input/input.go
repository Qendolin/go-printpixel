package input

import (
	"fmt"
	"strings"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	JoystickDeadzoneLeft         = float32(0.1)
	JoystickDeadzoneRight        = float32(0.1)
	JoystickDeadzoneLeftTrigger  = float32(-0.5)
	JoystickDeadzoneRightTrigger = float32(-0.5)
)

var joystickAxesDeadzones = [...]*float32{&JoystickDeadzoneLeft, &JoystickDeadzoneLeft, &JoystickDeadzoneRight, &JoystickDeadzoneRight}

const (
	MouseButtonPrefix   = "mouse"
	MouseWheelPrefix    = "wheel"
	GamepadButtonPrefix = "gamepad"
	GamepadAxisPrefix   = "axis"
	ModifierKeyPrefix   = "mod"
	KeyPrefix           = "key"
)

var (
	MouseButtonNames = map[string]uint8{
		"1":      uint8(glfw.MouseButton1),
		"2":      uint8(glfw.MouseButton2),
		"3":      uint8(glfw.MouseButton3),
		"4":      uint8(glfw.MouseButton4),
		"5":      uint8(glfw.MouseButton5),
		"6":      uint8(glfw.MouseButton6),
		"7":      uint8(glfw.MouseButton7),
		"8":      uint8(glfw.MouseButton8),
		"left":   uint8(glfw.MouseButtonLeft),
		"middle": uint8(glfw.MouseButtonMiddle),
		"right":  uint8(glfw.MouseButtonRight),
	}
	MouseWheelNames = map[string]int8{
		"up":   1,
		"down": -1,
	}
	GamepadButtonNames = map[string]uint16{
		"a":           uint16(glfw.ButtonA),
		"b":           uint16(glfw.ButtonB),
		"x":           uint16(glfw.ButtonX),
		"y":           uint16(glfw.ButtonY),
		"leftbumper":  uint16(glfw.ButtonLeftBumper),
		"rightbumper": uint16(glfw.ButtonRightBumper),
		"back":        uint16(glfw.ButtonBack),
		"start":       uint16(glfw.ButtonStart),
		"guide":       uint16(glfw.ButtonGuide),
		"leftthumb":   uint16(glfw.ButtonLeftThumb),
		"rightthumb":  uint16(glfw.ButtonRightThumb),
		"dpadup":      uint16(glfw.ButtonDpadUp),
		"dpadright":   uint16(glfw.ButtonDpadRight),
		"dpaddown":    uint16(glfw.ButtonDpadDown),
		"dpadleft":    uint16(glfw.ButtonDpadLeft),
		"cross":       uint16(glfw.ButtonCross),
		"circle":      uint16(glfw.ButtonCircle),
		"square":      uint16(glfw.ButtonSquare),
		"triangle":    uint16(glfw.ButtonTriangle),
	}
	GamepadAxesNames = map[string]uint8{
		"leftx-":        0 + 2*uint8(glfw.AxisLeftX),
		"leftx+":        1 + 2*uint8(glfw.AxisLeftX),
		"lefty-":        0 + 2*uint8(glfw.AxisLeftY),
		"lefty+":        1 + 2*uint8(glfw.AxisLeftY),
		"rightx-":       0 + 2*uint8(glfw.AxisRightX),
		"rightx+":       1 + 2*uint8(glfw.AxisRightX),
		"righty-":       0 + 2*uint8(glfw.AxisRightY),
		"righty+":       1 + 2*uint8(glfw.AxisRightY),
		"lefttrigger-":  0 + 2*uint8(glfw.AxisLeftTrigger),
		"lefttrigger+":  1 + 2*uint8(glfw.AxisLeftTrigger),
		"righttrigger-": 0 + 2*uint8(glfw.AxisRightTrigger),
		"righttrigger+": 1 + 2*uint8(glfw.AxisRightTrigger),
	}
	ModifierKeyNames = map[string]uint8{
		"shift":    uint8(glfw.ModShift),
		"control":  uint8(glfw.ModControl),
		"alt":      uint8(glfw.ModAlt),
		"super":    uint8(glfw.ModSuper),
		"capslock": uint8(glfw.ModCapsLock),
		"numlock":  uint8(glfw.ModNumLock),
	}
	KeyNames = map[string]int{
		"unknown":      int(glfw.KeyUnknown),
		"space":        int(glfw.KeySpace),
		"apostrophe":   int(glfw.KeyApostrophe),
		"comma":        int(glfw.KeyComma),
		"minus":        int(glfw.KeyMinus),
		"period":       int(glfw.KeyPeriod),
		"slash":        int(glfw.KeySlash),
		"0":            int(glfw.Key0),
		"1":            int(glfw.Key1),
		"2":            int(glfw.Key2),
		"3":            int(glfw.Key3),
		"4":            int(glfw.Key4),
		"5":            int(glfw.Key5),
		"6":            int(glfw.Key6),
		"7":            int(glfw.Key7),
		"8":            int(glfw.Key8),
		"9":            int(glfw.Key9),
		"semicolon":    int(glfw.KeySemicolon),
		"equal":        int(glfw.KeyEqual),
		"a":            int(glfw.KeyA),
		"b":            int(glfw.KeyB),
		"c":            int(glfw.KeyC),
		"d":            int(glfw.KeyD),
		"e":            int(glfw.KeyE),
		"f":            int(glfw.KeyF),
		"g":            int(glfw.KeyG),
		"h":            int(glfw.KeyH),
		"i":            int(glfw.KeyI),
		"j":            int(glfw.KeyJ),
		"k":            int(glfw.KeyK),
		"l":            int(glfw.KeyL),
		"m":            int(glfw.KeyM),
		"n":            int(glfw.KeyN),
		"o":            int(glfw.KeyO),
		"p":            int(glfw.KeyP),
		"q":            int(glfw.KeyQ),
		"r":            int(glfw.KeyR),
		"s":            int(glfw.KeyS),
		"t":            int(glfw.KeyT),
		"u":            int(glfw.KeyU),
		"v":            int(glfw.KeyV),
		"w":            int(glfw.KeyW),
		"x":            int(glfw.KeyX),
		"y":            int(glfw.KeyY),
		"z":            int(glfw.KeyZ),
		"leftbracket":  int(glfw.KeyLeftBracket),
		"backslash":    int(glfw.KeyBackslash),
		"rightbracket": int(glfw.KeyRightBracket),
		"graveaccent":  int(glfw.KeyGraveAccent),
		"world1":       int(glfw.KeyWorld1),
		"world2":       int(glfw.KeyWorld2),
		"escape":       int(glfw.KeyEscape),
		"enter":        int(glfw.KeyEnter),
		"tab":          int(glfw.KeyTab),
		"backspace":    int(glfw.KeyBackspace),
		"insert":       int(glfw.KeyInsert),
		"delete":       int(glfw.KeyDelete),
		"right":        int(glfw.KeyRight),
		"left":         int(glfw.KeyLeft),
		"down":         int(glfw.KeyDown),
		"up":           int(glfw.KeyUp),
		"pageup":       int(glfw.KeyPageUp),
		"pagedown":     int(glfw.KeyPageDown),
		"home":         int(glfw.KeyHome),
		"end":          int(glfw.KeyEnd),
		"capslock":     int(glfw.KeyCapsLock),
		"scrolllock":   int(glfw.KeyScrollLock),
		"numlock":      int(glfw.KeyNumLock),
		"printscreen":  int(glfw.KeyPrintScreen),
		"pause":        int(glfw.KeyPause),
		"f1":           int(glfw.KeyF1),
		"f2":           int(glfw.KeyF2),
		"f3":           int(glfw.KeyF3),
		"f4":           int(glfw.KeyF4),
		"f5":           int(glfw.KeyF5),
		"f6":           int(glfw.KeyF6),
		"f7":           int(glfw.KeyF7),
		"f8":           int(glfw.KeyF8),
		"f9":           int(glfw.KeyF9),
		"f10":          int(glfw.KeyF10),
		"f11":          int(glfw.KeyF11),
		"f12":          int(glfw.KeyF12),
		"f13":          int(glfw.KeyF13),
		"f14":          int(glfw.KeyF14),
		"f15":          int(glfw.KeyF15),
		"f16":          int(glfw.KeyF16),
		"f17":          int(glfw.KeyF17),
		"f18":          int(glfw.KeyF18),
		"f19":          int(glfw.KeyF19),
		"f20":          int(glfw.KeyF20),
		"f21":          int(glfw.KeyF21),
		"f22":          int(glfw.KeyF22),
		"f23":          int(glfw.KeyF23),
		"f24":          int(glfw.KeyF24),
		"f25":          int(glfw.KeyF25),
		"kp0":          int(glfw.KeyKP0),
		"kp1":          int(glfw.KeyKP1),
		"kp2":          int(glfw.KeyKP2),
		"kp3":          int(glfw.KeyKP3),
		"kp4":          int(glfw.KeyKP4),
		"kp5":          int(glfw.KeyKP5),
		"kp6":          int(glfw.KeyKP6),
		"kp7":          int(glfw.KeyKP7),
		"kp8":          int(glfw.KeyKP8),
		"kp9":          int(glfw.KeyKP9),
		"kpdecimal":    int(glfw.KeyKPDecimal),
		"kpdivide":     int(glfw.KeyKPDivide),
		"kpmultiply":   int(glfw.KeyKPMultiply),
		"kpsubtract":   int(glfw.KeyKPSubtract),
		"kpadd":        int(glfw.KeyKPAdd),
		"kpenter":      int(glfw.KeyKPEnter),
		"kpequal":      int(glfw.KeyKPEqual),
		"leftshift":    int(glfw.KeyLeftShift),
		"leftcontrol":  int(glfw.KeyLeftControl),
		"leftalt":      int(glfw.KeyLeftAlt),
		"leftsuper":    int(glfw.KeyLeftSuper),
		"rightshift":   int(glfw.KeyRightShift),
		"rightcontrol": int(glfw.KeyRightControl),
		"rightalt":     int(glfw.KeyRightAlt),
		"rightsuper":   int(glfw.KeyRightSuper),
		"menu":         int(glfw.KeyMenu),
	}
)

type window interface {
	SetKeyCallback(glfw.KeyCallback) glfw.KeyCallback
	SetMouseButtonCallback(glfw.MouseButtonCallback) glfw.MouseButtonCallback
	SetScrollCallback(glfw.ScrollCallback) glfw.ScrollCallback
}

type Trigger struct {
	MouseButtons   uint8
	MouseWheel     int8
	GamepadButtons uint16
	GamepadAxes    uint16
	ModifierKeys   uint8
	Key            int
	Error          string
}

func (tr Trigger) Valid() bool {
	return tr.Error == ""
}

func Combo(inputs ...string) (trigger Trigger) {
	for _, rawInput := range inputs {
		input := strings.ToLower(rawInput)
		if strings.HasPrefix(input, MouseButtonPrefix) {
			value, ok := MouseButtonNames[input[len(MouseButtonPrefix):]]
			if !ok {
				trigger.Error = fmt.Sprintf("unknown mouse button %q", rawInput)
				return
			}
			trigger.MouseButtons |= 1 << value
			continue
		}
		if strings.HasPrefix(input, MouseWheelPrefix) {
			value, ok := MouseWheelNames[input[len(MouseWheelPrefix):]]
			if !ok {
				trigger.Error = fmt.Sprintf("unknown mouse wheel direction %q", rawInput)
				return
			}
			trigger.MouseWheel |= 1 << value
			continue
		}
		if strings.HasPrefix(input, GamepadButtonPrefix) {
			value, ok := GamepadButtonNames[input[len(GamepadButtonPrefix):]]
			if !ok {
				trigger.Error = fmt.Sprintf("unknown gamepad button %q", rawInput)
				return
			}
			trigger.GamepadButtons |= 1 << value
			continue
		}
		if strings.HasPrefix(input, GamepadAxisPrefix) {
			value, ok := GamepadAxesNames[input[len(GamepadAxisPrefix):]]
			if !ok {
				trigger.Error = fmt.Sprintf("unknown gamepad axis %q", rawInput)
				return
			}
			trigger.GamepadAxes |= 1 << value
			continue
		}
		if strings.HasPrefix(input, ModifierKeyPrefix) {
			value, ok := ModifierKeyNames[input[len(ModifierKeyPrefix):]]
			if !ok {
				trigger.Error = fmt.Sprintf("unknown modifier key %q", rawInput)
				return
			}
			trigger.ModifierKeys |= value
			continue
		}
		if strings.HasPrefix(input, KeyPrefix) {
			value, ok := KeyNames[input[len(KeyPrefix):]]
			if !ok {
				trigger.Error = fmt.Sprintf("unknown key %q", rawInput)
				return
			}
			trigger.Key = value
			continue
		}
		trigger.Error = fmt.Sprintf("unknown input %q", rawInput)
		return
	}
	return
}

type InputState struct {
	Keyboard KeyboardState
	Gamepad  GamepadState
	Mouse    MouseState
}

type KeyboardState struct {
	Keys      [glfw.KeyLast]bool
	Modifiers uint8
}

type GamepadState struct {
	Buttons uint16
	Axes    uint16
}

type MouseState struct {
	Buttons uint8
	Wheel   int8
}

type Action string

type ActionEvent struct {
	Action Action
}

type Handler func(ActionEvent)

type Manager struct {
	triggers       map[Trigger]Action
	override       map[Trigger]Action
	stickyOverride bool
	handlers       map[Action][]Handler
	activeTrigger  Trigger
	state          InputState
}

var Default = NewManager()

func NewManager() Manager {
	return Manager{
		triggers:       map[Trigger]Action{},
		override:       nil,
		stickyOverride: false,
		handlers:       map[Action][]Handler{},
	}
}

func (m *Manager) AddTrigger(trigger Trigger, action Action) error {
	if !trigger.Valid() {
		return fmt.Errorf("%v", trigger.Error)
	}
	m.triggers[trigger] = action
	return nil
}

func (m *Manager) RemoveTrigger(trigger Trigger) {
	delete(m.triggers, trigger)
}

// Set a temporary override for all triggers
// The next active input will clear the override unless sticky is true
// This can be used to create more complex combinations. E.g.: Ctrl + V + S
func (m *Manager) SetOverride(triggers map[Trigger]Action, sticky bool) {
	m.override = triggers
	m.stickyOverride = sticky
}

func (m *Manager) ClearOverride() {
	m.override = nil
	m.stickyOverride = false
}

func (m *Manager) On(action Action, handler Handler) {
	handlers := m.handlers[action]
	handlers = append(handlers, handler)
	m.handlers[action] = handlers
}

func (m *Manager) Bind(w window) {
	w.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		m.state.Keyboard.Keys[key] = action == glfw.Press
		m.state.Keyboard.Modifiers = uint8(mods)
		m.activeTrigger.ModifierKeys = uint8(mods)
		activeInput := false
		// Above 340 are modifier keys
		if key < 340 {
			if action != glfw.Press {
				m.activeTrigger.Key = 0
			} else {
				m.activeTrigger.Key = int(key)
				activeInput = true
			}
		}
		m.Check(activeInput)
	})
	w.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		activeInput := false
		if action != glfw.Press {
			m.state.Mouse.Buttons &^= 1 << button
			m.activeTrigger.MouseButtons &^= 1 << button
		} else {
			m.state.Mouse.Buttons |= 1 << button
			m.activeTrigger.MouseButtons |= 1 << button
			activeInput = true
		}
		m.activeTrigger.ModifierKeys = uint8(mods)
		m.state.Keyboard.Modifiers = uint8(mods)
		m.Check(activeInput)
	})
	w.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
		if yoff > 0 {
			m.state.Mouse.Wheel = 1
			m.activeTrigger.MouseWheel = 1
		} else if yoff < 0 {
			m.state.Mouse.Wheel = -1
			m.activeTrigger.MouseWheel = -1
		}
		m.Check(true)
	})
}

func (m *Manager) Update() {
	glfw.PollEvents()

	m.state.Mouse.Wheel = 0
	m.activeTrigger.MouseWheel = 0
	m.activeTrigger.Key = 0

	activeInput := false
	if glfw.Joystick1.IsGamepad() {
		state := glfw.Joystick1.GetGamepadState()
		var buttonState uint16
		for i := 0; i < 15; i++ {
			if state.Buttons[i] == glfw.Press {
				mask := uint16(1 << i)
				activeInput = activeInput || buttonState&mask == 0
				buttonState |= mask
			}
		}
		m.state.Gamepad.Buttons = buttonState
		m.activeTrigger.GamepadButtons = buttonState
		var axesState uint16
		for i := 0; i < 4; i++ {
			if state.Axes[i] > *joystickAxesDeadzones[i] {
				mask := uint16(1 << (i * 2))
				activeInput = activeInput || m.activeTrigger.GamepadButtons&mask == 0
				axesState |= mask
			} else if state.Axes[i] < -*joystickAxesDeadzones[i] {
				mask := uint16(1 << (i*2 + 1))
				activeInput = activeInput || m.activeTrigger.GamepadButtons&mask == 0
				axesState |= mask
			}
		}
		if state.Axes[4] > JoystickDeadzoneLeftTrigger {
			const mask = 1 << 8
			activeInput = activeInput || axesState&mask == 0
			axesState |= mask
		}
		if state.Axes[5] > JoystickDeadzoneRightTrigger {
			const mask = 1 << 9
			activeInput = activeInput || axesState&mask == 0
			axesState |= mask
		}
		m.state.Gamepad.Axes = axesState
		m.activeTrigger.GamepadAxes = axesState
	}

	m.Check(activeInput)
}

// An 'active' input is a change away from the default state, i.e. a key press, a mouse scroll.
// A 'passive' input would be a repeat or a release
func (m *Manager) Check(activeInput bool) {
	activeTriggers := m.triggers
	if m.override != nil {
		activeTriggers = m.override
		if !m.stickyOverride && activeInput {
			m.ClearOverride()
		}
	}

	act, found := activeTriggers[m.activeTrigger]
	if !found {
		return
	}

	ev := ActionEvent{
		Action: act,
	}
	for _, handle := range m.handlers[act] {
		handle(ev)
	}
}
