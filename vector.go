package main

import "math"

type Vector struct {
	X1, X2, X3 float64
}

func (v *Vector) Init(x1, x2, x3 float64) {
	v.X1 = x1
	v.X2 = x2
	v.X3 = x3
}

func (v Vector) GetDotProduct(vec1 Vector) float64 {
	return v.X1*vec1.X1 + v.X2*vec1.X2 + v.X3*vec1.X3
}

func (v Vector) GetVecSquared() float64 {
	return v.X1*v.X1 + v.X2*v.X2 + v.X3*v.X3
}

func (v Vector) GetAbs() float64 {
	return math.Sqrt(v.X1*v.X1 + v.X2*v.X2 + v.X3*v.X3)
}

func (v *Vector) Scale(s float64) {
	v.X1 *= s
	v.X2 *= s
	v.X3 *= s
}

func (v Vector) GetNormalWith(vec2 Vector) Vector {
	vec := Vector{
		X1: v.X2*vec2.X3 - v.X3*vec2.X2,
		X2: v.X3*vec2.X1 - v.X1*vec2.X3,
		X3: v.X1*vec2.X2 - v.X2*vec2.X1,
	}
	return vec
}

func (v Vector) DifferenceVector(vec1 Vector) Vector {
	vec := Vector{
		X1: v.X1 - vec1.X1,
		X2: v.X2 - vec1.X2,
		X3: v.X3 - vec1.X3,
	}
	return vec
}
