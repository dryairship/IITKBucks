package controllers

import (
	"github.com/dryairship/IITKBucks/config"
	"github.com/dryairship/IITKBucks/logger"
	"github.com/dryairship/IITKBucks/models"
)

var currentMinerChannel chan bool

func performNoobInitialization() {
	c := make(chan bool)
	go tryToAddPeers(c)
	addedAPeer := <-c
	if !addedAPeer {
		logger.Fatal("[Controllers/Init] [FATAL] No peers added.")
	}
}

func performProInitialization() {
	genesisBlock := models.NewGenesisBlock()
	currentMinerChannel = make(chan bool)
	go mineBlock(genesisBlock)
	signal := <-currentMinerChannel
	if signal {
		logger.Println(logger.MajorEvent, "[Controllers/Init] [INFO] Genesis block mined successfully.")
	} else {
		logger.Fatal("[Controllers/Init] [FATAL] Could not mine genesis block.")
	}
}

func PerformInitialization() {
	if config.IS_PRO {
		performProInitialization()
	} else {
		performNoobInitialization()
	}
}
