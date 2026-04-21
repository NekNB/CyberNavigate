package specs

import (
	"io/fs"
	"net/http"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/assets"
)

func NewSpecs() http.Handler {
	// Открываем подсистему папки docs
	specsDir, err := fs.Sub(assets.SpecsFS, "docs/proto")
	if err != nil {
		panic(err)
	}
	specServer := http.FileServer(http.FS(specsDir))

	return specServer
}
