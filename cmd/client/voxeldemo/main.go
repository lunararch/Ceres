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
    "Ceres/pkg/voxel"

    "github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
    window, err := graphics.NewWindow(1280, 720, "Step 8: Voxel Data Structure Demo")
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

    cam := camera.NewCamera(ceresmath.NewVector3(8, 6, 12))
    cam.LookAt(ceresmath.NewVector3(0, 0, 0))

    inputHandler := input.NewInputHandler(window.GetHandle())
    inputHandler.SetCursorMode(glfw.CursorDisabled)

    cubeMesh := mesh.NewCubeMesh(1.0)
    cubeRenderer := graphics.NewCubeRenderer(cubeMesh)
    defer cubeRenderer.Delete()

    voxelWorld := createDemoVoxelWorld()

    fmt.Println("✓ Voxel data structure demo initialized")
    fmt.Println("\nVoxel World Contents:")
    printVoxelWorldStats(voxelWorld)
    fmt.Println("\nControls:")
    fmt.Println("  WASD - Move horizontally")
    fmt.Println("  Space/Shift - Move up/down")
    fmt.Println("  Mouse - Look around")
    fmt.Println("  F1 - Toggle wireframe mode")
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

        // Lighting uniforms
        shader.SetVec3("lightPos", 10.0, 15.0, 10.0)
        shader.SetVec3("viewPos", cam.Position.X, cam.Position.Y, cam.Position.Z)
        shader.SetVec3("lightColor", 1.0, 1.0, 1.0)
        shader.SetInt("useTexture", 0)

        // Render all voxels
        renderVoxelWorld(voxelWorld, shader, cubeRenderer)

        if int(currentFrame.Unix())%2 == 0 {
            modeStr := "Solid"
            switch renderMode {
            case graphics.RenderModeWireframe:
                modeStr = "Wireframe"
            case graphics.RenderModeBoth:
                modeStr = "Both"
            }

            fmt.Printf("\rPos: (%.1f, %.1f, %.1f) | Voxels: %d | Mode: %s | FPS: %.0f    ",
                cam.Position.X, cam.Position.Y, cam.Position.Z,
                len(voxelWorld), modeStr, 1.0/deltaTime)
        }

        window.SwapBuffers()
        window.PollEvents()
    }

    fmt.Println("\n✓ Voxel data structure demo completed")
}

type VoxelInstance struct {
    Position voxel.VoxelPosition
    Voxel    voxel.Voxel
}

func createDemoVoxelWorld() []VoxelInstance {
    world := make([]VoxelInstance, 0)

    
    // Stone base (5x5 platform)
    for x := int32(-2); x <= 2; x++ {
        for z := int32(-2); z <= 2; z++ {
            world = append(world, VoxelInstance{
                Position: voxel.NewVoxelPosition(x, 0, z),
                Voxel:    voxel.NewVoxel(voxel.VoxelTypeStone),
            })
        }
    }

    // Grass layer on top
    for x := int32(-2); x <= 2; x++ {
        for z := int32(-2); z <= 2; z++ {
            world = append(world, VoxelInstance{
                Position: voxel.NewVoxelPosition(x, 1, z),
                Voxel:    voxel.NewVoxel(voxel.VoxelTypeGrass),
            })
        }
    }

    // Showcase tower - different materials on each level
    materials := []voxel.VoxelType{
        voxel.VoxelTypeDirt,
        voxel.VoxelTypeSand,
        voxel.VoxelTypeWood,
        voxel.VoxelTypeBrick,
    }

    for i, material := range materials {
        y := int32(2 + i)
        world = append(world, VoxelInstance{
            Position: voxel.NewVoxelPosition(0, y, 0),
            Voxel:    voxel.NewVoxel(material),
        })
    }

    // Glass window at the top
    world = append(world, VoxelInstance{
        Position: voxel.NewVoxelPosition(0, 6, 0),
        Voxel:    voxel.NewVoxel(voxel.VoxelTypeGlass),
    })

    // Small tree
    // Trunk (wood)
    for y := int32(2); y <= 4; y++ {
        world = append(world, VoxelInstance{
            Position: voxel.NewVoxelPosition(-2, y, -2),
            Voxel:    voxel.NewVoxel(voxel.VoxelTypeWood),
        })
    }

    leafPositions := []voxel.VoxelPosition{
        voxel.NewVoxelPosition(-2, 5, -2),
        voxel.NewVoxelPosition(-3, 5, -2),
        voxel.NewVoxelPosition(-1, 5, -2),
        voxel.NewVoxelPosition(-2, 5, -3),
        voxel.NewVoxelPosition(-2, 5, -1),
    }

    for _, pos := range leafPositions {
        world = append(world, VoxelInstance{
            Position: pos,
            Voxel:    voxel.NewVoxel(voxel.VoxelTypeLeaves),
        })
    }

    for x := int32(1); x <= 2; x++ {
        for z := int32(1); z <= 2; z++ {
            world = append(world, VoxelInstance{
                Position: voxel.NewVoxelPosition(x, 2, z),
                Voxel:    voxel.NewVoxel(voxel.VoxelTypeWater),
            })
        }
    }

    return world
}

