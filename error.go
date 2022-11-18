package goerrorpath

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/*
..................................................................................................
. ErrorItem 错误接口定义
. 实现错误的分类管理
. BaseError   基础错误用来作为详细错误和路径错误的前缀
. ShortError  简短错误描述，用来进行对外输出，避免代码核心信息暴露
. DetailError 是对代码中多层调用产生出错误进行按执行顺讯从里到外拼接，暴露的错误产生的完整路线
. PathError   它是将代码执行过程中出现错误的包的基础错误与包路径层层拼接，暴露出程序的执行路径及执行出现错误的路径
..................................................................................................
*/
type ErrorItem interface {
	IsNil() bool
	BaseError() error
	ShortError() error
	DetailError() error
	PathError() error
}

/*
....................................................
. GoPathErrorItem go-error-path 的输出错误单元定义
. baseError   存储的是基础错误定义
. shortError  存储的简短错误定义即允许对外输出的错误信息
. detailError 存储每层代码调用的具体错误拼接
. pathError   存储每层代码调用包路径及其基础错误信息
.....................................................
*/
type GoPathErrorItem struct {
	baseError   error // 基础错误信息
	shortError  error // 对外输出错误信息
	detailError error // 详细错误信息
	pathError   error // 错误携带路径
}

// IsNil 是否为空
func (e *GoPathErrorItem) IsNil() bool {
	return e.detailError == nil && e.pathError == nil
}

// BError 获取基础错误信息
func (e *GoPathErrorItem) BaseError() error {
	return e.baseError
}

// ShortError 获取简短错误信息，一般用于外部错误返回，不暴露程序信息
func (e *GoPathErrorItem) ShortError() error {
	return e.shortError
}

// DetailError 获取详细错误包含完整路径
func (e *GoPathErrorItem) DetailError() error {
	return e.detailError
}

// PathError 获取包路径错误
func (e *GoPathErrorItem) PathError() error {
	return e.pathError
}

/*
..........................................................................................................................
. GoPathErrorType go-err-path 错误类型的定义
. 这个结构体是为了用户基于 GoPathErrorType 定义自己的错误类型时能够嵌入该类型定义，进而能够方便的使用 GoPathErrorType 下定义的方法。
. 同时用户在自定义错误类型中嵌入该结构体能够帮助 go-err-path 快速锁定错误类型所在的包。
.
. baseErr     用户设定的基础错误
. shortError  用户设定的对外输出的简短错误
. pkgPath     用户所定义的错误所在的包路径，有go-path-err自动解析出来
.
. Init           初始化方法，用户在定义完错误类型后，应立即在当前文件下的 init() 方法中完成对Init()的调用，已完成错误信息及包路径初始化
. ParseError     从一个给定错误中生成一个 GoPathErrorItem 错误
. CombineErrors  采用 | 将多个错误合并为1个，如果错误都为空则返回 nil
. IsNilErr       判断 GoPathErrorType 是否为空错误
. MergeError     多级错误合并
..........................................................................................................................
*/
type GoPathErrorType struct {
	baseErr    error  // 错误信息
	shortError error  // 对外输出错误信息
	pkgPath    string // 包路径
}

// BError 获取基础错误信息
func (g *GoPathErrorType) BError() error {
	return g.baseErr
}

// Init 初始化包路径
func (g *GoPathErrorType) Init(err interface{}, bserr string, sherr string) {
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
	g.shortError = errors.New(sherr)
	g.baseErr = errors.New(bserr)
}

// ParseError 格式化错误
func (g *GoPathErrorType) ParseError(dwt string, errs ...error) ErrorItem {
	err := g.CombineErrors(errs...)
	if err == nil {
		return nil
	}

	err, perr := g.parseError(dwt, err)

	return &GoPathErrorItem{
		baseError:   g.baseErr,
		shortError:  g.shortError,
		detailError: err,
		pathError:   perr,
	}
}

// CombineErrors 合并错误
func (g *GoPathErrorType) CombineErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	msgs := []string{}
	for _, err := range errs {
		if err == nil {
			continue
		}

		msgs = append(msgs, err.Error())
	}

	if len(msgs) == 0 {
		return nil
	}

	return errors.New(strings.Join(msgs, "|"))
}

// parseError 格式化错误
func (g *GoPathErrorType) parseError(dwt string, errs ...error) (error, error) {
	err := g.CombineErrors(errs...)
	if err == nil {
		return nil, nil
	}

	berr := g.baseErr.Error()

	return fmt.Errorf("%s%s ", berr+" : ", err.Error()), fmt.Errorf("%s%s : %s ", g.pkgPath, dwt, berr)
}

// parseItemError 格式化ErrorItem
func (g *GoPathErrorType) parseItemError(dwt string, err ErrorItem) (error, error) {
	berr := g.baseErr.Error()

	return fmt.Errorf("%s%s ", berr+" : ", err.DetailError()), fmt.Errorf("%s%s : %s : %s", g.pkgPath, dwt, g.baseErr.Error(), err.PathError().Error())
}

// IsNilErr 是否是空错误
func (g *GoPathErrorType) IsNilErr(err ErrorItem) bool {
	return err == nil || err.IsNil()
}

// MergeError 合并错误
func (g *GoPathErrorType) MergeError(dwt string, err ErrorItem) ErrorItem {
	if g.IsNilErr(err) {
		return nil
	}

	dterr, perr := g.parseItemError(dwt, err)

	return &GoPathErrorItem{
		baseError:   g.baseErr,
		shortError:  g.shortError,
		detailError: dterr,
		pathError:   perr,
	}
}
