package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Map(data []int, action func(int) int) []int {
	for i := range len(data) {
		data[i] = action(data[i])
	}
	return data
}

func Filter(data []int, action func(int) bool) []int {
	// вот тут не знаю как покрасивше сделать
	var resultData []int
	if data != nil {
		resultData = make([]int, 0)
	}
	/////////////////////////////////////////
	for _, element := range data {
		if action(element) {
			resultData = append(resultData, element)
		}
	}
	return resultData
}

func Reduce(data []int, initial int, action func(int, int) int) int {
	result := initial
	for _, element := range data {
		result = action(result, element)
	}
	return result
}

func TestMap(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) int
		result []int
	}{
		"nil numbers": {
			action: func(number int) int {
				return -number
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) int {
				return -number
			},
			result: []int{},
		},
		"inc numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number + 1
			},
			result: []int{2, 3, 4, 5, 6},
		},
		"double numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number * number
			},
			result: []int{1, 4, 9, 16, 25},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Map(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestFilter(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) bool
		result []int
	}{
		"nil numbers": {
			action: func(number int) bool {
				return number == 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) bool {
				return number == 1
			},
			result: []int{},
		},
		"even numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) bool {
				return number%2 == 0
			},
			result: []int{2, 4},
		},
		"positive numbers": {
			data: []int{-1, -2, 1, 2},
			action: func(number int) bool {
				return number > 0
			},
			result: []int{1, 2},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Filter(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestReduce(t *testing.T) {
	tests := map[string]struct {
		initial int
		data    []int
		action  func(int, int) int
		result  int
	}{
		"nil numbers": {
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"sum of numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 15,
		},
		"sum of numbers with initial value": {
			initial: 10,
			data:    []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 25,
		},
		// доп тесты
		"multiply": {
			initial: 2,
			data:    []int{1, 2, 1},
			action: func(lhs, rhs int) int {
				return lhs * rhs
			},
			result: 4,
		},
		"divide zero": {
			data: []int{10, 100, 23},
			action: func(lhs, rhs int) int {
				return lhs / rhs
			},
			result: 0,
		},
		"divide": {
			initial: 36,
			data:    []int{6, 3, 2},
			action: func(lhs, rhs int) int {
				return lhs / rhs
			},
			result: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Reduce(test.data, test.initial, test.action)
			assert.Equal(t, test.result, result)
		})
	}
}
