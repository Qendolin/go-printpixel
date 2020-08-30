package glcontext

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Error struct {
	Severity string
	Id       uint32
	Type     string
	Message  string
	Fatal    bool
	Stack    string
	Source   string
}

func (glerr Error) Error() string {
	return fmt.Sprintf("[%s] %v/%v from %v: %v", glerr.Severity, glerr.Id, glerr.Type, glerr.Source, glerr.Message)
}

type Config struct {
	//Enables DEBUG_OUTPUT and DEBUG_OUTPUT_SYNCHRONOUS. Also sets DebugMessageCallback.
	Debug         bool
	Multisampling bool
	Errors        <-chan Error
	errors        chan<- Error
}

func NewGlConfig(errorChanSize int) Config {
	errorChan := make(chan Error, errorChanSize)
	return Config{
		Debug:         false,
		Multisampling: true,
		Errors:        errorChan,
		errors:        errorChan,
	}
}

func (cfg Config) apply() error {
	if cfg.Debug {
		gl.DebugMessageCallback(debugMessageCallback(cfg.errors), nil)
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
	}
	if cfg.Multisampling {
		gl.Enable(gl.MULTISAMPLE)
	}
	return nil
}

func debugMessageCallback(errorChan chan<- Error) gl.DebugProc {
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
		var sourceStr string
		switch source {
		case gl.DEBUG_SOURCE_API:
			sourceStr = "API"
		case gl.DEBUG_SOURCE_WINDOW_SYSTEM:
			sourceStr = "WINDOW SYSTEM"
		case gl.DEBUG_SOURCE_SHADER_COMPILER:
			sourceStr = "SHADER COMPILER"
		case gl.DEBUG_SOURCE_THIRD_PARTY:
			sourceStr = "THIRD PARTY"
		case gl.DEBUG_SOURCE_APPLICATION:
			sourceStr = "APPLICATION"
		case gl.DEBUG_SOURCE_OTHER:
			sourceStr = "OTHER"
		default:
			sourceStr = "UNKNOWN"
		}

		stack := debug.Stack()

		err := Error{
			Severity: severityStr,
			Id:       id,
			Type:     gltypeStr,
			Message:  message,
			Fatal:    fatal,
			Stack:    string(stack),
			Source:   sourceStr,
		}
		select {
		case errorChan <- err:
		case <-time.After(500 * time.Millisecond):
			o := log.Writer()
			log.SetOutput(os.Stderr)
			log.Printf("Error stuck for 500ms: %v", err)
			log.SetOutput(o)
		}
	}
}
