package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Root []string
	Module
	Directory
}

type Content []struct {
	Name string
	Source
}

type Module []struct {
	Name string
	Content
}

type Directory []struct {
	Path string
	Content
}

type Source []struct {
	EC2    *EC2
	Bucket *Bucket
}

type AWSObject interface {
	generateContent() string
}

// Terraform object structs

func main() {
	repo := flag.String("name", "", "Repository name")
	tplFile := flag.String("template", "default", "Scaffold template file")
	flag.Parse()

	if len(*repo) > 0 {
		createDirectory(strings.Join([]string{*repo}, "/"))
		fmt.Printf("Generating %s repository\n", *repo)
	} else {
		log.Fatal("Please provide repository name by setting '-name' arg. Aborting...")
	}

	if *tplFile == "default" {
		folders := []string{
			strings.Join([]string{*repo, "modules"}, "/"),
		}

		files := []string{
			strings.Join([]string{*repo, ".gitignore.tf"}, "/"),
			strings.Join([]string{*repo, ".travis.yml"}, "/"),
			strings.Join([]string{*repo, ".travis_script.sh"}, "/"),
		}

		for _, folder := range folders {
			createDirectory(folder)
			fmt.Printf("Generating %s directory\n", folder)
		}

		for _, file := range files {
			createFile(file, "")
			fmt.Printf("Generating %s file\n", file)
		}
	} else {
		template, err := ioutil.ReadFile(*tplFile)
		check(err)
		config := Config{}

		err = yaml.Unmarshal(template, &config)
		check(err)

		fmt.Printf("%#v\n", config.Module)

		if len(config.Root) > 0 {
			for _, file := range config.Root {
				createFile(strings.Join([]string{*repo, file}, "/"), "")
				fmt.Printf("Generating %s file\n", file)
			}
		}

		if len(config.Module) > 0 {
			createDirectory(strings.Join([]string{*repo, "modules"}, "/"))
			fmt.Printf("Generating modules directory\n")
			for _, module := range config.Module {
				moduleName := module.Name
				createDirectory(strings.Join([]string{*repo, "modules", moduleName}, "/"))
				fmt.Printf("Generating modules %s directory\n", moduleName)
				// Generating module content
				for _, content := range module.Content {
					fileName := content.Name
					source := content.Source
					var buffer bytes.Buffer

					for _, sourceType := range source {
						buffer.WriteString(sourceType.EC2.generateContent())
						buffer.WriteString(sourceType.Bucket.generateContent())
					}

					createFile(strings.Join([]string{*repo, "modules", moduleName, fileName}, "/"), buffer.String())
					fmt.Printf("Generating modules %s content: %s \n", moduleName, fileName)

				}
			}
		}

		if len(config.Directory) > 0 {
			for _, dir := range config.Directory {
				dirPath := dir.Path
				createPath(strings.Join([]string{*repo, dirPath}, "/"))
				fmt.Printf("Generating directory %s\n", dirPath)
				// Generating directory content
				for _, content := range dir.Content {
					fileName := content.Name
					createFile(strings.Join([]string{*repo, dirPath, fileName}, "/"), "")
					fmt.Printf("Generating direcotry %s content: %s \n", dirPath, fileName)
				}
			}
		}
		// Lint all terraform files
		lintFiles()
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createDirectory(directoryPath string) {
	//choose your permissions well
	err := os.Mkdir(directoryPath, 0777)
	check(err)
}

func createPath(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	check(err)
}

func createFile(filePath string, content string) {
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	check(err)
}

func lintFiles() {
	cmd := "terraform"
	args := []string{"fmt", "."}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Linted terraform files")
}
