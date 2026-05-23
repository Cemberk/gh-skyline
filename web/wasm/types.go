package main

// ContributionDay represents a single day of GitHub contributions.
// This mirrors the types.ContributionDay from the main package but is defined here
// to avoid importing the entire internal package structure into WASM.
type ContributionDay struct {
	ContributionCount int    `json:"contributionCount"`
	Date              string `json:"date"`
}

// Point3D represents a point in 3D space using float64 for accuracy in calculations.
type Point3D struct {
	X, Y, Z float64
}

// Point3DFloat32 represents a point in 3D space using float32 for STL output.
type Point3DFloat32 struct {
	X, Y, Z float32
}

// ToFloat32 converts a Point3D to Point3DFloat32.
func (p Point3D) ToFloat32() Point3DFloat32 {
	return Point3DFloat32{
		X: float32(p.X),
		Y: float32(p.Y),
		Z: float32(p.Z),
	}
}

// Triangle represents a triangle in 3D space using float64 coordinates.
type Triangle struct {
	Normal     Point3D
	V1, V2, V3 Point3D
}

// TriangleFloat32 represents a triangle in 3D space using float32 coordinates for STL output.
type TriangleFloat32 struct {
	Normal     Point3DFloat32
	V1, V2, V3 Point3DFloat32
}

// ToFloat32 converts a Triangle to TriangleFloat32.
func (t Triangle) ToFloat32() TriangleFloat32 {
	return TriangleFloat32{
		Normal: t.Normal.ToFloat32(),
		V1:     t.V1.ToFloat32(),
		V2:     t.V2.ToFloat32(),
		V3:     t.V3.ToFloat32(),
	}
}
