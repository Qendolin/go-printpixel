package context

import (
	"fmt"
	"io/ioutil"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type Logger interface {
	Print(v ...interface{})
	Fatal(v ...interface{})
}

type DebugMessageCallbackFactory func(Logger) gl.DebugProc

type glConfig struct {
	//Enables DEBUG_OUTPUT and DEBUG_OUTPUT_SYNCHRONOUS. Also sets DebugMessageCallback.
	Debug  bool
	Logger Logger
	DMC    DebugMessageCallbackFactory
}

func NewGlConfig() glConfig {
	return glConfig{
		Debug:  false,
		Logger: log.New(ioutil.Discard, "", 0),
		DMC:    defaultDebugMessageCallback,
	}
}

func (cfg glConfig) Apply() error {
	if cfg.Debug {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
		gl.DebugMessageCallback(cfg.DMC(cfg.Logger), nil)
	}
	return nil
}

func defaultDebugMessageCallback(log Logger) gl.DebugProc {
	return func(source uint32,
		gltype uint32,
		id uint32,
		severity uint32,
		length int32,
		message string,
		userParam unsafe.Pointer) {
		var gltypeStr string
		switch gltype {
		case gl.DEBUG_TYPE_ERROR:
			gltypeStr = "ERROR"
		case gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR:
			gltypeStr = "DEPRECATED_BEHAVIOR"
		case gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR:
			gltypeStr = "UNDEFINED_BEHAVIOR"
		case gl.DEBUG_TYPE_PORTABILITY:
			gltypeStr = "PORTABILITY"
		case gl.DEBUG_TYPE_PERFORMANCE:
			gltypeStr = "PERFORMANCE"
		case gl.DEBUG_TYPE_OTHER:
			gltypeStr = "OTHER"
		case gl.DEBUG_TYPE_MARKER:
			gltypeStr = "MARKER"
		case gl.DEBUG_TYPE_POP_GROUP:
			gltypeStr = "POP_GROUP"
		case gl.DEBUG_TYPE_PUSH_GROUP:
			gltypeStr = "PUSH_GTOUP"
		}
		var severityStr string
		fatal := false
		switch severity {
		case gl.DEBUG_SEVERITY_LOW:
			severityStr = "LIGHT ERROR"
		case gl.DEBUG_SEVERITY_MEDIUM:
			severityStr = "MEDIUM ERROR"
		case gl.DEBUG_SEVERITY_HIGH:
			severityStr = "HEAVY ERROR"
			fatal = true
		case gl.DEBUG_SEVERITY_NOTIFICATION:
			severityStr = "WARNING"
		}

		logStr := fmt.Sprintf("[%-12s] %v/%v: %v\n", severityStr, id, gltypeStr, message)

		if fatal {
			log.Fatal(logStr)
		} else {
			log.Print(logStr)
		}
	}
}
