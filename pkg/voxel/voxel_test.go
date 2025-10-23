package voxel

import (
	"testing"

	ceresmath "Ceres/pkg/math"
)

func TestNewVoxel(t *testing.T) {
	voxel := NewVoxel(VoxelTypeStone)
	if voxel.Type != VoxelTypeStone {
		t.Errorf("Expected voxel type Stone, got %v", voxel.Type)
	}
}

func TestVoxelIsAir(t *testing.T) {
	air := NewVoxel(VoxelTypeAir)
	if !air.IsAir() {
		t.Error("Air voxel should return true for IsAir()")
	}

	stone := NewVoxel(VoxelTypeStone)
	if stone.IsAir() {
		t.Error("Stone voxel should return false for IsAir()")
	}
}

func TestVoxelIsSolid(t *testing.T) {
	air := NewVoxel(VoxelTypeAir)
	if air.IsSolid() {
		t.Error("Air voxel should not be solid")
	}

	stone := NewVoxel(VoxelTypeStone)
	if !stone.IsSolid() {
		t.Error("Stone voxel should be solid")
	}
}

func TestVoxelTransparency(t *testing.T) {
	tests := []struct {
		voxelType    VoxelType
		transparent  bool
	}{
		{VoxelTypeAir, true},
		{VoxelTypeWater, true},
		{VoxelTypeGlass, true},
		{VoxelTypeStone, false},
		{VoxelTypeDirt, false},
	}

	for _, tt := range tests {
		voxel := NewVoxel(tt.voxelType)
		if voxel.IsTransparent() != tt.transparent {
			t.Errorf("Voxel type %v: expected transparent=%v, got %v",
				tt.voxelType, tt.transparent, voxel.IsTransparent())
		}
		if voxel.IsOpaque() == tt.transparent {
			t.Errorf("Voxel type %v: IsOpaque should be opposite of IsTransparent", tt.voxelType)
		}
	}
}

func TestVoxelGetName(t *testing.T) {
	tests := []struct {
		voxelType VoxelType
		name      string
	}{
		{VoxelTypeAir, "Air"},
		{VoxelTypeStone, "Stone"},
		{VoxelTypeDirt, "Dirt"},
		{VoxelTypeGrass, "Grass"},
		{VoxelTypeSand, "Sand"},
		{VoxelTypeWater, "Water"},
		{VoxelTypeWood, "Wood"},
		{VoxelTypeLeaves, "Leaves"},
		{VoxelTypeGlass, "Glass"},
		{VoxelTypeBrick, "Brick"},
	}

	for _, tt := range tests {
		voxel := NewVoxel(tt.voxelType)
		if voxel.GetName() != tt.name {
			t.Errorf("Expected name %s for voxel type %v, got %s",
				tt.name, tt.voxelType, voxel.GetName())
		}
	}
}

func TestVoxelPosition(t *testing.T) {
	pos := NewVoxelPosition(10, 20, 30)
	if pos.X != 10 || pos.Y != 20 || pos.Z != 30 {
		t.Errorf("Expected position (10, 20, 30), got (%d, %d, %d)", pos.X, pos.Y, pos.Z)
	}
}

func TestVoxelPositionToWorldSpace(t *testing.T) {
	pos := NewVoxelPosition(5, 10, 15)
	worldPos := pos.ToWorldSpace()

	if worldPos.X != 5.0 || worldPos.Y != 10.0 || worldPos.Z != 15.0 {
		t.Errorf("Expected world position (5.0, 10.0, 15.0), got (%f, %f, %f)",
			worldPos.X, worldPos.Y, worldPos.Z)
	}
}

func TestVoxelPositionToWorldSpaceCenter(t *testing.T) {
	pos := NewVoxelPosition(5, 10, 15)
	worldPos := pos.ToWorldSpaceCenter()

	if worldPos.X != 5.5 || worldPos.Y != 10.5 || worldPos.Z != 15.5 {
		t.Errorf("Expected world position (5.5, 10.5, 15.5), got (%f, %f, %f)",
			worldPos.X, worldPos.Y, worldPos.Z)
	}
}

