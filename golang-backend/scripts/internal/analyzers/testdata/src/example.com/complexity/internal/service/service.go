package service

import "context"

func TooManyParams(ctx context.Context, a, b, c, d, e int) int { // want `function "TooManyParams" has 5 parameters excluding context.Context\.`
	return a + b + c + d + e
}
