package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kjhch/go-common/restful"
)

func TestResponse(t *testing.T) {
	resp := &restful.Response{
		Code:    "KJ0000",
		Message: "ok",
		Data:    make([]int, 1),
	}
	respJson, _ := json.Marshal(resp)
	fmt.Println(string(respJson))
}

func TestPagination(t *testing.T) {
	p := restful.Pagination{
		CurrentPage: new(int),
		PageSize:    new(int),
		PageCount:   new(int),
		Total:       new(int),
		Data:        []any{1, 2, 3},
	}

	j, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(j))
}

func TestXxx(t *testing.T) {
	// 查看『中文』的utf8编码
	fmt.Printf("%x\n", []byte("中文"))
	// 查看『中』的unicode编码
	fmt.Printf("%d %x\n", '中', '中')
}
