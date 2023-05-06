package api

import "net/http"

type Router struct{}

func (g *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
