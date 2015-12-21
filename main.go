package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/erbridge/gotwit"
	"github.com/erbridge/gotwit/twitter"
	"github.com/erbridge/wikipaedian/wiki"
)

func main() {
	var (
		con twitter.ConsumerConfig
		acc twitter.AccessConfig
	)

	f := "secrets.json"
	if _, err := os.Stat(f); err == nil {
		con, acc, _ = twitter.LoadConfigFile(f)
	} else {
		con, acc, _ = twitter.LoadConfigEnv()
	}

	b := gotwit.NewBot("wikipaedian", con, acc)

	go func() {
		if err := b.Start(); err != nil {
			panic(err)
		}
	}()

	now := time.Now()

	rand.Seed(now.UnixNano())

	next := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour()+1,
		0,
		0,
		0,
		now.Location(),
	)

	sleep := next.Sub(now)

	fmt.Printf("%v until first tweet\n", sleep)

	time.Sleep(sleep)

	if c, err := wiki.NewClient(&b); err != nil {
		panic(err)
	} else {
		c.Start(1 * time.Hour)
	}
}
