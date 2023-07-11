package cmd

import (
	"github.com/danilmir/cloudphoto/utils"
	"github.com/danilmir/cloudphoto/yandex_cloud"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	album string
	path  string
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "To send photos to cloud storage, use the command",
	Long:  `To send photos to cloud storage, use the command`,
	Run: func(cmd *cobra.Command, args []string) {
		if album == "" {
			println("Album not given")
			os.Exit(1)
		}

		if path == "" {
			ex, err := os.Executable()
			if err != nil {
				panic(err)
			}
			path = filepath.Dir(ex)
		}

		imgs := utils.GetImagePaths(path)

		if len(imgs) == 0 {
			println("Warning: Photos not found in directory", path)
			os.Exit(1)
		}

		yandex_cloud.LoadFromConfigFile()
		yandex_cloud.Init()
		if !yandex_cloud.IsBucketExist() {
			println("Bucket not exist")
		}

		yandex_cloud.CreateFolderIfNotExists(album)
		yandex_cloud.UploadImagesToAlbum(album, imgs)
	},
}

func init() {
	uploadCmd.Flags().StringVar(&album, "album", "", "Name of photo album")
	uploadCmd.Flags().StringVar(&path, "path", "", "Absolute or relative path to the photo directory.")
	rootCmd.AddCommand(uploadCmd)
}
