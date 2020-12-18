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

	// Here you put your JailID
	jailid, err := jail.Update(5)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("JailID: %d\n", jailid)

}
