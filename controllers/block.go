package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/models"
)

func newBlockHandler(c *gin.Context) {
	var body []byte
	numBytes, err := c.Request.Body.Read(body)
	if err != nil || numBytes == 0 {
		_ = c.AbortWithError(400, err)
	}

	block, err := models.BlockFromByteArray(body)
	if err != nil {
		_ = c.AbortWithError(400, err)
	}

	models.Blockchain().AppendBlock(block)
	c.Status(200)
}
