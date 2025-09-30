package main

import (
    "fmt"
    "log"
    "time"

    "Ceres/pkg/graphics"

    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
    window, err := graphics.NewWindow(800, 600, "Window Test")
    if err != nil {
        log.Fatal(err)
    }
    defer window.Close()

    fmt.Println("✓ Window created successfully")
    fmt.Println("✓ OpenGL context initialized")
    fmt.Println("Press ESC to close the window")

    var hue float32 = 0.0
    lastTime := time.Now()
    frames := 0

    for !window.ShouldClose() {
        frames++
        if time.Since(lastTime) >= time.Second {
            fmt.Printf("FPS: %d\n", frames)
            frames = 0
            lastTime = time.Now()
        }

        hue += 0.001
        if hue > 1.0 {
            hue = 0.0
        }
        r, g, b := hsvToRgb(hue, 0.5, 0.8)
        gl.ClearColor(r, g, b, 1.0)

        window.Clear()
        window.SwapBuffers()
        window.PollEvents()

        if window.GetKey(glfw.KeyEscape) == glfw.Press {
            break
        }
    }

    fmt.Println("✓ Window closed successfully")
}

func hsvToRgb(h, s, v float32) (float32, float32, float32) {
    c := v * s
    x := c * (1 - abs(mod(h*6, 2)-1))
    m := v - c

    var r, g, b float32
    switch {
    case h < 1.0/6.0:
        r, g, b = c, x, 0
    case h < 2.0/6.0:
        r, g, b = x, c, 0
    case h < 3.0/6.0:
        r, g, b = 0, c, x
    case h < 4.0/6.0:
        r, g, b = 0, x, c
    case h < 5.0/6.0:
        r, g, b = x, 0, c
    default:
        r, g, b = c, 0, x
    }

    return r + m, g + m, b + m
}

func abs(x float32) float32 {
    if x < 0 {
        return -x
    }
    return x
}

func mod(x, y float32) float32 {
    return x - y*float32(int(x/y))
}