package restful

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestResponse(t *testing.T) {
	resp := &Response{
		Code:    "KJ0000",
		Message: "ok",
		Data:    make([]int, 1),
	}
	respJson, _ := json.Marshal(resp)
	fmt.Println(string(respJson))
}

func TestPagination(t *testing.T) {
	p := Pagination{
		CurrentPage: new(int),
		PageSize:    new(int),
		PageCount:   new(int),
		Total:       new(int),
		Data:        []any{1, 2, 3},
	}

	j, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(j))
}

func TestSucceeded(t *testing.T) {
	t.Log(Succeeded(struct {
		a int
		b string
		c []string
	}{a: 1, b: "2", c: []string{"test", "111"}}))
}

func TestFailed(t *testing.T) {
	t.Log(Failed(ResourceNotFound, map[string]any{
		"a": 1,
		"b": "2",
		"c": []string{"test", "111"},
	}))
}
