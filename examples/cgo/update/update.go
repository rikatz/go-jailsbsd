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

	// Here you put your JailID
	jailid, err := jail.Update(5)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("JailID: %d\n", jailid)

}