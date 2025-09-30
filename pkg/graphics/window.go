package graphics

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	DefaultWidth = 1280
	DefaultHeight = 720
	DefaultTitle  = "Ceres Voxel Engine"
)

type Window struct {
	handle *glfw.Window
	Width int
	Height int
}

func NewWindow(width, height int, title string) (*Window, error){
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize GLFW: %w", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	glfw.WindowHint(glfw.DepthBits, 24)

	glfw.WindowHint(glfw.Samples, 4)

	handle, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
        glfw.Terminate()
        return nil, fmt.Errorf("failed to create window: %w", err)
    }

	handle.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		handle.Destroy()
		glfw.Terminate()
		return nil, fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	w := &Window {
		handle: handle,
		Width: width,
		Height: height,
	}

	handle.SetFramebufferSizeCallback(w.onResize)

	glfw.SwapInterval(1)

	w.initializeOpenGLState()

	return w, nil
}

func (w *Window) initializeOpenGLState() {
    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)

    gl.Enable(gl.CULL_FACE)
    gl.CullFace(gl.BACK)
    gl.FrontFace(gl.CCW)

    gl.Enable(gl.MULTISAMPLE)

    gl.ClearColor(0.53, 0.81, 0.92, 1.0)

    gl.Viewport(0, 0, int32(w.Width), int32(w.Height))

    fmt.Println("OpenGL Version:", gl.GoStr(gl.GetString(gl.VERSION)))
    fmt.Println("GLSL Version:", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
    fmt.Println("Renderer:", gl.GoStr(gl.GetString(gl.RENDERER)))
}

func (w *Window) onResize(window *glfw.Window, width, height int) {
    w.Width = width
    w.Height = height
    gl.Viewport(0, 0, int32(width), int32(height))
    fmt.Printf("Window resized to %dx%d\n", width, height)
}

func (w *Window) ShouldClose() bool {
    return w.handle.ShouldClose()
}

func (w *Window) Clear() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (w *Window) SwapBuffers() {
    w.handle.SwapBuffers()
}

func (w *Window) PollEvents() {
    glfw.PollEvents()
}

func (w *Window) GetAspectRatio() float32 {
    return float32(w.Width) / float32(w.Height)
}

func (w *Window) SetKeyCallback(cb glfw.KeyCallback) {
    w.handle.SetKeyCallback(cb)
}

func (w *Window) SetCursorPosCallback(cb glfw.CursorPosCallback) {
    w.handle.SetCursorPosCallback(cb)
}

func (w *Window) SetMouseButtonCallback(cb glfw.MouseButtonCallback) {
    w.handle.SetMouseButtonCallback(cb)
}

func (w *Window) SetCursorMode(mode int) {
    w.handle.SetInputMode(glfw.CursorMode, mode)
}

func (w *Window) GetKey(key glfw.Key) glfw.Action {
    return w.handle.GetKey(key)
}

func (w *Window) Close() {
    w.handle.Destroy()
    glfw.Terminate()
}

func (w *Window) GetHandle() *glfw.Window {
    return w.handle
}