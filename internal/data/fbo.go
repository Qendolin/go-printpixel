package data

import (
	"fmt"

	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

//The highest bit indicates this is a stencil attachment.
//The second highest bit indicates this is a depth attachment.
//If both are set this is s a depth stencil attachment
//Otherwise this is gl.COLOR_ATTACHMENTi where i is the value
type FboAttachment uint

const (
	FboAttachmentColor        = FboAttachment(0)
	FboAttachmentStencil      = FboAttachment(^(^uint(0) >> 1))
	FboAttachmentDepth        = FboAttachmentStencil >> 1
	FboAttachmentDepthStencil = FboAttachmentStencil | FboAttachmentDepth
)

func (a FboAttachment) ToGlEnum() uint32 {
	switch a & FboAttachmentDepthStencil {
	case FboAttachmentDepthStencil:
		return gl.DEPTH_STENCIL_ATTACHMENT
	case FboAttachmentStencil:
		return gl.STENCIL_ATTACHMENT
	case FboAttachmentDepth:
		return gl.DEPTH_ATTACHMENT
	default:
		return uint32(gl.COLOR_ATTACHMENT0 + a)
	}
}

type FboTarget int

const (
	FboTargetRead      = FboTarget(gl.READ_FRAMEBUFFER)
	FboTargetWrite     = FboTarget(gl.DRAW_FRAMEBUFFER)
	FboTargetReadWrite = FboTarget(gl.FRAMEBUFFER)
)

type FboStatusError struct {
	Status   string
	GlStatus string
	Target   FboTarget
	Id       uint32
}

func (err FboStatusError) Error() string {
	target := "draw"
	if err.Target == FboTargetRead {
		target = "read"
	}
	return fmt.Sprintf("Framebuffer %v bound as %v target is incomplete. Reason: %v. (%v)", err.Id, target, err.Status, err.GlStatus)
}

var DefaultFbo = Fbo{
	uint32: new(uint32),
}

type Fbo struct {
	*uint32
}

func NewFbo() *Fbo {
	return &Fbo{uint32: &NewId}
}

func (fbo *Fbo) Id() uint32 {
	if fbo.uint32 == &NewId {
		id := new(uint32)
		gl.GenFramebuffers(1, id)
		fbo.uint32 = id
	}
	return *fbo.uint32
}

func (fbo *Fbo) Bind(target FboTarget) {
	gl.BindFramebuffer(uint32(target), fbo.Id())
}

func (fbo *Fbo) Unbind(target FboTarget) {
	gl.BindFramebuffer(uint32(target), 0)
}

func (fbo *Fbo) BindFor(target FboTarget, context utils.BindingClosure) {
	fbo.Bind(target)
	context()
	fbo.Unbind(target)
}

func (fbo *Fbo) Destroy() {
	gl.DeleteFramebuffers(1, fbo.uint32)
	*fbo.uint32 = 0
}

func FboCheck(target FboTarget) error {
	status := gl.CheckFramebufferStatus(uint32(target))
	var glStatus string
	var statusStr string
	switch status {
	case gl.FRAMEBUFFER_COMPLETE:
		return nil
	case gl.FRAMEBUFFER_UNDEFINED:
		glStatus = "GL_FRAMEBUFFER_UNDEFINED"
		statusStr = "Framebuffer does not exist"
	case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
		glStatus = "FRAMEBUFFER_INCOMPLETE_ATTACHMENT"
		statusStr = "An attachment is incomplete"
	case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
		glStatus = "FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT"
		statusStr = "No attachments"
	case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
		glStatus = "FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER"
		statusStr = "A draw buffer is incomplete"
	case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
		glStatus = "FRAMEBUFFER_INCOMPLETE_READ_BUFFER"
		statusStr = "A read buffer is incomplete"
	case gl.FRAMEBUFFER_UNSUPPORTED:
		glStatus = "FRAMEBUFFER_UNSUPPORTED"
		statusStr = "An attachment has an unspported internal format"
	case gl.FRAMEBUFFER_INCOMPLETE_MULTISAMPLE:
		glStatus = "FRAMEBUFFER_INCOMPLETE_MULTISAMPLE"
		statusStr = "Different attachment samples"
	case gl.FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS:
		glStatus = "FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS"
		statusStr = "Different or invalid attachment layers"
	case 0:
		glStatus = "0"
		statusStr = "Unexpected error"
	default:
		glStatus = fmt.Sprintf("%X", status)
		statusStr = "Unknown"
	}

	enum := gl.DRAW_FRAMEBUFFER_BINDING
	if target == gl.READ_FRAMEBUFFER {
		enum = gl.READ_FRAMEBUFFER_BINDING
	}
	var id int32
	gl.GetIntegerv(uint32(enum), &id)

	return FboStatusError{
		GlStatus: glStatus,
		Status:   statusStr,
		Id:       uint32(id),
		Target:   target,
	}
}

func (fbo *Fbo) AttachTexture(tex GLTexture, attachment FboAttachment, level int) {
	gl.FramebufferTexture(gl.FRAMEBUFFER, attachment.ToGlEnum(), tex.Id(), int32(level))
}

func (fbo *Fbo) AttachTextureLayer(tex GLTexture, attachment FboAttachment, level, layer int) {
	gl.FramebufferTextureLayer(gl.FRAMEBUFFER, attachment.ToGlEnum(), tex.Id(), int32(level), int32(layer))
}

func (fbo *Fbo) AttachRenderbuffer(bufId uint32, attachment FboAttachment) {
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, attachment.ToGlEnum(), gl.RENDERBUFFER, bufId)
}
