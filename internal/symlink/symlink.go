package symlink

import (
	"fmt"
	"os"
	"path/filepath"
)

func SymlinkPackage(name, version string) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("unable to get user home directory: %w", err)
    }

    cacheDir := filepath.Join(homeDir, ".grog", "cache", name, version)
    nodeModulesDir := filepath.Join(".", "node_modules")

    if err := os.MkdirAll(nodeModulesDir, os.ModePerm); err != nil {
        return fmt.Errorf("failed to create node_modules directory: %w", err)
    }

    symlinkPath := filepath.Join(nodeModulesDir, name)

    _, err = os.Lstat(symlinkPath)
    if err != nil {
        if !os.IsNotExist(err) {
            return fmt.Errorf("failed to stat the symlink: %w", err)
        }
    } else {
        if err := os.Remove(symlinkPath); err != nil {
            return fmt.Errorf("failed to remove existing symlink: %w", err)
        }
    }

    if err := os.Symlink(cacheDir, symlinkPath); err != nil {
        return fmt.Errorf("failed to create symlink: %w", err)
    }

    fmt.Printf("Symlinked %s to node_modules\n", name)
    return nil
}