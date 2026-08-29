package app

import "github.com/vandad1901/p3s/apps/api/internal/post"

func initializeServices(a *App) {
	a.PostService = post.NewService(a.db)
}
