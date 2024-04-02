package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
    "path/filepath"
)

type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Scripts         map[string]string `json:"scripts"`
	Keywords        []string          `json:"keywords"`
	Author          string            `json:"author"`
	License         string            `json:"license"`
	Private         bool              `json:"private"`
	DevDependencies map[string]struct {
		Version string `json:"version"`
	} `json:"devDependencies"`
	Dependencies map[string]struct {
		Version string `json:"version"`
	} `json:"dependencies"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new grog project.",
	Run:   initProject,
}

func initProject(cmd *cobra.Command, args []string) {
    packageJSONPath := "package.json"

    nameOfProject, err := os.Getwd()
    if err != nil {
        fmt.Println("Error getting current working directory")
        return
    }
    nameOfProject = filepath.Base(nameOfProject)

	packageJSON := PackageJSON{
		Name:        nameOfProject,
		Version:     "1.0.0",
		Description: "",
		Main:        "index.js",
		Scripts: map[string]string{
			"start": "node index.js",
		},
		Keywords: []string{},
		Author:   "",
		License:  "",
		Private:  true,
		DevDependencies: map[string]struct {
			Version string `json:"version"`
		}{},
		Dependencies: map[string]struct {
			Version string `json:"version"`
		}{},
	}

    packageJSONBytes, err := json.MarshalIndent(packageJSON, "", "  ")   
    if err != nil {
        fmt.Println("Error creating package.json")
        return
    }

    if _, err := os.Stat(packageJSONPath); err == nil {
        fmt.Println("package.json already exists")
        return
    }
    
    err = os.WriteFile(packageJSONPath, packageJSONBytes, 0644)
    if err != nil {
        fmt.Println("Error creating package.json")
        return
    }

    fmt.Println("grog project initialized.")
}
