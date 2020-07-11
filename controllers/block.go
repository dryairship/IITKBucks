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
		return
	}

	block, err := models.BlockFromByteArray(body)
	if err != nil {
		_ = c.AbortWithError(400, err)
		return
	}

	isValid := models.Blockchain().IsBlockValid(&block)
	if !isValid {
		_ = c.AbortWithError(400, models.ERROR_INVALID_BLOCK)
		return
	}

	performPostNewBlockSteps(block)
	c.Status(200)
}

func performPostNewBlockSteps(newBlock models.Block) {
	models.Blockchain().ProcessBlock(newBlock)
	models.Blockchain().AppendBlock(newBlock)

	close(currentMinerChannel)
	currentMinerChannel = make(chan bool)

	go startMining()
}
