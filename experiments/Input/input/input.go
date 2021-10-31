package input

import (
	"sort"
	"strings"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// TODO: Localized key names

var (
	Gamepad1AxesConfig GamepadAxesConfig
)

func init() {
	Gamepad1AxesConfig = GamepadAxesConfig{
		Indices:        [6]int{0, 1, 3, 4, 2, 2},
		UpperDeadzones: [10]float32{0.95, 0.95, 0.95, 0.95, 0.95, 0.95, 0.95, 0.95, 0.95, 0.95},
		LowerDeadzones: [10]float32{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
	}
	Gamepad1AxesConfig.CalculateHalfwayPoint()
}

// GamepadAxesConfig
// Order:
// AxisLeftX- (Left), AxisLeftX+ (Right), AxisLeftY- (Up), AxisLeftY+ (Down),
// AxisRightX- (Left), AxisRightX+ (Right), AxisRightY- (Up), AxisRightY+ (Down),
// AxisLeftTrigger (1), AxisRightTrigger (-1)
type GamepadAxesConfig struct {
	Indices        [6]int
	UpperDeadzones [10]float32
	LowerDeadzones [10]float32
	HalfwayPoints  [10]float32
}

func (c *GamepadAxesConfig) CalculateHalfwayPoint() {
	for i := range c.LowerDeadzones {
		l, u := c.LowerDeadzones[i], c.UpperDeadzones[i]
		c.HalfwayPoints[i] = (u - l) / 2
	}
}

// var joystickAxesDeadzones = [...]*float32{&JoystickDeadzoneLeftLower, &JoystickDeadzoneLeftLower, &JoystickDeadzoneRightLower, &JoystickDeadzoneRightLower}

type window interface {
	SetKeyCallback(glfw.KeyCallback) glfw.KeyCallback
	SetMouseButtonCallback(glfw.MouseButtonCallback) glfw.MouseButtonCallback
	SetScrollCallback(glfw.ScrollCallback) glfw.ScrollCallback
	SetCursorPosCallback(glfw.CursorPosCallback) glfw.CursorPosCallback
}

// 8 inputs is enough for everyone, right?
const MaxTriggerSources = 8

type Trigger struct {
	Sources     [MaxTriggerSources]SourceId
	SourceCount uint8
}

func (tr Trigger) String() string {
	names := []string{}

	for i := tr.SourceCount - 1; i != 0xff; i-- {
		id := tr.Sources[i]
		names = append(names, Sources[int(id)].String())
	}

	return strings.Join(names, " + ")
}

func (tr *Trigger) Add(source Source) {
	if int(tr.SourceCount) == cap(tr.Sources) || source.Id == 0 {
		return
	}
	// Keep order, lower ids at the start
	for i := 0; i < int(tr.SourceCount); i++ {
		if source.Id <= tr.Sources[i] {
			if source.Id == tr.Sources[i] {
				return
			}
			copy(tr.Sources[i+1:], tr.Sources[i:])
			tr.Sources[i] = source.Id
			tr.SourceCount++
			return
		}
	}
	tr.Sources[tr.SourceCount] = source.Id
	tr.SourceCount++
}

func (tr *Trigger) Remove(source Source) {
	if int(tr.SourceCount) == 0 || source.Id == 0 {
		return
	}
	// Keep order, lower ids at the start
	for i := 0; i < int(tr.SourceCount); i++ {
		if source.Id == tr.Sources[i] {
			copy(tr.Sources[i:], tr.Sources[i+1:])
			tr.SourceCount--
			tr.Sources[tr.SourceCount] = 0
			return
		}
	}
}

func (tr *Trigger) Sort() (trigger Trigger) {
	sort.Slice(tr.Sources[:], func(i, j int) bool {
		return tr.Sources[i] < tr.Sources[j]
	})
	return
}

func (tr *Trigger) Has(source Source) bool {
	if int(tr.SourceCount) == 0 || source.Id == 0 {
		return false
	}
	// Keep order, lower ids at the start
	for i := 0; i < int(tr.SourceCount); i++ {
		if source.Id == tr.Sources[i] {
			return true
		} else if source.Id < tr.Sources[i] {
			return false
		}
	}
	return false
}

func Single(input string) (trigger Trigger) {
	// trigger.Add(input)
	trigger.Add(NamedSources[strings.ToLower(input)])
	return
}

func Combo(inputs ...string) (trigger Trigger) {
	for _, input := range inputs {
		// trigger.Add(input)
		trigger.Add(NamedSources[strings.ToLower(input)])
	}
	return
}

type InputState struct {
	Keyboard KeyboardState
	Gamepad  GamepadState
	Mouse    MouseState
}

type KeyboardState struct {
	Keys      [glfw.KeyLast + 1]bool
	Modifiers uint8
}

type GamepadState struct {
	Buttons [glfw.ButtonLast + 1]bool
	Axes    [glfw.AxisLast + 1]float32
}

type MouseState struct {
	Buttons uint8
	WheelY  int8
	WheelX  int8
	PosX    float32
	PosY    float32
}

type Action string

type ActionEvent struct {
	Action Action
	Active bool
	InputState
}

type Handler func(ActionEvent)

// TODO: implement repeat lifecycle
// TODO: implement deactivate lifecycle
// TODO: test lifecycle logic
type ActionLifecycle uint8

// Given        | Required
// OnActivate   | OnActivate
// OnActivate   | WhileActive
// OnActivate   | WhileRepeat
// WhileActive  | WhileActive
// WhileRepeat  | WhileActive
// WhileRepeat  | WhileRepeat
// OnDeactivate | OnDeactivate
const (
	OnActivate   ActionLifecycle = 0b00010001
	WhileActive  ActionLifecycle = 0b00100111
	WhileRepeat  ActionLifecycle = 0b01000101
	OnDeactivate ActionLifecycle = 0b10001000
)

type HandlerRegistration struct {
	handler   Handler
	action    Action
	lifecycle ActionLifecycle
}

func (reg HandlerRegistration) Handler() Handler {
	return reg.handler
}

func (reg HandlerRegistration) Invoke(ev ActionEvent) {
	reg.handler(ev)
}

func (reg HandlerRegistration) Lifecycle() ActionLifecycle {
	return reg.lifecycle
}

func (reg HandlerRegistration) Action() Action {
	return reg.action
}

type ActionState uint8

const (
	Inactive  = 0
	Triggered = 1
	Active    = 2
)

type Manager struct {
	triggers       map[Trigger]Action
	actionState    map[Action]ActionState
	override       map[Trigger]Action
	stickyOverride bool
	handlers       map[Action][]*HandlerRegistration
	activeTrigger  Trigger
	state          InputState
}

var Default = NewManager()

func NewManager() Manager {
	return Manager{
		triggers:       map[Trigger]Action{},
		override:       nil,
		stickyOverride: false,
		handlers:       map[Action][]*HandlerRegistration{},
		actionState:    map[Action]ActionState{},
	}
}

func (m *Manager) AddTrigger(trigger Trigger, action Action) {
	m.triggers[trigger] = action
}

func (m *Manager) AddTriggers(triggers map[Trigger]Action) {
	for tr, act := range triggers {
		m.triggers[tr] = act
	}
}

func (m *Manager) RemoveTrigger(trigger Trigger) {
	delete(m.triggers, trigger)
}

func (m *Manager) Triggers() map[Trigger]Action {
	return m.triggers
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

func (m *Manager) On(action Action, handler Handler) *HandlerRegistration {
	return m.Listen(action, OnActivate, handler)
}

func (m *Manager) While(action Action, handler Handler) *HandlerRegistration {
	return m.Listen(action, WhileActive, handler)
}

func (m *Manager) OnRepeat(action Action, handler Handler) *HandlerRegistration {
	return m.Listen(action, WhileRepeat, handler)
}

func (m *Manager) OnDeactivate(action Action, handler Handler) *HandlerRegistration {
	return m.Listen(action, OnDeactivate, handler)
}

func (m *Manager) Listen(action Action, lifecycle ActionLifecycle, handler Handler) *HandlerRegistration {
	handlers := m.handlers[action]
	reg := &HandlerRegistration{handler, action, lifecycle}
	handlers = append(handlers, reg)
	m.handlers[action] = handlers
	return reg
}

func (m *Manager) Off(reg *HandlerRegistration) {
	handlers := m.handlers[reg.action]
	for i := range handlers {
		if handlers[i] == reg {
			handlers[i] = handlers[len(handlers)-1]
			m.handlers[reg.action] = handlers[:len(handlers)-1]
			break
		}
	}
}

func (m *Manager) State() InputState {
	return m.state
}

func (m *Manager) Bind(w window) {
	w.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		m.state.Keyboard.Keys[key] = action == glfw.Press
		m.state.Keyboard.Modifiers = uint8(mods)
		// m.activeTrigger.ModifierKeys = uint8(mods)
		activeInput := false
		// Above 340 are modifier keys
		// if key < 340 {
		if action == glfw.Press {
			m.activeTrigger.Add(KeyboardKeySources[key])
			// m.activeTrigger.Key = int(key)
			activeInput = true
		} else if action == glfw.Release {
			m.activeTrigger.Remove(KeyboardKeySources[key])
			// m.activeTrigger.Key = 0
		}
		// }
		m.Check(activeInput)
	})
	w.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		activeInput := false
		if action == glfw.Press {
			m.activeTrigger.Add(MouseButtonSources[button])
			// m.state.Mouse.Buttons |= 1 << button
			// m.activeTrigger.MouseButtons |= 1 << button
			activeInput = true
		} else if action == glfw.Release {
			m.activeTrigger.Remove(MouseButtonSources[button])
			// m.state.Mouse.Buttons &^= 1 << button
			// m.activeTrigger.MouseButtons &^= 1 << button
		}
		// m.activeTrigger.ModifierKeys = uint8(mods)
		m.state.Keyboard.Modifiers = uint8(mods)
		m.Check(activeInput)
	})
	var (
		MouseWheelYUp    Source = NamedMouseWheels["mousewheelyup"]
		MouseWheelYDown  Source = NamedMouseWheels["mousewheelydown"]
		MouseWheelXLeft  Source = NamedMouseWheels["mousewheelxleft"]
		MouseWheelXRight Source = NamedMouseWheels["mousewheelxright"]
		MouseWheelX      Source = NamedMouseWheels["mousewheelx"]
		MouseWheelY      Source = NamedMouseWheels["mousewheely"]
	)
	w.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
		// FIXME: A offest greater / less than 1 / -1 is lost as handler does not have the info
		if yoff != 0 {
			m.state.Mouse.WheelY += int8(yoff)
			m.activeTrigger.Add(MouseWheelY)
			m.Check(true)
			m.activeTrigger.Remove(MouseWheelY)
		}
		if yoff > 0 {
			// m.activeTrigger.MouseWheel = 1
			// if !m.activeTrigger.Has(MouseWheelYUp) {
			m.activeTrigger.Add(MouseWheelYUp)
			for ; yoff != 0; yoff-- {
				m.Check(true)
			}
			m.activeTrigger.Remove(MouseWheelYUp)
			// }
		} else if yoff < 0 {
			// m.activeTrigger.MouseWheel = -1
			// if !m.activeTrigger.Has(MouseWheelYDown) {
			m.activeTrigger.Add(MouseWheelYDown)
			for ; yoff != 0; yoff++ {
				m.Check(true)
			}
			m.activeTrigger.Remove(MouseWheelYDown)
			// }
		}
		// FIXME: Left and right might be the other way around
		if xoff != 0 {
			m.state.Mouse.WheelX += int8(yoff)
			m.activeTrigger.Add(MouseWheelX)
			m.Check(true)
			m.activeTrigger.Remove(MouseWheelX)
		}
		if xoff > 0 {
			// m.activeTrigger.MouseWheel = 1
			// if !m.activeTrigger.Has(MouseWheelXRight) {
			m.activeTrigger.Add(MouseWheelXRight)
			for ; xoff != 0; xoff-- {
				m.Check(true)
			}
			m.activeTrigger.Remove(MouseWheelXRight)
			// }
		} else if xoff < 0 {
			// m.activeTrigger.MouseWheel = -1
			// if !m.activeTrigger.Has(MouseWheelXLeft) {
			m.activeTrigger.Add(MouseWheelXLeft)
			for ; xoff != 0; xoff++ {
				m.Check(true)
			}
			m.activeTrigger.Remove(MouseWheelXLeft)
			// }
		}
	})
	w.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		m.state.Mouse.PosX = float32(xpos)
		m.state.Mouse.PosY = float32(ypos)
	})
}

