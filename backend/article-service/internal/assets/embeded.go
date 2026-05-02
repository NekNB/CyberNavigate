package assets

import "embed"

// Внедряем сам UI
//go:embed swagger/*
var SwaggerUI embed.FS
