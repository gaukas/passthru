package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gaukas/passthru/config"
	"github.com/gaukas/passthru/handler"
	"github.com/gaukas/passthru/protocol"
	"github.com/gaukas/passthru/protocol/tls"
)

var (
	supportedProtocols = []protocol.Protocol{
		&tls.Protocol{},
	}
	serverVersion *config.Version = &config.Version{
		Major: 0,
		Minor: 2,
		Patch: 0,
	}
	conf *config.Config
)

// STOP EDITING! OR YOU ARE HACKING THE PROJECT.

func main() {
	configFile := flag.String("c", "", "path to config file")
	workerCountPerServer := flag.Int("w", 10, "number of workers (default 10, 0 for unlimited) assigned for each server")
	workerTimeout := flag.Duration("t", 5*time.Second, "worker timeout in seconds (default 5)")
	flag.Parse()

	// Load config
	var err error
	conf, err = config.LoadConfig(*configFile)
	if err != nil {
		panic(err)
	}

	// Check version
	switch conf.Version.CanFitInServer(serverVersion) {
	case config.WONT_FIT:
		panic("[FATAL] config version is too new for the server.")
	case config.MAY_FIT:
		fmt.Println("[WARNING] config version is newer than the server. Some features may not work.")
	case config.SHOULD_FIT:
		fmt.Println("[INFO] config version is better patched than the server. There could be unintended bahaviors.")
	}

	bufServer := make(chan *handler.Server, len(conf.Servers))
	workerWg := &sync.WaitGroup{}

	for serverAddr, protoGroup := range conf.Servers {
		// Create Protocol Manager
		protoMgr := protocol.NewProtocolManager()

		// Register supported protocols
		for _, supportedProtocol := range supportedProtocols {
			protoMgr.RegisterProtocol(supportedProtocol)
		}

		// Import protocol group
		err := protoMgr.ImportProtocolGroup(protoGroup)
		if err != nil {
			panic(err)
		}

		var server *handler.Server
		if *workerCountPerServer <= 0 {
			// Create unlimited server
			server = handler.NewServer(serverAddr, protoMgr, handler.SERVER_MODE_UNLIMITED)
			server.Start()
		} else {
			// Create worker-based server
			server = handler.NewServer(serverAddr, protoMgr, handler.SERVER_MODE_WORKER)
			server.Start()
			// spawn workers
			for i := 0; i < *workerCountPerServer; i++ {
				workerWg.Add(1)
				go func(server *handler.Server, wg *sync.WaitGroup) {
					defer wg.Done()

					for {
						ctxTimeOut, cancel := context.WithTimeout(context.Background(), *workerTimeout*time.Second)
						defer cancel()
						err := server.HandleNextConn(ctxTimeOut)
						switch err {
						case handler.ErrUnknownAction:
							fmt.Println("[WARNING] unknown action from a protocol.Protocol")
							continue
						case handler.ErrServerStopped:
							return
						default:
							if err != nil && err != context.DeadlineExceeded {
								fmt.Printf("[ERROR] error while handling connection: %v\n", err)
							}
						}
					}
				}(server, workerWg)
			}
		}

		fmt.Printf("[INFO] server %s started\n", serverAddr)
		bufServer <- server
	}
	close(bufServer)

	// Capture Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		for {
			server := <-bufServer
			if server == nil {
				break
			}
			server.Stop()
		}

		// wait for all workers to finish
		workerWg.Wait()
		os.Exit(0)
	}()

	select {}
}
