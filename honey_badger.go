package main

import(
	"fmt"
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"encoding/json"
	"io/ioutil"
	"log"
)

type Exts struct {
	Extarr	[]Extensions	`json:"exts"`
}

type Extensions struct {
	Ext		string	`json:"ext"`
	Skip	string	`json:"skip"`
}

type Dirs struct {
	Dirarr	[]Directories	`json:"dirs"`
}

type Directories struct {
	Dir		string	`json:"dir"`
	Skip	string	`json:"skip"`
}

type Testing struct {
	Exts	map[string]string
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please give me the root directory:")

	scanner.Scan()

	root := scanner.Text()

	fmt.Println("What word am I looking for?")

	scanner.Scan()

	word := scanner.Text()

	diggin(root, word)
	// testFirst(root, word)
}

func testFirst(root string, word string) {
	configJson, err1 := ioutil.ReadFile("excludes.json") 
	if err1 != nil {
		fmt.Println("error reading json: ", err1)
	}

	var exts Exts

	err2 := json.Unmarshal(configJson, &exts)
	if err2 != nil {
		fmt.Println("error unmarshalling json: ", err2)
	}

	var dirs Dirs

	err3 := json.Unmarshal(configJson, &dirs)
	if err3 != nil {
		fmt.Println("error unmarshalling json: ", err3)
	}

	log.Printf("ext %s\n", exts)
	log.Printf("dir %s\n",dirs)

	fmt.Println(len(exts.Extarr))

	var rgx []string

	for i := 0; i < len(exts.Extarr); i++ {
		if exts.Extarr[i].Skip == "y" {
			rgx = append(rgx, "^[A-Za-z0-9]*[.]"+ exts.Extarr[i].Ext+"$")
		}
	}

	for _, val := range rgx {
		fmt.Println(val)
	}
	
}

func diggin(root string, word string) {

	var skip bool = false

	fmt.Println("I'm diggin a hole")

	configJson, err1 := ioutil.ReadFile("excludes.json") 
	if err1 != nil {
		fmt.Println("error reading json: ", err1)
	}

	var exts Exts

	err2 := json.Unmarshal(configJson, &exts)
	if err2 != nil {
		fmt.Println("error unmarshalling json: ", err2)
	}

	var dirs Dirs

	err3 := json.Unmarshal(configJson, &dirs)
	if err3 != nil {
		fmt.Println("error unmarshalling json: ", err3)
	}

	thisWord := "(?i)" + word
	regWord, err4 := regexp.Compile(thisWord)
	if err4 != nil {
		fmt.Println(err4)
	}

	var rgx []string

	for i := 0; i < len(exts.Extarr); i++ {
		if exts.Extarr[i].Skip == "y" {
			rgx = append(rgx, "^[A-Za-z0-9]*[.]"+ exts.Extarr[i].Ext+"$")
		}
	}

	var drx []string

	for i := 0; i < len(dirs.Dirarr); i++ {
		if dirs.Dirarr[i].Skip == "y" {
			drx = append(drx, dirs.Dirarr[i].Dir)
		}
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		
		if info.IsDir() {
			for _, dirs := range drx {
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
					skip = true
					break
				} else {
					skip = false
					continue
				}


			}
			
			if skip == false {
				p, err := os.Open(path)
					if err != nil {
						fmt.Println("messed up opening file")
						return nil
					}
					scanner := bufio.NewScanner(p)

					defer p.Close()
		
					for scanner.Scan() {
						if regWord.MatchString(scanner.Text()) == true {
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