package swagger

import "embed"

// Внедряем папки со спецификациями
//
//go:embed docs/article-service/*
//go:embed docs/user-service/*
var SpecsFS embed.FS
