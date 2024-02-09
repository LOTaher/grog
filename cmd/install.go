package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/LOTaher/grog/pkg/util"
	"github.com/spf13/cobra"
)

var install = &cobra.Command{
    Use: "install",
    Short: "Install a package.",
    Long: `Install a package.
    Example:

    grog install express`,

    Run: installPackage,
}

// Represents the package that will be installed.
type Installer struct {
    package_name string
    package_version string
}

// Parse the inputted package details and initializes the Installer's data.
func (i *Installer) parsePackageDetails(pkg string) {
    
    parts := strings.Split(pkg, "@")
    name := parts[0]
    version := "latest"

    if len(parts) > 1 {
        version = parts[1]
    }
    
    if version == "latest" {
        i.package_version = version
        i.package_name = name
        return
    }


    ok, err := util.ValidVersion(version);
    
    if ok {
        i.package_version = version
    } else if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    i.package_name = name

    return
}

// The function that executes when the cli command is entered.
func installPackage(cmd *cobra.Command, args []string) {
    if len(args) < 1 {
        fmt.Printf("Please specify a package name to install.\n")
        return
    }
    
    i := &Installer{}

    pkg := args[0]

    i.parsePackageDetails(pkg)
    
    fmt.Printf("Installing package: %s, with version: %s\n", i.package_name, i.package_version)

}
