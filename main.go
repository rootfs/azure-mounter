package main

import (
	"fmt"
	"os"
	"os/exec"

	azure "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/codegangsta/cli"
)

func main() {
	cmd := cli.NewApp()
	cmd.Name = "azurefile-mounter"
	cmd.Version = "0.1"
	cmd.Usage = "Mount Azure File Service"
	var mountpoint, share string
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "account-name",
			Usage:  "Azure storage account name",
			EnvVar: "AZURE_STORAGE_ACCOUNT",
		},
		cli.StringFlag{
			Name:   "account-key",
			Usage:  "Azure storage account key",
			EnvVar: "AZURE_STORAGE_ACCOUNT_KEY",
		},
		cli.StringFlag{
			Name:  "mountpoint",
			Usage: "Host path where volumes are mounted at",
			Value: mountpoint,
		},
		cli.StringFlag{
			Name:  "share",
			Usage: "Share Name",
			Value: share,
		},
	}
	cmd.Action = func(c *cli.Context) {
		accountName := c.String("account-name")
		accountKey := c.String("account-key")
		mountpoint := c.String("mountpoint")
		share := c.String("share")
		if accountName == "" || accountKey == "" {
			fmt.Println("azure storage account name and key must be provided.")
			return
		}
		if mountpoint == "" || share == "" {
			fmt.Println("mountpoint and share must be provided.")
			return

		}
		storageClient, err := azure.NewBasicClient(accountName, accountKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating azure client: %v", err)
			return
		}
		cl := storageClient.GetFileService()
		// Create azure file share
		if _, err := cl.CreateShareIfNotExists(share); err != nil {
			fmt.Fprintf(os.Stderr, "error creating azure file share: %v", err)
			return
		}
		if err := os.MkdirAll(mountpoint, 0700); err != nil {
			fmt.Fprintf(os.Stderr, "could not create mount point: %v", err)
			return
		}
		cmd := exec.Command("mount", "-t", "cifs", fmt.Sprintf("//%s.file.core.windows.net/%s", accountName, share), mountpoint, "-o", fmt.Sprintf("vers=3.0,username=%s,password=%s,dir_mode=0777,file_mode=0777", accountName, accountKey), "--verbose")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "mount failed: %v \noutput=%s", err, string(out))
			return
		}
		return

	}
	cmd.Run(os.Args)
}
