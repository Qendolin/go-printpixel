package input

import (
	"strings"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type SourceType int

const (
	Unknown       = SourceType(0)
	KeyboardKey   = SourceType(1)
	ModifierKey   = SourceType(2)
	GamepadButton = SourceType(3)
	GamepadAxis   = SourceType(4)
	MouseButton   = SourceType(5)
	MouseWheel    = SourceType(6)
)

type SourceId uint8

type Source struct {
	Id          SourceId
	GlfwId      int
	Type        SourceType
	Name        string
	PrettyName  string
	Description string
	Analog      bool
}

func (s Source) String() string {
	if s.Type != KeyboardKey {
		return s.PrettyName
	}
	name := glfw.GetKeyName(glfw.Key(s.GlfwId), 0)
	// TODO: Spacial cases for keypad keys and so on
	if name == "" {
		return s.PrettyName
	}

	if s.GlfwId >= int(glfw.KeyKP0) && s.GlfwId <= int(glfw.KeyKPEqual) {
		name = "KP‑" + name
	}

	return name
}

func init() {
	NamedSources = make(map[string]Source, len(Sources))
	NamedKeyboardKeys = make(map[string]Source, 120)
	KeyboardKeySources = make(map[glfw.Key]Source, 120)
	NamedModifierKeys = make(map[string]Source, 6)
	ModifierKeySources = make(map[glfw.ModifierKey]Source, 6)
	NamedGamepadButtons = make(map[string]Source, 19)
	NamedGamepadAxes = make(map[string]Source, 18)
	NamedMouseButtons = make(map[string]Source, 13)
	NamedMouseWheels = make(map[string]Source, 6)
	for i := 0; i < len(Sources); i++ {
		if i > 255 {
			panic("Too many sources")
		}
		s := Sources[i]
		s.Id = SourceId(i)
		Sources[i] = s
		name := strings.ToLower(s.Name)
		NamedSources[name] = Sources[i]
		switch s.Type {
		case KeyboardKey:
			NamedKeyboardKeys[name] = s
			KeyboardKeySources[glfw.Key(s.GlfwId)] = s
		case ModifierKey:
			NamedModifierKeys[name] = s
			ModifierKeySources[glfw.ModifierKey(s.GlfwId)] = s
		case GamepadButton:
			NamedGamepadButtons[name] = s
			GamepadButtonSources[s.GlfwId] = s
		case GamepadAxis:
			NamedGamepadAxes[name] = s
			if s.Analog {
				AnalogAxesSources[s.GlfwId] = s
			} else {
				if DigitalAxesSources[s.GlfwId*2].GlfwId != s.GlfwId {
					DigitalAxesSources[s.GlfwId*2] = s
				} else {
					DigitalAxesSources[s.GlfwId*2+1] = s
				}
			}
		case MouseButton:
			NamedMouseButtons[name] = s
			MouseButtonSources[s.GlfwId] = s
		case MouseWheel:
			NamedMouseWheels[name] = s
		}
	}
}

// Sources by lowercase name
var NamedSources map[string]Source

// Keyboard key sources by lowercase name
var NamedKeyboardKeys map[string]Source

// Keyboard key sources by glfw key
var KeyboardKeySources map[glfw.Key]Source

// Modifier key sources by lowercase name
var NamedModifierKeys map[string]Source

// Modifier key sources by glfw modifier key
var ModifierKeySources map[glfw.ModifierKey]Source

// Mouse buttons sources by lowercase name
var NamedMouseButtons map[string]Source

// Keyboard key sources by lowercase name
var MouseButtonSources [glfw.MouseButtonLast + 1]Source

// Mouse wheel sources by lowercase name
var NamedMouseWheels map[string]Source

// Gamepad buttons sources by lowercase name
var NamedGamepadButtons map[string]Source

// Gamepad buttons sources by index
var GamepadButtonSources [glfw.ButtonLast + 1]Source

// Gamepad axes sources by lowercase name
var NamedGamepadAxes map[string]Source

// Gamepad buttons sources by index
var AnalogAxesSources [glfw.AxisLast + 1]Source

// Gamepad buttons sources by index
var DigitalAxesSources [(glfw.AxisLast + 1) * 2]Source

var Sources = [...]Source{
	{GlfwId: -1, Name: "Unknown", Description: "Unknown", PrettyName: "Unknown", Type: Unknown},
	// Keyboard Keys
	{GlfwId: int(glfw.KeySpace), Name: "KeySpace", Description: "Key Space", PrettyName: "⎵", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyApostrophe), Name: "KeyApostrophe", Description: "Key Apostrophe", PrettyName: "'", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyComma), Name: "KeyComma", Description: "Key Comma", PrettyName: ",", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyMinus), Name: "KeyMinus", Description: "Key Minus", PrettyName: "-", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyPeriod), Name: "KeyPeriod", Description: "Key Period", PrettyName: ".", Type: KeyboardKey},
	{GlfwId: int(glfw.KeySlash), Name: "KeySlash", Description: "Key Slash", PrettyName: "/", Type: KeyboardKey},
	{GlfwId: int(glfw.Key0), Name: "Key0", Description: "Key 0", PrettyName: "0", Type: KeyboardKey},
	{GlfwId: int(glfw.Key1), Name: "Key1", Description: "Key 1", PrettyName: "1", Type: KeyboardKey},
	{GlfwId: int(glfw.Key2), Name: "Key2", Description: "Key 2", PrettyName: "2", Type: KeyboardKey},
	{GlfwId: int(glfw.Key3), Name: "Key3", Description: "Key 3", PrettyName: "3", Type: KeyboardKey},
	{GlfwId: int(glfw.Key4), Name: "Key4", Description: "Key 4", PrettyName: "4", Type: KeyboardKey},
	{GlfwId: int(glfw.Key5), Name: "Key5", Description: "Key 5", PrettyName: "5", Type: KeyboardKey},
	{GlfwId: int(glfw.Key6), Name: "Key6", Description: "Key 6", PrettyName: "6", Type: KeyboardKey},
	{GlfwId: int(glfw.Key7), Name: "Key7", Description: "Key 7", PrettyName: "7", Type: KeyboardKey},
	{GlfwId: int(glfw.Key8), Name: "Key8", Description: "Key 8", PrettyName: "8", Type: KeyboardKey},
	{GlfwId: int(glfw.Key9), Name: "Key9", Description: "Key 9", PrettyName: "9", Type: KeyboardKey},
	{GlfwId: int(glfw.KeySemicolon), Name: "KeySemicolon", Description: "Key Semicolon", PrettyName: ";", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyEqual), Name: "KeyEqual", Description: "Key Equal", PrettyName: "=", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyA), Name: "KeyA", Description: "Key A", PrettyName: "A", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyB), Name: "KeyB", Description: "Key B", PrettyName: "B", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyC), Name: "KeyC", Description: "Key C", PrettyName: "C", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyD), Name: "KeyD", Description: "Key D", PrettyName: "D", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyE), Name: "KeyE", Description: "Key E", PrettyName: "E", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF), Name: "KeyF", Description: "Key F", PrettyName: "F", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyG), Name: "KeyG", Description: "Key G", PrettyName: "G", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyH), Name: "KeyH", Description: "Key H", PrettyName: "H", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyI), Name: "KeyI", Description: "Key I", PrettyName: "I", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyJ), Name: "KeyJ", Description: "Key J", PrettyName: "J", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyK), Name: "KeyK", Description: "Key K", PrettyName: "K", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyL), Name: "KeyL", Description: "Key L", PrettyName: "L", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyM), Name: "KeyM", Description: "Key M", PrettyName: "M", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyN), Name: "KeyN", Description: "Key N", PrettyName: "N", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyO), Name: "KeyO", Description: "Key O", PrettyName: "O", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyP), Name: "KeyP", Description: "Key P", PrettyName: "P", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyQ), Name: "KeyQ", Description: "Key Q", PrettyName: "Q", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyR), Name: "KeyR", Description: "Key R", PrettyName: "R", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyS), Name: "KeyS", Description: "Key S", PrettyName: "S", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyT), Name: "KeyT", Description: "Key T", PrettyName: "T", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyU), Name: "KeyU", Description: "Key U", PrettyName: "U", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyV), Name: "KeyV", Description: "Key V", PrettyName: "V", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyW), Name: "KeyW", Description: "Key W", PrettyName: "W", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyX), Name: "KeyX", Description: "Key X", PrettyName: "X", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyY), Name: "KeyY", Description: "Key Y", PrettyName: "Y", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyZ), Name: "KeyZ", Description: "Key Z", PrettyName: "Z", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyLeftBracket), Name: "KeyLeftBracket", Description: "Key Left Bracket", PrettyName: "[", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyBackslash), Name: "KeyBackslash", Description: "Key Backslash", PrettyName: "\\", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyRightBracket), Name: "KeyRightBracket", Description: "Key Right Bracket", PrettyName: "]", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyGraveAccent), Name: "KeyGraveAccent", Description: "Key Grave Accent", PrettyName: "`", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyWorld1 + 0), Name: "KeyWorld1", Description: "Key World 1", PrettyName: "Key World 1", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyWorld1 + 1), Name: "KeyWorld2", Description: "Key World 2", PrettyName: "Key World 2", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyWorld1 + 2), Name: "KeyWorld3", Description: "Key World 3", PrettyName: "Key World 3", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyWorld1 + 3), Name: "KeyWorld4", Description: "Key World 4", PrettyName: "Key World 4", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyWorld1 + 4), Name: "KeyWorld5", Description: "Key World 5", PrettyName: "Key World 5", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyWorld1 + 5), Name: "KeyWorld6", Description: "Key World 6", PrettyName: "Key World 6", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyWorld1 + 6), Name: "KeyWorld7", Description: "Key World 7", PrettyName: "Key World 7", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyWorld1 + 7), Name: "KeyWorld8", Description: "Key World 8", PrettyName: "Key World 8", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyWorld1 + 8), Name: "KeyWorld9", Description: "Key World 9", PrettyName: "Key World 9", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyEnter), Name: "KeyEnter", Description: "Key Enter", PrettyName: "⏎", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyTab), Name: "KeyTab", Description: "Key Tab", PrettyName: "↹", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyBackspace), Name: "KeyBackspace", Description: "Key Backspace", PrettyName: "←", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyInsert), Name: "KeyInsert", Description: "Key Insert", PrettyName: "Ins", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyDelete), Name: "KeyDelete", Description: "Key Delete", PrettyName: "Del", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyRight), Name: "KeyRight", Description: "Key Right", PrettyName: "→", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyLeft), Name: "KeyLeft", Description: "Key Left", PrettyName: "←", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyDown), Name: "KeyDown", Description: "Key Down", PrettyName: "↓", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyUp), Name: "KeyUp", Description: "Key Up", PrettyName: "↑", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyPageUp), Name: "KeyPageUp", Description: "Key Page Up", PrettyName: "PgUp", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyPageDown), Name: "KeyPageDown", Description: "Key Page Down", PrettyName: "PgDn", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyHome), Name: "KeyHome", Description: "Key Home", PrettyName: "Home", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyEnd), Name: "KeyEnd", Description: "Key End", PrettyName: "End", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyCapsLock), Name: "KeyCapsLock", Description: "Key Caps Lock", PrettyName: "Caps‑⇪", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyScrollLock), Name: "KeyScrollLock", Description: "Key Scroll Lock", PrettyName: "Scroll‑⇩", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyNumLock), Name: "KeyNumLock", Description: "Key Num Lock", PrettyName: "Num‑⇩", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyPrintScreen), Name: "KeyPrintScreen", Description: "Key Print Screen", PrettyName: "PrtSc", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyPause), Name: "KeyPause", Description: "Key Pause", PrettyName: "Pause", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF1), Name: "KeyF1", Description: "Key F1", PrettyName: "F1", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF2), Name: "KeyF2", Description: "Key F2", PrettyName: "F2", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF3), Name: "KeyF3", Description: "Key F3", PrettyName: "F3", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF4), Name: "KeyF4", Description: "Key F4", PrettyName: "F4", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF5), Name: "KeyF5", Description: "Key F5", PrettyName: "F5", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF6), Name: "KeyF6", Description: "Key F6", PrettyName: "F6", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF7), Name: "KeyF7", Description: "Key F7", PrettyName: "F7", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF8), Name: "KeyF8", Description: "Key F8", PrettyName: "F8", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF9), Name: "KeyF9", Description: "Key F9", PrettyName: "F9", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF10), Name: "KeyF10", Description: "Key F10", PrettyName: "F10", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF11), Name: "KeyF11", Description: "Key F11", PrettyName: "F11", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF12), Name: "KeyF12", Description: "Key F12", PrettyName: "F12", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF13), Name: "KeyF13", Description: "Key F13", PrettyName: "F13", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF14), Name: "KeyF14", Description: "Key F14", PrettyName: "F14", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF15), Name: "KeyF15", Description: "Key F15", PrettyName: "F15", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF16), Name: "KeyF16", Description: "Key F16", PrettyName: "F16", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF17), Name: "KeyF17", Description: "Key F17", PrettyName: "F17", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF18), Name: "KeyF18", Description: "Key F18", PrettyName: "F18", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF19), Name: "KeyF19", Description: "Key F19", PrettyName: "F19", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF20), Name: "KeyF20", Description: "Key F20", PrettyName: "F20", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF21), Name: "KeyF21", Description: "Key F21", PrettyName: "F21", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF22), Name: "KeyF22", Description: "Key F22", PrettyName: "F22", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF23), Name: "KeyF23", Description: "Key F23", PrettyName: "F23", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF24), Name: "KeyF24", Description: "Key F24", PrettyName: "F24", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyF25), Name: "KeyF25", Description: "Key F25", PrettyName: "F25", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP0), Name: "KeyKP0", Description: "Keypad 0", PrettyName: "KP‑0", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP1), Name: "KeyKP1", Description: "Keypad 1", PrettyName: "KP‑1", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP2), Name: "KeyKP2", Description: "Keypad 2", PrettyName: "KP‑2", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP3), Name: "KeyKP3", Description: "Keypad 3", PrettyName: "KP‑3", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP4), Name: "KeyKP4", Description: "Keypad 4", PrettyName: "KP‑4", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP5), Name: "KeyKP5", Description: "Keypad 5", PrettyName: "KP‑5", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP6), Name: "KeyKP6", Description: "Keypad 6", PrettyName: "KP‑6", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP7), Name: "KeyKP7", Description: "Keypad 7", PrettyName: "KP‑7", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP8), Name: "KeyKP8", Description: "Keypad 8", PrettyName: "KP‑8", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKP9), Name: "KeyKP9", Description: "Keypad 9", PrettyName: "KP‑9", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKPDecimal), Name: "KeyKPDecimal", Description: "Keypad Decimal", PrettyName: "KP‑.", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKPDivide), Name: "KeyKPDivide", Description: "Keypad Divide", PrettyName: "KP‑/", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKPMultiply), Name: "KeyKPMultiply", Description: "Keypad Multiply", PrettyName: "KP‑*", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKPSubtract), Name: "KeyKPSubtract", Description: "Keypad Subtract", PrettyName: "KP‑-", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKPAdd), Name: "KeyKPAdd", Description: "Keypad Add", PrettyName: "KP‑+", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKPEnter), Name: "KeyKPEnter", Description: "Keypad Enter", PrettyName: "KP‑Enter", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyKPEqual), Name: "KeyKPEqual", Description: "Keypad Equal", PrettyName: "KP‑=", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyLeftShift), Name: "KeyLeftShift", Description: "Key Left Shift", PrettyName: "L‑⇧", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyLeftControl), Name: "KeyLeftControl", Description: "Key Left Control", PrettyName: "L‑Ctrl", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyLeftAlt), Name: "KeyLeftAlt", Description: "Key Left Alt", PrettyName: "L‑Alt/⌥", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyLeftSuper), Name: "KeyLeftSuper", Description: "Key Left Super", PrettyName: "L‑⊞/⌘", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyRightShift), Name: "KeyRightShift", Description: "Key Right Shift", PrettyName: "R‑⇧", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyRightControl), Name: "KeyRightControl", Description: "Key Right Control", PrettyName: "R‑Ctrl", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyRightAlt), Name: "KeyRightAlt", Description: "Key Right Alt", PrettyName: "R‑Alt/⌥", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyRightSuper), Name: "KeyRightSuper", Description: "Key Right Super", PrettyName: "R‑⊞/⌘", Type: KeyboardKey},
	{GlfwId: int(glfw.KeyMenu), Name: "KeyMenu", Description: "Key Menu", PrettyName: "≣", Type: KeyboardKey},
	// Mods
	{GlfwId: int(glfw.ModShift), Name: "ModShift", Description: "Modifier Shift", PrettyName: "⇧", Type: ModifierKey},
	{GlfwId: int(glfw.ModControl), Name: "ModControl", Description: "Modifier Control", PrettyName: "Ctrl", Type: ModifierKey},
	{GlfwId: int(glfw.ModAlt), Name: "ModAlt", Description: "Modifier Alt", PrettyName: "Alt/⌥", Type: ModifierKey},
	{GlfwId: int(glfw.ModSuper), Name: "ModSuper", Description: "Modifier Super", PrettyName: "⊞/⌘", Type: ModifierKey},
	{GlfwId: int(glfw.ModCapsLock), Name: "ModCapsLock", Description: "Modifier Caps Lock", PrettyName: "Caps‑⇪", Type: ModifierKey},
	{GlfwId: int(glfw.ModNumLock), Name: "ModNumLock", Description: "Modifier Num Lock", PrettyName: "Num‑⇩", Type: ModifierKey},
	// Gamepad Buttons
	{GlfwId: int(glfw.ButtonA), Name: "ButtonA", Description: "Gamepad Xbox A", PrettyName: "Ⓐ/🟢", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonB), Name: "ButtonB", Description: "Gamepad Xbox B", PrettyName: "Ⓑ/🔴", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonX), Name: "ButtonX", Description: "Gamepad Xbox X", PrettyName: "Ⓧ/🔵", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonY), Name: "ButtonY", Description: "Gamepad Xbox Y", PrettyName: "Ⓨ/🟡", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonLeftBumper), Name: "ButtonLeftBumper", Description: "Gamepad Left Bumper", PrettyName: "L‑⁽🎮", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonRightBumper), Name: "ButtonRightBumper", Description: "Gamepad Right Bumper", PrettyName: "R‑🎮⁾", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonBack), Name: "ButtonBack", Description: "Gamepad Back", PrettyName: "🎮‑Back", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonStart), Name: "ButtonStart", Description: "Gamepad Start", PrettyName: "🎮‑Start", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonGuide), Name: "ButtonGuide", Description: "Gamepad Guide", PrettyName: "🎮‑Guide", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonLeftThumb), Name: "ButtonLeftThumb", Description: "Gamepad Left Stick", PrettyName: "L‑◉", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonRightThumb), Name: "ButtonRightThumb", Description: "Gamepad Right Stick", PrettyName: "R‑◉", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonDpadUp), Name: "ButtonDpadUp", Description: "D-pad Up", PrettyName: "🞧⇧", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonDpadRight), Name: "ButtonDpadRight", Description: "D-pad Right", PrettyName: "🞧⇨", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonDpadDown), Name: "ButtonDpadDown", Description: "D-pad Down", PrettyName: "🞧⇩", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonDpadLeft), Name: "ButtonDpadLeft", Description: "D-pad Left", PrettyName: "🞧⇦", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonCross), Name: "ButtonCross", Description: "Gamepad PlayStation Cross", PrettyName: "✕", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonCircle), Name: "ButtonCircle", Description: "Gamepad PlayStation Circle", PrettyName: "⭘", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonSquare), Name: "ButtonSquare", Description: "Gamepad PlayStation Square", PrettyName: "◻", Type: GamepadButton},
	{GlfwId: int(glfw.ButtonTriangle), Name: "ButtonTriangle", Description: "Gamepad PlayStation Triangle", PrettyName: "△", Type: GamepadButton},
	// Gamepad Axes Analog
	{GlfwId: int(glfw.AxisLeftX), Name: "AxisLeftX", Description: "Gamepad Left Stick X-Axis", PrettyName: "L‑◉ ⇄", Type: GamepadAxis, Analog: true},
	{GlfwId: int(glfw.AxisLeftY), Name: "AxisLeftY", Description: "Gamepad Left Stick Y-Axis", PrettyName: "L‑◉ ⇅", Type: GamepadAxis, Analog: true},
	{GlfwId: int(glfw.AxisRightX), Name: "AxisRightX", Description: "Gamepad Right Stick X-Axis", PrettyName: "R‑◉ ⇄", Type: GamepadAxis, Analog: true},
	{GlfwId: int(glfw.AxisRightY), Name: "AxisRightY", Description: "Gamepad Right Stick Y-Axis", PrettyName: "R‑◉ ⇅", Type: GamepadAxis, Analog: true},
	{GlfwId: int(glfw.AxisLeftTrigger), Name: "AxisLeftTrigger", Description: "Gamepad Left Trigger", PrettyName: "L‑Trigger ⬒", Type: GamepadAxis, Analog: true},
	{GlfwId: int(glfw.AxisRightTrigger), Name: "AxisRightTrigger", Description: "Gamepad Right Trigger", PrettyName: "R‑Trigger ⬒", Type: GamepadAxis, Analog: true},
	// Gamepad Axes Digital
	// FIXME: Left / Right might be switched
	{GlfwId: int(glfw.AxisLeftX), Name: "AxisLeftXLeft", Description: "Gamepad Left Stick Left", PrettyName: "L‑⇦◉", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisLeftX), Name: "AxisLeftXRight", Description: "Gamepad Left Stick Right", PrettyName: "L‑◉⇨", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisLeftY), Name: "AxisLeftYUp", Description: "Gamepad Left Stick Up", PrettyName: "L‑⇧◉", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisLeftY), Name: "AxisLeftYDown", Description: "Gamepad Left Stick Down", PrettyName: "L‑◉⇩", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisRightX), Name: "AxisRightXLeft", Description: "Gamepad Right Stick Left", PrettyName: "R‑⇦◉", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisRightX), Name: "AxisRightXRight", Description: "Gamepad Right Stick right", PrettyName: "R‑◉⇨", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisRightY), Name: "AxisRightYUp", Description: "Gamepad Right Stick Up", PrettyName: "R‑⇧◉", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisRightY), Name: "AxisRightYDown", Description: "Gamepad Right Stick Down", PrettyName: "R‑◉⇩", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisLeftTrigger), Name: "AxisLeftTriggerHalf", Description: "Gamepad Left Trigger Half", PrettyName: "L‑Trigger ⬒", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisLeftTrigger), Name: "AxisLeftTriggerFull", Description: "Gamepad Left Trigger Full", PrettyName: "L‑Trigger ◼", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisRightTrigger), Name: "AxisRightTriggerHalf", Description: "Gamepad Right Trigger Half", PrettyName: "R‑Trigger ⬒", Type: GamepadAxis},
	{GlfwId: int(glfw.AxisRightTrigger), Name: "AxisRightTriggerFull", Description: "Gamepad Right Trigger Full", PrettyName: "R‑Trigger ◼", Type: GamepadAxis},
	// Mouse Buttons
	{GlfwId: int(glfw.MouseButton1), Name: "MouseButton1", Description: "Mouse Button 1", PrettyName: "🖱️‑1", Type: MouseButton},
	{GlfwId: int(glfw.MouseButton2), Name: "MouseButton2", Description: "Mouse Button 2", PrettyName: "🖱️‑2", Type: MouseButton},
	{GlfwId: int(glfw.MouseButton3), Name: "MouseButton3", Description: "Mouse Button 3", PrettyName: "🖱️‑3", Type: MouseButton},
	{GlfwId: int(glfw.MouseButton4), Name: "MouseButton4", Description: "Mouse Button 4", PrettyName: "🖱️‑4", Type: MouseButton},
	{GlfwId: int(glfw.MouseButton5), Name: "MouseButton5", Description: "Mouse Button 5", PrettyName: "🖱️‑5", Type: MouseButton},
	{GlfwId: int(glfw.MouseButton6), Name: "MouseButton6", Description: "Mouse Button 6", PrettyName: "🖱️‑6", Type: MouseButton},
	{GlfwId: int(glfw.MouseButton7), Name: "MouseButton7", Description: "Mouse Button 7", PrettyName: "🖱️‑7", Type: MouseButton},
	{GlfwId: int(glfw.MouseButton8), Name: "MouseButton8", Description: "Mouse Button 8", PrettyName: "🖱️‑8", Type: MouseButton},
	{GlfwId: int(glfw.MouseButtonLeft), Name: "MouseButtonLeft", Description: "Mouse Button Left", PrettyName: "🖱️‑◖", Type: MouseButton},
	{GlfwId: int(glfw.MouseButtonRight), Name: "MouseButtonRight", Description: "Mouse Button Right", PrettyName: "🖱️‑◗", Type: MouseButton},
	{GlfwId: int(glfw.MouseButtonMiddle), Name: "MouseButtonMiddle", Description: "Mouse Button Middle", PrettyName: "🖱️‑⯊", Type: MouseButton},
	{GlfwId: int(glfw.MouseButton4), Name: "MouseButtonBack", Description: "Mouse Button Back", PrettyName: "🖱️‑⮨", Type: MouseButton},
	{GlfwId: int(glfw.MouseButton5), Name: "MouseButtonForward", Description: "Mouse Button Forward", PrettyName: "🖱️‑⮫", Type: MouseButton},
	// Mouse Wheel Analog
	{GlfwId: -1, Name: "MouseWheelX", Description: "Mouse Wheel X", PrettyName: "🖱️‑⊚⇄", Type: MouseWheel, Analog: true},
	{GlfwId: -1, Name: "MouseWheelY", Description: "Mouse Wheel Y", PrettyName: "🖱️‑⊚⇅", Type: MouseWheel, Analog: true},
	// Mouse Wheel Digital
	{GlfwId: -1, Name: "MouseWheelXLeft", Description: "Mouse Wheel Left", PrettyName: "🖱️‑⇇⊚", Type: MouseWheel},
	{GlfwId: -1, Name: "MouseWheelXRight", Description: "Mouse Wheel Right", PrettyName: "🖱️‑⊚⇉", Type: MouseWheel},
	{GlfwId: -1, Name: "MouseWheelYUp", Description: "Mouse Wheel Up", PrettyName: "🖱️‑⊚⇈", Type: MouseWheel},
	{GlfwId: -1, Name: "MouseWheelYDown", Description: "Mouse Wheel Down", PrettyName: "🖱️‑⊚⇊", Type: MouseWheel},
}
