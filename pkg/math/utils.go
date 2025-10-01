package math

import (
    "math"
)

const (
    Pi = float32(math.Pi)
    
    Epsilon = 1e-6
    
    DegToRad = Pi / 180.0
    
    RadToDeg = 180.0 / Pi
)

func Deg2Rad(degrees float32) float32 {
    return degrees * DegToRad
}

func Rad2Deg(radians float32) float32 {
    return radians * RadToDeg
}

func Clamp(value, min, max float32) float32 {
    return clampFloat(value, min, max)
}

func Lerp(a, b, t float32) float32 {
    return a + (b-a)*t
}

func InverseLerp(a, b, value float32) float32 {
    if a == b {
        return 0
    }
    return (value - a) / (b - a)
}

func Abs(x float32) float32 {
    return float32(math.Abs(float64(x)))
}

func Sign(x float32) float32 {
    if x > 0 {
        return 1
    }
    if x < 0 {
        return -1
    }
    return 0
}

func Min(a, b float32) float32 {
    if a < b {
        return a
    }
    return b
}

func Max(a, b float32) float32 {
    if a > b {
        return a
    }
    return b
}

func Floor(x float32) float32 {
    return float32(math.Floor(float64(x)))
}

func Ceil(x float32) float32 {
    return float32(math.Ceil(float64(x)))
}

func Round(x float32) float32 {
    return float32(math.Round(float64(x)))
}

func Sqrt(x float32) float32 {
    return float32(math.Sqrt(float64(x)))
}

func Pow(a, b float32) float32 {
    return float32(math.Pow(float64(a), float64(b)))
}

func Sin(x float32) float32 {
    return float32(math.Sin(float64(x)))
}

func Cos(x float32) float32 {
    return float32(math.Cos(float64(x)))
}

func Tan(x float32) float32 {
    return float32(math.Tan(float64(x)))
}

func ApproxEqual(a, b float32) bool {
    return Abs(a-b) < Epsilon
}