package tnyrtr

import (
	"net/http"
	"os"
	"path"
	"syscall"

	"github.com/razonyang/fastrouter"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

/* Router is a wrapper around a net/http server and the razonyang/fastrouter to tie it together
with a logger and the shutdown channel */
type Router struct {
	root     *fastrouter.Router
	omux     http.Handler
	base     string
	shutdown chan os.Signal
	mws      []Middleware
}

// New returns a new router.
func New(shutdown chan os.Signal, mws ...Middleware) *Router {
	mux := fastrouter.New()
	return &Router{
		root:     mux,
		omux:     otelhttp.NewHandler(mux, "request"),
		base:     "/",
		shutdown: shutdown,
		mws:      mws,
	}
}

// SignalShutdown sends the app wide terminate signal
func (r *Router) SignalShutdown() {
	r.shutdown <- syscall.SIGTERM
}

// Group returns a new group
func (r *Router) Group(sub string, mws ...Middleware) *Router {
	return &Router{
		mws:      append(mws, r.mws...),
		base:     path.Join(r.base, sub),
		root:     r.root,
		shutdown: r.shutdown,
	}
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.omux.ServeHTTP(w, req)
}

// DELETE is a shortcut for group.Handle("DELETE", path, handler)
func (r *Router) DELETE(pth string, handler Handler, mws ...Middleware) {
	r.root.Delete(path.Join(r.base, pth), newHandle(r, handler, mws...))
}

// GET is a shortcut for group.Handle("GET", path, handler)
func (r *Router) GET(pth string, handler Handler, mws ...Middleware) {
	r.root.Get(path.Join(r.base, pth), newHandle(r, handler, mws...))
}

// PATCH is a shortcut for group.Handle("PATCH", path, handler)
func (r *Router) PATCH(pth string, handler Handler, mws ...Middleware) {
	r.root.Handle(http.MethodPatch, path.Join(r.base, pth), newHandle(r, handler, mws...))
}

// PUT is a shortcut for group.Handle("PUT", path, handler)
func (r *Router) PUT(pth string, handler Handler, mws ...Middleware) {
	r.root.Put(path.Join(r.base, pth), newHandle(r, handler, mws...))
}

// POST is a shortcut for group.Handle("POST", path, handler)
func (r *Router) POST(pth string, handler Handler, mws ...Middleware) {
	r.root.Post(path.Join(r.base, pth), newHandle(r, handler, mws...))
}

// Prepare perpares the registered routes
func (r *Router) Prepare() {
	r.root.Prepare()
}

// GetParams retrieves the named params from a request
func GetParams(r *http.Request) map[string]string {
	return fastrouter.Params(r)
}
