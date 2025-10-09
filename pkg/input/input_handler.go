package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type InputHandler struct {
	window *glfw.Window

	firstMouse bool
	lastX float64
	lastY float64
	mouseMovedX float64
	mouseMovedY float64

	keys map[glfw.Key]bool
}

func NewInputHandler(window *glfw.Window) *InputHandler {
	h := &InputHandler{
		window: window,
		firstMouse: true,
		lastX:      0,
        lastY:      0,
        keys:       make(map[glfw.Key]bool),
	}

	window.SetCursorPosCallback(h.mouseCallback)
	window.SetKeyCallback(h.keyCallback)

	return h
}

func (h *InputHandler) mouseCallback(w *glfw.Window, xpos, ypos float64) {
	if h.firstMouse {
		h.lastX = xpos
		h.lastY = ypos
		h.firstMouse = false
	}

	h.mouseMovedX = xpos - h.lastX
	h.mouseMovedY = h.lastY - ypos

	h.lastX = xpos
	h.lastY = ypos
}

func (h *InputHandler) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    if action == glfw.Press {
        h.keys[key] = true
    } else if action == glfw.Release {
        h.keys[key] = false
    }
}

func (h *InputHandler) GetMouseMovement() (float64, float64) {
    x, y := h.mouseMovedX, h.mouseMovedY
    h.mouseMovedX = 0
    h.mouseMovedY = 0
    return x, y
}

func (h *InputHandler) IsKeyPressed(key glfw.Key) bool {
    return h.keys[key]
}

func (h *InputHandler) SetCursorMode(mode int) {
    h.window.SetInputMode(glfw.CursorMode, mode)
}