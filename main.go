package main

import (
	"bytes"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var directory string
var importPath string

func main()  {
	if os.Args[1] == "init" {

		// Clone repository : clone letsGo
		importPath = os.Args[2]
		directory = os.Args[3]

		pwd := exec.Command("pwd")

		var out bytes.Buffer
		pwd.Stdout = &out
		err := pwd.Run()
		if err != nil {
			log.Fatal(err)
		}
		path := strings.TrimSuffix(out.String(), "\n")+"/"+directory

		fmt.Println("Cloning letsGo into : " + directory)

		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL:      "https://github.com/letsgo-framework/letsGo",
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Cloning complete")

		// Change package name : change package name in glide.yaml to your package name
		read, err := ioutil.ReadFile(path+"/glide.yaml")
		if err != nil {
			panic(err)
		}

		newContents := strings.Replace(string(read), "letsGo", importPath+"/"+directory, -1)

		err = ioutil.WriteFile(path+"/glide.yaml", []byte(newContents), 0)
		if err != nil {
			panic(err)
		}

		fmt.Println("gilde updated")
		
		// change the internal package (controllers, tests, helpers etc.) paths as per your requirement
		err = filepath.Walk(path+"/controllers", Visit)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Controllers refracted")
		}
		
		err = filepath.Walk(path+"/gql", Visit)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("graphql refracted")
		}

		err = filepath.Walk(path+"/middlewares", Visit)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("middlewares refracted")
		}
		err = filepath.Walk(path+"/routes", Visit)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("routes refracted")
		}
		
		err = filepath.Walk(path+"/tests", Visit)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("tests refracted")
		}

		err = filepath.Walk(path, Visit)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("main refracted")
		}

		fmt.Println("refraction Done")

		// setup .env and .env.testing
		_, _ = exec.Command("cp", path+"/.env.example", path+"/.env").Output()
		_, _ = exec.Command("cp", path+"/.env.example", path+"/.env.testing").Output()

		// Glide
		gl := exec.Command("glide", "install")
		gl.Dir = path
		err = gl.Run()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("env updated")

		// remove .git
		gr := exec.Command("rm", "-rf", ".git")
		gr.Dir = path
		err = gr.Run()
		if err != nil {
			log.Fatal(err)
		}

	} else if os.Args[1] == "generate" || os.Args[1] == "g" {
		switch os.Args[2] {
			case "component", "c":
				fmt.Println("Generating component : "+os.Args[3])
				componentContent := []byte("package controllers")
				err := ioutil.WriteFile("./controllers/"+os.Args[3]+".go", componentContent, 0644)
				if err != nil {
					log.Fatal(err)
				}
				break
			case "type", "t":
				fmt.Println("Generating type : "+os.Args[3])
				typeContent := []byte("package types")
				err := ioutil.WriteFile("./types/"+os.Args[3]+".go", typeContent, 0644)
				if err != nil {
					log.Fatal(err)
				}
				break
			default:
				fmt.Println("Invalid argument")
				break
		}
	}
}

func Visit(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil //
	}

	matched, err := filepath.Match("*.go", fi.Name())

	if err != nil {
		panic(err)
		return err
	}

	if matched {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		newContents := strings.Replace(string(read), "github.com/letsGo", importPath+"/"+directory, -1)
		
		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			panic(err)
		}

	}

	return nil
}