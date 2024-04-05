package cmd

import (
	"path/filepath"
	"fmt"
	"time"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/itrepablik/kopy"
	"github.com/itrepablik/itrlog"
	"strings"
)

var copydirCmd = &cobra.Command{
	Use:   "copydir",
	Short: "Copy the entire folder or a directory without a compression",
	Long: `copydir command is to copy the entire directory or folder including its sub-folders and sub-directories contents.
Take note that, it will replace any existing files and its contents to the destination directory or a folder.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src:=filepath.FromSlash(args[0])
		dst:=filepath.FromSlash(args[1])

		IgnoreFileTypes = viper.Get("ignore.file_type_or_folder_name")
		IGFT := fmt.Sprint(IgnoreFileTypes)
		IgnoreFT = strings.Split(IGFT, ",")

		msg := `Starts copying the entire directory or a folder: `
		fmt.Println(msg, src)
		Sugar.Sugar.Infow(msg, "src", src, "log_time", time.Now().Format(itrlog.LogTimeFormat))

		// Starts copying the entire directory or a folder.
		// filesCopied, foldersCopied, err := kopy.CopyDir(src, dst, IsLogCopiedFile, IgnoreFT, Sugar)
		filesCopied, foldersCopied, err := kopy.CopyDir(src, dst, IsLogCopiedFile, IgnoreFT)
		if err != nil {
			fmt.Println(err)
			Sugar.Sugar.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			return
		}

		// Give some info back to the user's console and the logs as well.
		msg = `Successfully copied the entire directory or a folder: `
		fmt.Println(msg, src, ", Number of Folders Copied: ", filesCopied, " Number of Files Copied: ", foldersCopied)
		Sugar.Sugar.Infow(msg, "src", src, "dst", dst, "folder_copied", filesCopied, "files_copied", foldersCopied, "log_time", time.Now().Format(itrlog.LogTimeFormat))
	},
}

func init() {
	rootCmd.AddCommand(copydirCmd)
}