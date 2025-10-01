package math

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Vector3 struct {
	X, Y, Z float32
}

func NewVector3(x, y, z float32) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

func Zero() Vector3 {
	return Vector3{0, 0, 0}
}

func One() Vector3 {
	return Vector3{1, 1, 1}
}

func Up() Vector3 {
	return Vector3{0, 1, 0}
}

func Down() Vector3 {
	return Vector3{0, -1, 0}
}

func Forward() Vector3 {
	return Vector3{0, 0, -1}
}

func Back() Vector3 {
	return Vector3{0, 0, 1}
}

func Right() Vector3 {
	return Vector3{1, 0, 0}
}

func Left() Vector3 {
	return Vector3{-1, 0, 0}
}

func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

func (v Vector3) Mul(scalar float32) Vector3 {
	return Vector3{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

func (v Vector3) Div(scalar float32) Vector3 {
	return Vector3{
		X: v.X / scalar,
		Y: v.Y / scalar,
		Z: v.Z / scalar,
	}
}

func (v Vector3) Dot(other Vector3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

func (v Vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vector3) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector3) Normalize() Vector3 {
	len := v.Length()
	if len == 0 {
		return Zero()
	}
	return v.Div(len)
}

func (v Vector3) Distance(other Vector3) float32 {
	return v.Sub(other).Length()
}

func (v Vector3) DistanceSquared(other Vector3) float32 {
	return v.Sub(other).LengthSquared()
}

func (v Vector3) Lerp(other Vector3, t float32) Vector3 {
	return Vector3{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
		Z: v.Z + (other.Z-v.Z)*t,
	}
}

func (v Vector3) Clamp(min, max Vector3) Vector3 {
	return Vector3{
		clampFloat(v.X, min.X, max.X),
		clampFloat(v.Y, min.Y, max.Y),
		clampFloat(v.Z, min.Z, max.Z),
	}
}

func FromMgl32(v mgl32.Vec3) Vector3 {
	return Vector3{v.X(), v.Y(), v.Z()}
}

func (v Vector3) ToMgl32() mgl32.Vec3 {
	return mgl32.Vec3{v.X, v.Y, v.Z}
}

func clampFloat(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
