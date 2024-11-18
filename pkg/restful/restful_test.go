package restful

import (
	"encoding/json"
	"fmt"
	"os"
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

type MyError struct {
	code    string
	message string
}

func (myError *MyError) Code() string {
	return myError.code
}
func (myError *MyError) Message() string {
	return myError.message
}
func (myError *MyError) Status() int {
	return 200
}

func (myError *MyError) Error() string {
	return fmt.Sprintf("restful.Error(%v-%v):%v", myError.Status(), myError.Code(), myError.Message())
}

func TestFailed(t *testing.T) {
	e := &MyError{
		code:    "mdm-1000",
		message: "sth wrong",
	}
	t.Log(Failed(e, map[string]any{
		"a": 1,
		"b": "2",
		"c": []string{"test", "111"},
	}))

	t.Log(Failed(os.ErrClosed, nil))
}
