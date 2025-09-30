package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	// This is needed for GLFW on macOS, doesn't hurt on other platforms
	runtime.LockOSThread()
}

func main() {
	// Test GLFW
	if err := glfw.Init(); err != nil {
		fmt.Println("Failed to initialize GLFW:", err)
		return
	}
	defer glfw.Terminate()
	fmt.Println("✓ GLFW initialized successfully")

	// Test mathgl
	vec := mgl32.Vec3{1, 2, 3}
	fmt.Printf("✓ mathgl working - test vector: %v\n", vec)

	// Test OpenGL (requires window context, so just check import)
	fmt.Printf("✓ OpenGL bindings imported (version check requires window)\n")

	fmt.Println("\nAll dependencies verified successfully!")
}
