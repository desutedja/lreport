package router

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func RewriteURL(router *mux.Router, baseURLPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = func(url string) string {
			// log.Println("url path: ", url)
			if strings.Index(url, baseURLPath) == 0 {
				url = url[len(baseURLPath):]
			}
			// log.Println("after rewrite: ", url)
			return url
		}(r.URL.Path)

		router.ServeHTTP(w, r)
	}
}
