package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

const GoFileExtension = ".go"

func processFile(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	lineCount := countLinesInFile(file)
	functionCount := countFunctionsInFile(path)

	return lineCount, functionCount, nil
}

func countLinesInFile(file *os.File) int {
	scanner := bufio.NewScanner(file)
	var lineCount int
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return 0
	}
	return lineCount
}

func countFunctionsInFile(path string) int {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return 0
	}
	functionCount := 0
	ast.Inspect(node, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncDecl:
			functionCount++
		}
		return true
	})
	return functionCount
}

func countLinesAndFunctions(path string) (int, int, error) {
	var totalLineCount, totalFunctionCount int
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == GoFileExtension {
			lineCount, functionCount, processFileError := processFile(path)
			if processFileError != nil {
				return processFileError
			}
			totalLineCount += lineCount
			totalFunctionCount += functionCount
			fmt.Printf("Lines in %s: %d; Functions: %d\n", filepath.Base(path), lineCount, functionCount)
		}
		return nil
	})
	return totalLineCount, totalFunctionCount, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path>")
		return
	}
	path := os.Args[1]

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	projectName := filepath.Base(path)
	if path == "." {
		projectName = filepath.Base(currentDir)
	}
	fmt.Println("Project Name:", projectName)
	fmt.Printf("-------------\n")

	totalLineCount, totalFunctionCount, err := countLinesAndFunctions(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("-------------\n")
	fmt.Printf("Total lines in.go files: %d\n", totalLineCount)
	fmt.Printf("Total functions in.go files: %d\n", totalFunctionCount)
}
