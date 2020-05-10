package shader

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Uniform struct {
	*int32
}

type UniformLinkError struct {
	Name    string
	Program uint32
}

func (ulerr UniformLinkError) Error() string {
	return fmt.Sprintf("Failed to get the location of uniform '%v'. Program: %v", ulerr.Name, ulerr.Program)
}

func NewUniform(prog Program, name string) (uni *Uniform, err error) {
	loc := gl.GetUniformLocation(*prog.uint32, gl.Str(name))
	if loc == -1 {
		err = UniformLinkError{
			Name:    name,
			Program: *prog.uint32,
		}
	}

	uni = &Uniform{&loc}
	return
}

func (u *Uniform) Set(value interface{}) {
	for refVal := reflect.ValueOf(value); refVal.Kind() == reflect.Ptr; refVal = reflect.ValueOf(value) {
		value = refVal.Elem().Interface()
	}

	switch v := value.(type) {
	case float64:
		gl.Uniform1d(*u.int32, v)
	case float32:
		gl.Uniform1f(*u.int32, v)
	case int:
	case int64:
	case int32:
	case int16:
	case int8:
		gl.Uniform1i(*u.int32, int32(v))
	case uint:
	case uint64:
	case uint32:
	case uint16:
	case uint8:
		gl.Uniform1ui(*u.int32, uint32(v))
	case mgl32.Vec2:
		gl.Uniform2f(*u.int32, v.X(), v.Y())
	case mgl64.Vec2:
		gl.Uniform2d(*u.int32, v.X(), v.Y())
	case mgl32.Vec3:
		gl.Uniform3f(*u.int32, v.X(), v.Y(), v.Z())
	case mgl64.Vec3:
		gl.Uniform3d(*u.int32, v.X(), v.Y(), v.Z())
	case mgl32.Vec4:
		gl.Uniform4f(*u.int32, v.X(), v.Y(), v.Z(), v.W())
	case mgl64.Vec4:
		gl.Uniform4d(*u.int32, v.X(), v.Y(), v.Z(), v.W())
	case mgl32.Mat3:
		gl.UniformMatrix3fv(*u.int32, 1, false, &v[0])
	case mgl64.Mat3:
		gl.UniformMatrix3dv(*u.int32, 1, false, &v[0])
	case mgl32.Mat4:
		gl.UniformMatrix4fv(*u.int32, 1, false, &v[0])
	case mgl64.Mat4:
		gl.UniformMatrix4dv(*u.int32, 1, false, &v[0])
	default:
		reflectType := reflect.TypeOf(value)
		dataType := reflectType.String()
		log.Printf("Unsupported type %v", dataType)
	}
}
