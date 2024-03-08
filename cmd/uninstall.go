package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	//"github.com/LOTaher/grog/internal/cache"
	ver "github.com/LOTaher/grog/internal/version"
	"github.com/spf13/cobra"
)

var uninstall = &cobra.Command{
	Use:   "uninstall [package]",
	Short: "Uninstall a package.",
	Long:  `Uninstall a package. Example: grog uninstall express`,
	Run:   uninstallPackage,
}

type Uninstaller struct {
	Name    string
	Version string
	Global  bool
}

func (u *Uninstaller) parsePackageDetails(pkg string) error {
	var packageName, version string
	atCount := strings.Count(pkg, "@")

	if atCount == 0 {
		packageName = pkg
		version = ""
	} else if atCount == 1 {
		parts := strings.SplitN(pkg, "@", 2)
		packageName = parts[0]
		version = parts[1]
	} else if atCount == 2 {
		parts := strings.SplitN(pkg, "@", 3)
		packageName = parts[0] + "@" + parts[1]
		version = parts[2]
	} else {
		return fmt.Errorf("invalid package format")
	}

	if version == "" {
		version = "latest"
	} else {
		ok, err := ver.ValidVersion(version)
		if err != nil {
			return fmt.Errorf("invalid version '%s': %w", version, err)
		}
		if !ok {
			return fmt.Errorf("version '%s' is not a valid version", version)
		}
	}

	u.Name = packageName
	u.Version = version

	return nil
}

func uninstallPackage(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please specify a package name to uninstall.")
		return
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(args))

	for _, arg := range args {
		wg.Add(1)
		go func(arg string) {
			defer wg.Done()

			uninstaller := Uninstaller{}
			if err := uninstaller.parsePackageDetails(arg); err != nil {
				errChan <- fmt.Errorf("error parsing package details for %s: %w", arg, err)
				return
			}

			global, _ := cmd.Flags().GetString("g")
			uninstaller.Global = global != ""

			scope := "package"
			if uninstaller.Global {
				scope = "package globally"
			}
			fmt.Printf("Preparing to uninstall %s: %s@%s\n", scope, uninstaller.Name, uninstaller.Version)

			if err := performUninstallation(uninstaller.Name, uninstaller.Version, uninstaller.Global); err != nil {
				errChan <- fmt.Errorf("uninstallation failed for %s@%s: %w", uninstaller.Name, uninstaller.Version, err)
			}

		}(arg)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func performUninstallation(name, version string, global bool) error {

	// fmt.Println(InstalledPackages)
	//
	// versions, ok := InstalledPackages[name]
	// if !ok {
	// 	return fmt.Errorf("package %s is not installed", name)
	// }
	//
	// fmt.Printf("Uninstalling package: %s@%s\n", name, version)
	//
	// for i, v := range versions {
	// 	if v == version || version == "latest" {
	// 		InstalledPackages[name] = append(versions[:i], versions[i+1:]...)
	//
	//            if global {
	//                os.RemoveAll(cache.Cache + "/" + name + "/" + version)
	//            }
	//
	//            os.RemoveAll("./node_modules/" + name)
	//
	// 		if len(InstalledPackages[name]) == 0 {
	// 			delete(InstalledPackages, name)
	// 		}
	//
	// 		fmt.Printf("Uninstalled package: %s@%s\n", name, version)
	// 		return nil
	// 	}
	// }

	// if global {
	// 	os.RemoveAll(cache.Cache + "/" + name + "/" + version)
	// }

	os.RemoveAll("./node_modules/" + name)

    fmt.Printf("Uninstalled package: %s@%s\n", name, version)

	return nil
}
