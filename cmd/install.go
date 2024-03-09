package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/LOTaher/grog/internal/cache"
	"github.com/LOTaher/grog/internal/symlink"
	"github.com/LOTaher/grog/internal/tarball"
	ver "github.com/LOTaher/grog/internal/version"
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

/* Put in new file */
type Response struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
	Dist         struct {
		Tarball string `json:"tarball"`
	} `json:"dist"`
}
/* Put in new file */

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

	// Check if the package name is located in the cache. If it is, read it's lock file and check if the version is already installed.
	// If it is, skip the installation. If it is labeled as "isLatest" in the cache, then do the symlink, if it is not the latest, proceed
	// with the installation.

	// if exists, err := cache.PackageVersionCached(name, version); exists {
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	lockFile, err := cache.ReadLockFile(name)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	if lockFile.isLatest {
	// 		fmt.Printf("Package %s@%s already exists in the cache. Skipping installation.\n", name, version)
	// 		if err := symlink.SymlinkPackage(name, version); err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	}
	// }
    
    /* Put in new file */
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var packageInfo Response
	if err := json.Unmarshal(body, &packageInfo); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return err
	}
    /* Put in new file */

	exists, err := cache.IsVersionCached(packageInfo.Name, packageInfo.Version)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Package %s@%s already exists in the cache. Skipping installation.\n", packageInfo.Name, packageInfo.Version)
	} else {
		targetDir := filepath.Join(cacheDir, packageInfo.Name, packageInfo.Version)

		if err := tarball.DownloadTarball(packageInfo.Dist.Tarball, targetDir); err != nil {
			return fmt.Errorf("failed to download tarball: %w", err)
		}
		fmt.Printf("Successfully installed %s@%s\n", packageInfo.Name, packageInfo.Version)
	}

	if err := symlink.SymlinkPackage(packageInfo.Name, packageInfo.Version); err != nil {
		return err
	}

    // Add the packageinfo.Dependencies to the lock file, if it doesn't already exist.

	for depName, depVersion := range packageInfo.Dependencies {
		fmt.Printf("Installing dependency %s@%s\n", depName, depVersion)
		if err := performInstallation(depName, depVersion); err != nil {
			return fmt.Errorf("failed to install dependency %s@%s: %w", depName, depVersion, err)
		}
	}

	return nil
}
