package math

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Matrix4 struct {
	mgl32.Mat4
}

func Identity() Matrix4 {
	return Matrix4{mgl32.Ident4()}
}

func Translate(x, y, z float32) Matrix4 {
	return Matrix4{mgl32.Translate3D(x, y, z)}
}

func TranslateVec(v Vector3) Matrix4 {
	return Matrix4{mgl32.Translate3D(v.X, v.Y, v.Z)}
}

func Scale(x, y, z float32) Matrix4 {
	return Matrix4{mgl32.Scale3D(x, y, z)}
}

func ScaleVec(v Vector3) Matrix4 {
	return Matrix4{mgl32.Scale3D(v.X, v.Y, v.Z)}
}

func ScaleUniform(s float32) Matrix4 {
	return Matrix4{mgl32.Scale3D(s, s, s)}
}

func RotateX(angle float32) Matrix4 {
	return Matrix4{mgl32.HomogRotate3DX(angle)}
}

func RotateY(angle float32) Matrix4 {
	return Matrix4{mgl32.HomogRotate3DY(angle)}
}

func RotateZ(angle float32) Matrix4 {
	return Matrix4{mgl32.HomogRotate3DZ(angle)}
}

func Rotate(angle float32, axis Vector3) Matrix4 {
	return Matrix4{mgl32.HomogRotate3D(angle, axis.ToMgl32())}
}

func Perspective(fovY, aspect, near, far float32) Matrix4 {
	return Matrix4{mgl32.Perspective(fovY, aspect, near, far)}
}

func Ortho(left, right, bottom, top, near, far float32) Matrix4 {
	return Matrix4{mgl32.Ortho(left, right, bottom, top, near, far)}
}

func LookAt(eye, center, up Vector3) Matrix4 {
	return Matrix4{mgl32.LookAtV(eye.ToMgl32(), center.ToMgl32(), up.ToMgl32())}
}

func (m Matrix4) Mul(other Matrix4) Matrix4 {
	return Matrix4{m.Mat4.Mul4(other.Mat4)}
}

func (m Matrix4) MulVec(v Vector3) Vector3 {
	vec4 := m.Mat4.Mul4x1(mgl32.Vec4{v.X, v.Y, v.Z, 1.0})
	return Vector3{vec4.X(), vec4.Y(), vec4.Z()}
}

func (m Matrix4) MulDir(v Vector3) Vector3 {
    vec4 := m.Mat4.Mul4x1(mgl32.Vec4{v.X, v.Y, v.Z, 0.0})
    return Vector3{vec4.X(), vec4.Y(), vec4.Z()}
}

func (m Matrix4) Inverse() Matrix4 {
	return Matrix4{m.Mat4.Inv()}
}

func (m Matrix4) Transpose() Matrix4 {
	return Matrix4{m.Mat4.Transpose()}
}

func (m Matrix4) ToPtr() *float32 {
	return &m.Mat4[0]
}