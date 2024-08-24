package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
	"github.com/spf13/viper"
)

func main() {
	// Load configuration
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("ini")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Open serial port
	port, err := os.OpenFile(config.GetString("FLIPDOTSIGN.USB"), os.O_RDWR, 0)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer port.Close()

	// Create controller
	controller, err := goflipdot.NewController(port)
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}

	// Add sign
	width := config.GetInt("FLIPDOTSIGN.COLUMNS")
	height := config.GetInt("FLIPDOTSIGN.ROWS")
	err = controller.AddSign("main",
		config.GetInt("FLIPDOTSIGN.ADDRESS"),
		width,
		height,
		false)
	if err != nil {
		log.Fatalf("Failed to add sign: %v", err)
	}

	// Create image
	img, err := controller.CreateImage("main")
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	// Draw checkerboard pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	// Display image
	fmt.Println("Displaying checkerboard pattern...")
	err = controller.DrawImage(img, "main")
	if err != nil {
		log.Fatalf("Failed to draw image: %v", err)
	}

	fmt.Println("Checkerboard pattern displayed successfully!")
}
