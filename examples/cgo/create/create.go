package main

import (
	"fmt"
	"github.com/rikatz/go-jailsbsd/pkg/jailcgo"
	"log"
)

func main() {
	jail := jailcgo.Jail{}
	jail.Params = make(map[string]string)
	jail.Params["name"] = "bloblo"
	jail.Params["host.hostname"] = "bloblo.com"
	jail.Params["path"] = "/jails/katz"
	jail.Params["persist"] = "true"
	jail.Params["securelevel"] = "3"
	jail.Params["ip4.addr"] = "192.168.0.222"

	jailid, err := jail.Set()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("JailID: %d\n", jailid)

}
