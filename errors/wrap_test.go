/**
* @Author: tc.lam
* @Date: 2022/7/6 下午11:12
* @Software : GoLand
* @File: wrap_test
* @Description:
* @Return:
**/

package errors_test

import (
	"fmt"
	"github.com/tagol-org/dev/errors"
	"runtime"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	foo := errors.New("foo")
	bar := errors.Wrap(foo, "bar")
	if actual := bar.Error(); "bar: foo" != actual {
		t.Errorf("bar should be wrapped")
	}
}

func TestWrapNil(t *testing.T) {
	err := errors.Wrapf(nil, "something else")
	if err != nil {
		t.Errorf("err should be nil")
	}
}

type callInfo struct {
	File string
	Line int
}

func saveInfo(info *callInfo) {
	_, info.File, info.Line, _ = runtime.Caller(1)
	info.Line++
}
func foo(i int, info *callInfo) error {
	saveInfo(info)
	return errors.Errorf("foo: %d", i)
}

func TestWrapFrame(t *testing.T) {
	var inner, outer callInfo

	// wrap error returned from foo using `Wrap`.
	saveInfo(&outer)
	err := errors.Wrap(foo(42, &inner), "wrap foo")

	actual := errors.ErrorStack(err)
	t.Logf("error stack: %v", actual)

	innerMessage := "foo: 42"
	if !strings.Contains(actual, innerMessage) {
		t.Errorf("err should contain inner error message")
	}
	outerMessage := "wrap foo"
	if !strings.Contains(actual, outerMessage) {
		t.Errorf("err should contain outer error message")
	}
	innerInfo := fmt.Sprintf("%s:%d", inner.File, inner.Line)
	if !strings.Contains(actual, innerInfo) {
		t.Errorf("err should contain inner error frame info")
	}
	outerInfo := fmt.Sprintf("%s:%d", outer.File, outer.Line)
	if !strings.Contains(actual, outerInfo) {
		t.Errorf("err should contain outer error frame info")
	}

	// wrap err again using `Errorf`.
	var third callInfo
	saveInfo(&third)
	err = errors.Errorf("third: %w", err)
	actual = errors.ErrorStack(err)
	t.Logf("third error stack: %v", actual)

	if !strings.Contains(actual, innerMessage) {
		t.Errorf("err should contain inner error message")
	}
	if !strings.Contains(actual, outerMessage) {
		t.Errorf("err should contain outer error message")
	}
	if !strings.Contains(actual, innerInfo) {
		t.Errorf("err should contain inner error frame info")
	}
	if !strings.Contains(actual, outerInfo) {
		t.Errorf("err should contain outer error frame info")
	}
	thirdMessage := "third"
	if !strings.Contains(actual, thirdMessage) {
		t.Errorf("err should contain third error message")
	}
	thirdInfo := fmt.Sprintf("%s:%d", third.File, third.Line)
	if !strings.Contains(actual, thirdInfo) {
		t.Errorf("err should contain third error frame info")
	}
}
