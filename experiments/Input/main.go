package main

import (
	"fmt"
	"runtime"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/core/shader"
	"github.com/Qendolin/go-printpixel/experiments/Input/input"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	runtime.LockOSThread()
	glwConf := glw.BasicConfig("Input Test", 1600/2, 900/2, glw.DontCare, glw.DontCare)
	glwConf.DebugContext = true
	win, err := glw.New(glwConf)
	panicIf(err)

	input.Default.Bind(win.GetGLFWWindow())
	input.Default.AddTriggers(map[input.Trigger]input.Action{
		input.Combo("KeyLeftControl", "KeySpace"):    "print_hello",
		input.Combo("KeyLeftControl", "KeyK"):        "switch_mode",
		input.Single("MouseWheelYUp"):                "inventory_prev",
		input.Single("MouseWheelYDown"):              "inventory_next",
		input.Single("MouseButtonLeft"):              "draw",
		input.Single("KeyEqual"):                     "brush_size+",
		input.Single("KeyMinus"):                     "brush_size-",
		input.Combo("KeyLeftControl", "KeyUp"):       "print_fail",
		input.Single("KeyKPDecimal"):                 "kp_decimal",
		input.Single("KeyKPEnter"):                   "kp_enter",
		input.Single("KeyKP0"):                       "kp_0",
		input.Single("KeyKPDivide"):                  "kp_divide",
		input.Single("KeyKPEqual"):                   "kp_equal",
		input.Single("KeyKPAdd"):                     "kp_add",
		input.Single("KeyKPMultiply"):                "kp_multiply",
		input.Combo("KeyLeftControl", "MouseWheelY"): "zoom",
	})

	altMode := map[input.Trigger]input.Action{
		input.Combo("KeyLeftControl", "KeyX"): "win_exit",
		input.Combo("KeyLeftControl", "KeyQ"): "win_exit",
		input.Combo("KeyLeftControl", "KeyF"): "win_fullscreen",
	}

	var (
		isFullscreen                   = false
		drawing                        = false
		brushSize                      = float32(4)
		canvasTexture                  = core.MustNewTexture2D(core.InitEmpty(1920, 1080, 2), data.RGBA)
		canvasQuad                     = core.Quad()
		canvasProgram                  = shader.MustNewProgramFromPaths("@mod/assets/shaders/quad_tex_transform.vert", "@mod/assets/shaders/quad_tex_transform.frag")
		canvasTransformUniform         = canvasProgram.MustGetUniform("u_transform")
		brushPattern, brushPatternSize = generateBrushPattern(brushSize, true)
	)

	input.Default.While("inventory_prev", func(ae input.ActionEvent) {
		fmt.Println("Next inventory slot")
	})
	input.Default.While("inventory_next", func(ae input.ActionEvent) {
		fmt.Println("Previous inventory slot")
	})

	input.Default.On("draw", func(ae input.ActionEvent) {
		fmt.Println("Draw start")
		drawing = true
	})

	input.Default.OnDeactivate("draw", func(ae input.ActionEvent) {
		fmt.Println("Draw end")
		drawing = false
	})

	input.Default.On("brush_size+", func(ae input.ActionEvent) {
		brushSize++
		brushPattern, brushPatternSize = generateBrushPattern(brushSize, true)
		fmt.Printf("Brush Size: %v\n", brushSize)
	})
	input.Default.On("brush_size-", func(ae input.ActionEvent) {
		brushSize--
		if brushSize < 1 {
			brushSize = 1
		}
		brushPattern, brushPatternSize = generateBrushPattern(brushSize, true)
		fmt.Printf("Brush Size: %v\n", brushSize)
	})

	failReg := input.Default.On("print_fail", func(ae input.ActionEvent) {
		fmt.Println("fail")
	})
	input.Default.Off(failReg)

	input.Default.On("print_hello", func(ae input.ActionEvent) {
		fmt.Println("Hello World!")
	})
	input.Default.On("switch_mode", func(ae input.ActionEvent) {
		input.Default.SetOverride(altMode, false)
	})
	input.Default.On("win_exit", func(ae input.ActionEvent) {
		win.Close()
	})
	input.Default.On("win_fullscreen", func(ae input.ActionEvent) {
		if isFullscreen {
			win.SetMonitor(nil, 100, 100, 1600/2, 900/2, 0)
		} else {
			win.SetFullscreen(glfw.GetPrimaryMonitor())
		}
		isFullscreen = !isFullscreen
	})
	input.Default.While("zoom", func(ae input.ActionEvent) {
		fmt.Printf("Zoom: %d\n", ae.Mouse.WheelY)
	})

	for tr, act := range input.Default.Triggers() {
		fmt.Printf("%-10s : %s\n", tr, act)
	}

	// glfw.UpdateGamepadMappings("03000000a30600001af5000000000000,Saitek Cyborg,a:b1,b:b2,x:b0,y:b3,back:b8,guide:b12,start:b9,leftstick:b10,rightstick:b11,leftshoulder:b4,rightshoulder:b5,dpup:h0.1,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,leftx:a0,lefty:a1,rightx:a3,righty:a4,lefttrigger:+a2,righttrigger:-a2~,platform:Windows,")
	fmt.Printf("Joystick1: %s | %s | %s\n", glfw.Joystick1.GetName(), glfw.Joystick1.GetGamepadName(), glfw.Joystick1.GetGUID())

	canvasQuad.Bind()
	canvasProgram.Bind()
	canvasTransformUniform.Set(mgl32.Ortho2D(-0.5, 0.5, -0.5, 0.5))
	canvasTexture.Bind(0)
	for !win.ShouldClose() {
		input.Default.Update()

		if drawing {
			mouse := input.Default.State().Mouse
			if mouse.PosX >= 0 && mouse.PosY >= 0 && mouse.PosX < float32(win.GetWidth()) && mouse.PosY < float32(win.GetHeight()) {
				cx := mouse.PosX * 1920 / float32(win.GetWidth())
				cy := mouse.PosY * 1080 / float32(win.GetHeight())

				canvasTexture.WriteBytes(0, int32(cx-float32(brushPatternSize)/2), int32(cy-float32(brushPatternSize)/2), int32(brushPatternSize), int32(brushPatternSize), data.RGBA.PixelFormatEnum(), brushPattern)
			}
		}

		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		if glfw.Joystick1.Present() && glfw.Joystick1.IsGamepad() {
			glfw.Joystick1.GetGamepadState()
			fmt.Printf("%v\n", glfw.Joystick1.GetGamepadState().Axes)
		}

		win.SwapBuffers()
	}
}

func generateBrushPattern(radius float32, aa bool) ([]byte, int) {
	size := int(radius+0.5) * 2
	p := make([]byte, 0, size*size*4)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			p = append(p, 255, 255, 255, 255)
		}
	}

	return p, size
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
