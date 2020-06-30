package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/config"
)

type newPeerRequestBody struct {
	Url string `json:"url" binding:"required"`
}

type getPeersResponseBody struct {
	Peers []string `json:"peers"`
}

var peers []string
var potentialPeers []string = config.POTENTIAL_PEERS

func getPeersHandler(c *gin.Context) {
	if peers != nil {
		c.JSON(200, gin.H{"peers": peers})
	} else {
		c.JSON(200, gin.H{})
	}
}

func newPeerHandler(c *gin.Context) {
	if len(peers) == config.MAX_PEERS {
		c.AbortWithStatus(500)
		return
	}

	var body newPeerRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(400, err)
	}

	peers = append(peers, body.Url)
	c.Status(200)
}

func makeGetPeersRequest(peer string) {
	response, err := http.Get(fmt.Sprintf("%s/getPeers", peer))
	if err != nil {
		log.Println("[ERROR] go HTTP error while asking for peers. Peer: ", peer, ", Error: ", err)
		return
	}

	defer response.Body.Close()
	var bodyBytes []byte
	_, err = response.Body.Read(bodyBytes)
	if err != nil {
		log.Println("[ERROR] Cannot read getPeers response body. Peer:  ", peer, "Error: ", err)
		return
	}

	var body getPeersResponseBody
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.Println("[ERROR] Cannot unmarshal getPeers response body. Peer:  ", peer, "Error: ", err)
		return
	}

	potentialPeers = append(potentialPeers, body.Peers...)
}

func makeNewPeerRequest(peer string) {
	var jsonStr = []byte(fmt.Sprintf(`"url":"%s"}`, config.MY_URL))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/newPeer", peer), bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("[ERROR] go HTTP error while builing newPeer request. Peer: ", peer, ", Error: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] go HTTP error while making newPeer request. Peer: ", peer, ", Error: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		peers = append(peers, peer)
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

		makeNewPeerRequest(nextPeer)
		makeGetPeersRequest(nextPeer)

		if !added && len(peers) > 1 {
			added = true
			c <- true
		}
	}

	c <- false
}