func (m *Manager) Update() {
	glfw.PollEvents()

	m.state.Mouse.WheelY = 0
	m.state.Mouse.WheelX = 0

	for act, state := range m.actionState {
		if state&Triggered != 0 {
			m.actionState[act] = Active
		} else if state != Inactive {
			ev := ActionEvent{
				Action: act,
			}
			for _, handler := range m.handlers[act] {
				if handler.Lifecycle() == OnDeactivate {
					handler.Invoke(ev)
				}
			}
			m.actionState[act] = Inactive
		}
	}

	m.activeTrigger.Remove(NamedMouseWheels["mousewheelxright"])
	m.activeTrigger.Remove(NamedMouseWheels["mousewheelxleft"])
	m.activeTrigger.Remove(NamedMouseWheels["mousewheelyup"])
	m.activeTrigger.Remove(NamedMouseWheels["mousewheelydown"])
	// m.activeTrigger.Key = 0

	activeInput := false
	// TODO: implement gamepad logic
	if glfw.Joystick1.IsGamepad() {
		state := glfw.Joystick1.GetGamepadState()
		// var buttonState uint16
		for i := 0; i < 15; i++ {
			if state.Buttons[i] == glfw.Press {
				// mask := uint16(1 << i)
				// activeInput = activeInput || buttonState&mask == 0
				// buttonState |= mask
				m.activeTrigger.Add(GamepadButtonSources[i])
				m.state.Gamepad.Buttons[i] = true
			} else {
				m.activeTrigger.Remove(GamepadButtonSources[i])
				m.state.Gamepad.Buttons[i] = false
			}
		}
		// m.state.Gamepad.Buttons = buttonState
		// m.activeTrigger.GamepadButtons = buttonState
		// var axesState uint16
		// TODO: Analog Sources???
		m.state.Gamepad.Axes = state.Axes
		for i := 0; i < 8; i++ {
			value := state.Axes[i/2]
			if value < 0 {
				value *= -1
			}
			dzLower := Gamepad1AxesConfig.LowerDeadzones[i]
			if value >= dzLower {
				m.activeTrigger.Add(DigitalAxesSources[i])
			} else {
				m.activeTrigger.Remove(DigitalAxesSources[i])
			}
		}
		// TODO: Implement analog and digital gamepad triggers (digital threshold is defined by the driver)
		for i := 0; i < 2; i++ {
			value := state.Axes[Gamepad1AxesConfig.Indices[i+4]]
			if value < 0 {
				value *= -1
			}
			dzLower := Gamepad1AxesConfig.LowerDeadzones[i*2]
			halfway := Gamepad1AxesConfig.HalfwayPoints[i*2]
			if value >= halfway {
				m.activeTrigger.Add(DigitalAxesSources[i*2+8+1])
			} else if value > dzLower {
				m.activeTrigger.Remove(DigitalAxesSources[i*2+8+1])
				m.activeTrigger.Add(DigitalAxesSources[i*2+8])
			} else {
				m.activeTrigger.Remove(DigitalAxesSources[i*2+8+1])
				m.activeTrigger.Remove(DigitalAxesSources[i*2+8])
			}
		}
		// for i := 0; i < 4; i++ {
		// 	if state.Axes[i] > *joystickAxesDeadzones[i] {
		// 		// mask := uint16(1 << (i * 2))
		// 		// activeInput = activeInput || m.activeTrigger.GamepadButtons&mask == 0
		// 		// axesState |= mask
		// 		m.activeTrigger.Remove(DigitalAxesSources[i*2+1])
		// 		m.activeTrigger.Add(DigitalAxesSources[i*2])
		// 		m.state.Gamepad.Axes[i*2] = int8(state.Axes[i])
		// 	} else if state.Axes[i] < -*joystickAxesDeadzones[i] {
		// 		// mask := uint16(1 << (i*2 + 1))
		// 		// activeInput = activeInput || m.activeTrigger.GamepadButtons&mask == 0
		// 		// axesState |= mask
		// 		m.activeTrigger.Remove(DigitalAxesSources[i*2])
		// 		m.activeTrigger.Add(DigitalAxesSources[i*2+1])
		// 		m.state.Gamepad.Axes[i*2+1] = int8(state.Axes[i])
		// 	} else {
		// 		m.activeTrigger.Remove(DigitalAxesSources[i*2+1])
		// 		m.activeTrigger.Remove(DigitalAxesSources[i*2])
		// 	}
		// }
		// if state.Axes[4] > TriggerDeadzoneLeftLower {
		// 	// const mask = 1 << 8
		// 	// activeInput = activeInput || axesState&mask == 0
		// 	// axesState |= mask
		// 	m.activeTrigger.Add(DigitalAxesSources[glfw.AxisLeftTrigger])
		// } else {
		// 	m.activeTrigger.Remove(DigitalAxesSources[glfw.AxisLeftTrigger])
		// }
		// if state.Axes[5] > TriggerDeadzoneRightLower {
		// 	// const mask = 1 << 9
		// 	// activeInput = activeInput || axesState&mask == 0
		// 	// axesState |= mask
		// 	m.activeTrigger.Add(DigitalAxesSources[glfw.AxisRightTrigger])
		// } else {
		// 	m.activeTrigger.Remove(DigitalAxesSources[glfw.AxisRightTrigger])
		// }
		// m.state.Gamepad.Axes = axesState
		// m.activeTrigger.GamepadAxes = axesState
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
		Action:     act,
		Active:     activeInput,
		InputState: m.state,
	}
	actionState := m.actionState[act]
	actionLifecycle := OnActivate
	if actionState != Inactive {
		actionLifecycle = WhileActive
	}
	actionLifecycle = (actionLifecycle >> 4) & 0xf

	m.actionState[act] |= Triggered

	for _, handler := range m.handlers[act] {
		handlerLifecycle := handler.Lifecycle() & 0xf
		if handlerLifecycle&actionLifecycle != 0 {
			handler.Invoke(ev)
		}
	}
}
