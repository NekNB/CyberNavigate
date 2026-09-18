package assets

import "embed"

//go:embed articles
var ArticlesFS embed.FS

//go:embed simulator/*.yaml
var SimulatorFS embed.FS
