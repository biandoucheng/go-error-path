# go-error-path

## 概述
```
它将运行中的一些列可知错误已路径+错误信息的形式表达出来，方便错误排查
```

## 使用方法
[示例代码 <https://github.com/biandoucheng/open-example/tree/main/go-error-path-example>](https://github.com/biandoucheng/open-example/tree/main/go-error-path-example)
```
// 在某个包中添加 error.go 内容如下
package pkga

import (
	gperr "github.com/biandoucheng/go-error-path"
)

type PkgAErrorType struct {
	gperr.GoPathErrorType
}

var pkgAError *PkgAErrorType

func init() {
	pkgAError = &PkgAErrorType{}
	pkgAError.Init(pkgAError, "Wrong A")
}

// 在包内的其他文件中可直接使用 pkgAError 对应的错误信息格式化方法。如：
package pkgb

import (
	"go-error-path-example/pkga"
)

func FuncB() error {
	err := pkga.FuncA()
	if err != nil {
		return pkgBError.ParsePkgDwtErr("FuncB", err)
	}

	return nil
}

// 输出的错误信息形式如：
// pkgb.FuncB Error: Wrong B : pkga.FuncA Error: Wrong A : faild run func A
```