package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	AppURL          string
	Debug           bool
	UseCache        bool
	Port            string
	ServerName      string
	Secure          bool
	DatabaseDriver  string
	DatabaseHost    string
	DatabasePort    string
	DatabaseUser    string
	DatabasePass    string
	DatabaseName    string
	DatabaseSSLMode string
	MongoDBURI      string
	MongoDBName     string
	RedisHost       string
	RedisPassword   string
	RedisPrefix     string
	Cache           string
	CookieName      string
	CookieLifetime  int
	CookiePersist   bool
	CookieSecure    bool
	CookieDomain    string
	SessionStore    string
	SMTPHost        string
	SMTPUsername    string
	SMTPPassword    string
	SMTPPort        int
	SMTPEncryption  string
	SMTPFrom        string
	MailerAPI       string
	MailerKey       string
	MailerURL       string
	Renderer        string
	EncryptionKey   string
	AdditionalKeys  map[string]interface{}
}

var Env Config

func Load(root string, filename string, additionalKeys ...string) {
	// Load .env file
	err := godotenv.Load(root + "/env/" + filename + ".env")
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	Env = Config{
		AppName:         getEnv("APP_NAME", ""),
		AppURL:          getEnv("APP_URL", "http://localhost:4000"),
		Debug:           getEnvAsBool("DEBUG", true),
		UseCache:        getEnvAsBool("USE_CACHE", true),
		Port:            getEnv("PORT", "4000"),
		ServerName:      getEnv("SERVER_NAME", "localhost"),
		Secure:          getEnvAsBool("SECURE", false),
		DatabaseDriver:  getEnv("DATABASE_DRIVER", ""),
		DatabaseHost:    getEnv("DATABASE_HOST", ""),
		DatabasePort:    getEnv("DATABASE_PORT", ""),
		DatabaseUser:    getEnv("DATABASE_USER", ""),
		DatabasePass:    getEnv("DATABASE_PASS", ""),
		DatabaseName:    getEnv("DATABASE_NAME", ""),
		DatabaseSSLMode: getEnv("DATABASE_SSL_MODE", ""),
		MongoDBURI:      getEnv("MONGODB_URI", ""),
		MongoDBName:     getEnv("MONGODB_DB", ""),
		RedisHost:       getEnv("REDIS_HOST", ""),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisPrefix:     getEnv("REDIS_PREFIX", ""),
		Cache:           getEnv("CACHE", ""),
		CookieName:      getEnv("COOKIE_NAME", ""),
		CookieLifetime:  getEnvAsInt("COOKIE_LIFETIME", 1440),
		CookiePersist:   getEnvAsBool("COOKIE_PERSIST", true),
		CookieSecure:    getEnvAsBool("COOKIE_SECURE", false),
		CookieDomain:    getEnv("COOKIE_DOMAIN", "localhost"),
		SessionStore:    getEnv("SESSION_STORE", "cookie"),
		SMTPHost:        getEnv("SMTP_HOST", ""),
		SMTPUsername:    getEnv("SMTP_USERNAME", ""),
		SMTPPassword:    getEnv("SMTP_PASSWORD", ""),
		SMTPPort:        getEnvAsInt("SMTP_PORT", 1025),
		SMTPEncryption:  getEnv("SMTP_ENCRYPTION", ""),
		SMTPFrom:        getEnv("SMTP_FROM", ""),
		MailerAPI:       getEnv("MAILER_API", ""),
		MailerKey:       getEnv("MAILER_KEY", ""),
		MailerURL:       getEnv("MAILER_URL", ""),
		Renderer:        getEnv("RENDERER", "jet"),
		EncryptionKey:   getEnv("KEY", ""),
		AdditionalKeys:  make(map[string]interface{}),
	}

	// Set values that depend on other config values
	Env.RedisPrefix = getEnv("REDIS_PREFIX", Env.AppName)
	Env.CookieName = getEnv("COOKIE_NAME", Env.AppName)

	// Load additional keys
	for _, key := range additionalKeys {
		Env.AdditionalKeys[key] = getEnvAsInterface(key)
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultValue
}

func getEnvAsInterface(key string) interface{} {
	value := os.Getenv(key)

	// Try to parse as bool
	if boolValue, err := strconv.ParseBool(value); err == nil {
		return boolValue
	}

	// Try to parse as int
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	// Try to parse as float
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}

	// Check if it's a slice (comma-separated values)
	if strings.Contains(value, ",") {
		return strings.Split(value, ",")
	}

	// Return as string if no other type matches
	return value
}

// GetAdditionalKey retrieves the value of an additional key
func GetAdditionalKey(key string) interface{} {
	return Env.AdditionalKeys[key]
}

// Type-specific getter functions
func GetAdditionalKeyAsString(key string) (string, bool) {
	value, ok := Env.AdditionalKeys[key]
	if !ok {
		return "", false
	}
	strValue, ok := value.(string)
	return strValue, ok
}

func GetAdditionalKeyAsBool(key string) (bool, bool) {
	value, ok := Env.AdditionalKeys[key]
	if !ok {
		return false, false
	}
	boolValue, ok := value.(bool)
	return boolValue, ok
}

func GetAdditionalKeyAsInt(key string) (int, bool) {
	value, ok := Env.AdditionalKeys[key]
	if !ok {
		return 0, false
	}
	intValue, ok := value.(int)
	return intValue, ok
}

func GetAdditionalKeyAsFloat(key string) (float64, bool) {
	value, ok := Env.AdditionalKeys[key]
	if !ok {
		return 0, false
	}
	floatValue, ok := value.(float64)
	return floatValue, ok
}

func GetAdditionalKeyAsSlice(key string) ([]string, bool) {
	value, ok := Env.AdditionalKeys[key]
	if !ok {
		return nil, false
	}
	sliceValue, ok := value.([]string)
	return sliceValue, ok
}
