package app

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

// RegisterPprof registers Go runtime profiling endpoints under /debug/pprof/
// with strict production access control. Profiling endpoints are disabled by default
// and must be explicitly enabled via environment configuration.
func RegisterPprof(router *gin.Engine, enablePprof bool, adminToken string) {
	group := router.Group("/debug/pprof")
	group.Use(pprofSecurityMiddleware(enablePprof, adminToken))
	{
		group.GET("/", gin.WrapF(pprof.Index))
		group.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		group.GET("/profile", gin.WrapF(pprof.Profile))
		group.POST("/symbol", gin.WrapF(pprof.Symbol))
		group.GET("/symbol", gin.WrapF(pprof.Symbol))
		group.GET("/trace", gin.WrapF(pprof.Trace))
		group.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		group.GET("/block", gin.WrapH(pprof.Handler("block")))
		group.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		group.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		group.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		group.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}

// pprofSecurityMiddleware enforces production access protection on profiling endpoints.
// 1. Returns 404 Not Found if profiling is disabled in configuration.
// 2. Returns 401 Unauthorized if an admin token is configured but missing/invalid in request.
func pprofSecurityMiddleware(enabled bool, adminToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enabled {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "pprof profiling endpoints are disabled in this environment",
			})
			return
		}

		if adminToken != "" {
			reqToken := c.GetHeader("X-Pprof-Token")
			if reqToken == "" {
				reqToken = c.Query("token")
			}
			if reqToken != adminToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized: invalid or missing pprof admin token",
				})
				return
			}
		}

		c.Next()
	}
}
