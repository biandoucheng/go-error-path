package goerrorpath

import (
	"fmt"
	"reflect"
)

// GoPathErrorInterface 接口定义
type GoPathErrorInterface interface {
	Error() string
	Init()
}

// GoPathErrorType 类型定义
type GoPathErrorType struct {
	baseError string // 基础错误信息
	pkgPath   string // 包路径
}

// Error 获取基础错误信息
func (g *GoPathErrorType) Error() string {
	return g.baseError
}

// Init 初始化包路径
func (g *GoPathErrorType) Init(err interface{}, bserr string) {
	typ := reflect.TypeOf(err)
	if typ.Kind().String() == "ptr" {
		g.pkgPath = typ.Elem().PkgPath()
	} else {
		g.pkgPath = typ.PkgPath()
	}

	g.baseError = bserr
}

// ParseDwtErr 格式化 "做什么" 错误
func (g *GoPathErrorType) ParsePkgDwtErr(dwt string, err error) error {
	berr := ""
	if len(g.baseError) > 0 {
		berr = g.baseError + " : "
	}
	return fmt.Errorf("%s.%s Error: %s%s ", g.pkgPath, dwt, berr, err.Error())
}

// ParseError 格式化错误
func (g *GoPathErrorType) ParsePkgError(err error) error {
	berr := ""
	if len(g.baseError) > 0 {
		berr = g.baseError + " : "
	}
	return fmt.Errorf("%s Error: %s%s ", g.pkgPath, berr, err.Error())
}

// ParseNormalError 格式化常规错误
func (g *GoPathErrorType) ParseNormalError(err error) error {
	berr := ""
	if len(g.baseError) > 0 {
		berr = g.baseError + " : "
	}
	return fmt.Errorf("%s%s ", berr, err.Error())
}
