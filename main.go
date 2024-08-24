package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
	"github.com/spf13/viper"
)

type FlipdotServer struct {
	controller *goflipdot.Controller
	config     *viper.Viper
	width      int
	height     int
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

	width := config.GetInt("FLIPDOTSIGN.COLUMNS")
	height := config.GetInt("FLIPDOTSIGN.ROWS")

	err = controller.AddSign("main",
		config.GetInt("FLIPDOTSIGN.ADDRESS"),
		width,
		height,
		false)
	if err != nil {
		return nil, fmt.Errorf("failed to add sign: %w", err)
	}

	return &FlipdotServer{
		controller: controller,
		config:     config,
		width:      width,
		height:     height,
	}, nil
}

func (s *FlipdotServer) handleDotArray(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("Received request to /api/dots")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Error: Method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("Error: Failed to read request body: %v", err)
		return
	}

	var imageArray [][]int
	err = json.Unmarshal(body, &imageArray)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		log.Printf("Error: Failed to parse JSON: %v", err)
		return
	}

	img, err := s.controller.CreateImage("main")
	if err != nil {
		http.Error(w, "Failed to create image", http.StatusInternalServerError)
		log.Printf("Error: Failed to create image: %v", err)
		return
	}

	for y, row := range imageArray {
		for x, pixel := range row {
			if pixel == 1 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	err = s.controller.DrawImage(img, "main")
	if err != nil {
		http.Error(w, "Failed to draw image", http.StatusInternalServerError)
		log.Printf("Error: Failed to draw image: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": "image displayed"})

	duration := time.Since(start)
	log.Printf("Request processed in %v", duration)
}

func (s *FlipdotServer) handleDocumentation(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "README.md")
}

func (s *FlipdotServer) displayTestPattern() error {
	log.Println("Displaying test pattern")
	img, err := s.controller.CreateImage("main")
	if err != nil {
		return fmt.Errorf("failed to create image: %w", err)
	}

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			if (x + y) % 2 == 0 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	err = s.controller.DrawImage(img, "main")
	if err != nil {
		return fmt.Errorf("failed to draw test pattern: %w", err)
	}

	log.Println("Test pattern displayed successfully")
	return nil
}

func main() {
	log.Println("Starting Flipdot Server")
	server, err := NewFlipdotServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	err = server.displayTestPattern()
	if err != nil {
		log.Printf("Warning: Failed to display test pattern: %v", err)
	}

	http.HandleFunc("/api/dots", server.handleDotArray)
	http.HandleFunc("/documentation", server.handleDocumentation)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/documentation", http.StatusFound)
	})

	addr := fmt.Sprintf("%s:%d",
		server.config.GetString("SERVER.HOST"),
		server.config.GetInt("SERVER.PORT"))

	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
