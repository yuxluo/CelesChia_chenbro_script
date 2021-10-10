package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	. "github.com/yuxluo/CelesChia_chenbro_script/client"
	. "github.com/yuxluo/CelesChia_chenbro_script/server"
)

var (
	port        = flag.Uint("port", 1337, "port to listen or connect to for rpc calls")
	isServer    = flag.Bool("server", false, "activates server mode")
	json        = flag.Bool("json", false, "whether it should use json-rpc")
	serverSleep = flag.Duration("server.sleep", 0, "time for the server to sleep on requests")
	http        = flag.Bool("http", false, "whether it should use HTTP")
)

// handleSignals is a blocking function that waits for termination/interrupt
// signals.
//
// Running it in the background (non-main goroutine) has the effect of keeping
// track of the desire of termination of the current execution and then responding
// accordingly.
//
// In this example we gracefully  close the server listener in the case
// of the server - in the case of the client, breaks the request by cancelling the
// context.
func handleSignals() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	log.Println("signal received")
}

// must panics in the case of error.
func must(err error) {
	if err == nil {
		return
	}

	log.Panicln(err)
}

// runServer sets up the server with the
// flags as they were parsed and then initiates
// the server listening.
func runServer() {
	server := &Server{
		UseHttp: *http,
		UseJson: *json,
		Sleep:   *serverSleep,
		Port:    *port,
	}
	defer server.Close()

	go func() {
		handleSignals()
		server.Close()
		os.Exit(0)
	}()

	must(server.Start())
	return
}

func findEmpty() string {
	out, err := exec.Command("df").Output()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	//fmt.Printf("df output is \n%s\n", out)
	var dfOutput = string(out)

	s := strings.Split(dfOutput, "\n")
	for index, partition := range s {
		if index == 0 || partition == "" {
			continue
		}
		items := strings.Split(partition, " ")
		var withoutSpace []string
		for _, item := range items {
			if item != "" {
				withoutSpace = append(withoutSpace, item)
			}
		}
		intCapacity, _ := strconv.Atoi(withoutSpace[1])
		intRemaining, _ := strconv.Atoi(withoutSpace[3])
		if intCapacity > 999999999 && intRemaining > 100999999 {
			return withoutSpace[5]
		}
	}
	return ""
}

// runClient sets up the client with the
// flags as they were parsed and then initiates
// the client execution.
func runClient(masterIP string) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	client := &Client{
		UseHttp: *http,
		UseJson: *json,
		Port:    *port,
	}
	defer cancel()
	defer client.Close()

	must(client.Init(masterIP))

	for true {
		//	找到空盤
		nextEmptyDrive := findEmpty()
		if nextEmptyDrive == "" {
			fmt.Println("Either 吃屎了 or 全部盤都已經裝滿")
			return
		} else {
			for true {
				plotFileName, _ := client.Execute(ctx, "request")
				if plotFileName == "" {
					fmt.Println("母雞沒有plot, 等待1分鐘....")
					time.Sleep(1 * time.Minute)
				} else {
					println("Transfering %s to %s", plotFileName, nextEmptyDrive)
					exec.Command("wget", masterIP+"/"+plotFileName, "-P", nextEmptyDrive).Output()
					println("Finished transfering, requesting delete")
					plotFileName, _ := client.Execute(ctx, "delete " + plotFileName)
					break
				}
			}
		}
	}
}

// main execution - validates flags and constructs the internal
// runtime configuration based on the flags supplied.
func main() {
	flag.Parse()

	if *isServer {
		log.Println("starting server")
		log.Printf("will listen on port %d\n", *port)

		runServer()
		return
	}

	log.Println("starting client")
	log.Printf("will connect to port %d\n", *port)
	fmt.Println("Enter Muji IP: ")
	var masterIP string
	fmt.Scanln(&masterIP)
	runClient(masterIP)
	return
}
