package mesh

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Mesh struct {
	VAO           uint32
	VBO           uint32
	EBO           uint32
	IndexCount    int32
	VertexCount   int32
	TriangleCount int32
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.VAO)
	gl.DrawElements(gl.TRIANGLES, m.IndexCount, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

func (m *Mesh) Delete() {
	gl.DeleteVertexArrays(1, &m.VAO)
	gl.DeleteBuffers(1, &m.VBO)
	gl.DeleteBuffers(1, &m.EBO)
}

func (m *Mesh) GetVertexCount() int32 {
	return m.VertexCount
}

func (m *Mesh) GetTriangleCount() int32 {
	return m.TriangleCount
}
