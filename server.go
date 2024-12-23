package main

import (
	"fmt"
	"server/configs"
	"server/endpoints"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	configs.ConnectDB()

	// Initialize Gin router
	router := gin.Default()
	router.MaxMultipartMemory = 256 << 20 // Set max upload size to 256 MiB

	// Register endpoint groups
	endpoints.TapeEndpoints(router)
	endpoints.AudiobookEndpoints(router)

	// Serve static files
	router.Static("/static/audio", "./Audiobooks")             // Serve MP3 files from /static/audio
	router.Static("/static/images", "./Audiobooks/thumbnails") // Serve images from /static/images

	// List of authors and corresponding image filenames
	authors := []struct {
		Name     string
		ImageURL string
	}{
		{"Marian Keyes", "/static/images/marian_keyes.png"},
		{"Máiréad Ní Ghráda", "/static/images/mairead_ni_ghrada.png"},
		{"Roddy Doyle", "/static/images/roddy_doyle.png"},
		{"Áine Ní Ghlinn", "/static/images/aine_ni_ghlinn.png"},
		{"John Connolly", "/static/images/john_connolly.png"},
		{"Patricia Scanlan", "/static/images/patricia_scanlan.png"},
		{"Julie Parsons", "/static/images/julie_parsons.png"},
		{"Deirdre Purcell", "/static/images/deirdre_purcell.png"},
		{"Jo O'Donoghue", "/static/images/jo_odonoghue.png"},
		{"Vincent Banville", "/static/images/vincent_banville.png"},
	}

	// Root route: plain text response
	router.GET("/", func(c *gin.Context) {
		// Print the welcome message ONCE
		c.Writer.Write([]byte("Message: Welcome to the Audiobook Server!\n"))

		// Loop through authors and write details to the response
		for i, author := range authors {
			audioPath := fmt.Sprintf("/static/audio/1/%d.mp3", i+1)

			// Write details for each file directly to the response
			c.Writer.Write([]byte(fmt.Sprintf(
				"Author: %s\nAudio: %s\nImage: %s\n-------------------\n",
				author.Name, audioPath, author.ImageURL,
			)))
		}
	})

	// Start the server on port 8080
	router.Run(":8080")
}
