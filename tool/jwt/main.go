package main

import (
	"fmt"
	"log"
	"os"
	"time"

	auth "gitlab.com/thefrol/notty/pkg/jwt"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage jwt <key>")
	}

	key := os.Args[1]
	jwtS, err := auth.BuildExpirable(time.Hour*300, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(jwtS)
}
