package glcontext

import (
	"fmt"
	"runtime/debug"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type GlError struct {
	Severity string
	Id       uint32
	Type     string
	Message  string
	Fatal    bool
	Stack    string
}

func (glerr GlError) Error() string {
	return fmt.Sprintf("[%s] %v/%v: %v\n", glerr.Severity, glerr.Id, glerr.Type, glerr.Message)
}

type GlConfig struct {
	//Enables DEBUG_OUTPUT and DEBUG_OUTPUT_SYNCHRONOUS. Also sets DebugMessageCallback.
	Debug  bool
	Errors <-chan GlError
	errors chan<- GlError
}

func NewGlConfig(errorChanSize int) GlConfig {
	errorChan := make(chan GlError, errorChanSize)
	return GlConfig{
		Debug:  false,
		Errors: errorChan,
		errors: errorChan,
	}
}

func (cfg GlConfig) apply() error {
	if cfg.Debug {
		gl.DebugMessageCallback(debugMessageCallback(cfg.errors), nil)
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
	}
	return nil
}

func debugMessageCallback(errorChan chan<- GlError) gl.DebugProc {
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

		stack := debug.Stack()

		errorChan <- GlError{
			Severity: severityStr,
			Id:       id,
			Type:     gltypeStr,
			Message:  message,
			Fatal:    fatal,
			Stack:    string(stack),
		}
	}
}