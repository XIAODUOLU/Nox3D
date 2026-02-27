package scene

import "math"

// Camera represents a 3D camera
type Camera struct {
	Position Vec3
	Target   Vec3
	Up       Vec3

	// Orbit camera parameters
	Distance float32
	Yaw      float32
	Pitch    float32

	// Projection parameters
	FOV    float32
	Aspect float32
	Near   float32
	Far    float32
}

// NewCamera creates a new camera with default settings
func NewCamera() *Camera {
	return &Camera{
		Position: Vec3{0, 0, 5},
		Target:   Vec3{0, 0, 0},
		Up:       Vec3{0, 1, 0},
		Distance: 5.0,
		Yaw:      0,
		Pitch:    0,
		FOV:      float32(math.Pi / 4), // 45 degrees
		Aspect:   1.0,
		Near:     0.1,
		Far:      100.0,
	}
}

// UpdateOrbit updates camera position based on orbit parameters
func (c *Camera) UpdateOrbit() {
	// Convert spherical coordinates to cartesian
	x := c.Distance * float32(math.Cos(float64(c.Pitch))) * float32(math.Sin(float64(c.Yaw)))
	y := c.Distance * float32(math.Sin(float64(c.Pitch)))
	z := c.Distance * float32(math.Cos(float64(c.Pitch))) * float32(math.Cos(float64(c.Yaw)))

	c.Position = Add(c.Target, Vec3{x, y, z})
}

// Zoom adjusts camera distance
func (c *Camera) Zoom(delta float32) {
	c.Distance += delta
	if c.Distance < 0.5 {
		c.Distance = 0.5
	}
	if c.Distance > 50 {
		c.Distance = 50
	}
	c.UpdateOrbit()
}

// Rotate rotates the camera
func (c *Camera) Rotate(deltaYaw, deltaPitch float32) {
	c.Yaw += deltaYaw
	c.Pitch += deltaPitch

	// Clamp pitch to avoid gimbal lock
	maxPitch := float32(math.Pi/2 - 0.01)
	if c.Pitch > maxPitch {
		c.Pitch = maxPitch
	}
	if c.Pitch < -maxPitch {
		c.Pitch = -maxPitch
	}

	c.UpdateOrbit()
}

// GetViewMatrix returns the view matrix
func (c *Camera) GetViewMatrix() Mat4 {
	return LookAt(c.Position, c.Target, c.Up)
}

// GetProjectionMatrix returns the projection matrix
func (c *Camera) GetProjectionMatrix() Mat4 {
	return Perspective(c.FOV, c.Aspect, c.Near, c.Far)
}
