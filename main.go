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
	// testFirst(root, word, exclude)
}

func testFirst(root string, word string, exclude []string) {
	fmt.Println("Root is: ", root)
	fmt.Println("Word is: ", word)
	for _, v := range exclude {
		fmt.Println("Exclude: ", v)
	}
}

func diggin(root string, word string, excludes []string) {

	//var skip bool = false
	rgx := []string{"^[A-Za-z0-9]*[.]java$", "^[A-Za-z0-9]*[.]cs$", "^[A-Za-z0-9]*[.]php$", "^[A-Za-z0-9]*[.]html$", "^[A-Za-z0-9]*[.]go$"}

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
						if strings.Contains(scanner.Text(), word) {
							fmt.Println(path, scanner.Text())
						} else {
							continue
						}
					}
				} else {
					fmt.Println("should skip: ", info.Name())
					return filepath.SkipDir
				}
			}

			// if skip == true {
			// 	fmt.Println("should skip: ", info.Name())
			// 	return filepath.SkipDir
			// } else {
			// 	p, err := os.Open(path)
			// 	if err != nil {
			// 		fmt.Println("messed up opening file")
			// 		return nil
			// 	}
			// 	scanner := bufio.NewScanner(p)
	
			// 	for scanner.Scan() {
			// 		if strings.Contains(scanner.Text(), word) {
			// 			fmt.Println(path, scanner.Text())
			// 		} else {
			// 			continue
			// 		}
			// 	}
			// }
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}