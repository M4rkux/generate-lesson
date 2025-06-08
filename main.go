package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <lesson_name> [chapter_folder]")
		os.Exit(1)
	}

	lessonName := os.Args[1]
	chapterFolder := ""
	if len(os.Args) > 2 {
		chapterFolder = os.Args[2]
	}

	var folderName string
	if chapterFolder != "" {
		err := os.MkdirAll(chapterFolder, 0755)
		if err != nil {
			fmt.Printf("Failed to create chapter directory: %s\n", chapterFolder)
			os.Exit(1)
		}

		err = os.Chdir(chapterFolder)
		if err != nil {
			fmt.Printf("Failed to enter chapter directory: %s\n", chapterFolder)
			os.Exit(1)
		}

		folderName = filepath.Base(chapterFolder)
	} else {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get current directory")
			os.Exit(1)
		}
		folderName = filepath.Base(currentDir)
	}

	chapterName := strings.SplitN(folderName, "-", 2)[0]

	err := os.Mkdir(lessonName, 0755)
	if err != nil {
		fmt.Printf("Failed to create lesson directory: %s\n", lessonName)
		os.Exit(1)
	}

	err = os.Chdir(lessonName)
	if err != nil {
		fmt.Printf("Failed to enter lesson directory: %s\n", lessonName)
		os.Exit(1)
	}

	fmt.Printf("Creating lesson files for: %s\n", lessonName)

	// Initialize Go module
	modPath := fmt.Sprintf("example.com/%s/%s", chapterName, lessonName)
	cmd := exec.Command("go", "mod", "init", modPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to initialize go module")
		os.Exit(1)
	}

	// Create main.go and main_test.go
	for _, fname := range []string{"main.go", "main_test.go"} {
		file, err := os.Create(fname)
		if err != nil {
			fmt.Printf("Failed to create file: %s\n", fname)
			os.Exit(1)
		}
		file.Close()
	}

	// Open files in VS Code
	cmd = exec.Command("code", "main.go", "main_test.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to open files in VS Code")
		os.Exit(1)
	}
}
