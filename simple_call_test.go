package math_calculation

import (
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {
	res, err := Calculate("(867255+-440375)-426878", nil, nil)

	fmt.Println(res.String())
	fmt.Println(err)
}
