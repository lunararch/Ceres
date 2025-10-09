package camera

import (
	"math"

	ceresmath "Ceres/pkg/math"
)

type Camera struct {
	Position ceresmath.Vector3

	// Orientation
	Yaw float32
	Pitch float32
	Roll  float32

	// Camera Vectors
	Front ceresmath.Vector3
	Right ceresmath.Vector3
	Up    ceresmath.Vector3

	// World Up vector
	WorldUp ceresmath.Vector3

	// Camera Settings
	MovementSpeed float32
	MouseSensitivity float32
	Zoom float32 // fov

	// Constraints
	MinPitch float32
	MaxPitch float32
}

func NewCamera(position ceresmath.Vector3) *Camera {
	c := &Camera{
		Position: position,
		Yaw:              -90.0,
        Pitch:            0.0,
        Roll:             0.0,
        WorldUp:          ceresmath.Up(),
        MovementSpeed:    5.0,
        MouseSensitivity: 0.1,
        Zoom:             45.0,
        MinPitch:         -89.0,
        MaxPitch:         89.0,
	}

	c.updateCameraVectors()
	return c
}

func (c *Camera) GetViewMatrix() ceresmath.Matrix4 {
	return ceresmath.LookAt(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *Camera) GetProjectionMatrix(aspectRatio, near, far float32) ceresmath.Matrix4 {
	return ceresmath.Perspective(ceresmath.Deg2Rad(c.Zoom), aspectRatio, near, far)
}

func (c *Camera) ProcessKeyboard(direction Direction, deltaTime float32) {
	velocity := c.MovementSpeed * deltaTime

	switch direction {
	case Forward:
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	case Backward:
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	case Left:
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	case Right:
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	case Up:
		c.Position = c.Position.Add(c.WorldUp.Mul(velocity))
	case Down:
		c.Position = c.Position.Sub(c.WorldUp.Mul(velocity))
	}
}

func (c *Camera) ProcessMouseMovement(xOffset, yOffset float32, constrainPitch bool) {
	xOffset *= c.MouseSensitivity
	yOffset *= c.MouseSensitivity

	c.Yaw += xOffset
	c.Pitch += yOffset

	if constrainPitch {
		if c.Pitch > c.MaxPitch {
			c.Pitch = c.MaxPitch
		}
		if c.Pitch < c.MinPitch {
			c.Pitch = c.MinPitch
		}
	}

	c.updateCameraVectors()
}

func (c *Camera) ProcessMouseScroll(yoffset float32) {
	c.Zoom -= yoffset

	if c.Zoom < 1.0 {
		c.Zoom = 1.0
	}
	if c.Zoom > 120.0 {
		c.Zoom = 120.0
	}
}

func (c *Camera) updateCameraVectors() {
	yawRad := ceresmath.Deg2Rad(c.Yaw)
	pitchRad := ceresmath.Deg2Rad(c.Pitch)

	front := ceresmath.Vector3{
		X: ceresmath.Cos(yawRad) * ceresmath.Cos(pitchRad),
		Y: ceresmath.Sin(pitchRad),
		Z: ceresmath.Sin(yawRad) * ceresmath.Cos(pitchRad),
	}

	c.Front = front.Normalize()
	
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}

func (c *Camera) GetFrustum(aspectRatio, near, far float32) Frustum {
	return Frustum{}
}

func (c *Camera) SetPosition(position ceresmath.Vector3) {
    c.Position = position
}

func (c *Camera) SetYaw(yaw float32) {
    c.Yaw = yaw
    c.updateCameraVectors()
}

func (c *Camera) SetPitch(pitch float32) {
    c.Pitch = ceresmath.Clamp(pitch, c.MinPitch, c.MaxPitch)
    c.updateCameraVectors()
}

func (c *Camera) LookAt(target ceresmath.Vector3) {
    direction := target.Sub(c.Position).Normalize()
    
    c.Yaw = ceresmath.Rad2Deg(float32(math.Atan2(float64(direction.Z), float64(direction.X))))
    c.Pitch = ceresmath.Rad2Deg(float32(math.Asin(float64(direction.Y))))
    
    c.updateCameraVectors()
}