package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/models"
)

func newBlockHandler(c *gin.Context) {
	var body []byte
	numBytes, err := c.Request.Body.Read(body)
	if err != nil || numBytes == 0 {
		c.AbortWithStatus(400)
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

	log.Println("[INFO] Valid new block received with index", block.Index)
	performPostNewBlockSteps(block)
	c.Status(200)
}

func propagateBlockToPeers(block models.Block) {
	blockBytes := block.ToByteArray()
	buffer := bytes.NewBuffer(blockBytes)
	client := &http.Client{}

	count := 0

	for _, peer := range peers {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/newBlock", peer), buffer)
		if err != nil {
			log.Println("[ERROR] go HTTP error while builing newBlock request. Peer: ", peer, ", Error: ", err)
			continue
		}
		req.Header.Set("Content-Type", "application/octet-stream")

		resp, err := client.Do(req)
		if err != nil {
			log.Println("[ERROR] go HTTP error while making newBlock request. Peer: ", peer, ", Error: ", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var reply []byte
			_, _ = resp.Body.Read(reply)
			log.Printf("[WARNING] Peer %s gave %d response on newBlock request. %s\n", peer, resp.StatusCode, reply)
		} else {
			count++
		}
	}
	log.Printf("[INFO] Successfully sent newBlock requests to %d peers.\n", count)
}

func performPostNewBlockSteps(newBlock models.Block) {
	log.Println("[INFO] Performing post new block steps.")

	close(currentMinerChannel)
	currentMinerChannel = make(chan bool)

	models.Blockchain().ProcessBlock(newBlock)
	models.Blockchain().AppendBlock(newBlock)

	go startMining()
	go propagateBlockToPeers(newBlock)
}
