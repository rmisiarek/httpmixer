package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedColor(t *testing.T) {
	c1 := Red("test")
	c2 := _color("test", red, "linux")
	c3 := _color("test", red, "windows")

	assert.Equal(t, "\033[31mtest\033[0m", c1)
	assert.Equal(t, c1, c2)
	assert.Equal(t, "test", c3)
}

func TestGreenColor(t *testing.T) {
	c1 := Green("test")
	c2 := _color("test", green, "linux")
	c3 := _color("test", green, "windows")

	assert.Equal(t, "\033[32mtest\033[0m", c1)
	assert.Equal(t, c1, c2)
	assert.Equal(t, "test", c3)
}

func TestYellowColor(t *testing.T) {
	c1 := Yellow("test")
	c2 := _color("test", yellow, "linux")
	c3 := _color("test", yellow, "windows")

	assert.Equal(t, "\033[33mtest\033[0m", c1)
	assert.Equal(t, c1, c2)
	assert.Equal(t, "test", c3)
}

func TestBlueColor(t *testing.T) {
	c1 := Blue("test")
	c2 := _color("test", blue, "linux")
	c3 := _color("test", blue, "windows")

	assert.Equal(t, "\033[34mtest\033[0m", c1)
	assert.Equal(t, c1, c2)
	assert.Equal(t, "test", c3)
}

func TestPurpleColor(t *testing.T) {
	c1 := Purple("test")
	c2 := _color("test", purple, "linux")
	c3 := _color("test", purple, "windows")

	assert.Equal(t, "\033[35mtest\033[0m", c1)
	assert.Equal(t, c1, c2)
	assert.Equal(t, "test", c3)
}

func TestCyanColor(t *testing.T) {
	c1 := Cyan("test")
	c2 := _color("test", cyan, "linux")
	c3 := _color("test", cyan, "windows")

	assert.Equal(t, "\033[36mtest\033[0m", c1)
	assert.Equal(t, c1, c2)
	assert.Equal(t, "test", c3)
}

func TestGrayColor(t *testing.T) {
	c1 := Gray("test")
	c2 := _color("test", gray, "linux")
	c3 := _color("test", gray, "windows")

	assert.Equal(t, "\033[37mtest\033[0m", c1)
	assert.Equal(t, c1, c2)
	assert.Equal(t, "test", c3)
}

func TestWhiteColor(t *testing.T) {
	c1 := White("test")
	c2 := _color("test", white, "linux")
	c3 := _color("test", white, "windows")

	assert.Equal(t, "\033[97mtest\033[0m", c1)
	assert.Equal(t, c1, c2)
	assert.Equal(t, "test", c3)
}
