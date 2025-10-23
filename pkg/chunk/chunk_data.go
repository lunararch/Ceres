package chunk

import (
	"sync"

	"Ceres/pkg/voxel"
)

const ChunkSize = 32

type Chunk struct {
	Position ChunkPosition

	voxels [ChunkSize * ChunkSize * ChunkSize]voxel.Voxel

	neighbors [6]*Chunk

	isDirty    bool
	isModified bool
	isEmpty    bool

	mutex sync.RWMutex
}

type ChunkPosition struct {
	X, Y, Z int32
}

// NewChunk creates a new empty chunk at the given position
func NewChunk(position ChunkPosition) *Chunk {
	chunk := &Chunk{
		Position:   position,
		isDirty:    true,
		isEmpty:    true,
		isModified: false,
	}

	// Initialize with air
	for i := range chunk.voxels {
		chunk.voxels[i] = voxel.NewVoxel(voxel.VoxelTypeAir)
	}

	return chunk
}

func (c *Chunk) GetVoxel(x, y, z int32) voxel.Voxel {
	if !isValidLocalCoord(x, y, z) {
		return voxel.NewVoxel(voxel.VoxelTypeAir)
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	index := localToIndex(x, y, z)
	return c.voxels[index]
}

func (c *Chunk) SetVoxel(x, y, z int32, v voxel.Voxel) {
	if !isValidLocalCoord(x, y, z) {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	index := localToIndex(x, y, z)
	oldVoxel := c.voxels[index]

	if oldVoxel.Type != v.Type {
		c.voxels[index] = v
		c.isDirty = true
		c.isModified = true

		// Update isEmpty flag
		if !v.IsAir() {
			c.isEmpty = false
		} else {
			c.checkIfEmpty()
		}
	}
}

func (c *Chunk) GetVoxelSafe(x, y, z int32) voxel.Voxel {
	// If within bounds, get from this chunk
	if isValidLocalCoord(x, y, z) {
		return c.GetVoxel(x, y, z)
	}

	// Otherwise, check neighbors
	var neighborIndex int
	var nx, ny, nz int32

	if x < 0 {
		neighborIndex = int(voxel.VoxelFaceLeft)
		nx, ny, nz = ChunkSize-1, y, z
	} else if x >= ChunkSize {
		neighborIndex = int(voxel.VoxelFaceRight)
		nx, ny, nz = 0, y, z
	} else if y < 0 {
		neighborIndex = int(voxel.VoxelFaceBottom)
		nx, ny, nz = x, ChunkSize-1, z
	} else if y >= ChunkSize {
		neighborIndex = int(voxel.VoxelFaceTop)
		nx, ny, nz = x, 0, z
	} else if z < 0 {
		neighborIndex = int(voxel.VoxelFaceBack)
		nx, ny, nz = x, y, ChunkSize-1
	} else {
		neighborIndex = int(voxel.VoxelFaceFront)
		nx, ny, nz = x, y, 0
	}

	c.mutex.RLock()
	neighbor := c.neighbors[neighborIndex]
	c.mutex.RUnlock()

	if neighbor != nil {
		return neighbor.GetVoxel(nx, ny, nz)
	}

	return voxel.NewVoxel(voxel.VoxelTypeAir)
}

func (c *Chunk) SetNeighbor(face voxel.VoxelFace, neighbor *Chunk) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.neighbors[face] = neighbor
	c.isDirty = true
}

func (c *Chunk) GetNeighbor(face voxel.VoxelFace) *Chunk {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.neighbors[face]
}

func (c *Chunk) IsDirty() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.isDirty
}

func (c *Chunk) SetDirty(dirty bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.isDirty = dirty
}

func (c *Chunk) IsEmpty() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.isEmpty
}

func (c *Chunk) checkIfEmpty() {
	for _, v := range c.voxels {
		if !v.IsAir() {
			c.isEmpty = false
			return
		}
	}
	c.isEmpty = true
}

func (c *Chunk) Fill(voxelType voxel.VoxelType) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	v := voxel.NewVoxel(voxelType)
	for i := range c.voxels {
		c.voxels[i] = v
	}

	c.isEmpty = voxelType == voxel.VoxelTypeAir
	c.isDirty = true
	c.isModified = true
}

func (c *Chunk) GetWorldPosition() voxel.VoxelPosition {
	return voxel.NewVoxelPosition(
		c.Position.X*ChunkSize,
		c.Position.Y*ChunkSize,
		c.Position.Z*ChunkSize,
	)
}

func isValidLocalCoord(x, y, z int32) bool {
	return x >= 0 && x < ChunkSize &&
		y >= 0 && y < ChunkSize &&
		z >= 0 && z < ChunkSize
}

func localToIndex(x, y, z int32) int {
	return int(x + y*ChunkSize + z*ChunkSize*ChunkSize)
}

func indexToLocal(index int) (x, y, z int32) {
	z = int32(index / (ChunkSize * ChunkSize))
	index -= int(z * ChunkSize * ChunkSize)
	y = int32(index / ChunkSize)
	x = int32(index % ChunkSize)
	return
}
