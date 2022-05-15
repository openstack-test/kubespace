package alert

import (
	"fmt"
	"testing"
)

func TestCalculateReversePolishNotation(t *testing.T) {
	t1 := "level=info&app=java"
	t2 := make(map[string]string)
	t2["level"] = "info"
	t2["app"] = "java"

	n, err := BuildTree(t1)
	if err != nil {
		fmt.Println("err ==>", err)
	}
	nodeStr := Converse2ReversePolishNotation(n)

	if CalculateReversePolishNotation(t2, nodeStr) {
		fmt.Println(true)
	} else {
		fmt.Println(false)
	}
}
