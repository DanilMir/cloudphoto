package cmd

import (
	"fmt"
	"github.com/danilmir/cloudphoto/yandex_cloud"
	"github.com/spf13/cobra"
	"os"
	"sort"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list alubums",
	Long:  `list alubums`,
	Run: func(cmd *cobra.Command, args []string) {
		yandex_cloud.LoadFromConfigFile()
		yandex_cloud.Init()
		if !yandex_cloud.IsBucketExist() {
			println("Bucket not exist")
		}

		albums, err := yandex_cloud.GetAllRootFolders()
		if err != nil {
			println(err)
		}

		if len(albums) == 0 {
			println("Photo albums not found")
			os.Exit(1)
		}

		sort.Strings(albums)

		for _, album := range albums {
			fmt.Println(album)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
