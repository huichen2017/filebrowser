//go:generate go-bindata -pkg assets -o assets/assets.go templates/ assets/css/ assets/js/ assets/fonts/

package hugo

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/hacdias/caddy-hugo/assets"
	"github.com/hacdias/caddy-hugo/browse"
	"github.com/hacdias/caddy-hugo/edit"
	"github.com/hacdias/caddy-hugo/settings"
	"github.com/hacdias/caddy-hugo/utils"
	"github.com/mholt/caddy/config/setup"
	"github.com/mholt/caddy/middleware"
	"github.com/spf13/hugo/commands"
)

// Setup function
func Setup(c *setup.Controller) (middleware.Middleware, error) {
	commands.Execute()

	return func(next middleware.Handler) middleware.Handler {
		return &handler{Next: next}
	}, nil
}

type handler struct{ Next middleware.Handler }

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if middleware.Path(r.URL.Path).Matches("/admin") {
		page := utils.ParseComponents(r)[1]

		log.Print(page)

		if page == "assets" {
			filename := strings.Replace(r.URL.Path, "/admin/", "", 1)
			file, err := assets.Asset(filename)

			if err != nil {
				return 404, nil
			}

			extension := filepath.Ext(filename)
			mime := mime.TypeByExtension(extension)

			header := w.Header()
			header.Set("Content-Type", mime)

			w.Write(file)
			return 200, nil
		} else if page == "browse" {
			return browse.Execute(w, r)
		} else if page == "edit" {
			return edit.Execute(w, r)
		} else if page == "settings" {
			return settings.Execute(w, r)
		}

		return 404, nil
	}

	return h.Next.ServeHTTP(w, r)
}
