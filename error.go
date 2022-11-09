package goerrorpath

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// GoPathErrorType 类型定义
type GoPathErrorType struct {
	baseErr    error  // 错误信息
	baseErrStr string // 基础错误信息
	pkgPath    string // 包路径
}

// BErrorStr 获取基础错误信息字符串
func (g *GoPathErrorType) BErrorStr() string {
	return g.baseErrStr
}

// BError 获取基础错误信息
func (g *GoPathErrorType) BError() error {
	return g.baseErr
}

// Init 初始化包路径
func (g *GoPathErrorType) Init(err interface{}, bserr string) {
	// 包路径查询
	typ := reflect.TypeOf(err)
	p := ""
	if typ.Kind().String() == "ptr" {
		p = typ.Elem().PkgPath()
	} else {
		p = typ.PkgPath()
	}

	// 包路径优化
	idx := strings.Index(p, "/")
	if idx == -1 {
		p = ""
	} else {
		p = strings.TrimLeft(p[idx:], "/") + "."
	}

	g.pkgPath = p
	g.baseErrStr = bserr
	g.baseErr = errors.New(bserr)
}

// ParseNormalError 格式化常规错误
func (g *GoPathErrorType) ParseError(err error) error {
	if err == nil {
		return nil
	}

	berr := ""
	if len(g.baseErrStr) > 0 {
		berr = g.baseErrStr + " : "
	}
	return fmt.Errorf("%s%s ", berr, err.Error())
}

// ParsePathError 格式化 "做什么" 错误
func (g *GoPathErrorType) ParsePathError(dwt string, err error) error {
	if err == nil {
		return nil
	}

	berr := ""
	if len(g.baseErrStr) > 0 {
		berr = g.baseErrStr + " : "
	}

	return fmt.Errorf("%s%s : %s%s ", g.pkgPath, dwt, berr, err.Error())
}

// 输出干净错误,不输出详细错误信息
func (g *GoPathErrorType) ClearError(err error) error {
	if err == nil || err == g.baseErr || !strings.HasPrefix(err.Error(), g.pkgPath+" : ") {
		return err
	}
	return errors.New(strings.TrimLeft(err.Error(), g.pkgPath+" : "))
}
