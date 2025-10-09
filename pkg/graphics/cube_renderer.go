package graphics

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"

	"Ceres/pkg/mesh"
)

type RenderMode int

const (
	RenderModeSolid RenderMode = iota
	RenderModeWireframe
	RenderModeBoth
)

type CubeRenderer struct {
	vao        uint32
	vbo        uint32
	ebo        uint32
	indexCount int32
	renderMode RenderMode
}

func NewCubeRenderer(cubeMesh *mesh.CubeMesh) *CubeRenderer {
	var vao, vbo, ebo uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	vertexData := cubeMesh.ToFloatArray()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(cubeMesh.Indices)*4, gl.Ptr(cubeMesh.Indices), gl.STATIC_DRAW)

	stride := cubeMesh.GetVertexStride()

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, unsafe.Pointer(uintptr(3*4)))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, stride, unsafe.Pointer(uintptr(6*4)))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	return &CubeRenderer{
		vao:        vao,
		vbo:        vbo,
		ebo:        ebo,
		indexCount: int32(len(cubeMesh.Indices)),
		renderMode: RenderModeSolid,
	}
}

func (cr *CubeRenderer) Render() {
	gl.BindVertexArray(cr.vao)

	switch cr.renderMode {
	case RenderModeSolid:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		gl.DrawElements(gl.TRIANGLES, cr.indexCount, gl.UNSIGNED_INT, nil)
	case RenderModeWireframe:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		gl.DrawElements(gl.TRIANGLES, cr.indexCount, gl.UNSIGNED_INT, nil)
	case RenderModeBoth:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		gl.DrawElements(gl.TRIANGLES, cr.indexCount, gl.UNSIGNED_INT, nil)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		gl.DrawElements(gl.TRIANGLES, cr.indexCount, gl.UNSIGNED_INT, nil)
	}

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	gl.BindVertexArray(0)
}

func (cr *CubeRenderer) SetRenderMode(mode RenderMode) {
	cr.renderMode = mode
}

func (cr *CubeRenderer) GetRenderMode() RenderMode {
	return cr.renderMode
}

func (cr *CubeRenderer) Delete() {
	gl.DeleteVertexArrays(1, &cr.vao)
	gl.DeleteBuffers(1, &cr.vbo)
	gl.DeleteBuffers(1, &cr.ebo)
}
