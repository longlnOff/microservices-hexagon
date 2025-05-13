package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type httpServer struct {
	configuration configuration.Configuration
}


func routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(time.Second * 60))

	r.Route("/v1", func(r chi.Router) {

		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthcheckHandler)

		// docsURL := fmt.Sprintf("%s/swagger/doc.json",
		// 	app.configuration.Server.SERVER_PORT)
		// r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL((docsURL))))
		// Post API
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)

				r.Get("/", app.getPostsHandler)
				r.Patch("/", app.checkPostownership("moderator", app.updatePostHandler))
				r.Delete("/", app.checkPostownership("admin", app.deletePostHandler))

				// Create comment for post
				r.Post("/comments", app.createCommentHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// User API
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				// Follow & Unfollow
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
		})

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}
