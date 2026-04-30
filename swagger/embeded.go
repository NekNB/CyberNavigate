package swagger

import "embed"

// Внедряем папки со спецификациями
//
//go:embed article-service/*
//go:embed user-service/*
var SpecsFS embed.FS
