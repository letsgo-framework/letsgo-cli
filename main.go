package main

import (
	"bytes"
	"flag"
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

	initCommand := flag.NewFlagSet("init", flag.ExitOnError)

	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)

	runCommand := flag.NewFlagSet("run", flag.ExitOnError)
	livereload := runCommand.Bool("livereload", false, "run with livereload")
	l := runCommand.Bool("l", false, "run with livereload")

	generateCommand := flag.NewFlagSet("generate", flag.ExitOnError)
	//logCommand := flag.NewFlagSet("log", flag.ExitOnError)

	switch os.Args[1] {
		case "init", "i":
			initCommand.Parse(os.Args[2:])
		case "run", "r":
			runCommand.Parse(os.Args[2:])
		case "help", "h":
			helpCommand.Parse(os.Args[2:])
		case "generate", "g":
			generateCommand.Parse(os.Args[2:])
		default:
			fmt.Printf("%q is not valid command.\n", os.Args[1])
			os.Exit(2)
	}

	if initCommand.Parsed() {

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

	}

	if generateCommand.Parsed() {
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

	if runCommand.Parsed() {
		runCommand.Parse(os.Args[2:])
		if *livereload == true || *l == true{
			_, _ = exec.Command("fresh").Output()
		} else {
			_, _ = exec.Command("go", "run", "main.go").Output()
		}

		// TODO: Show output on console
	}

	if helpCommand.Parsed() {
		Usage()
	}
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "init <PROJECT_NAME> : Create a new letsgo project\n")
	fmt.Fprintf(os.Stderr, "generate <FILE_TYPE> <FILE_NAME> : Generate file of controller of type\n")
	fmt.Fprintf(os.Stderr, "run : Run your project\n")
	fmt.Fprintf(os.Stderr, "\t -livereload or -l \t with livereload\n")
	fmt.Fprintf(os.Stderr, "log <ACTION> : Tail or Clear log\n")
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

		newContents := strings.Replace(string(read), "github.com/letsgo-framework/letsgo", importPath+"/"+directory, -1)
		
		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			panic(err)
		}

	}

	return nil
}