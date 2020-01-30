package shader

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type CompileErr struct {
	Log    string
	Shader uint32
}

func (cerr CompileErr) Error() string {
	return fmt.Sprintf("Failed to compile shader (id: %v). Info: \n\n%v", cerr.Shader, cerr.Log)
}

type Shader struct {
	*uint32
}

func NewVertexShader(source string) (*Shader, error) {
	id := gl.CreateShader(gl.VERTEX_SHADER)
	err := loadAndCompileShader(id, source)
	return &Shader{&id}, err
}

func NewFragmentShader(source string) (*Shader, error) {
	id := gl.CreateShader(gl.VERTEX_SHADER)
	err := loadAndCompileShader(id, source)
	return &Shader{&id}, err
}

func (shader *Shader) Id() uint32 {
	return *shader.uint32
}

func (shader *Shader) Destroy() {
	gl.DeleteShader(shader.Id())
}

func loadAndCompileShader(id uint32, source string) error {
	cStr := gl.Str(source + "\x00")
	gl.ShaderSource(id, 1, &cStr, nil)
	gl.CompileShader(id)

	var ok int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &ok)
	if ok == gl.FALSE {
		return CompileErr{
			Log:    readShaderInfoLog(id),
			Shader: id,
		}
	}
	return nil
}

func readShaderInfoLog(id uint32) string {
	var logLength int32
	gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &logLength)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(id, logLength, nil, gl.Str(log))
	return log
}
