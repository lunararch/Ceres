package mesh

import (
	ceresmath "Ceres/pkg/math"
)

type Vertex struct {
	Position ceresmath.Vector3
	Normal   ceresmath.Vector3
	TexCoord [2]float32
}

type CubeMesh struct {
	Vertices []Vertex
	Indices  []uint32
}

func NewCubeMesh(size float32) *CubeMesh {
	halfSize := size / 2.0

	vertices := []Vertex{
		// Front face (Z+)
        {ceresmath.NewVector3(-halfSize, -halfSize, halfSize), ceresmath.NewVector3(0, 0, 1), [2]float32{0, 0}},
        {ceresmath.NewVector3(halfSize, -halfSize, halfSize), ceresmath.NewVector3(0, 0, 1), [2]float32{1, 0}},
        {ceresmath.NewVector3(halfSize, halfSize, halfSize), ceresmath.NewVector3(0, 0, 1), [2]float32{1, 1}},
        {ceresmath.NewVector3(-halfSize, halfSize, halfSize), ceresmath.NewVector3(0, 0, 1), [2]float32{0, 1}},

        // Back face (Z-)
        {ceresmath.NewVector3(halfSize, -halfSize, -halfSize), ceresmath.NewVector3(0, 0, -1), [2]float32{0, 0}},
        {ceresmath.NewVector3(-halfSize, -halfSize, -halfSize), ceresmath.NewVector3(0, 0, -1), [2]float32{1, 0}},
        {ceresmath.NewVector3(-halfSize, halfSize, -halfSize), ceresmath.NewVector3(0, 0, -1), [2]float32{1, 1}},
        {ceresmath.NewVector3(halfSize, halfSize, -halfSize), ceresmath.NewVector3(0, 0, -1), [2]float32{0, 1}},

        // Top face (Y+)
        {ceresmath.NewVector3(-halfSize, halfSize, halfSize), ceresmath.NewVector3(0, 1, 0), [2]float32{0, 0}},
        {ceresmath.NewVector3(halfSize, halfSize, halfSize), ceresmath.NewVector3(0, 1, 0), [2]float32{1, 0}},
        {ceresmath.NewVector3(halfSize, halfSize, -halfSize), ceresmath.NewVector3(0, 1, 0), [2]float32{1, 1}},
        {ceresmath.NewVector3(-halfSize, halfSize, -halfSize), ceresmath.NewVector3(0, 1, 0), [2]float32{0, 1}},

        // Bottom face (Y-)
        {ceresmath.NewVector3(-halfSize, -halfSize, -halfSize), ceresmath.NewVector3(0, -1, 0), [2]float32{0, 0}},
        {ceresmath.NewVector3(halfSize, -halfSize, -halfSize), ceresmath.NewVector3(0, -1, 0), [2]float32{1, 0}},
        {ceresmath.NewVector3(halfSize, -halfSize, halfSize), ceresmath.NewVector3(0, -1, 0), [2]float32{1, 1}},
        {ceresmath.NewVector3(-halfSize, -halfSize, halfSize), ceresmath.NewVector3(0, -1, 0), [2]float32{0, 1}},

        // Right face (X+)
        {ceresmath.NewVector3(halfSize, -halfSize, halfSize), ceresmath.NewVector3(1, 0, 0), [2]float32{0, 0}},
        {ceresmath.NewVector3(halfSize, -halfSize, -halfSize), ceresmath.NewVector3(1, 0, 0), [2]float32{1, 0}},
        {ceresmath.NewVector3(halfSize, halfSize, -halfSize), ceresmath.NewVector3(1, 0, 0), [2]float32{1, 1}},
        {ceresmath.NewVector3(halfSize, halfSize, halfSize), ceresmath.NewVector3(1, 0, 0), [2]float32{0, 1}},

        // Left face (X-)
        {ceresmath.NewVector3(-halfSize, -halfSize, -halfSize), ceresmath.NewVector3(-1, 0, 0), [2]float32{0, 0}},
        {ceresmath.NewVector3(-halfSize, -halfSize, halfSize), ceresmath.NewVector3(-1, 0, 0), [2]float32{1, 0}},
        {ceresmath.NewVector3(-halfSize, halfSize, halfSize), ceresmath.NewVector3(-1, 0, 0), [2]float32{1, 1}},
        {ceresmath.NewVector3(-halfSize, halfSize, -halfSize), ceresmath.NewVector3(-1, 0, 0), [2]float32{0, 1}},
	}

	// Indices for indexed drawing (2 triangles per face, 6 faces)
	indices := []uint32{
		0, 1, 2, 2, 3, 0,       // Front face
		4, 5, 6, 6, 7, 4,       // Back face
		8, 9, 10, 10, 11, 8,    // Top face
		12, 13, 14, 14, 15, 12,  // Bottom face
		16, 17, 18, 18, 19, 16,  // Right face
		20, 21, 22, 22, 23, 20,  // Left face
	}

	return &CubeMesh{
		Vertices: vertices,
		Indices:  indices,
	}
}

func (cm *CubeMesh) ToFloatArray() []float32 {
	data := make([]float32, 0, len(cm.Vertices)*8)

	for _, v := range cm.Vertices{
		data = append(data, v.Position.X, v.Position.Y, v.Position.Z)
		data = append(data, v.Normal.X, v.Normal.Y, v.Normal.Z)
		data = append(data, v.TexCoord[0], v.TexCoord[1])
	}

	return data
}

func (cm *CubeMesh) GetVertexStride() int32 {
	return 8 * 4
}