package frontend

import (
	"embed"
)

//go:embed dist/*
var Static embed.FS
