package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"Ceres/pkg/chunk"
)

type ChunkRenderer struct {
	meshes map[chunk.ChunkPosition]*chunk.ChunkMesh

	renderedChunks int
	renderedFaces  int
}

func NewChunkRenderer() *ChunkRenderer {
	return &ChunkRenderer{
		meshes: make(map[chunk.ChunkPosition]*chunk.ChunkMesh),
	}
}

func (cr *ChunkRenderer) UpdateChunkMesh(c *chunk.Chunk) {
	mesh := c.GenerateMesh()

	if oldMesh, exists := cr.meshes[c.Position]; exists {
		cr.DeleteMesh(oldMesh)
	}

	if mesh.IsEmpty() {
		delete(cr.meshes, c.Position)
		return
	}

	cr.UploadMesh(mesh)

	cr.meshes[c.Position] = mesh
}

func (cr *ChunkRenderer) UploadMesh(mesh *chunk.ChunkMesh) {
	gl.GenVertexArrays(1, &mesh.VAO)
	gl.GenBuffers(1, &mesh.VBO)
	gl.GenBuffers(1, &mesh.EBO)

	gl.BindVertexArray(mesh.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(mesh.Vertices)*4, gl.Ptr(mesh.Vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mesh.Indices)*4, gl.Ptr(mesh.Indices), gl.STATIC_DRAW)

	stride := int32(11 * 4)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, stride, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, stride, gl.PtrOffset(8*4))
	gl.EnableVertexAttribArray(3)

	gl.BindVertexArray(0)
}

func (cr *ChunkRenderer) DeleteMesh(mesh *chunk.ChunkMesh) {
	if mesh.VAO != 0 {
		gl.DeleteVertexArrays(1, &mesh.VAO)
		mesh.VAO = 0
	}
	if mesh.VBO != 0 {
		gl.DeleteBuffers(1, &mesh.VBO)
		mesh.VBO = 0
	}
	if mesh.EBO != 0 {
		gl.DeleteBuffers(1, &mesh.EBO)
		mesh.EBO = 0
	}
}

func (cr *ChunkRenderer) RenderChunk(chunkPos chunk.ChunkPosition) {
	mesh, exists := cr.meshes[chunkPos]
	if !exists || mesh.IsEmpty() {
		return
	}

	gl.BindVertexArray(mesh.VAO)
	gl.DrawElements(gl.TRIANGLES, int32(mesh.IndexCount), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	cr.renderedFaces += mesh.IndexCount / 3
}

func (cr *ChunkRenderer) RenderAll() {
	cr.renderedChunks = 0
	cr.renderedFaces = 0

	for chunkPos := range cr.meshes {
		cr.RenderChunk(chunkPos)
		cr.renderedChunks++
	}
}

func (cr *ChunkRenderer) GetStats() (chunks, faces int) {
	return cr.renderedChunks, cr.renderedFaces
}

func (cr *ChunkRenderer) UpdateDirtyChunks(chunkManager *chunk.ChunkManager) int {
	dirtyChunks := chunkManager.GetDirtyChunks()

	for _, c := range dirtyChunks {
		cr.UpdateChunkMesh(c)
	}

	return len(dirtyChunks)
}

func (cr *ChunkRenderer) Clear() {
	for _, mesh := range cr.meshes {
		cr.DeleteMesh(mesh)
	}
	cr.meshes = make(map[chunk.ChunkPosition]*chunk.ChunkMesh)
}

func (cr *ChunkRenderer) GetMeshCount() int {
	return len(cr.meshes)
}
