package assets

import "embed"

// Внедряем папку со спецификациями
//
//go:embed docs/*
var SpecsFS embed.FS

// Внедряем сам UI
//
//go:embed swagger/*
var SwaggerUI embed.FS
