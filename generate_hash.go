package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"time"
)

func main() {
	println(generateHash("12"))
}

func generateHash(uid string) string {
	secret := []byte("secret")
	time := time.Now().Unix()
	ts := fmt.Sprintf("%d", time)

	newHash := hmac.New(sha256.New, secret)

	n, err := newHash.Write([]byte(ts + uid))
	if err != nil || n != len([]byte(ts+uid)) {
		panic(err)
	}

	return fmt.Sprintf("Hash: %x\nTimeStamp: %s\nUID: %s", newHash.Sum(nil), ts, uid)
}
