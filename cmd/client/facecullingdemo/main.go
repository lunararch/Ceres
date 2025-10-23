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

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	window, err := graphics.NewWindow(1280, 720, "Step 11: Face Culling Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer window.Close()

	if err := gl.Init(); err != nil {
		log.Fatal("Failed to initialize OpenGL:", err)
	}

	shaderManager := graphics.NewShaderManager()
	defer shaderManager.DeleteAll()

	if err := shaderManager.LoadDefaultShaders(); err != nil {
		log.Fatal(err)
	}

	shader, err := shaderManager.GetShader("basic")
	if err != nil {
		log.Fatal(err)
	}

	cam := camera.NewCamera(ceresmath.NewVector3(15, 10, 15))
	cam.LookAt(ceresmath.NewVector3(0, 0, 0))

	inputHandler := input.NewInputHandler(window.GetHandle())
	inputHandler.SetCursorMode(glfw.CursorDisabled)

	worldMesh := buildTestWorld()
	defer worldMesh.Delete()

	lastFrame := time.Now()
	frameCount := 0
	fpsTimer := time.Now()

	fmt.Println("=== Face Culling Demo ===")
	fmt.Println("Controls:")
	fmt.Println("  WASD: Move camera")
	fmt.Println("  Space/Shift: Move up/down")
	fmt.Println("  Mouse: Look around")
	fmt.Println("  ESC: Exit")
	fmt.Println()
	printMeshStats(worldMesh)

	for !window.ShouldClose() {
		currentFrame := time.Now()
		deltaTime := float32(currentFrame.Sub(lastFrame).Seconds())
		lastFrame = currentFrame

		if processInput(inputHandler, cam, deltaTime, window) {
			break
		}

		window.Clear()

		projection := cam.GetProjectionMatrix(window.GetAspectRatio(), 0.1, 100.0)
		view := cam.GetViewMatrix()
		model := ceresmath.Identity()

		shader.Use()
		shader.SetMat4("projection", projection.ToPtr())
		shader.SetMat4("view", view.ToPtr())
		shader.SetMat4("model", model.ToPtr())
		shader.SetVec3("lightPos", 20.0, 20.0, 20.0)
		shader.SetVec3("viewPos", cam.Position.X, cam.Position.Y, cam.Position.Z)
		shader.SetVec3("lightColor", 1.0, 1.0, 1.0)
		shader.SetFloat("ambientStrength", 0.3)
		shader.SetInt("useTexture", 0)

		worldMesh.Draw()

		window.SwapBuffers()
		glfw.PollEvents()

		frameCount++
		if time.Since(fpsTimer) >= time.Second {
			fmt.Printf("\rFPS: %d | Vertices: %d | Triangles: %d | Faces Culled: ~%.1f%%",
				frameCount, worldMesh.GetVertexCount(), worldMesh.GetTriangleCount(),
				calculateCullingPercentage(worldMesh.GetTriangleCount()))
			frameCount = 0
			fpsTimer = time.Now()
		}
	}
	fmt.Println()
}

func buildTestWorld() *mesh.Mesh {
	builder := mesh.NewVoxelMeshBuilder()

	size := int32(10)

	for x := int32(0); x < size; x++ {
		for y := int32(0); y < size; y++ {
			for z := int32(0); z < size; z++ {
				pos := voxel.NewVoxelPosition(x, y, z)
				voxelType := getVoxelTypeForPosition(x, y, z, size)

				if voxelType == voxel.VoxelTypeAir {
					continue
				}

				neighbors := getNeighbors(x, y, z, size)
				color := getVoxelColor(voxelType)

				builder.AddVoxelWithCulling(pos, voxelType, neighbors, color)
			}
		}
	}

	vertices := builder.GetVertices()
	indices := builder.GetIndices()

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	return &mesh.Mesh{
		VAO:           vao,
		VBO:           vbo,
		EBO:           ebo,
		IndexCount:    int32(len(indices)),
		VertexCount:   int32(len(vertices) / 9),
		TriangleCount: int32(len(indices) / 3),
	}
}

func getVoxelTypeForPosition(x, y, z, size int32) voxel.VoxelType {
	isEdge := x == 0 || x == size-1 || y == 0 || y == size-1 || z == 0 || z == size-1
	isCross := (x == size/2 || y == size/2 || z == size/2)

	if isEdge {
		return voxel.VoxelTypeStone
	} else if isCross {
		return voxel.VoxelTypeBrick
	}
	return voxel.VoxelTypeAir
}

func getNeighbors(x, y, z, size int32) mesh.VoxelNeighbors {
	return mesh.VoxelNeighbors{
		Top:    getVoxelAt(x, y+1, z, size),
		Bottom: getVoxelAt(x, y-1, z, size),
		Left:   getVoxelAt(x-1, y, z, size),
		Right:  getVoxelAt(x+1, y, z, size),
		Front:  getVoxelAt(x, y, z+1, size),
		Back:   getVoxelAt(x, y, z-1, size),
	}
}

func getVoxelAt(x, y, z, size int32) voxel.Voxel {
	if x < 0 || x >= size || y < 0 || y >= size || z < 0 || z >= size {
		return voxel.NewVoxel(voxel.VoxelTypeAir)
	}
	voxelType := getVoxelTypeForPosition(x, y, z, size)
	return voxel.NewVoxel(voxelType)
}

func getVoxelColor(voxelType voxel.VoxelType) ceresmath.Vector3 {
	switch voxelType {
	case voxel.VoxelTypeStone:
		return ceresmath.NewVector3(0.5, 0.5, 0.5)
	case voxel.VoxelTypeBrick:
		return ceresmath.NewVector3(0.7, 0.3, 0.2)
	default:
		return ceresmath.NewVector3(1.0, 1.0, 1.0)
	}
}

func printMeshStats(m *mesh.Mesh) {
	fmt.Printf("Mesh Statistics:\n")
	fmt.Printf("  Vertices: %d\n", m.GetVertexCount())
	fmt.Printf("  Triangles: %d\n", m.GetTriangleCount())
	fmt.Printf("  Memory (approx): %.2f KB\n", float32(m.GetVertexCount()*9*4+m.GetTriangleCount()*3*4)/1024.0)
	fmt.Printf("  Faces Culled: ~%.1f%%\n", calculateCullingPercentage(m.GetTriangleCount()))
	fmt.Println()
}

func calculateCullingPercentage(trianglesRendered int32) float32 {
	size := int32(10)
	totalVoxels := int32(0)

	for x := int32(0); x < size; x++ {
		for y := int32(0); y < size; y++ {
			for z := int32(0); z < size; z++ {
				if getVoxelTypeForPosition(x, y, z, size) != voxel.VoxelTypeAir {
					totalVoxels++
				}
			}
		}
	}

	maxTriangles := totalVoxels * 6 * 2
	culledTriangles := maxTriangles - trianglesRendered

	if maxTriangles == 0 {
		return 0
	}

	return float32(culledTriangles) / float32(maxTriangles) * 100.0
}

func processInput(inputHandler *input.InputHandler, cam *camera.Camera, deltaTime float32, window *graphics.Window) bool {
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
		return true
	}

	return false
}
