package routes

import (
	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api/v1")

	auth := api.Group("", middleware.ValidateRBAC()) // middleware.CheckFlag()
	{
		authPublicRoutes := []string{
			"/auth/student/register",
			"/auth/student/login",
			"/auth/teacher/register",
			"/auth/teacher/login",
			"/auth/admin/login",
			"/auth/verify-otp",
			"/auth/resend-otp",
			"/auth/forgot-password",
			"/auth/reset-password",
		}

		authProtectedRoutes := []string{
			"/auth/logout",
			"/auth/refresh",
			"/admin/users",
			"/users/me",
			"/users/me/password",
			"/users/me/image",
			"/users",
			"/users/students",
			"/users/teachers",
			"/users/lead-teachers",
			"/users/by-email",
			"/users/:user_id",
			"/users/:user_id/role",
			"/users/:user_id/reactivate",
		}

		for _, route := range authPublicRoutes {
			auth.Any(route, handlers.ForwardToAuth(cfg))
		}

		for _, route := range authProtectedRoutes {
			auth.Any(route, middleware.ValidateAccessJWT(cfg), handlers.ForwardToAuth(cfg))
		}
	}

	exam := api.Group("/exam",
		middleware.ValidateAccessJWT(cfg),
		middleware.ValidateRBAC(),
		//middleware.CheckFlag(),
	)
	{
		exam.Any("/*path", handlers.ForwardToExamService(cfg))
	}

	classes := api.Group("/classes",
		middleware.ValidateAccessJWT(cfg),
		middleware.ValidateRBAC(),
		//middleware.CheckFlag(),
	)
	{
		classes.Any("/*path", handlers.ForwardToClassesService(cfg))
	}

}
