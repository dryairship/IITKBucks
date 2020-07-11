package controllers

import (
	"github.com/gin-gonic/gin"
)

type aliasRequestBody struct {
	Alias     string `json:"alias"`
	PublicKey string `json:"publicKey"`
}

var aliasMap map[string]string = make(map[string]string)

func addAliasHandler(c *gin.Context) {
	var body aliasRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	_, exists := aliasMap[body.Alias]
	if exists {
		c.AbortWithStatus(400)
		return
	}
	aliasMap[body.Alias] = body.PublicKey
	c.Status(200)
}

func getPublicKeyHandler(c *gin.Context) {
	var body aliasRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	_, exists := aliasMap[body.Alias]
	if !exists {
		c.AbortWithStatus(404)
		return
	}
	c.JSON(200, gin.H{
		"publicKey": aliasMap[body.Alias],
	})
}
