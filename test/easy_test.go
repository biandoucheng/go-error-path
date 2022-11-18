package test

import (
	"errors"
	"fmt"
	"testing"

	goerrorpath "github.com/biandoucheng/go-error-path"
)

type A struct {
	goerrorpath.GoPathErrorType
}

func TestPath(t *testing.T) {
	a := &A{}
	a.Init(a, "功能测试", "内部错误")

	err := errors.New("执行错误")
	nerr := a.ParseError("TestPath", err)
	fmt.Println(nerr.BaseError())
	fmt.Println(nerr.ShortError())
	fmt.Println(nerr.DetailError())
	fmt.Println(nerr.PathError())
}
