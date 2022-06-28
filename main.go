package main

import(
	"fmt"
	//"io/ioutil"
	"bufio"
	"strings"
	"os"
	"path/filepath"
	"regexp"
)

func main() {

	// root := "C:/Users/robcampbell/crow_engine"

	var exclude []string
	var extensions []string

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

	fmt.Println("Extensions I should exclude:")

	for {
		scanner.Scan()

		text := scanner.Text()

		if text == "end" {
			if len(extensions) == 0 {
				extensions = append(extensions, "none")
			}
			break
		} else {
			extensions = append(extensions, text)
		}
	}

	diggin(root, word, exclude, extensions)
	// testFirst(root, word, exclude)
}

func testFirst(root string, word string, exclude []string) {
	fmt.Println("Root is: ", root)
	fmt.Println("Word is: ", word)
	for _, v := range exclude {
		fmt.Println("Exclude: ", v)
	}
}

func diggin(root string, word string, excludes []string, extensions []string) {

	var skip bool = false

	isJava := regexp.Compile("[.]java$")
	isCSharp := regexp.Compile("[.]cs$")
	isPhp := regexp.Compile("[.]php$")
	isHtml := regexp.Compile("[.]html$")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		// if info.IsDir() && info.Name() == ".git" {
		// 	return filepath.SkipDir
		// }
		
		if info.IsDir() {
			for _, dirs := range excludes {
				if info.Name() == dirs {
					return filepath.SkipDir
				}
			}
		}

		if !info.IsDir() {
			for _, exts := range extensions {
				if info.Name() == exts {
					skip = true
				}
			}

			if skip == true {
				fmt.Println("should skip: ", info.Name())
				return filepath.SkipDir
			} else {
				p, err := os.Open(path)
				if err != nil {
					fmt.Println("messed up opening file")
					return nil
				}
				scanner := bufio.NewScanner(p)
	
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), word) {
						fmt.Println(path, scanner.Text())
					} else {
						continue
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}