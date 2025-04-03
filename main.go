package main

import (
	"html/template"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	c := NewController(cfg)

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"urlencode": url.QueryEscape,
	})
	r.LoadHTMLGlob("templates/*.html")
	r.SetTrustedProxies(cfg.AppTrustedProxies)

	r.NoRoute(c.NotFound())
	r.GET("/health/v1/ping", c.Ping())
	r.GET("/assets/:type/:file", c.GetAsset())

	r.GET("/", c.Index())
	r.GET("/questions", c.Questions())
	r.GET("/results/:quiz_result", c.Results())

	r.Run(":" + cfg.AppPort)
}
