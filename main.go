package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"image/color"
	"os"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
	"github.com/spf13/viper"
)

type FlipdotServer struct {
	controller *goflipdot.Controller
	config     *viper.Viper
}

func NewFlipdotServer() (*FlipdotServer, error) {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("ini")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	port, err := os.OpenFile(config.GetString("FLIPDOTSIGN.USB"), os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}

	controller, err := goflipdot.NewController(port)
	if err != nil {
		return nil, fmt.Errorf("failed to create controller: %w", err)
	}

	err = controller.AddSign("main",
		config.GetInt("FLIPDOTSIGN.ADDRESS"),
		config.GetInt("FLIPDOTSIGN.COLUMNS"),
		config.GetInt("FLIPDOTSIGN.ROWS"),
		false)
	if err != nil {
		return nil, fmt.Errorf("failed to add sign: %w", err)
	}

	return &FlipdotServer{
		controller: controller,
		config:     config,
	}, nil
}

func (s *FlipdotServer) handleDotArray(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var imageArray [][]bool
	err = json.Unmarshal(body, &imageArray)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	img, err := s.controller.CreateImage("main")
	if err != nil {
		http.Error(w, "Failed to create image", http.StatusInternalServerError)
		return
	}

	for y, row := range imageArray {
		for x, pixel := range row {
			if pixel {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	err = s.controller.DrawImage(img, "main")
	if err != nil {
		http.Error(w, "Failed to draw image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": "image displayed"})
}

func (s *FlipdotServer) handleDocumentation(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "README.md")
}

func main() {
	server, err := NewFlipdotServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	http.HandleFunc("/api/dots", server.handleDotArray)
	http.HandleFunc("/documentation", server.handleDocumentation)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/documentation", http.StatusFound)
	})

	addr := fmt.Sprintf("%s:%d",
		server.config.GetString("SERVER.HOST"),
		server.config.GetInt("SERVER.PORT"))

	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
