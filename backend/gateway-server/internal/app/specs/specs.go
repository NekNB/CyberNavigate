package specs

import (
	"io/fs"
	"net/http"

	"github.com/NekNB/CyberNavigate/swagger"
)

func NewSpecs() http.Handler {

	// Открываем подсистему папки docs
	specsDir, err := fs.Sub(swagger.SpecsFS, "docs")
	if err != nil {
		panic(err)
	}

	// specsDir := os.DirFS("./internal/assets/docs")
	specServer := http.FileServer(http.FS(specsDir))

	return specServer
}
