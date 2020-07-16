package controllers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/logger"
	"github.com/dryairship/IITKBucks/models"
)

func newBlockHandler(c *gin.Context) {
	var body []byte
	numBytes, err := c.Request.Body.Read(body)
	if err != nil || numBytes == 0 {
		c.String(400, "Error while reading request body")
		return
	}

	block, err := models.BlockFromByteArray(body)
	if err != nil {
		c.String(400, "Given bytes could not be converted to a block")
		return
	}

	isValid, err := models.Blockchain().IsBlockValid(&block)
	if !isValid {
		c.String(400, err.Error())
		return
	}

	logger.Println(logger.MajorEvent, "[Controllers/Block] [INFO] Valid new block received with index", block.Index)
	performPostNewBlockSteps(block)
	c.String(200, "Block added to the blockchain")
}

func propagateBlockToPeers(blockBytes []byte) {
	buffer := bytes.NewBuffer(blockBytes)
	client := &http.Client{}

	count := 0

	for _, peer := range peers {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/newBlock", peer), buffer)
		if err != nil {
			logger.Println(logger.RareError, "[Controllers/Block] [WARN] go HTTP error while builing newBlock request. Peer: ", peer, ", Error: ", err)
			continue
		}
		req.Header.Set("Content-Type", "application/octet-stream")

		resp, err := client.Do(req)
		if err != nil {
			logger.Println(logger.CommonError, "[Controllers/Block] [WARN] go HTTP error while making newBlock request. Peer: ", peer, ", Error: ", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var reply []byte
			_, _ = resp.Body.Read(reply)
			logger.Printf(logger.RareError, "[Controllers/Block] [WARN] Peer %s gave %d response on newBlock request. %s\n", peer, resp.StatusCode, reply)
		} else {
			count++
		}
	}
	logger.Printf(logger.MajorEvent, "[Controllers/Block] [DEBUG] Successfully sent newBlock requests to %d out of %d peers.\n", count, len(peers))
}

func performPostNewBlockSteps(newBlock models.Block) {
	logger.Println(logger.MinorEvent, "[Controllers/Block] [TRACE] Performing post new block steps.")

	close(currentMinerChannel)
	currentMinerChannel = make(chan bool)

	blockBytes := newBlock.ToByteArray()

	models.Blockchain().ProcessBlock(newBlock)
	models.Blockchain().AppendBlock(newBlock)

	go startMining()
	go propagateBlockToPeers(blockBytes)
}
