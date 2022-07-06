/**
* @Author: tc.lam
* @Date: 2022/7/6 下午11:12
* @Software : GoLand
* @File: wrap
* @Description:
* @Return:
**/

package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

func wrap(err error, message string, skip int) error {
	if err == nil {
		return nil
	}
	return &wrapError{
		message: message,
		next:    err,
		frame:   xerrors.Caller(skip),
	}
}

// Wrap returns a error annotating `err` with `message` and the caller's frame.
// Wrap returns nil if `err` is nil.
func Wrap(err error, message string) error {
	return wrap(err, message, 2)
}

// Wrapf returns a error annotating `err` with `message` formatted and the caller's frame.
func Wrapf(err error, message string, args ...interface{}) error {
	return wrap(err, fmt.Sprintf(message, args...), 2)
}

// ErrorStack formats details for `err` (with frame info).
func ErrorStack(err error) string {
	return fmt.Sprintf("%+v", err)
}

var (
	Annotate  = Wrap
	Annotatef = Wrapf

	Unwrap = xerrors.Unwrap
	Opaque = xerrors.Opaque
	As     = xerrors.As
	Is     = xerrors.Is
	Errorf = xerrors.Errorf
	New    = xerrors.New
)
