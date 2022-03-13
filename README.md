# :package: tnyrtr

**tnyrtr** is a simple wrapper around
[razonyang/fastrouter](https://github.com/razonyang/fastrouter) to save some
boilerplate code.

It combines [razonyang/fastrouter](https://github.com/razonyang/fastrouter) with
convenience functions for sending and parsing JSON, a shutdown channel, a
context that offers access to the response status code, the start time of the
request, and a traceID (if provided on the request context), as well as the
ability to add middlewares.

## :computer: Example

```go
package main

import (
	"net/http"
	"os"

	"github.com/abecodes/tnyrtr"
	)

func main() {
	// create a shutdown channel
	// any handler can return a tnyrtr.
	shutdown := make(chan os.Signal)
	router := tnyrtr.New(shutdown)

	router.GET(
		"/hello",
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var user struct {
				name string
			}
			// parse the incoming body
			err := tnyrtr.ParseBodyJSON(r, &user)
			if err != nil {
				// smth went completly wrong and now we need to trigger a graceful server shutdown via the shutdown channel
				return tnyrtr.NewShutdownErr("oO")
			}

			// respond with a JSON
			return tnyrtr.RespondJSON(ctx, w, data, statusCode)
		},
	)

	// IMPORTANT: underlying fastrouter needs to prepare the registered routes
	tnyrtr.Prepare()

  // tnyrtr implements the http.Handler interface
	http.ListenAndServe("localhost:8080",r)

	// somewhere down the line handle the shutdown signal
	<-shutdown
}

```

## :compass: Named parameters

See the
[fastrouter documentation](https://pkg.go.dev/github.com/razonyang/fastrouter?utm_source=godoc#Parser.Parse)
for more detailed informations on parameter parsing.

```go
router.GET(
		"/hello/<world>",
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			name := tnyrtr.GetParams()["name"]

			fmt.Println(name)
			return nil
		},
	)


```

## :zap: Middlewares

A `router`, `group`, or `route` can recieve none or many _middlewares_. They
will be applied in the order they are passed.

_Middlewares_ are inherited: `router` ->`group` -> `route`

```go
func Logger(log *log.Logger) tnyrtr.Middleware {
	return func(next tnyrtr.Handler) tnyrtr.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			log.Info("REQUEST_START")
			defer func(start time.Time) {
				log.Infow(
					"REQUEST_DONE",
					"STATUS",
					tnyrtr.GetStatus(ctx),
					"DURATION",
					time.Since(start).Seconds(),
				)
			}(tnyrtr.GetTime(ctx))

			return next(ctx, w, r)
		}
	}
}

shutdown := make(chan os.Signal)
// create a router with the logging middleware
router := tnyrtr.New(shutdown, Logger)
// or apply it to a route directly
router.GET("/hello", handler, Logger)
```
