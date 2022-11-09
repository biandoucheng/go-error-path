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
	a.Init(a, "功能测试")

	err := errors.New("执行错误")
	fmt.Println(a.ParsePathError("TestPath", err))
}
