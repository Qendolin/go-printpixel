package shader

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
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

func readProgramInfoLog(id uint32) string {
	var logLength int32
	gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &logLength)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(id, logLength, nil, gl.Str(log))
	return log
}
