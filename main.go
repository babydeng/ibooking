package main

import (
	"Group4/ibooking-back/internal"
	"Group4/ibooking-back/internal/api"
	"Group4/ibooking-back/store"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"os"
	"os/signal"

	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		internal.GetAsString("DB_USER", "pi"),
		internal.GetAsString("DB_PASSWORD", "123456"),
		internal.GetAsString("DB_HOST", "10.177.29.226"),
		internal.GetAsInt("DB_PORT", 5432),
		internal.GetAsString("DB_NAME", "restapi"),
	)
	// Open the database
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalln("Error opening database:", err)
	}

	// Connectivity check
	if err := db.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
	}

	// Create our demo user
	createUserInDb(db)

	// Start our server
	server := api.NewServer(internal.GetAsInt("SERVER_PORT", 9002))
	server.MustStart()
	defer server.Stop()

	//defaultMiddleware := []mux.MiddlewareFunc{
	//	api.JSONMiddleware,
	//	api.CORSMiddleware(internal.GetAsSlice("CORS_WHITELIST",
	//		[]string{
	//			"http://localhost:9000",
	//			"http://0.0.0.0:9000",
	//		}, ","),
	//	),
	//}

	// Handlers
	// server.AddRoute("/login", handleLogin(db), http.MethodPost, defaultMiddleware...)
	// server.AddRoute("/logout", handleLogout(), http.MethodGet)

	// Our session protected middleware
	// protectedMiddleware := append(defaultMiddleware, validCookieMiddleware(db))
	// server.AddRoute("/checkSecret", checkSecret(db), http.MethodGet, protectedMiddleware...)

	// Workouts
	server.AddRoute("/order", handleCreateNewOrder(db), http.MethodPost)
	server.AddRoute("/order/{date}", handleListOrders(db), http.MethodGet)
	server.AddRoute("/order/{order_id}", handleDeleteOrder(db), http.MethodDelete)
	// server.AddRoute("/order/{order_id}", handleUpdateOrder(db), http.MethodPut)

	// Wait for CTRL-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// We block here until a CTRL-C / SigInt is received
	// Once received, we exit and the server is cleaned up
	<-sigChan

}

func createUserInDb(db *sql.DB) {

	ctx := context.Background()
	querier := store.New(db)

	log.Println("Creating user@user...")
	hashPwd := internal.HashPassword("123456")

	_, err := querier.CreateUsers(ctx, store.CreateUsersParams{
		UserNum:      "22212010007",
		PasswordHash: hashPwd,
		UserName:     "mixiaochao",
	})

	// This is interesting to look at, the sql/pq library recommends we use
	// this pattern to understand errors. We could use the ErrorCode directly
	// or look for the specific type. We know we'll be violating unique_violation
	// if our user already exists in the database
	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		log.Println("Dummy User already present")
		return
	}

	if err != nil {
		log.Println("Failed to create user:", err)
	}
}
