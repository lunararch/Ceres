package chunk

import (
	"Ceres/pkg/voxel"
)

type ChunkMesh struct {
	Vertices []float32

	Indices []uint32

	VertexCount int

	IndexCount int

	VAO uint32
	VBO uint32
	EBO uint32
}

func NewChunkMesh() *ChunkMesh {
	return &ChunkMesh{
		Vertices: make([]float32, 0, 4096),
		Indices:  make([]uint32, 0, 4096),
	}
}

func (cm *ChunkMesh) AddFace(position voxel.VoxelPosition, face voxel.VoxelFace, voxelType voxel.VoxelType) {
	normal := voxel.GetFaceNormal(face)
	vertices := getFaceVertices(face)

	color := getVoxelTypeColor(voxelType)

	baseIndex := uint32(len(cm.Vertices) / 11)

	for i := 0; i < 4; i++ {
		cm.Vertices = append(cm.Vertices,
			float32(position.X)+vertices[i][0],
			float32(position.Y)+vertices[i][1],
			float32(position.Z)+vertices[i][2],
		)

		cm.Vertices = append(cm.Vertices,
			float32(normal.X),
			float32(normal.Y),
			float32(normal.Z),
		)

		cm.Vertices = append(cm.Vertices,
			vertices[i][3],
			vertices[i][4],
		)

		cm.Vertices = append(cm.Vertices,
			color[0],
			color[1],
			color[2],
		)
	}

	cm.Indices = append(cm.Indices,
		baseIndex+0, baseIndex+1, baseIndex+2,
		baseIndex+2, baseIndex+3, baseIndex+0,
	)

	cm.VertexCount += 4
	cm.IndexCount += 6
}

func (cm *ChunkMesh) Clear() {
	cm.Vertices = cm.Vertices[:0]
	cm.Indices = cm.Indices[:0]
	cm.VertexCount = 0
	cm.IndexCount = 0
}

func (cm *ChunkMesh) IsEmpty() bool {
	return cm.IndexCount == 0
}

func getFaceVertices(face voxel.VoxelFace) [4][5]float32 {
	switch face {
	case voxel.VoxelFaceTop: // +Y
		return [4][5]float32{
			{0, 1, 0, 0, 0},
			{1, 1, 0, 1, 0},
			{1, 1, 1, 1, 1},
			{0, 1, 1, 0, 1},
		}
	case voxel.VoxelFaceBottom: // -Y
		return [4][5]float32{
			{0, 0, 1, 0, 0},
			{1, 0, 1, 1, 0},
			{1, 0, 0, 1, 1},
			{0, 0, 0, 0, 1},
		}
	case voxel.VoxelFaceLeft: // -X
		return [4][5]float32{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 1, 1, 1, 1},
			{0, 1, 0, 0, 1},
		}
	case voxel.VoxelFaceRight: // +X
		return [4][5]float32{
			{1, 0, 1, 0, 0},
			{1, 0, 0, 1, 0},
			{1, 1, 0, 1, 1},
			{1, 1, 1, 0, 1},
		}
	case voxel.VoxelFaceFront: // +Z
		return [4][5]float32{
			{0, 0, 1, 0, 0},
			{1, 0, 1, 1, 0},
			{1, 1, 1, 1, 1},
			{0, 1, 1, 0, 1},
		}
	case voxel.VoxelFaceBack: // -Z
		return [4][5]float32{
			{1, 0, 0, 0, 0},
			{0, 0, 0, 1, 0},
			{0, 1, 0, 1, 1},
			{1, 1, 0, 0, 1},
		}
	default:
		return [4][5]float32{}
	}
}

func getVoxelTypeColor(voxelType voxel.VoxelType) [3]float32 {
	switch voxelType {
	case voxel.VoxelTypeStone:
		return [3]float32{0.5, 0.5, 0.5} // Gray
	case voxel.VoxelTypeDirt:
		return [3]float32{0.55, 0.35, 0.2} // Brown
	case voxel.VoxelTypeGrass:
		return [3]float32{0.2, 0.8, 0.2} // Green
	case voxel.VoxelTypeSand:
		return [3]float32{0.95, 0.9, 0.6} // Sand yellow
	case voxel.VoxelTypeWater:
		return [3]float32{0.2, 0.4, 0.9} // Blue
	case voxel.VoxelTypeWood:
		return [3]float32{0.6, 0.4, 0.2} // Wood brown
	case voxel.VoxelTypeLeaves:
		return [3]float32{0.15, 0.6, 0.15} // Dark green
	case voxel.VoxelTypeGlass:
		return [3]float32{0.7, 0.9, 1.0} // Light blue
	case voxel.VoxelTypeBrick:
		return [3]float32{0.7, 0.3, 0.2} // Red brick
	default:
		return [3]float32{1.0, 1.0, 1.0} // White
	}
}

func (c *Chunk) GenerateMesh() *ChunkMesh {
	mesh := NewChunkMesh()

	if c.IsEmpty() {
		return mesh
	}

	for x := int32(0); x < ChunkSize; x++ {
		for y := int32(0); y < ChunkSize; y++ {
			for z := int32(0); z < ChunkSize; z++ {
				currentVoxel := c.GetVoxel(x, y, z)

				if currentVoxel.IsAir() {
					continue
				}

				for face := voxel.VoxelFaceTop; face <= voxel.VoxelFaceBack; face++ {
					if c.isFaceVisible(x, y, z, face) {
						position := c.GetWorldPosition().Add(voxel.NewVoxelPosition(x, y, z))
						mesh.AddFace(position, face, currentVoxel.Type)
					}
				}
			}
		}
	}

	c.SetDirty(false)

	return mesh
}

func (c *Chunk) isFaceVisible(x, y, z int32, face voxel.VoxelFace) bool {
	offset := voxel.GetFaceOffset(face)

	nx := x + offset.X
	ny := y + offset.Y
	nz := z + offset.Z

	neighbor := c.GetVoxelSafe(nx, ny, nz)

	return neighbor.IsAir() || neighbor.IsTransparent()
}
