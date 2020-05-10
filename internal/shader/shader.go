package shader

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type ShaderType int

const (
	TypeVertex   = ShaderType(gl.VERTEX_SHADER)
	TypeFragment = ShaderType(gl.FRAGMENT_SHADER)
)

type CompileErr struct {
	Log    string
	Shader uint32
}

func (cerr CompileErr) Error() string {
	return fmt.Sprintf("Failed to compile shader (id: %v). Compiler Log: \r\n\r\n%v\r\n\r\n", cerr.Shader, cerr.Log)
}

type Shader struct {
	*uint32
}

func NewShader(source string, shaderType ShaderType) (*Shader, error) {
	id := gl.CreateShader(uint32(shaderType))
	err := loadAndCompileShader(id, source)
	return &Shader{&id}, err
}

func NewVertexShader(source string) (*Shader, error) {
	return NewShader(source, TypeVertex)
}

func NewFragmentShader(source string) (*Shader, error) {
	return NewShader(source, TypeFragment)
}

func NewShaderFromModulePath(modulePath string, shaderType ShaderType) (*Shader, error) {
	absPath, err := utils.ResolveModulePath(modulePath)
	if err != nil {
		return nil, err
	}
	source, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	return NewShader(string(source), shaderType)
}

func (shader *Shader) Id() uint32 {
	return *shader.uint32
}

func (shader *Shader) Destroy() {
	gl.DeleteShader(shader.Id())
}

func loadAndCompileShader(id uint32, source string) error {
	cStrs, free := gl.Strs(source + "\x00")
	gl.ShaderSource(id, 1, cStrs, nil)
	free()
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
