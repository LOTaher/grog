package cmd

import (
	"fmt"
	"os"

	"github.com/LOTaher/grog/internal/cache"
	"github.com/spf13/cobra"
)

var clear = &cobra.Command{
	Use:   "clear",
	Short: "Clear the cache.",
	Run:   clearCache,
}

func clearCache(cmd *cobra.Command, args []string) {
    os.RemoveAll(cache.Cache)
    fmt.Println("Cache cleared.")
}
