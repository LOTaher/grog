package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/LOTaher/grog/pkg/util"
	"github.com/spf13/cobra"
)

var npmRegistryURL = "https://registry.npmjs.org"

var install = &cobra.Command{
    Use:   "install [package]",
    Short: "Install a package.",
    Long:  `Install a package. Example: grog install express`,
    Run:   installPackage,
}

type Installer struct {
    Name    string
    Version string
}

func (i *Installer) parsePackageDetails(pkg string) error {
    parts := strings.Split(pkg, "@")
    i.Name = parts[0]
    i.Version = "latest" 

    if len(parts) > 1 && parts[1] != "" {
        ok, err := util.ValidVersion(parts[1])
        if err != nil {
            return fmt.Errorf("invalid version '%s': %w", parts[1], err)
        }
        if ok {
            i.Version = parts[1]
        }
    }
    return nil
}

func installPackage(cmd *cobra.Command, args []string) {
    if len(args) < 1 {
        fmt.Println("Please specify a package name to install.")
        return
    }

    installer := Installer{}
    if err := installer.parsePackageDetails(args[0]); err != nil {
        fmt.Println("Error parsing package details:", err)
        os.Exit(1)
    }

    fmt.Printf("Preparing to install package: %s@%s\n", installer.Name, installer.Version)

    if err := performInstallation(installer.Name, installer.Version); err != nil {
        fmt.Println("Installation failed:", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully installed package: %s@%s\n", installer.Name, installer.Version)
}

func performInstallation(name, version string) error {
    url := fmt.Sprintf("%s/%s/%s", npmRegistryURL, name, version)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    req.Header.Set("Accept", "application/vnd.npm.install-v1+json; q=1.0, application/json; q=0.8, */*")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // TODO: Handle package response. Only printing right now.

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    fmt.Printf("Response from npm registry: %s\n", string(body))

    // TODO: Implement actual package download and installation logic here.

    return nil
}