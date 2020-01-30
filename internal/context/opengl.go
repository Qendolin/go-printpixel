package context

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

func Init(cfg glConfig) (err error) {
	cfg.Logger.Print(fmt.Sprintf("%[1]v Initializing OpenGL %[1]v", strings.Repeat("=", 5)))
	if err = gl.Init(); err != nil {
		return
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	cfg.Logger.Print(fmt.Sprintf("OpenGL version: %v", version))

	cfg.Logger.Print("Applying configuration")
	if err = cfg.Apply(); err != nil {
		return
	}

	cfg.Logger.Print(fmt.Sprintf("%[1]v Done %[1]v", strings.Repeat("=", 5)))
	return
}
