package main

import (
	"github.com/nfv-aws/wcafe-conductor/conductor"
	"github.com/nfv-aws/wcafe-conductor/db"
)

func main() {
	db.Init()
	conductor.GetMessage()
	db.Close()
}
