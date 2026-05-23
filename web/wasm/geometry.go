package main

import (
	"fmt"
	"math"
)

// Model dimension constants define the basic measurements for the 3D model.
const (
	BaseHeight    float64 = 10.0     // Height of the base in model units
	MaxHeight     float64 = 25.0     // Maximum height for contribution columns
	CellSize      float64 = 2.5      // Size of each contribution cell
	GridSize      int     = 53       // Number of weeks in a year
	BaseThickness float64 = 10.0     // Total thickness of the base
	MinHeight     float64 = CellSize // Minimum height for any contribution column
)

// YearOffset defines the depth spacing between successive years in a multi-year model.
const YearOffset float64 = 7.0 * CellSize

// NormalizeContribution converts a contribution count to a normalized height value.
// Returns 0 for no contributions, or a value between MinHeight and MaxHeight for active contributions.
func NormalizeContribution(count, maxCount int) float64 {
	if count == 0 {
		return 0 // No contribution means no column
	}
	if maxCount <= 0 {
		return MinHeight // Avoid division by zero, return minimum height
	}

	// Calculate the available height range for columns
	heightRange := MaxHeight - MinHeight

	// Use square root to create more visual variation in height
	normalizedValue := math.Sqrt(float64(count)) / math.Sqrt(float64(maxCount))

	// Scale to fit between MinHeight and MaxHeight
	return MinHeight + (normalizedValue * heightRange)
}

// CreateColumn generates a rectangular column (cuboid) with the given dimensions.
// x, y: position of the column
// height: height of the column
// size: width and depth of the column
func CreateColumn(x, y, height, size float64) ([]Triangle, error) {
	if height <= 0 {
		return nil, fmt.Errorf("height must be positive, got %f", height)
	}
	if size <= 0 {
		return nil, fmt.Errorf("size must be positive, got %f", size)
	}

	// Define the 8 vertices of the cuboid
	baseZ := BaseHeight
	topZ := BaseHeight + height

	v := []Point3D{
		{x, y, baseZ},           // 0: front-bottom-left
		{x + size, y, baseZ},    // 1: front-bottom-right
		{x + size, y + size, baseZ}, // 2: back-bottom-right
		{x, y + size, baseZ},    // 3: back-bottom-left
		{x, y, topZ},            // 4: front-top-left
		{x + size, y, topZ},     // 5: front-top-right
		{x + size, y + size, topZ}, // 6: back-top-right
		{x, y + size, topZ},     // 7: back-top-left
	}

	// Create triangles for each face of the cuboid (2 triangles per face)
	triangles := []Triangle{
		// Bottom face (pointing down: 0, 0, -1)
		{Point3D{0, 0, -1}, v[3], v[1], v[0]},
		{Point3D{0, 0, -1}, v[2], v[1], v[3]},
		// Top face (pointing up: 0, 0, 1)
		{Point3D{0, 0, 1}, v[4], v[5], v[6]},
		{Point3D{0, 0, 1}, v[4], v[6], v[7]},
		// Front face (pointing toward -Y: 0, -1, 0)
		{Point3D{0, -1, 0}, v[0], v[1], v[5]},
		{Point3D{0, -1, 0}, v[0], v[5], v[4]},
		// Back face (pointing toward +Y: 0, 1, 0)
		{Point3D{0, 1, 0}, v[3], v[7], v[6]},
		{Point3D{0, 1, 0}, v[3], v[6], v[2]},
		// Left face (pointing toward -X: -1, 0, 0)
		{Point3D{-1, 0, 0}, v[0], v[4], v[7]},
		{Point3D{-1, 0, 0}, v[0], v[7], v[3]},
		// Right face (pointing toward +X: 1, 0, 0)
		{Point3D{1, 0, 0}, v[1], v[2], v[6]},
		{Point3D{1, 0, 0}, v[1], v[6], v[5]},
	}

	return triangles, nil
}

