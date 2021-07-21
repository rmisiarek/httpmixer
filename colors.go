package main

import (
	"github.com/fatih/color"
)

var RedBold = color.New(color.FgHiRed).Add(color.Bold).SprintFunc()
var GreenBold = color.New(color.FgHiGreen).Add(color.Bold).SprintFunc()
var BlueBold = color.New(color.FgHiBlue).Add(color.Bold).SprintFunc()
var YellowBold = color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
var MagentaBold = color.New(color.FgHiMagenta).Add(color.Bold).SprintFunc()
var WhiteBold = color.New(color.FgWhite).Add(color.Bold).SprintFunc()

var Red = color.New(color.FgHiRed).SprintFunc()
var Green = color.New(color.FgHiGreen).SprintFunc()
var Blue = color.New(color.FgHiBlue).SprintFunc()
var Yellow = color.New(color.FgHiYellow).SprintFunc()
var Magenta = color.New(color.FgHiMagenta).SprintFunc()
var White = color.New(color.FgWhite).SprintFunc()
