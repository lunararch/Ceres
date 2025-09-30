package main

import (
    "fmt"
    "log"

    "Ceres/pkg/graphics"

    "github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
    // Create window
    window, err := graphics.NewWindow(
        graphics.DefaultWidth,
        graphics.DefaultHeight,
        graphics.DefaultTitle,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer window.Close()

    fmt.Println("Window created successfully!")
    fmt.Printf("Resolution: %dx%d\n", window.Width, window.Height)
    fmt.Printf("Aspect Ratio: %.2f\n", window.GetAspectRatio())

    window.SetKeyCallback(onKey)

    frameCount := 0
    for !window.ShouldClose() {
        window.Clear()


        window.SwapBuffers()
        window.PollEvents()

        frameCount++
        if frameCount%60 == 0 {
            fmt.Printf("Running... (Frame %d)\n", frameCount)
        }
    }

    fmt.Println("Shutting down gracefully...")
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    if key == glfw.KeyEscape && action == glfw.Press {
        w.SetShouldClose(true)
    }

    if action == glfw.Press {
        switch key {
        case glfw.KeyF11:
            fmt.Println("F11 pressed - fullscreen toggle not yet implemented")
        case glfw.KeyF1:
            fmt.Println("F1 pressed - wireframe toggle not yet implemented")
        }
    }
}