package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"

	"gchat/internal"
	"gchat/internal/env"
	"gchat/internal/ip"
	"gchat/internal/ws"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// in production use docker or k8s environment vars
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if env.GetEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	banlist := ip.NewBanList()
	hub := ws.NewHub(banlist)
	router := gin.Default()

	internal.ServeStatic(router)

	router.GET("/ws", hub.ServeWS)

	srv := &http.Server{
		Addr:    env.GetEnv("IP_ADDRESS", ":3000"),
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			args := strings.Fields(scanner.Text())

			if len(args) == 0 {
				return
			}

			cmd := args[0]
			switch cmd {
			case "block":
				if len(args) < 2 {
					log.Println("No IP address")
					continue
				}

				ip := args[1]
				hub.BanList.Block(ip)
				log.Printf("IP %s is blocked/n", ip)
			case "unblock":
				if len(args) < 2 {
					log.Println("No IP address")
					continue
				}

				ip := args[1]
				hub.BanList.Unblock(ip)
				log.Printf("IP %s is unblocked/n", ip)
			default:
				log.Println("Unknown command")
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Scanner error: %v/n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	hub.CloseAll()
	log.Println("Closing active Websocket connections...")
	hub.WG.Wait()

	log.Println("Server exiting")
}
