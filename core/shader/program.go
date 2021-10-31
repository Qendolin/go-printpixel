package shader

import (
	"fmt"
	"strings"

	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type LinkErr struct {
	Log     string
	Program uint32
}

func (lerr LinkErr) Error() string {
	return fmt.Sprintf("Failed to link shaders to program (id: %v). Info: \n\n%v\n\n", lerr.Program, lerr.Log)
}

type Program struct {
	*uint32
}

func NewProgram(vertShader *Shader, fragShader *Shader) (prog *Program, err error) {
	id := gl.CreateProgram()
	gl.AttachShader(id, vertShader.Id())
	gl.AttachShader(id, fragShader.Id())
	gl.LinkProgram(id)
	gl.DetachShader(id, vertShader.Id())
	gl.DetachShader(id, fragShader.Id())

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

func NewProgramFromPaths(vertex, fragment string) (prog *Program, err error) {
	vsh, err := NewShaderFromPath(vertex, TypeVertex)
	if err != nil {
		return nil, err
	}
	fsh, err := NewShaderFromPath(fragment, TypeFragment)
	if err != nil {
		return nil, err
	}
	prog, err = NewProgram(vsh, fsh)
	vsh.Destroy()
	fsh.Destroy()
	return prog, err
}

func MustNewProgramFromPaths(vertex, fragment string) (prog *Program) {
	vsh, err := NewShaderFromPath(vertex, TypeVertex)
	if err != nil {
		panic(err)
	}
	fsh, err := NewShaderFromPath(fragment, TypeFragment)
	if err != nil {
		panic(err)
	}
	prog, err = NewProgram(vsh, fsh)
	if err != nil {
		panic(err)
	}
	vsh.Destroy()
	fsh.Destroy()
	return prog
}

func (prog *Program) Validate() (ok bool, log string) {
	gl.ValidateProgram(*prog.uint32)

	log = readProgramInfoLog(*prog.uint32)

	var okInt int32
	gl.GetProgramiv(*prog.uint32, gl.VALIDATE_STATUS, &okInt)
	ok = okInt == gl.TRUE

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
	context()
	prog.Unbind()
}

func (prog *Program) Destroy() {
	gl.DeleteProgram(prog.Id())
	*prog.uint32 = 0
}

func readProgramInfoLog(id uint32) string {
	var logLength int32
	gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &logLength)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(id, logLength, nil, gl.Str(log))
	return log
}

func (prog *Program) GetUniform(name string) (*StatelessUniform, error) {
	return NewUniform(*prog, name)
}

func (prog *Program) MustGetUniform(name string) *StatelessUniform {
	u, e := prog.GetUniform(name)
	if e != nil {
		panic(e)
	}
	return u
}
