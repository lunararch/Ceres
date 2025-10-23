package mesh

import (
	ceresmath "Ceres/pkg/math"
	"Ceres/pkg/voxel"
)

type VoxelMeshBuilder struct {
	vertices []float32
	indices  []uint32
}

func NewVoxelMeshBuilder() *VoxelMeshBuilder {
	return &VoxelMeshBuilder{
		vertices: make([]float32, 0, 4096),
		indices:  make([]uint32, 0, 2048),
	}
}

type VoxelNeighbors struct {
	Top    voxel.Voxel
	Bottom voxel.Voxel
	Left   voxel.Voxel
	Right  voxel.Voxel
	Front  voxel.Voxel
	Back   voxel.Voxel
}

func ShouldRenderFace(currentVoxel voxel.Voxel, neighborVoxel voxel.Voxel) bool {
	if currentVoxel.IsAir() {
		return false
	}

	if neighborVoxel.IsAir() {
		return true
	}

	if neighborVoxel.IsTransparent() && currentVoxel.IsOpaque() {
		return true
	}

	return false
}

func (vmb *VoxelMeshBuilder) AddVoxelWithCulling(
	pos voxel.VoxelPosition,
	voxelType voxel.VoxelType,
	neighbors VoxelNeighbors,
	color ceresmath.Vector3,
) {
	currentVoxel := voxel.NewVoxel(voxelType)
	worldPos := pos.ToWorldSpace()

	if ShouldRenderFace(currentVoxel, neighbors.Top) {
		vmb.addTopFace(worldPos, color)
	}

	if ShouldRenderFace(currentVoxel, neighbors.Bottom) {
		vmb.addBottomFace(worldPos, color)
	}

	if ShouldRenderFace(currentVoxel, neighbors.Left) {
		vmb.addLeftFace(worldPos, color)
	}

	if ShouldRenderFace(currentVoxel, neighbors.Right) {
		vmb.addRightFace(worldPos, color)
	}

	if ShouldRenderFace(currentVoxel, neighbors.Front) {
		vmb.addFrontFace(worldPos, color)
	}

	if ShouldRenderFace(currentVoxel, neighbors.Back) {
		vmb.addBackFace(worldPos, color)
	}
}

func (vmb *VoxelMeshBuilder) addTopFace(pos ceresmath.Vector3, color ceresmath.Vector3) {
	baseIdx := uint32(len(vmb.vertices) / 9)

	vmb.addVertex(pos.X, pos.Y+1, pos.Z, 0, 1, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y+1, pos.Z, 0, 1, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y+1, pos.Z+1, 0, 1, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X, pos.Y+1, pos.Z+1, 0, 1, 0, color.X, color.Y, color.Z)

	vmb.addQuadIndices(baseIdx)
}

func (vmb *VoxelMeshBuilder) addBottomFace(pos ceresmath.Vector3, color ceresmath.Vector3) {
	baseIdx := uint32(len(vmb.vertices) / 9)

	vmb.addVertex(pos.X, pos.Y, pos.Z+1, 0, -1, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y, pos.Z+1, 0, -1, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y, pos.Z, 0, -1, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X, pos.Y, pos.Z, 0, -1, 0, color.X, color.Y, color.Z)

	vmb.addQuadIndices(baseIdx)
}

func (vmb *VoxelMeshBuilder) addLeftFace(pos ceresmath.Vector3, color ceresmath.Vector3) {
	baseIdx := uint32(len(vmb.vertices) / 9)

	vmb.addVertex(pos.X, pos.Y, pos.Z, -1, 0, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X, pos.Y+1, pos.Z, -1, 0, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X, pos.Y+1, pos.Z+1, -1, 0, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X, pos.Y, pos.Z+1, -1, 0, 0, color.X, color.Y, color.Z)

	vmb.addQuadIndices(baseIdx)
}

func (vmb *VoxelMeshBuilder) addRightFace(pos ceresmath.Vector3, color ceresmath.Vector3) {
	baseIdx := uint32(len(vmb.vertices) / 9)

	vmb.addVertex(pos.X+1, pos.Y, pos.Z+1, 1, 0, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y+1, pos.Z+1, 1, 0, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y+1, pos.Z, 1, 0, 0, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y, pos.Z, 1, 0, 0, color.X, color.Y, color.Z)

	vmb.addQuadIndices(baseIdx)
}

func (vmb *VoxelMeshBuilder) addFrontFace(pos ceresmath.Vector3, color ceresmath.Vector3) {
	baseIdx := uint32(len(vmb.vertices) / 9)

	vmb.addVertex(pos.X, pos.Y, pos.Z+1, 0, 0, 1, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X, pos.Y+1, pos.Z+1, 0, 0, 1, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y+1, pos.Z+1, 0, 0, 1, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y, pos.Z+1, 0, 0, 1, color.X, color.Y, color.Z)

	vmb.addQuadIndices(baseIdx)
}

func (vmb *VoxelMeshBuilder) addBackFace(pos ceresmath.Vector3, color ceresmath.Vector3) {
	baseIdx := uint32(len(vmb.vertices) / 9)

	vmb.addVertex(pos.X+1, pos.Y, pos.Z, 0, 0, -1, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X+1, pos.Y+1, pos.Z, 0, 0, -1, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X, pos.Y+1, pos.Z, 0, 0, -1, color.X, color.Y, color.Z)
	vmb.addVertex(pos.X, pos.Y, pos.Z, 0, 0, -1, color.X, color.Y, color.Z)

	vmb.addQuadIndices(baseIdx)
}

func (vmb *VoxelMeshBuilder) addVertex(x, y, z, nx, ny, nz, r, g, b float32) {
	vmb.vertices = append(vmb.vertices,
		x, y, z, 
		nx, ny, nz,
		r, g, b,
	)
}

func (vmb *VoxelMeshBuilder) addQuadIndices(baseIdx uint32) {
	vmb.indices = append(vmb.indices,
		baseIdx, baseIdx+1, baseIdx+2,
		baseIdx, baseIdx+2, baseIdx+3,
	)
}

func (vmb *VoxelMeshBuilder) Reset() {
	vmb.vertices = vmb.vertices[:0]
	vmb.indices = vmb.indices[:0]
}

func (vmb *VoxelMeshBuilder) GetVertices() []float32 {
	return vmb.vertices
}

func (vmb *VoxelMeshBuilder) GetIndices() []uint32 {
	return vmb.indices
}

func (vmb *VoxelMeshBuilder) GetVertexCount() int {
	return len(vmb.vertices) / 9
}

func (vmb *VoxelMeshBuilder) GetTriangleCount() int {
	return len(vmb.indices) / 3
}
