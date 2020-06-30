package controllers

import (
	"log"
)

func performNoobInitialization() {
	c := make(chan bool)
	go tryToAddPeers(c)
	addedAPeer := <-c
	if !addedAPeer {
		log.Fatal("[ERROR] No peers added.")
	}
}

func PerformInitialization() {
	performNoobInitialization()
}
