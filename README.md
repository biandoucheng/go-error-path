# go-error-path

## 概述
```
它将运行中的一些列可知错误已路径+错误信息的形式表达出来，方便错误排查
```

## 使用方法
[示例代码 <https://github.com/biandoucheng/open-example/tree/main/go-error-path-example>](https://github.com/biandoucheng/open-example/tree/main/go-error-path-example)
```
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

// 输出的错误信息形式如：
功能测试
内部错误
功能测试 : 执行错误
biandoucheng/go-error-path/test.TestPath : 功能测试
```