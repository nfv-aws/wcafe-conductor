package main

import (
	"github.com/nfv-aws/wcafe-api-controller/db"
	"github.com/nfv-aws/wcafe-conductor/conductor"
)

func main() {
	db.Init()
	conductor.GetMessage()
	db.Close()
}
