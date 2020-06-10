package main

import (
	"github.com/nfv-aws/wcafe-api-controller/db"
	"github.com/nfv-aws/wcafe-conductor/conductor"
)

func main() {
	db.Init()

	// 並列処理開始
	// それぞれとの連絡のためのchを作成する
	petsFin := make(chan bool)
	storesFin := make(chan bool)
	usersFin := make(chan bool)
	suppliesFin := make(chan bool)
	go func() {
		conductor.PetsGetMessage()
		petsFin <- true
	}()
	go func() {
		conductor.StoresGetMessage()
		storesFin <- true
	}()
	go func() {
		conductor.UsersGetMessage()
		usersFin <- true
	}()
	go func() {
		conductor.SuppliesGetMessage()
		suppliesFin <- true
	}()
	// 全部が終わるまでブロックし続ける
	<-petsFin
	<-storesFin
	<-usersFin
	<-suppliesFin

	db.Close()
}
