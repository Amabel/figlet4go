package figlet

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Font struct {
	Hardblank string
	Height    int
	FontSlice []string
}

type FontManager struct {
	// font library
	fontLib map[string]*Font

	// font name to path
	fontList map[string]string
}

func NewFontManager() *FontManager {
	this := &FontManager{}

	this.fontLib = make(map[string]*Font)
	this.fontList = make(map[string]string)
	this.loadBuildInFont()

	return this
}

// walk through the path, load all the *.flf font file
func (this *FontManager) LoadFont(fontPath string) error {

	return filepath.Walk(fontPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".flf") {
			return nil
		}

		fontName := strings.TrimSuffix(info.Name(), ".flf")
		this.fontList[fontName] = path
		return nil
	})
}

func (this *FontManager) loadBuildInFont() error {
	font, err := this.parseFontContent(BuildInFont())
	if err != nil {
		return err
	}
	this.fontLib["default"] = font
	return nil
}

func (this *FontManager) loadDiskFont(fontName, fontFilePath string) error {
	// read full file content
	fileBuf, err := ioutil.ReadFile(fontFilePath)
	if err != nil {
		return err
	}

	font, err := this.parseFontContent(string(fileBuf))
	if err != nil {
		return err
	}

	this.fontLib[fontName] = font
	return nil
}

func (this *FontManager) parseFontContent(cont string) (*Font, error) {
	lines := strings.Split(cont, "\n")
	if len(lines) < 1 {
		return nil, errors.New("font content error")
	}

	// flf2a$ 7 5 16 -1 12
	// Fender by Scooter 8/94 (jkratten@law.georgetown.edu)
	//
	// Explanation of first line:
	// flf2 - "magic number" for file identification
	// a    - should always be `a', for now
	// $    - the "hardblank" -- prints as a blank, but can't be smushed
	// 7    - height of a character
	// 5    - height of a character, not including descenders
	// 10   - max line length (excluding comment lines) + a fudge factor
	// -1   - default smushmode for this font (like "-m 15" on command line)
	// 12   - number of comment lines

	header := strings.Split(lines[0], " ")

	font := &Font{}
	font.Hardblank = header[0][len(header)-1:]
	font.Height, _ = strconv.Atoi(header[1])

	commentEndLine, _ := strconv.Atoi(header[5])
	font.FontSlice = lines[commentEndLine+1:]

	return font, nil
}
