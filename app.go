package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var version = "master"

type App struct {
	Config    *viper.Viper
	Wordlists map[string][]string
	Router    *mux.Router
	DB        DB
}

func (a *App) Init() {
	log.Printf("Starting asdfland version %s", version)
	rand.Seed(time.Now().UTC().UnixNano())
	a.InitConfig()
	a.InitWordlists()
	a.InitRouter()
	switch a.Config.GetString("db_kind") {
	case "redis":
		a.InitRedis()
	default:
		log.Fatalf("Database type not supported: %s",
			a.Config.GetString("db_kind"))
	}
}

func (a *App) InitConfig() {
	a.Config = viper.New()
	a.Config.AddConfigPath(".")
	a.Config.AddConfigPath("$HOME/.config/asdfland")
	a.Config.AddConfigPath("$/etc/asdfland")
	a.Config.SetConfigName("config")

	err := a.Config.ReadInConfig()
	if err != nil {
		log.Println("Config file couldn't be loaded; using defaults")
	} else {
		log.Println("Config file loaded")
	}

	a.Config.SetDefault("port", "9090")
	a.Config.SetDefault("db_kind", "redis")
	a.Config.SetDefault("redis_addr", "localhost:6379")
	a.Config.SetDefault("redis_db", 0)
	a.Config.SetDefault("redis_pass", "")

	a.Config.SetDefault("bcrypt_cost", 8)
	a.Config.SetDefault("word_delimiter", ".")
}

func (a *App) InitWordlists() {
	fname_regex := regexp.MustCompile(`^wordlists\/(.*)\.txt$`)
	a.Wordlists = make(map[string][]string)
	for _, wl_name := range AssetNames() {
		fname_match := fname_regex.FindStringSubmatch(wl_name)
		if fname_match == nil {
			continue
		}
		wl_bytes := MustAsset(wl_name)
		scanner := bufio.NewScanner(bytes.NewReader(wl_bytes))
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		a.Wordlists[fname_match[1]] = lines
		log.Printf("Wordlist loaded: %s", wl_name)
	}
}

func (a *App) GetReadableString(wl_name string, count int) (string, error) {
	var words []string
	wl, exists := a.Wordlists[wl_name]
	if !exists {
		return "", fmt.Errorf("GetReadableString: word list %s not found", wl)
	}
	wl_size := len(wl)
	for i := 0; i < count; i++ {
		words = append(words, wl[rand.Intn(wl_size)])
	}
	return strings.Join(words, a.Config.GetString("word_delimiter")), nil
}

func (a *App) InitRouter() {
	a.Router = mux.NewRouter().StrictSlash(true)

	// Setup the static Vue.js frontend routes
	frontendServer := Logger(http.FileServer(
		&assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    "frontend"}), "frontend")
	a.Router.Path("/").Handler(frontendServer)
	a.Router.PathPrefix("/static").Handler(frontendServer)

	// Setup the API routes
	routes := a.GetRoutes()
	for _, route := range *routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = a.SessionMiddleware(handler, route.EnsureSession)
		handler = Logger(handler, route.Name)

		a.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	log.Println("Loaded routes")
}

func (a *App) InitRedis() {
	rdb := RedisDB{}
	rdb.Init(
		a.Config.GetString("redis_addr"),
		a.Config.GetString("redis_pass"),
		a.Config.GetInt("redis_dbnum"))
	a.DB = &rdb
	log.Println("Loaded Redis")
}

func (a *App) Run() {
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         "localhost:" + a.Config.GetString("port"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Starting server at %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
