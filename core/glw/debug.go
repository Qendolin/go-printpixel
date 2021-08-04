package glw

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime/debug"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type DebugHandler = func(err DebugMessage)

type DebugMessageType int

const (
	// Events that generated an error
	DbgError = DebugMessageType(gl.DEBUG_TYPE_ERROR)
	// Behavior that has been marked for deprecation
	DbgDeprecatedBehavoir = DebugMessageType(gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR)
	// Behavior that is undefined according to the specification
	DbgUndefinedBehavoir = DebugMessageType(gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR)
	// Implementation-dependent performance warnings
	DbgPerformace = DebugMessageType(gl.DEBUG_TYPE_PERFORMANCE)
	// Use of extensions or shaders in a way that is highly vendor-specific
	DbgPortability = DebugMessageType(gl.DEBUG_TYPE_PORTABILITY)
	// Types of events that do not fit any of the ones listed above
	DbgOther = DebugMessageType(gl.DEBUG_TYPE_OTHER)
	// Annotation of the command stream
	DbgMarker = DebugMessageType(gl.DEBUG_TYPE_MARKER)
	// Entering a debug group
	DbgPushGroup = DebugMessageType(gl.DEBUG_TYPE_PUSH_GROUP)
	// Leaving a debug group
	DbgPopGroup = DebugMessageType(gl.DEBUG_TYPE_POP_GROUP)
)

var glErrorTypes = map[DebugMessageType]string{
	DbgError:              "ERROR",
	DbgDeprecatedBehavoir: "DEPRECATED_BEHAVIOR",
	DbgUndefinedBehavoir:  "UNDEFINED_BEHAVIOR",
	DbgPerformace:         "PERFORMANCE",
	DbgPortability:        "PORTABILITY",
	DbgOther:              "OTHER",
	DbgMarker:             "MARKER",
	DbgPushGroup:          "PUSH_GROUP",
	DbgPopGroup:           "POP_GROUP",
}

func (t DebugMessageType) String() string {
	if name, ok := glErrorTypes[t]; ok {
		return name
	}
	return "UNKNOWN"
}

type DebugMessageSeverity int

const (
	// Any GL error;
	// dangerous undefined behavior;
	// any GLSL or ARB shader compiler and linker errors
	DbgSeverityHigh = DebugMessageSeverity(gl.DEBUG_SEVERITY_HIGH)
	// Severe performance warnings;
	// GLSL or other shader compiler and linker warnings;
	// use of currently deprecatedbehavior
	DbgSeverityMedium = DebugMessageSeverity(gl.DEBUG_SEVERITY_MEDIUM)
	// Performance warnings from redundant state changes;
	// trivial undefined behavior
	DbgSeverityLow = DebugMessageSeverity(gl.DEBUG_SEVERITY_LOW)
	// Any message which is not an error or performance concern
	DbgSeverityNotification = DebugMessageSeverity(gl.DEBUG_SEVERITY_NOTIFICATION)
)

var glErrorSeverities = map[DebugMessageSeverity]string{
	DbgSeverityHigh:         "CRITICAL_ERROR",
	DbgSeverityMedium:       "ERROR",
	DbgSeverityLow:          "WARNING",
	DbgSeverityNotification: "INFO",
}

func (s DebugMessageSeverity) String() string {
	if name, ok := glErrorSeverities[s]; ok {
		return name
	}
	return "UNKNOWN"
}

type DebugMessageSource int

const (
	// The GL
	DbgSourceApi = DebugMessageSource(gl.DEBUG_SOURCE_API)
	// The GLSL shader compiler or compilers for other extension-provided languages
	DbgSourceShaderCompiler = DebugMessageSource(gl.DEBUG_SOURCE_SHADER_COMPILER)
	// The window system, such as WGL or GLX
	DbgSourceWindowSystem = DebugMessageSource(gl.DEBUG_SOURCE_WINDOW_SYSTEM)
	// External debuggers or third-party middleware libraries
	DbgSourceThirdParty = DebugMessageSource(gl.DEBUG_SOURCE_THIRD_PARTY)
	// The application
	DbgSourceApplication = DebugMessageSource(gl.DEBUG_SOURCE_APPLICATION)
	// Sources that do not fit to any of the ones listed above
	DbgSourceOther = DebugMessageSource(gl.DEBUG_SOURCE_OTHER)
)

