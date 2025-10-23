package voxel

import (
	ceresmath "Ceres/pkg/math"
)

// VoxelType represents different types of voxels in the world
type VoxelType uint8

const (
	VoxelTypeAir VoxelType = iota
	VoxelTypeStone
	VoxelTypeDirt
	VoxelTypeGrass
	VoxelTypeSand
	VoxelTypeWater
	VoxelTypeWood
	VoxelTypeLeaves
	VoxelTypeGlass
	VoxelTypeBrick
)

// Voxel represents a single voxel in the world
type Voxel struct {
	Type VoxelType
}

// NewVoxel creates a new voxel with the specified type
func NewVoxel(voxelType VoxelType) Voxel {
	return Voxel{
		Type: voxelType,
	}
}

// IsAir checks if the voxel is air (empty space)
func (v Voxel) IsAir() bool {
	return v.Type == VoxelTypeAir
}

// IsSolid checks if the voxel is solid (not air or transparent)
func (v Voxel) IsSolid() bool {
	return !v.IsAir()
}

// IsTransparent checks if the voxel is transparent
func (v Voxel) IsTransparent() bool {
	return v.Type == VoxelTypeAir || v.Type == VoxelTypeWater || v.Type == VoxelTypeGlass
}

// IsOpaque checks if the voxel is opaque (blocks light and visibility)
func (v Voxel) IsOpaque() bool {
	return !v.IsTransparent()
}

// GetName returns the name of the voxel type
func (v Voxel) GetName() string {
	switch v.Type {
	case VoxelTypeAir:
		return "Air"
	case VoxelTypeStone:
		return "Stone"
	case VoxelTypeDirt:
		return "Dirt"
	case VoxelTypeGrass:
		return "Grass"
	case VoxelTypeSand:
		return "Sand"
	case VoxelTypeWater:
		return "Water"
	case VoxelTypeWood:
		return "Wood"
	case VoxelTypeLeaves:
		return "Leaves"
	case VoxelTypeGlass:
		return "Glass"
	case VoxelTypeBrick:
		return "Brick"
	default:
		return "Unknown"
	}
}

// VoxelPosition represents a 3D position in voxel space (integer coordinates)
type VoxelPosition struct {
	X, Y, Z int32
}

// NewVoxelPosition creates a new voxel position
func NewVoxelPosition(x, y, z int32) VoxelPosition {
	return VoxelPosition{X: x, Y: y, Z: z}
}

// ToWorldSpace converts voxel coordinates to world space
// Each voxel is 1x1x1 in world units
func (vp VoxelPosition) ToWorldSpace() ceresmath.Vector3 {
	return ceresmath.Vector3{
		X: float32(vp.X),
		Y: float32(vp.Y),
		Z: float32(vp.Z),
	}
}

// ToWorldSpaceCenter converts voxel coordinates to the center of the voxel in world space
func (vp VoxelPosition) ToWorldSpaceCenter() ceresmath.Vector3 {
	return ceresmath.Vector3{
		X: float32(vp.X) + 0.5,
		Y: float32(vp.Y) + 0.5,
		Z: float32(vp.Z) + 0.5,
	}
}

// Add adds two voxel positions
func (vp VoxelPosition) Add(other VoxelPosition) VoxelPosition {
	return VoxelPosition{
		X: vp.X + other.X,
		Y: vp.Y + other.Y,
		Z: vp.Z + other.Z,
	}
}

// Sub subtracts another voxel position from this one
func (vp VoxelPosition) Sub(other VoxelPosition) VoxelPosition {
	return VoxelPosition{
		X: vp.X - other.X,
		Y: vp.Y - other.Y,
		Z: vp.Z - other.Z,
	}
}

// Distance calculates the Manhattan distance between two voxel positions
func (vp VoxelPosition) Distance(other VoxelPosition) int32 {
	dx := vp.X - other.X
	dy := vp.Y - other.Y
	dz := vp.Z - other.Z
	
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

// WorldToVoxelPosition converts world space coordinates to voxel position
func WorldToVoxelPosition(worldPos ceresmath.Vector3) VoxelPosition {
	return VoxelPosition{
		X: int32(worldPos.X),
		Y: int32(worldPos.Y),
		Z: int32(worldPos.Z),
	}
}

// VoxelFace represents the six faces of a voxel
type VoxelFace int

const (
	VoxelFaceTop VoxelFace = iota
	VoxelFaceBottom
	VoxelFaceLeft
	VoxelFaceRight
	VoxelFaceFront
	VoxelFaceBack
)

// GetFaceNormal returns the normal vector for a voxel face
func GetFaceNormal(face VoxelFace) ceresmath.Vector3 {
	switch face {
	case VoxelFaceTop:
		return ceresmath.Vector3{X: 0, Y: 1, Z: 0}
	case VoxelFaceBottom:
		return ceresmath.Vector3{X: 0, Y: -1, Z: 0}
	case VoxelFaceLeft:
		return ceresmath.Vector3{X: -1, Y: 0, Z: 0}
	case VoxelFaceRight:
		return ceresmath.Vector3{X: 1, Y: 0, Z: 0}
	case VoxelFaceFront:
		return ceresmath.Vector3{X: 0, Y: 0, Z: 1}
	case VoxelFaceBack:
		return ceresmath.Vector3{X: 0, Y: 0, Z: -1}
	default:
		return ceresmath.Vector3{X: 0, Y: 0, Z: 0}
	}
}

// GetFaceOffset returns the voxel position offset for a given face
func GetFaceOffset(face VoxelFace) VoxelPosition {
	switch face {
	case VoxelFaceTop:
		return VoxelPosition{X: 0, Y: 1, Z: 0}
	case VoxelFaceBottom:
		return VoxelPosition{X: 0, Y: -1, Z: 0}
	case VoxelFaceLeft:
		return VoxelPosition{X: -1, Y: 0, Z: 0}
	case VoxelFaceRight:
		return VoxelPosition{X: 1, Y: 0, Z: 0}
	case VoxelFaceFront:
		return VoxelPosition{X: 0, Y: 0, Z: 1}
	case VoxelFaceBack:
		return VoxelPosition{X: 0, Y: 0, Z: -1}
	default:
		return VoxelPosition{X: 0, Y: 0, Z: 0}
	}
}
