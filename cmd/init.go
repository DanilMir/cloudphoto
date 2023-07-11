package cmd

import (
	"fmt"
	"github.com/danilmir/cloudphoto/yandex_cloud"
	"github.com/spf13/cobra"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init program configuration",
	Long:  `init program configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		var bucket, keyId, accessKey string
		fmt.Print("Bucket name: ")
		fmt.Scanln(&bucket)
		fmt.Print("AWS access key id: ")
		fmt.Scanln(&keyId)
		fmt.Print("Access Key: ")
		fmt.Scanln(&accessKey)

		err := yandex_cloud.CreateConfigFile(bucket, keyId, accessKey)
		if err != nil {
			println("Error:", err)
			os.Exit(1)
		}

		yandex_cloud.LoadFromConfigFile()
		yandex_cloud.Init()
		if !yandex_cloud.IsBucketExist() {
			err = yandex_cloud.CreateBucket()
			if err != nil {
				println("Error:", err)
				os.Exit(1)
			}
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
