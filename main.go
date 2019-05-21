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
	"text/tabwriter"
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

	logCommand := flag.NewFlagSet("log", flag.ExitOnError)

	clearLog := logCommand.Bool("clear", false, "Clear log file")
	c := logCommand.Bool("c", false, "Clear log file")

	dockerizeCommand := flag.NewFlagSet("dockerize", flag.ExitOnError)

	buildCommand := flag.NewFlagSet("build", flag.ExitOnError)

	switch os.Args[1] {
		case "init", "i":
			initCommand.Parse(os.Args[2:])
		case "run", "r":
			runCommand.Parse(os.Args[2:])
		case "help", "h":
			helpCommand.Parse(os.Args[2:])
		case "generate", "g":
			generateCommand.Parse(os.Args[2:])
		case "log":
			logCommand.Parse(os.Args[2:])
		case "dockerize", "d":
			dockerizeCommand.Parse(os.Args[2:])
		case "build", "b":
			buildCommand.Parse(os.Args[2:])
		default:
			fmt.Printf("%q is not valid command.\n", os.Args[1])
			os.Exit(2)
	}

	if initCommand.Parsed() {

		// Clone repository : clone letsgo
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

		fmt.Println("Cloning letsgo into : " + directory)

		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL:      "https://github.com/letsgo-framework/letsGo",
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Cloning complete")

		// Checkout latest tag
		checkout := exec.Command("git", "checkout", "1.0.0-beta.1")
		checkout.Dir = path
		err = checkout.Run()
		if err != nil {
			log.Fatal(err)
		}

		// Change package name : change package name in glide.yaml to your package name
		read, err := ioutil.ReadFile(path+"/glide.yaml")
		if err != nil {
			panic(err)
		}

		newContents := strings.Replace(string(read), "github.com/letsgo-framework/letsgo", importPath+"/"+directory, -1)

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

		/* No middleware available
		err = filepath.Walk(path+"/middlewares", Visit)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("middlewares refracted")
		}
		*/

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
		if *livereload == true || *l == true {
			freshCommand := exec.Command("fresh")
			freshCommand.Stdout = os.Stdout
			freshCommand.Stderr = os.Stderr

			_ = freshCommand.Run()
		} else {
			runMainCommand := exec.Command("go", "run", "main.go")
			runMainCommand.Stdout = os.Stdout
			runMainCommand.Stderr = os.Stderr

			_ = runMainCommand.Run()
		}
	}

	if helpCommand.Parsed() {
		Usage()
	}

	if logCommand.Parsed() {
		if *clearLog == true || *c == true{
			emptyLogFile, err := os.Create("./log/letsgo.log")
			if err != nil {
				log.Fatal(err)
			}
			emptyLogFile.Close()

			fmt.Println("Log cleared")
		} else {
			lineCount := os.Args[2]
			tailCommand := exec.Command("tail", "-"+lineCount+"f", "log/letsgo.log")
			tailCommand.Stdout = os.Stdout
			tailCommand.Stderr = os.Stderr

			_ = tailCommand.Run()
		}
	}

	if dockerizeCommand.Parsed() {
		binaryName := os.Args[2]
		read, err := ioutil.ReadFile("./Dockerfile")
		if err != nil {
			panic(err)
		}

		newContents := strings.Replace(string(read), "letsgo", binaryName, -1)

		err = ioutil.WriteFile("./Dockerfile", []byte(newContents), 0)
		if err != nil {
			panic(err)
		}

		fmt.Println("Dockerized")
	}

	if buildCommand.Parsed() {
		build := exec.Command("go", "build")
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr

		_ = build.Run()

		fmt.Println("Binary is ready")
	}
}

func Usage() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintf(w, "log <ACTION> \t Tail or Clear log\n")
	fmt.Fprintf(w, "\t -clear or -c to clrear log\n")
	fmt.Fprintf(w, "\t <LINE_NUMBER> to tail log file\n")
	fmt.Fprintf(w, "init <IMPORT_PATH> <PROJECT_NAME> \t Create a new letsgo project \t short i\n")
	fmt.Fprintf(w, "generate <FILE_TYPE> <FILE_NAME> \t Generate file of controller of type \t short g\n")
	fmt.Fprintf(w, "build \t Build binary \t short b\n")
	fmt.Fprintf(w, "dockerize <BINARY_NAME> \t Dockerize your project \t short d\n")
	fmt.Fprintf(w, "help \t Print usage \t short h\n")
	fmt.Fprintf(w, "run \t Run your project \t short r\n")
	fmt.Fprintf(w, "\t -livereload or -l with livereload\n")
	w.Flush()
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