package main

import(
	"fmt"
	"bufio"
	"strings"
	"os"
	"path/filepath"
	"regexp"
)

func main() {

	var exclude []string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please give me the root directory:")

	scanner.Scan()

	root := scanner.Text()

	fmt.Println("What word am I looking for?")

	scanner.Scan()

	word := scanner.Text()

	fmt.Println("Directories I should exclude:")

	for {
		scanner.Scan()

		text := scanner.Text()

		if text == "end" {
			if len(exclude) == 0 {
				exclude = append(exclude, "none")
			}
			break
		} else {
			exclude = append(exclude, text)
		}
	}

	diggin(root, word, exclude)
}

func diggin(root string, word string, excludes []string) {

	thisWord := "(?i)(?<=\s|^|\W)" + word + "(?=\s|$|\W)"
	regWord, _ := regexp.Compile(thisWord)

	rgx := []string{"^[A-Za-z0-9]*[.]java$", "^[A-Za-z0-9]*[.]cs$", "^[A-Za-z0-9]*[.]php$", "^[A-Za-z0-9]*[.]html$", "^[A-Za-z0-9]*[.]cfm$", "^[A-Za-z0-9]*[.]js$", "^[A-Za-z0-9]*[.]xml$"}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		
		if info.IsDir() {
			for _, dirs := range excludes {
				if info.Name() == dirs {
					return filepath.SkipDir
				}
			}
		}

		if !info.IsDir() {
			for i := 0; i < len(rgx); i++ {
				reg, _ := regexp.Compile(rgx[i])
				match := reg.MatchString(info.Name())

				if match == true {
					p, err := os.Open(path)
					if err != nil {
						fmt.Println("messed up opening file")
						return nil
					}
					scanner := bufio.NewScanner(p)
		
					for scanner.Scan() {
						// if strings.Contains(scanner.Text(), word) {
						if regWord.MatchString(scanner.Text()) {
							fmt.Println(path, scanner.Text())
						} else {
							continue
						}
					}
				} else {
					continue
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}