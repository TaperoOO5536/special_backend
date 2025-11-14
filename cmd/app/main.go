package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/TaperoOO5536/special_backend/internal/app"
	"github.com/TaperoOO5536/special_backend/pkg/env"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	cfg := &app.Config{
		GrpcPort:     env.GetGRPCPort(),
		HttpPort:     env.GetHTTPPort(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Dsn:          env.GetDsn(),
	}

	app := app.New(cfg)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	GenInitData()
}

func GenInitData() error {
	user := User{
		ID:        1165017205,
		FirstName: "Олеся",
		LastName:  "",
		Username:  "@kakayatobeliberda",
	}
	token := env.GetToken()
	values := url.Values{}
	values.Set("query_id", "TEST_QUERY_123")
	values.Set("auth_date", fmt.Sprintf("%d", time.Now().Unix()))

	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}
	values.Set("user", string(userJSON))

	var dataCheckStrings []string
	for key := range values {
		dataCheckStrings = append(dataCheckStrings, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	sort.Strings(dataCheckStrings)
	dataCheckString := strings.Join(dataCheckStrings, "\n")

	h := sha256.New()
	h.Write([]byte("WebAppData"))
	secretKey := h.Sum(nil)
	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(token))
	hmacKey := h.Sum(nil)

	h = hmac.New(sha256.New, hmacKey)
	h.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(h.Sum(nil))

	values.Set("hash", hash)
	initData := values.Encode()
	log.Println("Test initData:", initData)
	return nil
}