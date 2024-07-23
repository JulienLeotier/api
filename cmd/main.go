package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"

	moduleRouter "quest/routers"
	utils "quest/utils"
)

func main() {
	utils.SetupDatabase()
	key := os.Getenv("SESSION_SECRET")
	maxAge := 86400 * 30
	isProd := false
	// Configure CORS

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), os.Getenv("GOOGLE_REDIRECT_URI")),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), os.Getenv("FACEBOOK_REDIRECT_URI")),
	)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Remplacez par l'IP ou l'origine que vous souhaitez autoriser
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	moduleRouter.SetupModuleRoutes(r, utils.DB)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
