package chunk

import (
	"fmt"
	"sync"

	"Ceres/pkg/voxel"
)

type ChunkManager struct {
	chunks map[ChunkPosition]*Chunk
	mutex  sync.RWMutex

	// Statistics
	totalChunks  int
	loadedChunks int
	dirtyChunks  int
}

func NewChunkManager() *ChunkManager {
	return &ChunkManager{
		chunks: make(map[ChunkPosition]*Chunk),
	}
}

func (cm *ChunkManager) GetChunk(pos ChunkPosition) *Chunk {
	cm.mutex.RLock()
	chunk, exists := cm.chunks[pos]
	cm.mutex.RUnlock()

	if exists {
		return chunk
	}

	return cm.CreateChunk(pos)
}

func (cm *ChunkManager) CreateChunk(pos ChunkPosition) *Chunk {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if chunk, exists := cm.chunks[pos]; exists {
		return chunk
	}

	chunk := NewChunk(pos)
	cm.chunks[pos] = chunk
	cm.totalChunks++
	cm.loadedChunks++

	cm.setupNeighbors(chunk)

	return chunk
}

func (cm *ChunkManager) setupNeighbors(chunk *Chunk) {
	for face := voxel.VoxelFaceTop; face <= voxel.VoxelFaceBack; face++ {
		neighborPos := chunk.Position.GetNeighborPosition(face)

		if neighbor, exists := cm.chunks[neighborPos]; exists {
			chunk.SetNeighbor(face, neighbor)

			oppositeFace := getOppositeFace(face)
			neighbor.SetNeighbor(oppositeFace, chunk)
		}
	}
}

func (cm *ChunkManager) GetVoxel(voxelPos voxel.VoxelPosition) voxel.Voxel {
	chunkPos := VoxelToChunkPosition(voxelPos)
	chunk := cm.GetChunk(chunkPos)

	x, y, z := VoxelToLocalPosition(voxelPos)
	return chunk.GetVoxel(x, y, z)
}

func (cm *ChunkManager) SetVoxel(voxelPos voxel.VoxelPosition, v voxel.Voxel) {
	chunkPos := VoxelToChunkPosition(voxelPos)
	chunk := cm.GetChunk(chunkPos)

	x, y, z := VoxelToLocalPosition(voxelPos)
	chunk.SetVoxel(x, y, z, v)

	cm.markAdjacentChunksDirty(voxelPos)
}

func (cm *ChunkManager) markAdjacentChunksDirty(voxelPos voxel.VoxelPosition) {
	x, y, z := VoxelToLocalPosition(voxelPos)
	chunkPos := VoxelToChunkPosition(voxelPos)

	if x == 0 {
		if neighbor := cm.GetChunkIfExists(chunkPos.GetNeighborPosition(voxel.VoxelFaceLeft)); neighbor != nil {
			neighbor.SetDirty(true)
		}
	}
	if x == ChunkSize-1 {
		if neighbor := cm.GetChunkIfExists(chunkPos.GetNeighborPosition(voxel.VoxelFaceRight)); neighbor != nil {
			neighbor.SetDirty(true)
		}
	}
	if y == 0 {
		if neighbor := cm.GetChunkIfExists(chunkPos.GetNeighborPosition(voxel.VoxelFaceBottom)); neighbor != nil {
			neighbor.SetDirty(true)
		}
	}
	if y == ChunkSize-1 {
		if neighbor := cm.GetChunkIfExists(chunkPos.GetNeighborPosition(voxel.VoxelFaceTop)); neighbor != nil {
			neighbor.SetDirty(true)
		}
	}
	if z == 0 {
		if neighbor := cm.GetChunkIfExists(chunkPos.GetNeighborPosition(voxel.VoxelFaceBack)); neighbor != nil {
			neighbor.SetDirty(true)
		}
	}
	if z == ChunkSize-1 {
		if neighbor := cm.GetChunkIfExists(chunkPos.GetNeighborPosition(voxel.VoxelFaceFront)); neighbor != nil {
			neighbor.SetDirty(true)
		}
	}
}

func (cm *ChunkManager) GetChunkIfExists(pos ChunkPosition) *Chunk {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	return cm.chunks[pos]
}

func (cm *ChunkManager) UnloadChunk(pos ChunkPosition) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if chunk, exists := cm.chunks[pos]; exists {
		for face := voxel.VoxelFaceTop; face <= voxel.VoxelFaceBack; face++ {
			if neighbor := chunk.GetNeighbor(face); neighbor != nil {
				oppositeFace := getOppositeFace(face)
				neighbor.SetNeighbor(oppositeFace, nil)
			}
		}

		delete(cm.chunks, pos)
		cm.loadedChunks--
	}
}

func (cm *ChunkManager) GetDirtyChunks() []*Chunk {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	dirtyChunks := make([]*Chunk, 0)
	for _, chunk := range cm.chunks {
		if chunk.IsDirty() && !chunk.IsEmpty() {
			dirtyChunks = append(dirtyChunks, chunk)
		}
	}

	return dirtyChunks
}

func (cm *ChunkManager) GetLoadedChunks() []*Chunk {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	chunks := make([]*Chunk, 0, len(cm.chunks))
	for _, chunk := range cm.chunks {
		chunks = append(chunks, chunk)
	}

	return chunks
}

func (cm *ChunkManager) GetStats() ChunkManagerStats {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	dirtyCount := 0
	emptyCount := 0

	for _, chunk := range cm.chunks {
		if chunk.IsDirty() {
			dirtyCount++
		}
		if chunk.IsEmpty() {
			emptyCount++
		}
	}

	return ChunkManagerStats{
		TotalChunks:  cm.totalChunks,
		LoadedChunks: cm.loadedChunks,
		DirtyChunks:  dirtyCount,
		EmptyChunks:  emptyCount,
	}
}

type ChunkManagerStats struct {
	TotalChunks  int
	LoadedChunks int
	DirtyChunks  int
	EmptyChunks  int
}

func getOppositeFace(face voxel.VoxelFace) voxel.VoxelFace {
	switch face {
	case voxel.VoxelFaceTop:
		return voxel.VoxelFaceBottom
	case voxel.VoxelFaceBottom:
		return voxel.VoxelFaceTop
	case voxel.VoxelFaceLeft:
		return voxel.VoxelFaceRight
	case voxel.VoxelFaceRight:
		return voxel.VoxelFaceLeft
	case voxel.VoxelFaceFront:
		return voxel.VoxelFaceBack
	case voxel.VoxelFaceBack:
		return voxel.VoxelFaceFront
	default:
		return face
	}
}

func (cm *ChunkManager) PrintStats() {
	stats := cm.GetStats()
	fmt.Printf("Chunk Manager Statistics:\n")
	fmt.Printf("  Total Chunks Created: %d\n", stats.TotalChunks)
	fmt.Printf("  Currently Loaded: %d\n", stats.LoadedChunks)
	fmt.Printf("  Dirty (need mesh): %d\n", stats.DirtyChunks)
	fmt.Printf("  Empty (all air): %d\n", stats.EmptyChunks)
}
