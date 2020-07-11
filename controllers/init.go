package controllers

import (
	"log"

	"github.com/dryairship/IITKBucks/config"
	"github.com/dryairship/IITKBucks/models"
)

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
	log.Printf("%+v\n", genesisBlock)
}

func PerformInitialization() {
	if config.IS_PRO {
		performProInitialization()
	} else {
		performNoobInitialization()
	}
}
