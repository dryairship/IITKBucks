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
		c.String(400, "Invalid JSON request body")
		return
	}
	_, exists := aliasMap[body.Alias]
	if exists {
		c.String(400, "Alias already exists")
		return
	}
	aliasMap[body.Alias] = body.PublicKey
	c.String(200, "Successfully added alias")
}

func getPublicKeyHandler(c *gin.Context) {
	var body aliasRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		c.String(400, "Invalid JSON request body")
		return
	}
	_, exists := aliasMap[body.Alias]
	if !exists {
		c.String(404, "Alias not related to any public key")
		return
	}
	c.JSON(200, gin.H{
		"publicKey": aliasMap[body.Alias],
	})
}
