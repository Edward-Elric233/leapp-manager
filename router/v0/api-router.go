package v0

import (
	"gin-template/controller"
	"gin-template/middleware"
	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.GlobalAPIRateLimit())
	{
		apiRouter.GET("/status", controller.GetStatus)
		apiRouter.GET("/notice", controller.GetNotice)
		apiRouter.GET("/verification", middleware.CriticalRateLimit(), controller.SendEmailVerification)
		apiRouter.GET("/reset_password", middleware.CriticalRateLimit(), controller.SendPasswordResetEmail)
		apiRouter.GET("/user/reset", controller.SendNewPasswordEmail)
		apiRouter.GET("/oauth/github", controller.GitHubOAuth)

		userRoute := apiRouter.Group("/user")
		{
			userRoute.POST("/register", middleware.CriticalRateLimit(), controller.Register)
			userRoute.POST("/login", middleware.CriticalRateLimit(), controller.Login)
			userRoute.GET("/logout", controller.Logout)

			selfRoute := userRoute.Group("/")
			selfRoute.Use(middleware.NoTokenAuth(), middleware.UserAuth())
			{
				selfRoute.GET("/self", controller.GetSelf)
				selfRoute.PUT("/self", controller.UpdateSelf)
				selfRoute.DELETE("/self", controller.DeleteSelf)
				selfRoute.GET("/token", controller.GenerateToken)
			}

			adminRoute := userRoute.Group("/")
			adminRoute.Use(middleware.NoTokenAuth(), middleware.AdminAuth())
			{
				adminRoute.GET("/", controller.GetAllUsers)
				adminRoute.GET("/:id", controller.GetUser)
				adminRoute.POST("/", controller.CreateUser)
				adminRoute.POST("/manage", controller.ManageUser)
				adminRoute.PUT("/", controller.UpdateUser)
				adminRoute.DELETE("/:id", controller.DeleteUser)
			}
		}
		optionRoute := apiRouter.Group("/option")
		optionRoute.Use(middleware.NoTokenAuth(), middleware.RootAuth())
		{
			optionRoute.GET("/", controller.GetOptions)
			optionRoute.PUT("/", controller.UpdateOption)
		}
		fileRoute := apiRouter.Group("/file")
		{
			fileRoute.GET("/:id", middleware.DownloadRateLimit(), controller.DownloadFile)
			fileRoute.POST("/", middleware.UserAuth(), middleware.UploadRateLimit(), controller.UploadFile)
			fileRoute.DELETE("/:id", middleware.UserAuth(), controller.DeleteFile)
		}
	}
}
