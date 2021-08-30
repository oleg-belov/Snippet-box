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
	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	// Update these routes to use the new dynamic middleware chain followed
	// by the appropriate handler function.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// Add the five new routes.
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	mux.Get("/user/profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userProfile))
	mux.Get("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePasswordForm))
	mux.Post("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePassword))

	// Add a new GET /ping route.
	mux.Get("/ping", http.HandlerFunc(ping))
	mux.Get("/healthcheck", http.HandlerFunc(app.healthcheckHandler))

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
