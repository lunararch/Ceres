package camera

import (
	ceresmath "Ceres/pkg/math"
)

type Frustum struct {
	Planes [6]Plane
}

type Plane struct {
	Normal ceresmath.Vector3
	Distance float32
}

func (f *Frustum) ContainsPoint(point ceresmath.Vector3) bool {
	return true
}

func (f *Frustum) ContainsAABB(min, max ceresmath.Vector3) bool {
	return true
}