package shader

import (
	"fmt"
	"strings"

	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type LinkErr struct {
	Log     string
	Program uint32
}

func (lerr LinkErr) Error() string {
	return fmt.Sprintf("Failed to link shaders to program (id: %v). Info: \n\n%v", lerr.Program, lerr.Log)
}

type Program struct {
	*uint32
}

func NewProgram(vertShader *Shader, fragShader *Shader) (prog *Program, err error) {
	id := gl.CreateProgram()
	gl.AttachShader(id, vertShader.Id())
	gl.AttachShader(id, fragShader.Id())
	gl.LinkProgram(id)

	var ok int32
	gl.GetProgramiv(id, gl.LINK_STATUS, &ok)
	if ok == gl.FALSE {
		err = LinkErr{
			Log:     readProgramInfoLog(id),
			Program: id,
		}
	}
	prog = &Program{&id}
	return
}

func (prog *Program) Id() uint32 {
	return *prog.uint32
}

func (prog *Program) Bind() {
	gl.UseProgram(prog.Id())
}

func (prog *Program) Unbind() {
	gl.UseProgram(0)
}

func (prog *Program) BindFor(context utils.BindingClosure) {
	prog.Bind()
	defered := context()
	prog.Unbind()
	for _, deferedFunc := range defered {
		deferedFunc()
	}
}

func (prog *Program) Destroy() {
	gl.DeleteProgram(prog.Id())
}

func readProgramInfoLog(id uint32) string {
	var logLength int32
	gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &logLength)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(id, logLength, nil, gl.Str(log))
	return log
}
