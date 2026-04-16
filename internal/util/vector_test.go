package util

import (
	"math"
	"testing"
)

func TestCosineSimilarity_IdenticalVectors(t *testing.T) {
	a := []float64{1.0, 2.0, 3.0}
	b := []float64{1.0, 2.0, 3.0}
	result := CosineSimilarity(a, b)
	if math.Abs(result-1.0) > 1e-9 {
		t.Errorf("expected 1.0, got %f", result)
	}
}

func TestCosineSimilarity_OrthogonalVectors(t *testing.T) {
	a := []float64{1.0, 0.0}
	b := []float64{0.0, 1.0}
	result := CosineSimilarity(a, b)
	if math.Abs(result) > 1e-9 {
		t.Errorf("expected 0.0, got %f", result)
	}
}

func TestCosineSimilarity_OppositeVectors(t *testing.T) {
	a := []float64{1.0, 2.0, 3.0}
	b := []float64{-1.0, -2.0, -3.0}
	result := CosineSimilarity(a, b)
	if math.Abs(result-(-1.0)) > 1e-9 {
		t.Errorf("expected -1.0, got %f", result)
	}
}

func TestCosineSimilarity_ZeroVector(t *testing.T) {
	a := []float64{0.0, 0.0, 0.0}
	b := []float64{1.0, 2.0, 3.0}
	result := CosineSimilarity(a, b)
	if result != 0.0 {
		t.Errorf("expected 0.0 for zero vector, got %f", result)
	}
}

func TestCosineSimilarity_EmptyVectors(t *testing.T) {
	result := CosineSimilarity([]float64{}, []float64{})
	if result != 0.0 {
		t.Errorf("expected 0.0 for empty vectors, got %f", result)
	}
}

func TestCosineSimilarity_DimensionMismatch(t *testing.T) {
	a := []float64{1.0, 2.0}
	b := []float64{1.0, 2.0, 3.0}
	result := CosineSimilarity(a, b)
	if result != 0.0 {
		t.Errorf("expected 0.0 for dimension mismatch, got %f", result)
	}
}

func TestCosineSimilarity_NilVectors(t *testing.T) {
	result := CosineSimilarity(nil, nil)
	if result != 0.0 {
		t.Errorf("expected 0.0 for nil vectors, got %f", result)
	}
}

func TestCosineSimilarity_PartialSimilarity(t *testing.T) {
	a := []float64{1.0, 1.0, 0.0}
	b := []float64{1.0, 0.0, 0.0}
	result := CosineSimilarity(a, b)
	// cos(45°) ≈ 0.7071
	expected := 1.0 / math.Sqrt(2.0)
	if math.Abs(result-expected) > 1e-9 {
		t.Errorf("expected %f, got %f", expected, result)
	}
}
