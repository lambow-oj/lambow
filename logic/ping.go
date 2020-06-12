package logic

import (
	"fmt"
	"git/lambow-oj/lambow/cinex"
)

func Add(ctx *cinex.APIContext, a int64, b int64) (int64, error) {
	c := a + b
	if c == 0 {
		err := fmt.Errorf("invaild num %v %v", a, b)

		return 0, err
	}

	return c, nil
}
