package chunk

import (
	"Ceres/pkg/voxel"
	"fmt"
)

func NewChunkPosition(x, y, z int32) ChunkPosition {
	return ChunkPosition{X: x, Y: y, Z: z}
}

func (cp ChunkPosition) Add(other ChunkPosition) ChunkPosition {
	return ChunkPosition{
		X: cp.X + other.X,
		Y: cp.Y + other.Y,
		Z: cp.Z + other.Z,
	}
}

func (cp ChunkPosition) Sub(other ChunkPosition) ChunkPosition {
	return ChunkPosition{
		X: cp.X - other.X,
		Y: cp.Y - other.Y,
		Z: cp.Z - other.Z,
	}
}

func (cp ChunkPosition) Distance(other ChunkPosition) int32 {
	dx := cp.X - other.X
	dy := cp.Y - other.Y
	dz := cp.Z - other.Z

	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	if dz < 0 {
		dz = -dz
	}

	return dx + dy + dz
}

func VoxelToChunkPosition(voxelPos voxel.VoxelPosition) ChunkPosition {
	return ChunkPosition{
		X: floorDiv(voxelPos.X, ChunkSize),
		Y: floorDiv(voxelPos.Y, ChunkSize),
		Z: floorDiv(voxelPos.Z, ChunkSize),
	}
}

func VoxelToLocalPosition(voxelPos voxel.VoxelPosition) (x, y, z int32) {
	x = mod(voxelPos.X, ChunkSize)
	y = mod(voxelPos.Y, ChunkSize)
	z = mod(voxelPos.Z, ChunkSize)
	return
}

func (cp ChunkPosition) GetNeighborPosition(face voxel.VoxelFace) ChunkPosition {
	offset := voxel.GetFaceOffset(face)
	return ChunkPosition{
		X: cp.X + offset.X,
		Y: cp.Y + offset.Y,
		Z: cp.Z + offset.Z,
	}
}

func floorDiv(a, b int32) int32 {
	if a < 0 {
		return (a - b + 1) / b
	}
	return a / b
}

func mod(a, b int32) int32 {
	result := a % b
	if result < 0 {
		result += b
	}
	return result
}

func (cp ChunkPosition) String() string {
	return fmt.Sprintf("(%d, %d, %d)", cp.X, cp.Y, cp.Z)
}
