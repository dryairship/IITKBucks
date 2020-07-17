package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/config"
	"github.com/dryairship/IITKBucks/logger"
)

type newPeerRequestBody struct {
	Url string `json:"url" binding:"required"`
}

type getPeersResponseBody struct {
	Peers []string `json:"peers"`
}

var peers []string
var peerIPs = make(map[string]string)
var potentialPeers []string = config.POTENTIAL_PEERS

func getPeersHandler(c *gin.Context) {
	if peers != nil {
		c.JSON(200, gin.H{"peers": peers})
	} else {
		c.JSON(200, gin.H{"peers": make([]int, 0)})
	}
}

func getPeerIPsHandler(c *gin.Context) {
	c.JSON(200, peerIPs)
}

func newPeerHandler(c *gin.Context) {
	if len(peers) == config.MAX_PEERS {
		c.String(500, "Max peer limit reached")
		return
	}

	var body newPeerRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		c.String(400, "Invalid JSON request body")
		return
	}

	peerIPs[c.ClientIP()] = body.Url
	logger.Println(logger.MinorEvent, "[Controllers/Peers] [INFO] Received peer add request. URL:", body.Url, ", IP:", c.ClientIP())

	if isAlreadyAPeer(body.Url) {
		c.String(200, "Peer has already been added")
		return
	}

	if strings.Contains(body.Url, "127.0.0.1") || strings.Contains(body.Url, "localhost") || strings.Contains(config.MY_URL, body.Url) {
		c.String(400, "Invalid Peer URL")
		return
	}

	peers = append(peers, body.Url)
	c.String(200, "Successfully added peer")
}

func makeGetPeersRequest(peer string) {
	if strings.Contains(peer, "127.0.0.1") || strings.Contains(peer, "localhost") || strings.Contains(config.MY_URL, peer) {
		return
	}

	response, err := http.Get(fmt.Sprintf("%s/getPeers", peer))
	if err != nil {
		logger.Println(logger.RareError, "[Controllers/Peers] [WARN] go HTTP error while asking for peers. Peer:", peer, ", Error:", err)
		return
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Println(logger.RareError, "[Controllers/Peers] [WARN] Cannot read getPeers response body. Peer:", peer, ". Error:", err)
		return
	}

	var body getPeersResponseBody
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		logger.Println(logger.RareError, "[Controllers/Peers] [WARN] Cannot unmarshal getPeers response body. Peer:", peer, ", Error:", err)
		return
	}

	potentialPeers = append(potentialPeers, body.Peers...)
}

func makeNewPeerRequest(peer string) bool {
	if strings.Contains(peer, "127.0.0.1") || strings.Contains(peer, "localhost") || strings.Contains(config.MY_URL, peer) {
		return false
	}

	var jsonStr = []byte(fmt.Sprintf(`{"url":"%s"}`, config.MY_URL))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/newPeer", peer), bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Println(logger.RareError, "[Controllers/Peers] [WARN] go HTTP error while builing newPeer request. Peer:", peer, ", Error:", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(logger.CommonError, "[Controllers/Peers] [WARN] go HTTP error while making newPeer request. Peer:", peer, ", Error:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		peers = append(peers, peer)
		logger.Println(logger.CommonError, "[Controllers/Peers] [INFO] Added new peer by sending request:", peer)
		return true
	} else {
		reason, _ := ioutil.ReadAll(resp.Body)
		logger.Println(logger.CommonError, "[Controllers/Peers] [ERROR] newPeer request rejected. Reason:", string(reason), ", Peer:", peer, ", Code:", resp.StatusCode)
		return false
	}
}

func isAlreadyAPeer(peer string) bool {
	for i := range peers {
		if peer == peers[i] {
			return true
		}
	}
	return false
}

func tryToAddPeers(c chan bool) {
	if len(peers) == config.MAX_PEERS {
		return
	}

	added := false

	for len(peers) < config.MAX_PEERS && len(potentialPeers) > 0 {
		nextPeer := potentialPeers[0]
		potentialPeers = potentialPeers[1:]

		if isAlreadyAPeer(nextPeer) {
			continue
		}

		if makeNewPeerRequest(nextPeer) && !added {
			added = true
			c <- true
		}

		makeGetPeersRequest(nextPeer)
	}

	c <- false
}
