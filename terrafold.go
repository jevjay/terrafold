package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
}

type Module []struct {
	Name string
	Content
}

type Directory []struct {
	Path string
	Content
}

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

		for i := range folders {
			f := folders[i]
			createDirectory(f)
			fmt.Printf("Generating %s directory\n", f)
		}

		for i := range files {
			f := files[i]
			createFile(f)
			fmt.Printf("Generating %s file\n", f)
		}
	} else {
		template, err := ioutil.ReadFile(*tplFile)
		var config Config

		check(err)

		err = yaml.Unmarshal(template, &config)
		check(err)

		fmt.Printf("%#v\n", config.Module)

		if len(config.Root) > 0 {
			for i := range config.Root {
				f := config.Root[i]
				createFile(strings.Join([]string{*repo, f}, "/"))
				fmt.Printf("Generating %s file\n", f)
			}
		}

		if len(config.Module) > 0 {
			createDirectory(strings.Join([]string{*repo, "modules"}, "/"))
			fmt.Printf("Generating modules directory\n")
			for i := range config.Module {
				moduleName := config.Module[i].Name
				createDirectory(strings.Join([]string{*repo, "modules", moduleName}, "/"))
				fmt.Printf("Generating modules %s directory\n", moduleName)
				// Generating module content
				for j := range config.Module[i].Content {
					fileName := config.Module[i].Content[j].Name
					createFile(strings.Join([]string{*repo, "modules", moduleName, fileName}, "/"))
					fmt.Printf("Generating modules %s content: %s \n", moduleName, fileName)
				}
			}
		}

		if len(config.Directory) > 0 {
			for i := range config.Directory {
				dirPath := config.Directory[i].Path
				createPath(strings.Join([]string{*repo, dirPath}, "/"))
				fmt.Printf("Generating directory %s\n", dirPath)
				// Generating directory content
				for j := range config.Directory[i].Content {
					fileName := config.Directory[i].Content[j].Name
					createFile(strings.Join([]string{*repo, dirPath, fileName}, "/"))
					fmt.Printf("Generating direcotry %s content: %s \n", dirPath, fileName)
				}
			}
		}
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

func createFile(filePath string) {
	newFile, err := os.Create(filePath)
	check(err)
	newFile.Close()
}
