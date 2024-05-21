package config

import (
	_ "embed"
)

//go:embed common.txt
var Dictfile string
