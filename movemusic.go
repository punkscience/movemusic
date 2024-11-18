package movemusic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MoveMusic(sourceFileFullPath string, destFolderPath string, useFolders bool) (string, error) {

	// Check if the source file exists
	if _, err := os.Stat(sourceFileFullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("source file does not exist")
	}

	// Check if the destination folder exists
	if _, err := os.Stat(destFolderPath); os.IsNotExist(err) {
		return "", fmt.Errorf("destination folder does not exist")
	}

	// Get the mp3, flac or wav file details from the file
	file, err := os.Open(sourceFileFullPath)
	if err != nil {
		return "", fmt.Errorf("error opening the file")
	}
	defer file.Close()

	// Get the file extension
	ext := strings.ToLower(filepath.Ext(sourceFileFullPath))
	if ext != ".mp3" && ext != ".flac" && ext != ".wav" {
		return "", fmt.Errorf("unsupported file type")
	}

	// Open the spirce file
	m, err := tag.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("error reading the file")
	}

	// Get the artist, album, track and track number
	artist := m.Artist()
	album := m.Album()
	track := m.Title()
	trackNumber, _ := m.Track()

	// Check if the artist, album, track and track number are empty
	if artist == "" {
		artist = "Unknown"
	}

	if album == "" {
		album = "Unknown"
	}

	if track == "" {
		track = "Unknown"
	}

	if trackNumber == 0 {
		trackNumber = 1
	}

	// Build a name
	var newName string
	if useFolders {
		newName = makeFileNameFolders(artist, album, track, trackNumber, ext)
	} else {
		newName = makeFileName(artist, album, track, trackNumber, ext)
	}

	// Build the destination file path
	destFileFullPath := filepath.Join(destFolderPath, newName)

	// Check if the destination file exists
	if _, err := os.Stat(destFileFullPath); err == nil {
		return "", fmt.Errorf("destination file already exists")
	}

	// Move the file
	err = os.Rename(sourceFileFullPath, destFileFullPath)
	if err != nil {
		return "", fmt.Errorf("error moving the file")
	}

	return destFileFullPath, nil
}

func makeFileNameFolders(artist string, album string, track string, trackNumber int, ext string) string {
	// Remove the invalid characters from the artist, album and track
	artist = cleanup(artist)
	album = cleanup(album)
	track = cleanup(track)

	newName := filepath.Join(artist, album)

	// Build the file name
	newName += fmt.Sprintf("%02d - %s%s", trackNumber, track, ext)

	return newName
}

func makeFileName(artist string, album string, track string, trackNumber int, ext string) string {
	// Remove the invalid characters from the artist, album and track
	artist = cleanup(artist)
	album = cleanup(album)
	track = cleanup(track)

	// Build the file name
	newName := fmt.Sprintf("%s - %s - %02d - %s%s", artist, album, trackNumber, track, ext)

	return newName
}

func cleanup(s string) string {

	s = strings.Trim(s, " \t\n\r\"'")

	// Remove the invalid characters from the artist, album and track
	s = strings.Replace(s, "/", "-", -1)
	s = strings.Replace(s, "\\", "-", -1)
	s = strings.Replace(s, ":", "-", -1)
	s = strings.Replace(s, "*", "-", -1)
	s = strings.Replace(s, "?", "-", -1)
	s = strings.Replace(s, "\"", "-", -1)
	s = strings.Replace(s, "<", "-", -1)
	s = strings.Replace(s, ">", "-", -1)
	s = strings.Replace(s, "|", "-", -1)
	s = strings.Replace(s, "  ", " ", -1)

	s = strings.Replace(s, "feat.", "ft.", -1)
	s = strings.Replace(s, "Feat.", "ft.", -1)
	s = strings.Replace(s, "Featuring", "ft.", -1)
	s = strings.Replace(s, "&", " and ", -1)

	// Remove any special characters
	s = strings.Map(func(r rune) rune {
		if r >= 32 && r <= 126 {
			return r
		}
		return -1
	}, s)

	// Finally, fix the capitalization to proper english title case
	s = cases.Title(language.English).String(s)

	return s
}