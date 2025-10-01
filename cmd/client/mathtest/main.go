package main

import (
    "fmt"

    ceresmath "Ceres/pkg/math"
)

func main() {
    fmt.Println("=== Ceres Math Library Test ===\n")

    // Test Vector3 operations
    fmt.Println("--- Vector3 Tests ---")
    v1 := ceresmath.NewVector3(1, 2, 3)
    v2 := ceresmath.NewVector3(4, 5, 6)

    fmt.Printf("v1: %+v\n", v1)
    fmt.Printf("v2: %+v\n", v2)
    fmt.Printf("v1 + v2: %+v\n", v1.Add(v2))
    fmt.Printf("v1 - v2: %+v\n", v1.Sub(v2))
    fmt.Printf("v1 * 2: %+v\n", v1.Mul(2))
    fmt.Printf("v1 · v2 (dot): %.2f\n", v1.Dot(v2))
    fmt.Printf("v1 × v2 (cross): %+v\n", v1.Cross(v2))
    fmt.Printf("Length of v1: %.2f\n", v1.Length())
    fmt.Printf("Normalized v1: %+v\n", v1.Normalize())
    fmt.Printf("Distance v1 to v2: %.2f\n", v1.Distance(v2))

    // Test directional vectors
    fmt.Println("\n--- Direction Vectors ---")
    fmt.Printf("Up: %+v\n", ceresmath.Up())
    fmt.Printf("Forward: %+v\n", ceresmath.Forward())
    fmt.Printf("Right: %+v\n", ceresmath.Right())

    // Test Matrix operations
    fmt.Println("\n--- Matrix4 Tests ---")
    _ = ceresmath.Identity()
    fmt.Printf("Identity matrix created\n")

    translation := ceresmath.Translate(5, 10, 15)
    fmt.Printf("Translation matrix created\n")

    rotation := ceresmath.RotateY(ceresmath.Deg2Rad(45))
    fmt.Printf("Rotation matrix (45° around Y) created\n")

    scale := ceresmath.Scale(2, 2, 2)
    fmt.Printf("Scale matrix (2x) created\n")

    // Combine transformations
    combined := translation.Mul(rotation).Mul(scale)
    fmt.Printf("Combined transformation matrix created\n")

    // Transform a point
    point := ceresmath.NewVector3(1, 0, 0)
    transformed := combined.MulVec(point)
    fmt.Printf("Point %+v transformed to %+v\n", point, transformed)

    // Test Transform struct
    fmt.Println("\n--- Transform Tests ---")
    transform := ceresmath.NewTransform()
    transform.Position = ceresmath.NewVector3(10, 5, 0)
    transform.Rotation = ceresmath.NewVector3(0, ceresmath.Deg2Rad(90), 0)
    transform.Scale = ceresmath.NewVector3(1, 1, 1)

    fmt.Printf("Transform position: %+v\n", transform.Position)
    fmt.Printf("Transform rotation (deg): %.0f, %.0f, %.0f\n",
        ceresmath.Rad2Deg(transform.Rotation.X),
        ceresmath.Rad2Deg(transform.Rotation.Y),
        ceresmath.Rad2Deg(transform.Rotation.Z))
    fmt.Printf("Transform forward: %+v\n", transform.Forward())
    fmt.Printf("Transform right: %+v\n", transform.Right())
    fmt.Printf("Transform up: %+v\n", transform.Up())

    // Test projection matrices
    fmt.Println("\n--- Projection Tests ---")
    _ = ceresmath.Perspective(
        ceresmath.Deg2Rad(60), // 60° FOV
        16.0/9.0,              // 16:9 aspect ratio
        0.1,                   // Near plane
        1000.0,                // Far plane
    )
    fmt.Printf("Perspective projection matrix created (60° FOV, 16:9)\n")

    // Test view matrix
    eye := ceresmath.NewVector3(0, 10, 20)
    center := ceresmath.NewVector3(0, 0, 0)
    up := ceresmath.Up()
    _ = ceresmath.LookAt(eye, center, up)
    fmt.Printf("View matrix created (eye at %+v, looking at %+v)\n", eye, center)

    // Test utility functions
    fmt.Println("\n--- Utility Tests ---")
    fmt.Printf("Deg2Rad(180): %.4f (should be ~%.4f)\n",
        ceresmath.Deg2Rad(180), ceresmath.Pi)
    fmt.Printf("Rad2Deg(π): %.1f (should be 180)\n",
        ceresmath.Rad2Deg(ceresmath.Pi))
    fmt.Printf("Lerp(0, 10, 0.5): %.1f\n", ceresmath.Lerp(0, 10, 0.5))
    fmt.Printf("Clamp(15, 0, 10): %.1f\n", ceresmath.Clamp(15, 0, 10))
    fmt.Printf("Min(5, 10): %.1f\n", ceresmath.Min(5, 10))
    fmt.Printf("Max(5, 10): %.1f\n", ceresmath.Max(5, 10))

    fmt.Println("\n✓ All math tests completed successfully!")
}