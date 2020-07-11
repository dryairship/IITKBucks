package controllers

import (
	"log"

	"github.com/dryairship/IITKBucks/config"
	"github.com/dryairship/IITKBucks/models"
)

var currentMinerChannel chan bool

func performNoobInitialization() {
	c := make(chan bool)
	go tryToAddPeers(c)
	addedAPeer := <-c
	if !addedAPeer {
		log.Fatal("[ERROR] No peers added.")
	}
}

func performProInitialization() {
	genesisBlock := models.NewGenesisBlock()
	currentMinerChannel = make(chan bool)
	go mineBlock(genesisBlock)
	signal := <-currentMinerChannel
	if signal {
		log.Println("[INFO] Genesis block mined successfully.")
	} else {
		log.Fatal("[ERROR] Could not mine genesis block.")
	}
}

func PerformInitialization() {
	if config.IS_PRO {
		performProInitialization()
	} else {
		performNoobInitialization()
	}
}
