package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryCache(t *testing.T) {
	cache := newCategoryCache()

	cache.set(200, "description")
	description, exist := cache.get(200)

	assert.Equal(t, "description", description)
	assert.Equal(t, true, exist)

	description, exist = cache.get(404)
	assert.Equal(t, "", description)
	assert.Equal(t, false, exist)
}

func TestInSlice(t *testing.T) {
	s := []int{1, 2, 3}

	exist := _inSlice(s, 1)
	assert.Equal(t, true, exist)

	exist = _inSlice(s, 0)
	assert.Equal(t, false, exist)
}

func TestAggregateCodes(t *testing.T) {
	m := map[int]string{1: "test1", 2: "test2", 3: "test3"}

	result := _aggregateCodes(m)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestResolveCodeDescription(t *testing.T) {
	_t := true
	_f := false

	// Options for filtering all categories
	opts := &statusFilter{
		showAll:       &_t,
		onlyInfo:      &_f,
		onlySuccess:   &_f,
		onlyClientErr: &_f,
		onlyServerErr: &_f,
	}

	want := map[int]string{
		100: "Continue",
		200: "OK",
		300: "Multiple Choices",
		400: "Bad Request",
		500: "Internal Server Error",
		999: "No such code",
	}

	for k, v := range want {
		if k != 999 {
			description, exist := resolveCodeDescription(k, opts)
			assert.Equal(t, v, description)
			assert.Equal(t, true, exist)
		} else {
			// There is no such code, UnknownCategory should be returned
			description, exist := resolveCodeDescription(k, opts)
			assert.Equal(t, "", description)
			assert.Equal(t, false, exist)
		}
	}

	// As 200 code will be resolved second time, then cache should be used
	description, exist := resolveCodeDescription(200, opts)
	assert.Equal(t, "OK", description)
	assert.Equal(t, true, exist)

	// Clear cache
	cache = newCategoryCache()

	// Options for filtering only success category
	opts = &statusFilter{
		showAll:       &_f,
		onlyInfo:      &_f,
		onlySuccess:   &_t,
		onlyClientErr: &_f,
		onlyServerErr: &_f,
	}

	for k, v := range want {
		if k == 200 {
			// Only 200 code should be found
			description, exist := resolveCodeDescription(k, opts)
			assert.Equal(t, v, description)
			assert.Equal(t, true, exist)
		} else {
			description, exist := resolveCodeDescription(k, opts)
			assert.Equal(t, "", description)
			assert.Equal(t, false, exist)
		}
	}
}
