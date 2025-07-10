package src

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type Frontmatter struct {
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	WPM     int    `yaml:"wpm"`
	Draft   bool   `yaml:"draft"`
	Created string `yaml:"created"`
	Updated string `yaml:"updated"`
	Author  string `yaml:"author"`
}

func ScanFrontmatter(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var frontmatterLines []string

	inFrontmatter := false

	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Line:", line)

		if strings.TrimSpace(line) == "---" {
			if inFrontmatter {
				inFrontmatter = false
				break // End of frontmatter
			}
			inFrontmatter = true
			continue
		}

		if inFrontmatter {
			frontmatterLines = append(frontmatterLines, line)
		}
	}

	if inFrontmatter {
		return "", errors.New("frontmatter not closed with '---'")
	} else if len(frontmatterLines) == 0 {
		return "", errors.New("no frontmatter found")
	} else if err := scanner.Err(); err != nil {
		return "", err
	}

	frontmatter := strings.Join(frontmatterLines, "\n")

	return frontmatter, nil
}

func ParseFrontmatter(content string) (*Frontmatter, error) {
	var fm Frontmatter

	if err := yaml.Unmarshal([]byte(content), &fm); err != nil {
		return nil, err
	}

	return &fm, nil
}
