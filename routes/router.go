package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"main/controller"
	"main/global"
	"main/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.HandleMethodNotAllowed = true
	router.Use(gin.Recovery())
	router.Use(middleware.AddTraceId())

	// monitor
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(router)

	router.Use(middleware.LoggerMiddleware())

	rPing := router.Group("/ping")
	rPing.GET("/", controller.Ping)

	// p0 related, client
	if global.DeployParty() == global.DEPLOY_PARTY_P0 || global.DeployParty() == global.DEPLOY_PARTY_BOTH {
		rP0 := router.Group("/p0")
		{
			rP0.GET("keygen", controller.GetP0KeyAggByUserId)
			rP0.POST("keygen_round1", controller.P0KeyGenRound1)
			rP0.POST("sign_round1", controller.P0SignRound1)
			rP0.POST("verify", controller.P0Verify)
		}
	}

	// p1 related, server
	if global.DeployParty() == global.DEPLOY_PARTY_P1 || global.DeployParty() == global.DEPLOY_PARTY_BOTH {
		rP1 := router.Group("/p1")
		{
			rP1.POST("keygen_round1", controller.P1KeyGenRound1)
			rP1.POST("sign_round1", controller.P1SignRound1)
			rP1.POST("sign_round2", controller.P1SignRound2)
		}
	}

	return router
}
