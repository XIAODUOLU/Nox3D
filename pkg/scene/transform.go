package scene

import "math"

// Mat4 represents a 4x4 matrix
type Mat4 [16]float32

// Identity returns an identity matrix
func Identity() Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Multiply multiplies two matrices
func (m Mat4) Multiply(other Mat4) Mat4 {
	var result Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i*4+j] = 0
			for k := 0; k < 4; k++ {
				result[i*4+j] += m[i*4+k] * other[k*4+j]
			}
		}
	}
	return result
}

// TransformVec3 transforms a Vec3 by this matrix (treating it as a point with w=1)
func (m Mat4) TransformVec3(v Vec3) Vec3 {
	x := m[0]*v.X + m[1]*v.Y + m[2]*v.Z + m[3]
	y := m[4]*v.X + m[5]*v.Y + m[6]*v.Z + m[7]
	z := m[8]*v.X + m[9]*v.Y + m[10]*v.Z + m[11]
	w := m[12]*v.X + m[13]*v.Y + m[14]*v.Z + m[15]

	if w != 0 {
		return Vec3{x / w, y / w, z / w}
	}
	return Vec3{x, y, z}
}

// TransformVec3Dir transforms a Vec3 as a direction (w=0)
func (m Mat4) TransformVec3Dir(v Vec3) Vec3 {
	x := m[0]*v.X + m[1]*v.Y + m[2]*v.Z
	y := m[4]*v.X + m[5]*v.Y + m[6]*v.Z
	z := m[8]*v.X + m[9]*v.Y + m[10]*v.Z
	return Vec3{x, y, z}
}

// Perspective creates a perspective projection matrix
func Perspective(fov, aspect, near, far float32) Mat4 {
	f := float32(1.0 / math.Tan(float64(fov)/2.0))

	return Mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / (near - far), (2 * far * near) / (near - far),
		0, 0, -1, 0,
	}
}

// LookAt creates a view matrix
func LookAt(eye, target, up Vec3) Mat4 {
	zaxis := Normalize(Sub(eye, target))
	xaxis := Normalize(Cross(up, zaxis))
	yaxis := Cross(zaxis, xaxis)

	return Mat4{
		xaxis.X, xaxis.Y, xaxis.Z, -Dot(xaxis, eye),
		yaxis.X, yaxis.Y, yaxis.Z, -Dot(yaxis, eye),
		zaxis.X, zaxis.Y, zaxis.Z, -Dot(zaxis, eye),
		0, 0, 0, 1,
	}
}

// RotationY creates a rotation matrix around Y axis
func RotationY(angle float32) Mat4 {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))

	return Mat4{
		c, 0, s, 0,
		0, 1, 0, 0,
		-s, 0, c, 0,
		0, 0, 0, 1,
	}
}

// RotationX creates a rotation matrix around X axis
func RotationX(angle float32) Mat4 {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))

	return Mat4{
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
		0, 0, 0, 1,
	}
}

// Translation creates a translation matrix
func Translation(x, y, z float32) Mat4 {
	return Mat4{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	}
}

// Scale creates a scale matrix
func Scale(x, y, z float32) Mat4 {
	return Mat4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

// Vector operations

// Add adds two vectors
func Add(a, b Vec3) Vec3 {
	return Vec3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Sub subtracts two vectors
func Sub(a, b Vec3) Vec3 {
	return Vec3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Dot computes dot product
func Dot(a, b Vec3) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Cross computes cross product
func Cross(a, b Vec3) Vec3 {
	return Vec3{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

// Length computes vector length
func Length(v Vec3) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Normalize normalizes a vector
func Normalize(v Vec3) Vec3 {
	l := Length(v)
	if l == 0 {
		return Vec3{0, 0, 0}
	}
	return Vec3{v.X / l, v.Y / l, v.Z / l}
}

// Scale3 scales a vector
func Scale3(v Vec3, s float32) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}
