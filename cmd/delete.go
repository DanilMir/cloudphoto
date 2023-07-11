package cmd

import (
	"github.com/danilmir/cloudphoto/yandex_cloud"
	"github.com/spf13/cobra"
	"os"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "To delete an album in cloud storage, use the command",
	Long:  `To delete an album in cloud storage, use the command`,
	Run: func(cmd *cobra.Command, args []string) {
		if album == "" {
			println("Album not given")
			os.Exit(1)
		}

		yandex_cloud.LoadFromConfigFile()
		yandex_cloud.Init()
		if !yandex_cloud.IsBucketExist() {
			println("Bucket not exist")
			os.Exit(1)
		}

		if !yandex_cloud.IsFolderExists(album) {
			println("Warning: Photo album not found", album)
			os.Exit(1)
		}
		yandex_cloud.DeleteAlbum(album)
	},
}

func init() {
	deleteCmd.Flags().StringVar(&album, "album", "", "Name of photo album")
	rootCmd.AddCommand(deleteCmd)
}
