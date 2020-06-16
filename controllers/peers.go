package controllers

import (
	"github.com/gin-gonic/gin"
)

type newPeerRequestBody struct {
	Url string `json:"url" binding:"required"`
}

var peers []string

func getPeersHandler(c *gin.Context) {
	if peers != nil {
		c.JSON(200, gin.H{"peers": peers})
	} else {
		c.JSON(200, gin.H{})
	}
}

func newPeerHandler(c *gin.Context) {
	var body newPeerRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(400, err)
	}

	peers = append(peers, body.Url)
	c.Status(200)
}