// CreateCuboidBase generates the base platform for the skyline model.
func CreateCuboidBase(innerWidth, innerDepth float64) ([]Triangle, error) {
	if innerWidth <= 0 || innerDepth <= 0 {
		return nil, fmt.Errorf("invalid base dimensions: width=%f, depth=%f", innerWidth, innerDepth)
	}

	// Base vertices
	v := []Point3D{
		{0, 0, 0},                          // 0: front-bottom-left
		{innerWidth, 0, 0},                 // 1: front-bottom-right
		{innerWidth, innerDepth, 0},        // 2: back-bottom-right
		{0, innerDepth, 0},                 // 3: back-bottom-left
		{0, 0, BaseHeight},                 // 4: front-top-left
		{innerWidth, 0, BaseHeight},        // 5: front-top-right
		{innerWidth, innerDepth, BaseHeight}, // 6: back-top-right
		{0, innerDepth, BaseHeight},        // 7: back-top-left
	}

	triangles := []Triangle{
		// Bottom face
		{Point3D{0, 0, -1}, v[3], v[1], v[0]},
		{Point3D{0, 0, -1}, v[2], v[1], v[3]},
		// Top face
		{Point3D{0, 0, 1}, v[4], v[5], v[6]},
		{Point3D{0, 0, 1}, v[4], v[6], v[7]},
		// Front face
		{Point3D{0, -1, 0}, v[0], v[1], v[5]},
		{Point3D{0, -1, 0}, v[0], v[5], v[4]},
		// Back face
		{Point3D{0, 1, 0}, v[3], v[7], v[6]},
		{Point3D{0, 1, 0}, v[3], v[6], v[2]},
		// Left face
		{Point3D{-1, 0, 0}, v[0], v[4], v[7]},
		{Point3D{-1, 0, 0}, v[0], v[7], v[3]},
		// Right face
		{Point3D{1, 0, 0}, v[1], v[2], v[6]},
		{Point3D{1, 0, 0}, v[1], v[6], v[5]},
	}

	return triangles, nil
}

// CreateContributionGeometry generates geometry for a single year's contributions
func CreateContributionGeometry(contributions [][]ContributionDay, yearIndex int, maxContrib int) ([]Triangle, error) {
	var triangles []Triangle

	// Base Y offset includes padding and positions each year accordingly
	baseYOffset := 2*CellSize + float64(yearIndex)*YearOffset

	for weekIdx, week := range contributions {
		for dayIdx, day := range week {
			if day.ContributionCount > 0 {
				height := NormalizeContribution(day.ContributionCount, maxContrib)
				x := 2*CellSize + float64(weekIdx)*CellSize
				y := baseYOffset + float64(dayIdx)*CellSize

				columnTriangles, err := CreateColumn(x, y, height, CellSize)
				if err != nil {
					return nil, err
				}
				triangles = append(triangles, columnTriangles...)
			}
		}
	}

	return triangles, nil
}

// CalculateMultiYearDimensions calculates dimensions for multiple years
func CalculateMultiYearDimensions(yearCount int) (width, depth float64) {
	// Total width: grid size + padding on both sides
	width = float64(GridSize)*CellSize + 4*CellSize
	// Total depth: (7 days * number of years) + padding on both sides
	depth = float64(7*yearCount)*CellSize + 4*CellSize
	return width, depth
}

// findMaxContributions finds the maximum contribution count in the data
func findMaxContributions(contributions [][]ContributionDay) int {
	maxContrib := 0
	for _, week := range contributions {
		for _, day := range week {
			if day.ContributionCount > maxContrib {
				maxContrib = day.ContributionCount
			}
		}
	}
	return maxContrib
}

// findMaxContributionsAcrossYears finds the maximum contribution count across all years
func findMaxContributionsAcrossYears(contributionsPerYear [][][]ContributionDay) int {
	maxContrib := 0
	for _, yearContributions := range contributionsPerYear {
		yearMax := findMaxContributions(yearContributions)
		if yearMax > maxContrib {
			maxContrib = yearMax
		}
	}
	return maxContrib
}
