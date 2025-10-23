package main

import (
	"fmt"
	"log"
	"time"

	"Ceres/pkg/camera"
	"Ceres/pkg/chunk"
	"Ceres/pkg/graphics"
	"Ceres/pkg/input"
	ceresmath "Ceres/pkg/math"
	"Ceres/pkg/mesh"
	"Ceres/pkg/voxel"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
    window, err := graphics.NewWindow(1280, 720, "Step 9: Chunk System Demo")
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

    cam := camera.NewCamera(ceresmath.NewVector3(chunk.ChunkSize*1.5, chunk.ChunkSize, chunk.ChunkSize*2))
    cam.LookAt(ceresmath.NewVector3(0, chunk.ChunkSize/2, 0))

    inputHandler := input.NewInputHandler(window.GetHandle())
    inputHandler.SetCursorMode(glfw.CursorDisabled)

    cubeMesh := mesh.NewCubeMesh(1.0)
    cubeRenderer := graphics.NewCubeRenderer(cubeMesh)
    defer cubeRenderer.Delete()

    chunkManager := chunk.NewChunkManager()

    generateDemoWorld(chunkManager)

    fmt.Println("✓ Chunk system demo initialized")
    fmt.Println("\nWorld Configuration:")
    fmt.Printf("  Chunk Size: %dx%dx%d voxels\n", chunk.ChunkSize, chunk.ChunkSize, chunk.ChunkSize)
    chunkManager.PrintStats()
    
    fmt.Println("\nControls:")
    fmt.Println("  WASD - Move horizontally")
    fmt.Println("  Space/Shift - Move up/down")
    fmt.Println("  Mouse - Look around")
    fmt.Println("  F1 - Toggle wireframe mode")
    fmt.Println("  ESC - Exit")

    lastFrame := time.Now()
    renderMode := graphics.RenderModeSolid
    voxelCount := 0

    for !window.ShouldClose() {
        currentFrame := time.Now()
        deltaTime := float32(currentFrame.Sub(lastFrame).Seconds())
        lastFrame = currentFrame

        if processInput(inputHandler, cam, deltaTime, window, &renderMode, cubeRenderer) {
            break
        }

        window.Clear()

        projection := cam.GetProjectionMatrix(window.GetAspectRatio(), 0.1, 500.0)
        view := cam.GetViewMatrix()

        shader.Use()
        shader.SetMat4("projection", projection.ToPtr())
        shader.SetMat4("view", view.ToPtr())

        shader.SetVec3("lightPos", 100.0, 100.0, 100.0)
        shader.SetVec3("viewPos", cam.Position.X, cam.Position.Y, cam.Position.Z)
        shader.SetVec3("lightColor", 1.0, 1.0, 1.0)
        shader.SetInt("useTexture", 0)

        voxelCount = renderChunks(chunkManager, shader, cubeRenderer)

        if int(currentFrame.Unix())%2 == 0 {
            stats := chunkManager.GetStats()
            modeStr := "Solid"
            switch renderMode {
            case graphics.RenderModeWireframe:
                modeStr = "Wireframe"
            case graphics.RenderModeBoth:
                modeStr = "Both"
            }

            fmt.Printf("\rPos: (%.0f, %.0f, %.0f) | Chunks: %d | Voxels: %d | Mode: %s | FPS: %.0f    ",
                cam.Position.X, cam.Position.Y, cam.Position.Z,
                stats.LoadedChunks, voxelCount, modeStr, 1.0/deltaTime)
        }

        window.SwapBuffers()
        window.PollEvents()
    }

    fmt.Println("\n✓ Chunk system demo completed")
}

func generateDemoWorld(cm *chunk.ChunkManager) {
    for cx := int32(-1); cx <= 1; cx++ {
        for cz := int32(-1); cz <= 1; cz++ {
            chunkPos := chunk.NewChunkPosition(cx, 0, cz)
            c := cm.CreateChunk(chunkPos)

            for x := int32(0); x < chunk.ChunkSize; x++ {
                for z := int32(0); z < chunk.ChunkSize; z++ {
                    worldX := cx*chunk.ChunkSize + x
                    worldZ := cz*chunk.ChunkSize + z
                    height := int32(8 + ((worldX+worldZ)%8))

                    for y := int32(0); y < height && y < chunk.ChunkSize; y++ {
                        var voxelType voxel.VoxelType
                        
                        if y == height-1 {
                            voxelType = voxel.VoxelTypeGrass
                        } else if y >= height-4 {
                            voxelType = voxel.VoxelTypeDirt
                        } else {
                            voxelType = voxel.VoxelTypeStone
                        }
                        
                        c.SetVoxel(x, y, z, voxel.NewVoxel(voxelType))
                    }
                }
            }

            if cx == 0 && cz == 0 {
                for y := int32(10); y < 20; y++ {
                    c.SetVoxel(16, y, 16, voxel.NewVoxel(voxel.VoxelTypeBrick))
                }
            }
        }
    }
}

func renderChunks(cm *chunk.ChunkManager, shader *graphics.Shader, renderer *graphics.CubeRenderer) int {
    voxelCount := 0
    chunks := cm.GetLoadedChunks()

    for _, c := range chunks {
        if c.IsEmpty() {
            continue
        }

        worldPos := c.GetWorldPosition()

        for x := int32(0); x < chunk.ChunkSize; x++ {
            for y := int32(0); y < chunk.ChunkSize; y++ {
                for z := int32(0); z < chunk.ChunkSize; z++ {
                    v := c.GetVoxel(x, y, z)
                    if v.IsAir() {
                        continue
                    }

                    voxelWorldPos := worldPos.Add(voxel.NewVoxelPosition(x, y, z))
                    model := ceresmath.TranslateVec(voxelWorldPos.ToWorldSpace())

                    color := getVoxelColor(v.Type)
                    shader.SetVec3("objectColor", color.X, color.Y, color.Z)

                    shader.SetMat4("model", model.ToPtr())
                    renderer.Render()

                    voxelCount++
                }
            }
        }
    }

    return voxelCount
}

func getVoxelColor(voxelType voxel.VoxelType) ceresmath.Vector3 {
    switch voxelType {
    case voxel.VoxelTypeStone:
        return ceresmath.NewVector3(0.5, 0.5, 0.5)
    case voxel.VoxelTypeDirt:
        return ceresmath.NewVector3(0.55, 0.35, 0.2)
    case voxel.VoxelTypeGrass:
        return ceresmath.NewVector3(0.2, 0.8, 0.2)
    case voxel.VoxelTypeBrick:
        return ceresmath.NewVector3(0.7, 0.3, 0.2)
    default:
        return ceresmath.NewVector3(1.0, 1.0, 1.0)
    }
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

    if inputHandler.IsKeyPressed(glfw.KeyEscape) {
        return true
    }

    return false
}