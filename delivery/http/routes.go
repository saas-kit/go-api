package http

import (
	"net/http"

	"github.com/xujiajun/gorouter"
)

func setupRoutes(r *gorouter.Router) {
	// .. custom routes
	r.GET("/test", testHandler)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK test func"))
}
