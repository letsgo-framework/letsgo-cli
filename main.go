package main

import (
	"bytes"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main()  {

	if os.Args[1] == "init" {

		// Clone repository : clone letsGo
		directory := os.Args[2]

		pwd := exec.Command("pwd")

		var out bytes.Buffer
		pwd.Stdout = &out
		err := pwd.Run()
		if err != nil {
			log.Fatal(err)
		}
		path := strings.TrimSuffix(out.String(), "\n")+"/"+directory

		fmt.Printf("Cloning letsGo into %v👍",directory)

		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL:      "https://github.com/Sab94/letsGo",
		})

		if err != nil {
			log.Fatal(err)
		}

		// Change package name : change package name in glide.yaml to your package name
		read, err := ioutil.ReadFile(path+"/glide.yaml")
		if err != nil {
			panic(err)
		}

		newContents := strings.Replace(string(read), "letsGo", directory, -1)

		err = ioutil.WriteFile(path+"/glide.yaml", []byte(newContents), 0)
		if err != nil {
			panic(err)
		}

		// TODO : change the internal package (controllers, tests, helpers etc.) paths as per your requirement

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
	}
}