package tnyrtr

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

/* Handler handles incoming requests */
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

func newHandle(r *Router, handler Handler, mws ...Middleware) http.HandlerFunc {
	handler = wrapMiddleware(mws, handler)
	handler = wrapMiddleware(r.mws, handler)

	h := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		// Capture the parent request span from the context.
		span := trace.SpanFromContext(ctx)

		// Then the ctx would be enriched with the context values
		v := Values{
			TraceID: span.SpanContext().TraceID().String(),
			Now:     time.Now().UTC(),
		}
		ctx = context.WithValue(ctx, key, &v)

		if err := handler(ctx, w, req); err != nil {
			// If the error has not be handled before
			// eg in a middleware
			// the app will be considered compromised.
			fmt.Println("ROUTER ERROR", err)
			r.SignalShutdown()
		}
	}

	return h
}
