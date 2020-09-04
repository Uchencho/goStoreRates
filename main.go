package main

import (
	"github.com/Uchencho/goStoreRates/config"
)

func main() {

	defer config.Db.Close()
}
