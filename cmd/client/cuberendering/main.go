package main

import (
    "fmt"
    "log"
    "time"

    "Ceres/pkg/camera"
    "Ceres/pkg/graphics"
    "Ceres/pkg/input"
    ceresmath "Ceres/pkg/math"
    "Ceres/pkg/mesh"

    "github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
    window, err := graphics.NewWindow(1280, 720, "Step 7: Basic Cube Rendering")
    if err != nil {
        log.Fatal(err)
    }
    defer window.Close()

    shaderManager := graphics.NewShaderManager()
    defer shaderManager.DeleteAll()

    if err := shaderManager.LoadDefaultShaders(); err != nil {
        log.Fatal(err)
    }

    shader, err := shaderManager.GetShader("basic")
    if err != nil {
        log.Fatal(err)
    }

    cam := camera.NewCamera(ceresmath.NewVector3(0, 2, 5))

    inputHandler := input.NewInputHandler(window.GetHandle())
    inputHandler.SetCursorMode(glfw.CursorDisabled)

    cubeMesh := mesh.NewCubeMesh(1.0)
    cubeRenderer := graphics.NewCubeRenderer(cubeMesh)
    defer cubeRenderer.Delete()

    cubePositions := []ceresmath.Vector3{
        ceresmath.NewVector3(0, 0, 0),
        ceresmath.NewVector3(2, 0, 0),
        ceresmath.NewVector3(-2, 0, 0),
        ceresmath.NewVector3(0, 2, 0),
        ceresmath.NewVector3(0, -2, 0),
        ceresmath.NewVector3(2, 2, 2),
        ceresmath.NewVector3(-2, -2, -2),
    }

    fmt.Println("✓ Cube rendering system initialized")
    fmt.Println("\nControls:")
    fmt.Println("  WASD - Move horizontally")
    fmt.Println("  Space/Shift - Move up/down")
    fmt.Println("  Mouse - Look around")
    fmt.Println("  F1 - Toggle wireframe mode")
    fmt.Println("  F2 - Toggle solid mode")
    fmt.Println("  F3 - Toggle both modes")
    fmt.Println("  ESC - Exit")

    lastFrame := time.Now()
    renderMode := graphics.RenderModeSolid

    for !window.ShouldClose() {
        currentFrame := time.Now()
        deltaTime := float32(currentFrame.Sub(lastFrame).Seconds())
        lastFrame = currentFrame

        if processInput(inputHandler, cam, deltaTime, window, &renderMode, cubeRenderer) {
            break
        }

        window.Clear()

        projection := cam.GetProjectionMatrix(window.GetAspectRatio(), 0.1, 100.0)
        view := cam.GetViewMatrix()

        shader.Use()
        shader.SetMat4("projection", projection.ToPtr())
        shader.SetMat4("view", view.ToPtr())

        shader.SetVec3("lightPos", 10.0, 10.0, 10.0)
        shader.SetVec3("viewPos", cam.Position.X, cam.Position.Y, cam.Position.Z)
        shader.SetVec3("lightColor", 1.0, 1.0, 1.0)
        shader.SetInt("useTexture", 0)

        for i, pos := range cubePositions {
            model := ceresmath.TranslateVec(pos)

            hue := float32(i) / float32(len(cubePositions))
            r, g, b := hsvToRgb(hue, 0.8, 0.9)
            shader.SetVec3("objectColor", r, g, b)

            shader.SetMat4("model", model.ToPtr())
            cubeRenderer.Render()
        }

        if int(currentFrame.Unix())%2 == 0 {
            modeStr := "Solid"
            switch renderMode {
            case graphics.RenderModeWireframe:
                modeStr = "Wireframe"
            case graphics.RenderModeBoth:
                modeStr = "Both"
            }

            fmt.Printf("\rPos: (%.1f, %.1f, %.1f) | Mode: %s | FPS: %.0f    ",
                cam.Position.X, cam.Position.Y, cam.Position.Z,
                modeStr, 1.0/deltaTime)
        }

        window.SwapBuffers()
        window.PollEvents()
    }

    fmt.Println("\n✓ Cube rendering test completed")
}

func processInput(inputHandler *input.InputHandler, cam *camera.Camera, deltaTime float32,
    window *graphics.Window, renderMode *graphics.RenderMode, cubeRenderer *graphics.CubeRenderer) bool {

    if inputHandler.IsKeyPressed(glfw.KeyW) {
        cam.ProcessKeyboard(camera.Forward, deltaTime)
    }
    if inputHandler.IsKeyPressed(glfw.KeyS) {
        cam.ProcessKeyboard(camera.Backward, deltaTime)
    }
    if inputHandler.IsKeyPressed(glfw.KeyA) {
        cam.ProcessKeyboard(camera.Left, deltaTime)
    }
    if inputHandler.IsKeyPressed(glfw.KeyD) {
        cam.ProcessKeyboard(camera.Right, deltaTime)
    }
    if inputHandler.IsKeyPressed(glfw.KeySpace) {
        cam.ProcessKeyboard(camera.Up, deltaTime)
    }
    if inputHandler.IsKeyPressed(glfw.KeyLeftShift) {
        cam.ProcessKeyboard(camera.Down, deltaTime)
    }

    xOffset, yOffset := inputHandler.GetMouseMovement()
    cam.ProcessMouseMovement(float32(xOffset), float32(yOffset), true)

    if inputHandler.IsKeyPressed(glfw.KeyF1) && *renderMode != graphics.RenderModeWireframe {
        *renderMode = graphics.RenderModeWireframe
        cubeRenderer.SetRenderMode(*renderMode)
        fmt.Println("\nSwitched to Wireframe mode")
    }
    if inputHandler.IsKeyPressed(glfw.KeyF2) && *renderMode != graphics.RenderModeSolid {
        *renderMode = graphics.RenderModeSolid
        cubeRenderer.SetRenderMode(*renderMode)
        fmt.Println("\nSwitched to Solid mode")
    }
    if inputHandler.IsKeyPressed(glfw.KeyF3) && *renderMode != graphics.RenderModeBoth {
        *renderMode = graphics.RenderModeBoth
        cubeRenderer.SetRenderMode(*renderMode)
        fmt.Println("\nSwitched to Both modes")
    }

    if inputHandler.IsKeyPressed(glfw.KeyEscape) {
        return true
    }

    return false
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