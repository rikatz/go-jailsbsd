package main

import (
	"github.com/rikatz/go-jailsbsd/pkg/jailcgo"
	"log"
	"fmt"
)

func main() {
	jail := jailcgo.JailParams{}
	jail.MapParams = make(map[string]string)
	jail.MapParams["name"] = "bloblo"
	jail.MapParams["host.hostname"] = "bloblo.com"
	jail.MapParams["path"] = "/jails/katz"
	jail.MapParams["persist"] = "true"
	jail.MapParams["securelevel"] = "3"
	jail.MapParams["ip4.addr"] = "192.168.0.222"


	jailid, err := jail.Set()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("JailID: %d\n", jailid)

}