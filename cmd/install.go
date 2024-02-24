package cmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

type Response struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
	Dist         struct {
		Tarball string `json:"tarball"`
	} `json:"dist"`
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

	fmt.Printf("Package: %s\nVersion: %s\n Tarball URL: %s\n", packageInfo.Name, packageInfo.Version, packageInfo.Dist.Tarball)

	// TODO: Implement actual package download and installation logic here.

	downloadTarball(packageInfo.Dist.Tarball, "test")

	return nil
}

func downloadTarball(url, targetDir string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		outputPath := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				return err
			}

		case tar.TypeReg:
			outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}

			outFile.Close()

		}
	}

	return nil

}
