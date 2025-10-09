package main

import (
    "fmt"
    "log"
    "time"

    "Ceres/pkg/camera"
    "Ceres/pkg/graphics"
    "Ceres/pkg/input"
    ceresmath "Ceres/pkg/math"

    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
    window, err := graphics.NewWindow(1280, 720, "Camera System Test")
    if err != nil {
        log.Fatal(err)
    }
    defer window.Close()

    shaderManager := graphics.NewShaderManager()
    defer shaderManager.DeleteAll()

    if err := shaderManager.LoadDefaultShaders(); err != nil {
        log.Fatal(err)
    }

    shader, err := shaderManager.GetShader("simple")
    if err != nil {
        log.Fatal(err)
    }

    cam := camera.NewCamera(ceresmath.NewVector3(0, 2, 5))

    inputHandler := input.NewInputHandler(window.GetHandle())
    inputHandler.SetCursorMode(glfw.CursorDisabled) // Capture cursor for FPS controls

    cubes := createCubeGrid(10, 10)

    var vao, vbo uint32
    vertices := createCubeVertices()
    gl.GenVertexArrays(1, &vao)
    gl.GenBuffers(1, &vbo)

    gl.BindVertexArray(vao)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
    gl.EnableVertexAttribArray(0)

    gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
    gl.EnableVertexAttribArray(1)

    gl.BindVertexArray(0)

    fmt.Println("✓ Camera system initialized")
    fmt.Println("Controls:")
    fmt.Println("  WASD - Move horizontally")
    fmt.Println("  Space/Shift - Move up/down")
    fmt.Println("  Mouse - Look around")
    fmt.Println("  ESC - Exit")

    lastFrame := time.Now()
    showControls := true

    for !window.ShouldClose() {
        currentFrame := time.Now()
        deltaTime := float32(currentFrame.Sub(lastFrame).Seconds())
        lastFrame = currentFrame

        processInput(inputHandler, cam, deltaTime, window)

        window.Clear()

        projection := cam.GetProjectionMatrix(window.GetAspectRatio(), 0.1, 100.0)
        view := cam.GetViewMatrix()

        shader.Use()
        shader.SetMat4("projection", projection.ToPtr())
        shader.SetMat4("view", view.ToPtr())

        gl.BindVertexArray(vao)
        for _, cubePos := range cubes {
            model := ceresmath.TranslateVec(cubePos)
            shader.SetMat4("model", model.ToPtr())
            gl.DrawArrays(gl.TRIANGLES, 0, 36)
        }
        gl.BindVertexArray(0)

        if showControls && int(currentFrame.Unix())%5 == 0 {
            fmt.Printf("\rPos: (%.1f, %.1f, %.1f) | Yaw: %.0f° | Pitch: %.0f° | FPS: %.0f",
                cam.Position.X, cam.Position.Y, cam.Position.Z,
                cam.Yaw, cam.Pitch, 1.0/deltaTime)
        }

        window.SwapBuffers()
        window.PollEvents()
    }

    gl.DeleteVertexArrays(1, &vao)
    gl.DeleteBuffers(1, &vbo)

    fmt.Println("\n✓ Camera test completed")
}

func processInput(inputHandler *input.InputHandler, cam *camera.Camera, deltaTime float32, window *graphics.Window) {
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

    if inputHandler.IsKeyPressed(glfw.KeyEscape) {
        window.GetHandle().SetShouldClose(true)
    }

    if inputHandler.IsKeyPressed(glfw.KeyTab) {
        inputHandler.SetCursorMode(glfw.CursorNormal)
    }
}

func createCubeGrid(width, depth int) []ceresmath.Vector3 {
    cubes := make([]ceresmath.Vector3, 0, width*depth)
    
    for x := 0; x < width; x++ {
        for z := 0; z < depth; z++ {
            cubes = append(cubes, ceresmath.NewVector3(
                float32(x*2-width),
                0,
                float32(z*2-depth),
            ))
        }
    }
    
    return cubes
}

func createCubeVertices() []float32 {
    return []float32{
        // Front face (Red)
        -0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
        0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
        0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
        0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
        -0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
        -0.5, -0.5, 0.5, 1.0, 0.0, 0.0,

        // Back face (Green)
        -0.5, -0.5, -0.5, 0.0, 1.0, 0.0,
        -0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
        0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
        0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
        0.5, -0.5, -0.5, 0.0, 1.0, 0.0,
        -0.5, -0.5, -0.5, 0.0, 1.0, 0.0,

        // Top face (Blue)
        -0.5, 0.5, -0.5, 0.0, 0.0, 1.0,
        -0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
        0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
        0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
        0.5, 0.5, -0.5, 0.0, 0.0, 1.0,
        -0.5, 0.5, -0.5, 0.0, 0.0, 1.0,

        // Bottom face (Yellow)
        -0.5, -0.5, -0.5, 1.0, 1.0, 0.0,
        0.5, -0.5, -0.5, 1.0, 1.0, 0.0,
        0.5, -0.5, 0.5, 1.0, 1.0, 0.0,
        0.5, -0.5, 0.5, 1.0, 1.0, 0.0,
        -0.5, -0.5, 0.5, 1.0, 1.0, 0.0,
        -0.5, -0.5, -0.5, 1.0, 1.0, 0.0,

        // Right face (Cyan)
        0.5, -0.5, -0.5, 0.0, 1.0, 1.0,
        0.5, 0.5, -0.5, 0.0, 1.0, 1.0,
        0.5, 0.5, 0.5, 0.0, 1.0, 1.0,
        0.5, 0.5, 0.5, 0.0, 1.0, 1.0,
        0.5, -0.5, 0.5, 0.0, 1.0, 1.0,
        0.5, -0.5, -0.5, 0.0, 1.0, 1.0,

        // Left face (Magenta)
        -0.5, -0.5, -0.5, 1.0, 0.0, 1.0,
        -0.5, -0.5, 0.5, 1.0, 0.0, 1.0,
        -0.5, 0.5, 0.5, 1.0, 0.0, 1.0,
        -0.5, 0.5, 0.5, 1.0, 0.0, 1.0,
        -0.5, 0.5, -0.5, 1.0, 0.0, 1.0,
        -0.5, -0.5, -0.5, 1.0, 0.0, 1.0,
    }
}