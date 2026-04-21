package swagger

import (
	"io/fs"
	"net/http"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/assets"
)

func NewSwagger() http.Handler {

	uiDir, err := fs.Sub(assets.SwaggerUI, "swagger")
	if err != nil {
		panic(err)
	}
	uiServer := http.FileServer(http.FS(uiDir))

	return uiServer
}
