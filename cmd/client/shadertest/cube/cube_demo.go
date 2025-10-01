package main

import (
	"fmt"
	"log"
	"math"

	"Ceres/pkg/graphics"
	ceresmath "Ceres/pkg/math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	window, err := graphics.NewWindow(1280, 720, "Cube with Lighting Test")
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

	cubeVertices := []float32{
		// Front face
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,

		// Back face
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,

		// Top face
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,

		// Bottom face
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,

		// Right face
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 0.0,

		// Left face
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
	}

	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	// Position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Normal
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// Texture coordinates
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	projection := ceresmath.Perspective(ceresmath.Deg2Rad(45.0), window.GetAspectRatio(), 0.1, 100.0)
	view := ceresmath.LookAt(
		ceresmath.NewVector3(3, 3, 3),
		ceresmath.NewVector3(0, 0, 0),
		ceresmath.Up(),
	)

	fmt.Println("✓ Cube with lighting initialized")
	fmt.Println("Press ESC to exit")

	rotation := float32(0)

	for !window.ShouldClose() {
		window.Clear()

		rotation += 0.005
		if rotation > 2*math.Pi {
			rotation = 0
		}

		model := ceresmath.RotateY(rotation).Mul(ceresmath.RotateX(rotation * 0.5))

		shader.Use()
		shader.SetMat4("model", model.ToPtr())
		shader.SetMat4("view", view.ToPtr())
		shader.SetMat4("projection", projection.ToPtr())

		shader.SetVec3("lightPos", 5.0, 5.0, 5.0)
		shader.SetVec3("viewPos", 3.0, 3.0, 3.0)
		shader.SetVec3("lightColor", 1.0, 1.0, 1.0)
		shader.SetVec3("objectColor", 0.2, 0.7, 1.0)
		shader.SetInt("useTexture", 0)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		gl.BindVertexArray(0)

		window.SwapBuffers()
		window.PollEvents()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			break
		}
	}

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)

	fmt.Println("✓ Cube test completed")
}
