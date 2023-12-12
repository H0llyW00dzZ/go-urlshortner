package handlers

import (
	"net/http"

	cloudDatastore "cloud.google.com/go/datastore"
	localDatastore "github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/shortid"
	"github.com/gin-gonic/gin"
)

func RegisterHandlersGin(router *gin.Engine, dsClient *cloudDatastore.Client) {
	router.GET("/:id", getURLHandlerGin(dsClient))
	router.POST("/", postURLHandlerGin(dsClient))
}

func getURLHandlerGin(dsClient *cloudDatastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		key := cloudDatastore.NameKey("urlz", id, nil)
		var url localDatastore.URL
		if err := dsClient.Get(c, key, &url); err != nil {
			if err == cloudDatastore.ErrNoSuchEntity {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "URL not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.Redirect(http.StatusFound, url.Original)
	}
}

func postURLHandlerGin(dsClient *cloudDatastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			URL string `json:"url"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		id, err := shortid.Generate(10) // Generate a 10-character ID
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
			return
		}

		key := cloudDatastore.NameKey("urlz", id, nil)
		url := localDatastore.URL{
			Original: req.URL,
			ID:       id,
		}
		if _, err := dsClient.Put(c, key, &url); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}
