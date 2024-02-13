package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// const tpl = `package restful

// type Error struct {
// 	Code    int
// 	Message string
// 	Status  int
// }

// var (
// {{range .}}
//     {{.VarName}} = &Error{Code: {{.Code}}, Message: "{{.Message}}", Status: {{.Status}}}
// {{end}}
// )

// `

type errordef struct {
	Code    string
	DescCn  string
	Comment string
	VarName string
	Status  string
	Message string
}

func main() {
	fileBytes, err := os.ReadFile("../configs/errorcode.csv")
	if err != nil {
		panic(err)
	}
	fileStr := string(fileBytes)
	fileLines := strings.Split(fileStr, "\n")

	errordefs := make([]errordef, 0)
	firstLine := true
	for _, line := range fileLines {
		if firstLine {
			firstLine = false
			continue
		}
		if line == "" {
			continue
		}
		lineSplits := strings.Split(line, ",")
		fmt.Println(lineSplits)
		errordefs = append(errordefs, errordef{
			Code:    strings.TrimSpace(lineSplits[0]),
			DescCn:  strings.TrimSpace(lineSplits[1]),
			Comment: strings.TrimSpace(lineSplits[2]),
			VarName: strings.TrimSpace(lineSplits[3]),
			Status:  strings.TrimSpace(lineSplits[4]),
			Message: strings.TrimSpace(lineSplits[5]),
		})
	}
	tpl, err := template.ParseFiles("../configs/errorcode.tpl")
	if err != nil {
		panic(err)
	}
	f, err := os.Create("../pkg/restful/error.go")
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(f, errordefs)
	if err != nil {
		panic(err)
	}
}
