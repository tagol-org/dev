/**
* @Author: tc.lam
* @Date: 2022/7/6 下午11:02
* @Software : GoLand
* @File: errors
* @Description:
* @Return:
**/

package errors

import (
	"fmt"
	"golang.org/x/xerrors"
)

type warpError struct {
	message string
	next    error
	frame   xerrors.Frame
}

func (e *warpError) Unwrap() error {
	return e.next
}

func (e *warpError) Error() string {
	if e.next == nil {
		return e.message
	}
	return fmt.Sprintf("%s:%v", e.message, e.next)
}

func (e *warpError) Format(f fmt.State, c rune) {
	xerrors.FormatError(e, f, c)
}

func (e *warpError) FormatError(p xerrors.Printer) error {
	p.Printf(e.message)
	if p.Detail() {
		e.frame.Format(p)
	}
	return e.next
}
