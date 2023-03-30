package slice

import (
	"fmt"
	"testing"
)

func TestIndexOf(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	if s.IndexOf(3) != 2 {
		t.Errorf("Expected index of 3 to be 2, got %d", s.IndexOf(3))
	}
	if s.IndexOf(6) != -1 {
		t.Errorf("Expected index of 6 to be -1, got %d", s.IndexOf(6))
	}
}

func TestHas(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	if !s.Has(3) {
		t.Errorf("Expected Has(3) to be true, got false")
	}
	if s.Has(6) {
		t.Errorf("Expected Has(6) to be false, got true")
	}
}

func TestEqual(t *testing.T) {
	s1 := New(1, 2, 3, 4, 5)
	s2 := New(1, 2, 3, 4, 5)
	if !s1.Equal(s2) {
		t.Errorf("Expected Equal() to be true, got false")
	}
	s3 := New(1, 2, 3, 4)
	if s1.Equal(s3) {
		t.Errorf("Expected Equal() to be false, got true")
	}
}

func TestAny(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	if !s.Any(func(v int) bool { return v == 3 }) {
		t.Errorf("Expected Any() to be true, got false")
	}
	if s.Any(func(v int) bool { return v == 6 }) {
		t.Errorf("Expected Any() to be false, got true")
	}
}

func TestAll(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	if !s.All(func(v int) bool { return v > 0 }) {
		t.Errorf("Expected All() to be true, got false")
	}
	if s.All(func(v int) bool { return v > 1 }) {
		t.Errorf("Expected All() to be false, got true")
	}
}

func TestFilter(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	filtered := s.Filter(func(v int) bool { return v > 3 })
	expected := New(4, 5)
	if !filtered.Equal(expected) {
		t.Errorf("Expected Filter() to return %v, got %v", expected, filtered)
	}
}

func TestFind(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	if s.Find(func(v int) bool { return v == 3 }) != 2 {
		t.Errorf("Expected Find() to return 2, got %d", s.Find(func(v int) bool { return v == 3 }))
	}
	if s.Find(func(v int) bool { return v == 6 }) != -1 {
		t.Errorf("Expected Find() to return -1, got %d", s.Find(func(v int) bool { return v == 6 }))
	}
}

func TestMap(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	mapped := Map(s, func(v int) string { return fmt.Sprintf("x%d", v) })
	expected := New("x1", "x2", "x3", "x4", "x5")
	if !mapped.Equal(expected) {
		t.Errorf("Expected Map() to return %v, got %v", expected, mapped)
	}
}
