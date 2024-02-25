package traefik_waiting_room

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

type Config struct {
	Enabled          bool          `json:"enabled"`
	WaitRoomFile     string        `json:"waitRoomFile"`
	WaitingTime      time.Duration `json:"waitingTime"`
	PurgeTime        time.Duration `json:"purgeTime"`
	MaxEntries       int           `json:"maxEntries"`
	HttpResponseCode int           `json:"httpResponseCode"`
	HttpContentType  string        `json:"httpContentType"`
}

func CreateConfig() *Config {
	return &Config{
		Enabled:          true,
		WaitRoomFile:     "waiting-room.html",
		WaitingTime:      1,
		PurgeTime:        5,
		MaxEntries:       5,
		HttpResponseCode: http.StatusTooManyRequests,
		//HttpResponseCode: http.StatusOK,
		HttpContentType: "text/html; charset=utf-8",
	}
}

type WaitingRoom struct {
	next             http.Handler
	enabled          bool
	waitRoomFile     string
	waitingTime      time.Duration
	purgeTime        time.Duration
	maxEntries       int
	httpResponseCode int
	httpContentType  string
	name             string
	config           *Config
	cache            cache.Cache
	template         *template.Template
}

const xTraefikID = "X-Trf-Id"

// Instantiate a new Waiting Room plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	// Log config
	log.Printf("waiting room plugin config %v", config)

	if len(config.WaitRoomFile) == 0 {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	return &WaitingRoom{
		enabled:          config.Enabled,
		waitRoomFile:     config.WaitRoomFile,
		waitingTime:      config.WaitingTime,
		purgeTime:        config.PurgeTime,
		maxEntries:       config.MaxEntries,
		httpResponseCode: config.HttpResponseCode,
		httpContentType:  config.HttpContentType,
		next:             next,
		name:             name,
		config:           config,
		cache:            *cache.New(config.WaitingTime*time.Minute, config.PurgeTime*time.Minute),
		template:         template.New("WaitRoom").Delims("[[", "]]"),
	}, nil
}

func (r *WaitingRoom) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	log.Println("Cache Items")
	for k, v := range r.cache.Items() {
		fmt.Println("k:", k, "v:", v)
	}
	log.Print("Item Count: " + strconv.Itoa(r.cache.ItemCount()))

	cookie, err := req.Cookie(xTraefikID)
	if err != nil {
		log.Println("Cookie not found!")
	} else {
		r.cache.Set(cookie.Value, time.Now().Unix(), cache.DefaultExpiration)
	}

	cookies := req.Cookies()
	for _, cookie := range cookies {
		log.Println("Cookie value: " + cookie.Value)
	}

	cSize := r.cache.ItemCount()

	if cSize >= r.config.MaxEntries {
		fileBytes, err := os.ReadFile(r.waitRoomFile)
		if err == nil {
			rw.Header().Add("Content-Type", r.httpContentType)
			rw.WriteHeader(r.httpResponseCode)
			_, err = rw.Write(fileBytes)
			if err != nil {
				log.Printf("Could not serve waiting room template %s: %s", r.waitRoomFile, err)
			} else {
				return
			}
		} else {
			log.Printf("Could not read waiting room template %s: %s", r.waitRoomFile, err)
		}

		r.next.ServeHTTP(rw, req)
	} else {
		r.next.ServeHTTP(rw, req)
	}
}