func TestVoxelPositionAdd(t *testing.T) {
	pos1 := NewVoxelPosition(1, 2, 3)
	pos2 := NewVoxelPosition(4, 5, 6)
	result := pos1.Add(pos2)

	if result.X != 5 || result.Y != 7 || result.Z != 9 {
		t.Errorf("Expected position (5, 7, 9), got (%d, %d, %d)", result.X, result.Y, result.Z)
	}
}

func TestVoxelPositionSub(t *testing.T) {
	pos1 := NewVoxelPosition(10, 20, 30)
	pos2 := NewVoxelPosition(3, 5, 7)
	result := pos1.Sub(pos2)

	if result.X != 7 || result.Y != 15 || result.Z != 23 {
		t.Errorf("Expected position (7, 15, 23), got (%d, %d, %d)", result.X, result.Y, result.Z)
	}
}

func TestVoxelPositionDistance(t *testing.T) {
	pos1 := NewVoxelPosition(0, 0, 0)
	pos2 := NewVoxelPosition(3, 4, 5)
	distance := pos1.Distance(pos2)

	// Manhattan distance: |3-0| + |4-0| + |5-0| = 12
	if distance != 12 {
		t.Errorf("Expected distance 12, got %d", distance)
	}

	// Test with negative coordinates
	pos3 := NewVoxelPosition(-2, -3, -4)
	distance2 := pos1.Distance(pos3)
	if distance2 != 9 {
		t.Errorf("Expected distance 9, got %d", distance2)
	}
}

func TestWorldToVoxelPosition(t *testing.T) {
	worldPos := ceresmath.Vector3{X: 10.7, Y: 20.3, Z: 30.9}
	voxelPos := WorldToVoxelPosition(worldPos)

	// Should truncate to integer coordinates
	if voxelPos.X != 10 || voxelPos.Y != 20 || voxelPos.Z != 30 {
		t.Errorf("Expected voxel position (10, 20, 30), got (%d, %d, %d)",
			voxelPos.X, voxelPos.Y, voxelPos.Z)
	}
}

func TestGetFaceNormal(t *testing.T) {
	tests := []struct {
		face   VoxelFace
		normal ceresmath.Vector3
	}{
		{VoxelFaceTop, ceresmath.Vector3{X: 0, Y: 1, Z: 0}},
		{VoxelFaceBottom, ceresmath.Vector3{X: 0, Y: -1, Z: 0}},
		{VoxelFaceLeft, ceresmath.Vector3{X: -1, Y: 0, Z: 0}},
		{VoxelFaceRight, ceresmath.Vector3{X: 1, Y: 0, Z: 0}},
		{VoxelFaceFront, ceresmath.Vector3{X: 0, Y: 0, Z: 1}},
		{VoxelFaceBack, ceresmath.Vector3{X: 0, Y: 0, Z: -1}},
	}

	for _, tt := range tests {
		normal := GetFaceNormal(tt.face)
		if normal.X != tt.normal.X || normal.Y != tt.normal.Y || normal.Z != tt.normal.Z {
			t.Errorf("Face %d: expected normal (%f, %f, %f), got (%f, %f, %f)",
				tt.face, tt.normal.X, tt.normal.Y, tt.normal.Z, normal.X, normal.Y, normal.Z)
		}
	}
}

func TestGetFaceOffset(t *testing.T) {
	tests := []struct {
		face   VoxelFace
		offset VoxelPosition
	}{
		{VoxelFaceTop, VoxelPosition{X: 0, Y: 1, Z: 0}},
		{VoxelFaceBottom, VoxelPosition{X: 0, Y: -1, Z: 0}},
		{VoxelFaceLeft, VoxelPosition{X: -1, Y: 0, Z: 0}},
		{VoxelFaceRight, VoxelPosition{X: 1, Y: 0, Z: 0}},
		{VoxelFaceFront, VoxelPosition{X: 0, Y: 0, Z: 1}},
		{VoxelFaceBack, VoxelPosition{X: 0, Y: 0, Z: -1}},
	}

	for _, tt := range tests {
		offset := GetFaceOffset(tt.face)
		if offset.X != tt.offset.X || offset.Y != tt.offset.Y || offset.Z != tt.offset.Z {
			t.Errorf("Face %d: expected offset (%d, %d, %d), got (%d, %d, %d)",
				tt.face, tt.offset.X, tt.offset.Y, tt.offset.Z, offset.X, offset.Y, offset.Z)
		}
	}
}