var glErrorSources = map[DebugMessageSource]string{
	DbgSourceApi:            "GRAPHICS_LIBRARY",
	DbgSourceShaderCompiler: "SHADER_COMPILER",
	DbgSourceWindowSystem:   "WINDOW_SYSTEM",
	DbgSourceThirdParty:     "THIRD_PARTY",
	DbgSourceApplication:    "APPLICATION",
	DbgSourceOther:          "OTHER",
}

func (s DebugMessageSource) String() string {
	if name, ok := glErrorSources[s]; ok {
		return name
	}
	return "UNKNOWN"
}

type DebugMessage struct {
	Severity DebugMessageSeverity
	Id       uint32
	Type     DebugMessageType
	Message  string
	Critical bool
	Stack    string
	Source   DebugMessageSource
}

func (err DebugMessage) String() string {
	return fmt.Sprintf("[%v] %v #%v from %v: %v", err.Severity, err.Type, err.Source, err.Id, err.Message)
}

func (err DebugMessage) Error() string {
	return err.String()
}

func LogCriticalf(format string, a ...interface{}) {
	insert(gl.DEBUG_SEVERITY_HIGH, fmt.Sprintf(format, a...))
}
func LogCritical(a ...interface{}) {
	insert(gl.DEBUG_SEVERITY_HIGH, fmt.Sprint(a...))
}

func LogErrorf(format string, a ...interface{}) {
	insert(gl.DEBUG_SEVERITY_MEDIUM, fmt.Sprintf(format, a...))
}
func LogError(a ...interface{}) {
	insert(gl.DEBUG_SEVERITY_MEDIUM, fmt.Sprint(a...))
}

func LogWarnf(format string, a ...interface{}) {
	insert(gl.DEBUG_SEVERITY_LOW, fmt.Sprintf(format, a...))
}
func LogWarn(a ...interface{}) {
	insert(gl.DEBUG_SEVERITY_LOW, fmt.Sprint(a...))
}

func LogInfof(format string, a ...interface{}) {
	insert(gl.DEBUG_SEVERITY_NOTIFICATION, fmt.Sprintf(format, a...))
}
func LogInfo(a ...interface{}) {
	insert(gl.DEBUG_SEVERITY_NOTIFICATION, fmt.Sprint(a...))
}

func insert(severity uint32, str string) {
	header := (*reflect.StringHeader)(unsafe.Pointer(&str))
	ptr := (*uint8)(unsafe.Pointer(header.Data))
	gl.DebugMessageInsert(gl.DEBUG_SOURCE_APPLICATION, gl.DEBUG_TYPE_ERROR, 0, severity, int32(len(str)), ptr)
}

// The ErrorTimeout is used when the error callback does not return in time
// 0 will disable the timeout
var ErrorTimeout = 100 * time.Millisecond

func DefaultDebugMessageCallback(win *extWindow) gl.DebugProc {
	return func(source uint32, gltype uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
		err := DebugMessage{
			Severity: DebugMessageSeverity(severity),
			Id:       id,
			Type:     DebugMessageType(gltype),
			Message:  message,
			Critical: severity == gl.DEBUG_SEVERITY_HIGH,
			Stack:    string(debug.Stack()),
			Source:   DebugMessageSource(source),
		}

		if ErrorTimeout != 0 {
			timeout := time.AfterFunc(ErrorTimeout, func() {
				logger := log.New(os.Stderr, log.Prefix(), log.Flags())
				logger.Printf("Error stuck for %vms: %v", ErrorTimeout.Milliseconds(), err)
			})
			win.debugHandler(err)
			timeout.Stop()
		} else {
			win.debugHandler(err)
		}
	}
}

// Prints error using default logger and calls log.Fatal on critical errors
func DefaultDebugHandler(err DebugMessage) {
	if err.Critical {
		log.Fatalf("%v\n%v", err, err.Stack)
	}
	log.Printf("%v\n", err)
}

func (w *extWindow) SetDebugHandler(cb DebugHandler) (previous DebugHandler) {
	p := w.debugHandler
	w.debugHandler = cb
	return p
}
