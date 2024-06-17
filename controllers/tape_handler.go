package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"server/database"
	"server/filemanager"
	"server/models"

	"github.com/gin-gonic/gin"
)

func ReturnCatalog(section string) []byte {
	result := database.SelectAll("Tape")
	var tapes []models.Tape
	err := json.Unmarshal(result, &tapes)
	if err != nil {
		print(err)
	}
	if len(tapes) == 0 {
		panic("Empty array returned from DB")
	}
	jsonString, _ := json.Marshal(tapes)
	return jsonString
}

func ReturnAudioFileData(id string) []byte {
	result := database.SelectById("Tape", "1")
	var value []models.Tape
	err := json.Unmarshal(result, &value)
	if err != nil {
		print(err)
	}
	if len(value) == 0 {
		panic("Empty array returned from DB")
	}
	jsonString, _ := json.Marshal(value[0])
	return jsonString
}

func DownloadFolder(c *gin.Context, audioId string) {
	folderPath := "./Audiobooks/" + audioId + "/"
	zipFileFolder := "./AudiobooksZipped/" + audioId + "/"
	zipFileName := "./AudiobooksZipped/" + audioId + "/files.zip"

	err := filemanager.ZipFolder(folderPath, zipFileFolder, zipFileName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating zip file: %v", err)
		return
	}
	// we can talk about this if you would like but I think we should reuse the zip files
	// and leave the original files there as a backup
	//defer os.Remove(zipFileName)

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFileName)
	c.Writer.Header().Set("Content-Type", "application/zip")

	file, err := os.Open(zipFileName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening zip file: %v", err)
		return
	}
	defer file.Close()

	io.Copy(c.Writer, file)
}

func UploadFile(c *gin.Context) {
	id := c.Param("id")

	form, err1 := c.MultipartForm()
	if err1 != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err1.Error())
		return
	}
	files := form.File["files"]

	uploadDir := "./AudioBooks/" + id
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not create upload directory")
		return
	}

	for _, file := range files {
		filePath := filepath.Join(uploadDir, filepath.Base(file.Filename))
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Could not save file %s", file.Filename))
			return
		}
	}

	c.String(http.StatusOK, fmt.Sprintf("Uploaded %d files successfully", len(files)))
}

func DownloadFile(c *gin.Context) {
	fileName := c.Param("fileName")
	folderType := c.Param("folderType")
	filePath := "./" + folderType + "/" + fileName
	c.FileAttachment(filePath, fileName)
}

func GetCatalogByType(c *gin.Context) {
	sectionType := c.Param("type")
	c.Header("Content-Type", "application/json")

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(ReturnCatalog(sectionType))
}