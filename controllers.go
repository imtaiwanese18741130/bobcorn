package main

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	*Config
}

func NewController(cfg *Config) *Controller {
	return &Controller{
		Config: cfg,
	}
}

func (ctrl Controller) Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"BaseURL": ctrl.AppBaseURL,
		})
	}
}

func (ctrl Controller) Questions() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "questions.html", gin.H{
			"BaseURL": ctrl.AppBaseURL,
		})
	}
}

func (ctrl Controller) Results() gin.HandlerFunc {
	return func(c *gin.Context) {
		quizResult := c.Param("quiz_result")
		isValid, isOpposed := ctrl.vaildateQuizResult(quizResult)
		if !isValid {
			ctrl.NotFound()(c)
			return
		}

		isOpposedHidden := ""
		if isOpposed {
			isOpposedHidden = "hidden"
		}

		c.HTML(http.StatusOK, "results.html", gin.H{
			"BaseURL":         ctrl.AppBaseURL,
			"QuizResult":      quizResult,
			"IsOpposedHidden": isOpposedHidden,
		})
	}
}

func (ctrl Controller) vaildateQuizResult(quizResult string) (bool, bool) {
	// format : {area}-{type}
	allowedList := []string{
		"1-A", "1-B", "1-C", "1-D", "1-E", "1-F", "1-G", "1-H",
		"2-A", "2-B", "2-C", "2-D", "2-E", "2-F", "2-G", "2-H",
		"3-A", "3-B", "3-C", "3-D", "3-E", "3-F", "3-G", "3-H",
		"4-A", "4-B", "4-C", "4-D", "4-E", "4-F", "4-G", "4-H",
	}

	isValid := slices.Contains(allowedList, quizResult)

	isOpposed := false
	if isValid {
		r := strings.Split(quizResult, "-")
		if len(r) == 2 {
			if r[1] == "A" || r[1] == "B" {
				isOpposed = true
			}
		}
	}

	return isValid, isOpposed
}

func (ctrl Controller) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "v0.0.1"})
	}
}

func (ctrl Controller) NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "4xx.html", GetViewHttpError(http.StatusNotFound, "您請求的頁面不存在", ctrl.AppBaseURL, ctrl.AppBaseURL))
	}
}

func (ctrl Controller) GetAsset() gin.HandlerFunc {
	return func(c *gin.Context) {
		up := RequestURIAsset{}
		if err := c.ShouldBindUri(&up); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		filePath := fmt.Sprintf("./assets/%s/%s", up.Type, up.File)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		if ctrl.AppEnv == "production" {
			c.Header("Cache-Control", "public, max-age=3600")
		} else {
			c.Header("Cache-Control", "no-cache")
		}

		c.File(filePath)
	}
}

type RequestURIAsset struct {
	Type string `uri:"type" binding:"required"`
	File string `uri:"file" binding:"required"`
}
