package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DurationString string
	RikonURL       string
	GmailPass      string
	EmailFrom      string
	PGHost         string
	PGPort         int32
	PGUser         string
	PGPassword     string
	PGDBName       string
	WebHost        string
}

func LoadConfig() Config {
	var gmailPass, pgPassword string
	awsSecrets := getEnv("AWS_SECRETS", "")
	if awsSecrets != "" {
		var secretMap map[string]string
		if err := json.Unmarshal([]byte(awsSecrets), &secretMap); err != nil {
			log.Fatal(err)
		}
		_, exists := secretMap["PG_PASS"]
		if !exists {
			log.Fatalf("aws secret does not contain pg password")
		}
		_, exists = secretMap["GMAIL_PASS"]
		if !exists {
			log.Fatalf("aws secret does not contain gmail password")
		}
		pgPassword = secretMap["PG_PASS"]
		gmailPass = secretMap["GMAIL_PASS"]
	} else {
		pgPassword = getEnv("PG_PASSWORD", "")
		gmailPass = getEnv("GMAIL_PASS", "")
	}
	durationString := getEnv("DURATION_STRING", "15s")
	rikonURL := getEnv("RIKON_URL", "https://rikontools.com/product/")
	emailFrom := getEnv("EMAIL_FROM", "")
	pgHost := getEnv("PG_HOST", "localhost")
	pgPort := getEnvInt("PG_PORT", 5432)
	pgUser := getEnv("PG_USER", "postgres")
	pgDBName := getEnv("PG_DBNAME", "postgres")
	webHost := getEnv("WEB_HOST", "http://localhost:5173/")

	return Config{
		DurationString: durationString,
		RikonURL:       rikonURL,
		GmailPass:      gmailPass,
		EmailFrom:      emailFrom,
		PGHost:         pgHost,
		PGPort:         pgPort,
		PGUser:         pgUser,
		PGPassword:     pgPassword,
		PGDBName:       pgDBName,
		WebHost:        webHost,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int32) int32 {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return int32(intValue)
		}
	}
	return fallback
}
