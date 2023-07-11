package cmd

import (
	"github.com/danilmir/cloudphoto/yandex_cloud"
	"github.com/spf13/cobra"
)

var mksiteCmd = &cobra.Command{
	Use:   "mksite",
	Short: "Generating and publishing photo archive web pages",
	Long:  `Generating and publishing photo archive web pages`,
	Run: func(cmd *cobra.Command, args []string) {
		yandex_cloud.LoadFromConfigFile()
		yandex_cloud.Init()
		if !yandex_cloud.IsBucketExist() {
			println("Bucket not exist")
		}

		yandex_cloud.CreateWebsite()
	},
}

func init() {
	rootCmd.AddCommand(mksiteCmd)
}
