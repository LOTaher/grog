package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/LOTaher/grog/internal/cache"
	"github.com/LOTaher/grog/internal/request"
	"github.com/LOTaher/grog/internal/symlink"
	"github.com/LOTaher/grog/internal/tarball"
	ver "github.com/LOTaher/grog/internal/version"
	"github.com/spf13/cobra"
)

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

	i.Name = packageName
	i.Version = version

	return nil
}

func installPackage(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please specify a package name to install.")
		return
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(args))

	for _, arg := range args {
		wg.Add(1)
		go func(arg string) {
			defer wg.Done()

			installer := Installer{}
			if err := installer.parsePackageDetails(arg); err != nil {
				errChan <- fmt.Errorf("error parsing package details for %s: %w", arg, err)
				return
			}

			fmt.Printf("Preparing to install package: %s@%s\n", installer.Name, installer.Version)

			if err := performInstallation(installer.Name, installer.Version); err != nil {
				errChan <- fmt.Errorf("installation failed for %s@%s: %w", installer.Name, installer.Version, err)
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

func performInstallation(name, version string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to get user home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, ".grog", "cache")

	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create cache directory: %w", err)
	}

	if len(version) == 1 {
		v := ver.GetMostRecentVersion(name)
		version = v
	}

	if strings.ContainsAny(version, "<>~^=") {
		foundVersion, err := ver.BestMatchingVersion(name, version)
		if err != nil {
			return fmt.Errorf("failed to resolve version constraint '%s' : %w", version, err)
		}

		version = foundVersion
	}
	if version == "latest" {
		if cached, err := cache.IsPackageCached(name); cached {
			if err != nil {
				return err
			}

			versions, err := cache.GetVersions(name)
			if err != nil {
				return err
			}

			for _, version := range versions {
				isLatest, err := cache.IsLockFileLatest(name, version)
				if err != nil {
					return err
				}
				if isLatest {
					if err := symlink.SymlinkPackage(name, version); err != nil {
						return err
					}

					lockfile, err := cache.ReadLockFile(name, version)
					if err != nil {
						return err
					}

					fmt.Printf("Package %s@%s is latest version and already exists in the cache. Skipping installation.\n", name, version)

					for depName, depVersion := range lockfile.Dependencies {
						fmt.Printf("Installing dependency %s@%s\n", depName, depVersion)
						if err := symlink.SymlinkPackage(depName, depVersion); err != nil {
							return err
						}
					}
					return nil
				}
			}
		}
	}

	packageInfo, err := request.FetchResponse(name, version)
	if err != nil {
		return err
	}

	if exists, err := cache.IsVersionCached(packageInfo.Name, packageInfo.Version); exists {
		if err != nil {
			return err
		}

		fmt.Printf("Package %s@%s already exists in the cache. Skipping installation.\n", packageInfo.Name, packageInfo.Version)

		if err := symlink.SymlinkPackage(packageInfo.Name, packageInfo.Version); err != nil {
			return err
		}

	} else {

		isLatest, err := ver.IsLatestVersion(packageInfo.Name, packageInfo.Version)
		if err != nil {
			return err
		}

		targetDir := filepath.Join(cacheDir, packageInfo.Name, packageInfo.Version)

		if err := tarball.DownloadTarball(packageInfo.Dist.Tarball, targetDir); err != nil {
			return fmt.Errorf("failed to download tarball: %w", err)
		}

		err = cache.CreateLockFile(packageInfo.Name, packageInfo.Version, isLatest, packageInfo.Dependencies)
		if err != nil {
			return err
		}

		if err := symlink.SymlinkPackage(packageInfo.Name, packageInfo.Version); err != nil {
			return err
		}

		fmt.Printf("Successfully installed %s@%s\n", packageInfo.Name, packageInfo.Version)
	}

	for depName, depVersion := range packageInfo.Dependencies {
		fmt.Printf("Installing dependency %s@%s\n", depName, depVersion)
		if err := performInstallation(depName, depVersion); err != nil {
			return fmt.Errorf("failed to install dependency %s@%s: %w", depName, depVersion, err)
		}
	}

	return nil
}
