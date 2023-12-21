package gorocks

import "net/http"

type Middleware func(HandlerFunc) HandlerFunc

func ApplicationJSON(next HandlerFunc) HandlerFunc {
	return func(app *App) {
		app.SetHeader("Content-Type", "application/json")
		next(app)
	}
}

func BasicAuth(username, password string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(app *App) {
			un, pw, ok := app.HttpRequest().BasicAuth()
			if !ok {
				app.SetStatusCode(http.StatusUnauthorized)
				app.ResponseWriter().Write([]byte("unautorized"))

				return
			}

			if un != username && pw != password {
				app.SetStatusCode(http.StatusUnauthorized)
				app.ResponseWriter().Write([]byte("wrong username/password"))

				return
			}

			next(app)
		}
	}
}

func ResponseWithHeaders(headers map[string]string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(app *App) {
			for key, value := range headers {
				app.SetHeader(key, value)
			}

			next(app)
		}
	}
}
