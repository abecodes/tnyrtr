package tnyrtr

// Middleware represents a function designed to run some code before/after another handler
type Middleware func(next Handler) Handler

func wrapMiddleware(mws []Middleware, handler Handler) Handler {
	/* looping backwards over the passed middlewares to ensure
	that the first passed middleware is executed first */
	for i := len(mws) - 1; i >= 0; i-- {
		if h := mws[i]; h != nil {
			handler = h(handler)
		}
	}

	return handler
}
