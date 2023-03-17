package main

import (
	"WeddingBackEnd/app/api"
	"WeddingBackEnd/ultilities"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewAPIv1(container *Container) http.Handler {
	router := api.NewRouter()
	v1 := router.Group("/api/v1")

	authRouter(v1)

	facebookRouter(v1)

	googleRouter(v1)

	v1.Use(
		ultilities.Midlleware(container.Config.JwtKey),
		func(handle httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				handle(w, r, p)
			}
		},
	)

	userRouter(v1)

	return router
}

func authRouter(parent *api.Router) {
	authHandler := api.AuthHandler{
		JwtSecret:          container.Config.JwtKey,
		AccountRepository:  container.AccountRepository,
		AccessIDRepository: container.UserAccessIDRepository,
		UserRepository:     container.UserRepository,
	}

	router := parent.Group("/auth")
	router.POST("/sign-in", authHandler.SignIn)
	router.POST("/sign-up", authHandler.SignUp)
	router.POST("/verify:headers", authHandler.VerifyToken)
}

func userRouter(parent *api.Router) {
	userHandler := api.UserHandler{
		UserRepository: container.UserRepository,
	}
	router := parent.Group("/user")
	router.GET("/:id", userHandler.GetByID)

	router.GET("", userHandler.GetAll)

	router.POST("", userHandler.Add)

	router.PUT("/update/:id", userHandler.UpdateByID)

	router.DELETE("/deletebyid/:id", userHandler.RemoveByID)
	router.DELETE("/deletebyphone/:phonenumber", userHandler.RemoveByPhoneNumber)

}

func facebookRouter(parent *api.Router) {
	facebookHandler := api.FacebookHandler{
		UserRepository: container.UserRepository,
	}
	router := parent.Group("/auth")
	router.POST("/facebook/sign-in", facebookHandler.FacebookLogin)
}

func googleRouter(parent *api.Router) {
	googleHandler := api.GoogleHandler{
		UserRepository: container.UserRepository,
	}
	router := parent.Group("/auth")
	router.POST("/google/sign-in", googleHandler.HandleGoogleCallBack)
}