func renderVoxelWorld(world []VoxelInstance, shader *graphics.Shader, renderer *graphics.CubeRenderer) {
    for _, instance := range world {
        // Skip air voxels
        if instance.Voxel.IsAir() {
            continue
        }

        worldPos := instance.Position.ToWorldSpace()

        model := ceresmath.TranslateVec(worldPos)

        color := getVoxelColor(instance.Voxel.Type)
        shader.SetVec3("objectColor", color.X, color.Y, color.Z)

        shader.SetMat4("model", model.ToPtr())
        renderer.Render()
    }
}

func getVoxelColor(voxelType voxel.VoxelType) ceresmath.Vector3 {
    switch voxelType {
    case voxel.VoxelTypeStone:
        return ceresmath.NewVector3(0.5, 0.5, 0.5) // Gray
    case voxel.VoxelTypeDirt:
        return ceresmath.NewVector3(0.55, 0.35, 0.2) // Brown
    case voxel.VoxelTypeGrass:
        return ceresmath.NewVector3(0.2, 0.8, 0.2) // Green
    case voxel.VoxelTypeSand:
        return ceresmath.NewVector3(0.95, 0.9, 0.6) // Sand yellow
    case voxel.VoxelTypeWater:
        return ceresmath.NewVector3(0.2, 0.4, 0.9) // Blue
    case voxel.VoxelTypeWood:
        return ceresmath.NewVector3(0.6, 0.4, 0.2) // Wood brown
    case voxel.VoxelTypeLeaves:
        return ceresmath.NewVector3(0.15, 0.6, 0.15) // Dark green
    case voxel.VoxelTypeGlass:
        return ceresmath.NewVector3(0.7, 0.9, 1.0) // Light blue
    case voxel.VoxelTypeBrick:
        return ceresmath.NewVector3(0.7, 0.3, 0.2) // Red brick
    default:
        return ceresmath.NewVector3(1.0, 1.0, 1.0) // White
    }
}

func printVoxelWorldStats(world []VoxelInstance) {
    typeCounts := make(map[voxel.VoxelType]int)

    for _, instance := range world {
        typeCounts[instance.Voxel.Type]++
    }

    fmt.Printf("Total voxels: %d\n", len(world))
    for voxelType, count := range typeCounts {
        v := voxel.NewVoxel(voxelType)
        fmt.Printf("  %s: %d\n", v.GetName(), count)
    }
}

func processInput(inputHandler *input.InputHandler, cam *camera.Camera, deltaTime float32,
    window *graphics.Window, renderMode *graphics.RenderMode, cubeRenderer *graphics.CubeRenderer) bool {

    // Movement
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

    if inputHandler.IsKeyPressed(glfw.KeyEscape) {
        return true
    }

    return false
}