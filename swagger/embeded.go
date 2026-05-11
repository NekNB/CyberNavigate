package swagger

import "embed"

// Внедряем папки со спецификациями
//
//go:embed docs/article-service/*
var SpecsFS embed.FS
