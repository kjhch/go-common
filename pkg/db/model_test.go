package db

import "testing"

func TestJsonArr(t *testing.T) {
	s := JsonArray[int]{1, 2, 3}
	v, _ := s.Value()
	t.Log(string(v.([]byte)))

	sa := JsonArray[string]{"a", "b", "c"}
	v, _ = sa.Value()
	t.Log(string(v.([]byte)))
}
