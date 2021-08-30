package main

import (
	"errors"
	"fmt"
	"net/http"
	"obelov.com/snippetbox/pkg/forms"
	"obelov.com/snippetbox/pkg/models"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because Pat matches the "/" path exactly, we can now remove the manual check
	// of r.URL.Path != "/" from this handler.
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create an instance of a templateData struct holding the slice of
	// snippets.
	data := &templateData{Snippets: s}

	// Use the new render helper.
	app.render(w, r, "home.page.tmpl", data)
}

// Add a showSnippet handler function.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function. If it can't
	// be converted to an integer, or the value is less than 1, we return a 404 page
	// not found response.
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{Snippet: s}

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl", data)
}

// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

// Add a createSnippet handler function.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Checking if the request method is a POST is now superfluous and can be
	// removed.

	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError helper to send
	// a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, we can use the Get() method to retrieve
	// the validated value for a particular form field.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Your snippet was saved
	// successfully!") and the corresponding key ("flash") to the session
	// data. Note that if there's no existing session for the current user
	// (or their session has expired) then a new, empty, session for them
	// will automatically be created by the session middleware.
	app.session.Put(r, "flash", "Snippet successfully created!")

	// Change the redirect to use the new semantic URL style of /snippet/:id
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using the form helper we made earlier.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	// Try to create a new user record in the database. If the email already exists
	// add an error message to the form and re-display it.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked and asking them to log in.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic error
	// message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Add the ID of the current user to the session, so that they are now 'logged in'.
	app.session.Put(r, "authenticatedUserID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.session.Remove(r, "authenticatedUserID")

	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "about.page.tmpl", nil)
}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")

	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "profile.page.tmpl", &templateData{
		User: user,
	})
}

func (app *application) changePasswordForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "password.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("currentPassword", "newPassword", "newPasswordConfirmation")
	form.MinLength("newPassword", 10)
	if form.Get("newPassword") != form.Get("newPasswordConfirmation") {
		form.Errors.Add("newPasswordConfirmation", "Passwords do not match")
	}

	if !form.Valid() {
		app.render(w, r, "password.page.tmpl", &templateData{Form: form})
		return
	}

	userID := app.session.GetInt(r, "authenticatedUserID")

	err = app.users.ChangePassword(userID, form.Get("currentPassword"), form.Get("newPassword"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("currentPassword", "Current password is incorrect")
			app.render(w, r, "password.page.tmpl", &templateData{Form: form})
		} else if err != nil {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your password has been updated!")
	http.Redirect(w, r, "/user/profile", 303)
}
