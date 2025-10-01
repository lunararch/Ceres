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
	window, err := graphics.NewWindow(800, 600, "Shader System Test")
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

	vertices := []float32{
		// positions        // colors
		-0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // bottom left - red
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // bottom right - green
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // top - blue
	}

	var vao, vbo uint32
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

	projection := ceresmath.Perspective(
		ceresmath.Deg2Rad(45.0),
		window.GetAspectRatio(),
		0.1,
		100.0,
	)

	view := ceresmath.LookAt(
		ceresmath.NewVector3(0, 0, 3),
		ceresmath.NewVector3(0, 0, 0),
		ceresmath.Up(),
	)

	fmt.Println("✓ Shader system initialized")
	fmt.Println("✓ Test triangle created")
	fmt.Println("Press ESC to exit")

	rotation := float32(0)

	for !window.ShouldClose() {
		window.Clear()

		rotation += 0.01
		if rotation > 2*math.Pi {
			rotation = 0
		}

		model := ceresmath.RotateY(rotation)

		shader.Use()
		shader.SetMat4("model", model.ToPtr())
		shader.SetMat4("view", view.ToPtr())
		shader.SetMat4("projection", projection.ToPtr())

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.BindVertexArray(0)

		window.SwapBuffers()
		window.PollEvents()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			break
		}
	}

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)

	fmt.Println("✓ Test completed successfully")
}
