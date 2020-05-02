package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, err := cfg.Build()
	if err != nil {
		log.Fatal("error creating log", err)
	}
	l.Info("starting service")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	err = router.Run(":" + port)
	if err != nil {
		l.Error("closing server", zap.Error(err))
	}
}
