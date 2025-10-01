package math

import (
	"math"
)

type Transform struct {
	Position Vector3
	Rotation Vector3
	Scale    Vector3
}

func NewTransform() Transform {
	return Transform{
		Position: Zero(),
		Rotation: Zero(),
		Scale:    One(),
	}
}

func (t Transform) Matrix() Matrix4 {
	translation := TranslateVec(t.Position)
	rotationX := RotateX(t.Rotation.X)
	rotationY := RotateY(t.Rotation.Y)
	rotationZ := RotateZ(t.Rotation.Z)
	scale := ScaleVec(t.Scale)

	return translation.Mul(rotationY).Mul(rotationX).Mul(rotationZ).Mul(scale)
}

func (t Transform) Forward() Vector3 {
	return t.Matrix().MulDir(Forward()).Normalize()
}

func (t Transform) Right() Vector3 {
	return t.Matrix().MulDir(Right()).Normalize()
}

func (t Transform) Up() Vector3 {
	return t.Matrix().MulDir(Up()).Normalize()
}

func (t *Transform) LookAt(target, worldUp Vector3) {
	forward := target.Sub(t.Position).Normalize()
    _ = forward.Cross(worldUp).Normalize() 

    t.Rotation.Y = float32(math.Atan2(float64(forward.X), float64(forward.Z)))
    t.Rotation.X = float32(math.Asin(float64(-forward.Y)))
    t.Rotation.Z = 0
}