package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	// Update these routes to use the new dynamic middleware chain followed
	// by the appropriate handler function.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(app.config.staticDir)})

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)
}
