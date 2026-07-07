package assets

import "embed"

// Внедряем сам UI
//
//go:embed swagger/*
var SwaggerUI embed.FS

//go:embed redoc/*
var RedocUI embed.FS
