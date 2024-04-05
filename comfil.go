package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/kopy"

	"github.com/spf13/cobra"
)

var comfileCmd = &cobra.Command{
	Use:   "comfile",
	Short: "Compress any single file",
	Long: `comfile command will compress any single file using .zip compression format.

Example of a valid directory path in Windows:
"C:\source_folder\filename.txt" "D:\backup_destination"`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		src := filepath.FromSlash(args[0])
		dst := filepath.FromSlash(args[1])

		msg := `Start compressing the file:`
		fmt.Println(msg, src)
		Sugar.Sugar.Errorw(msg, "src", src, "dst", dst, "log_time", time.Now().Format(itrlog.LogTimeFormat))

		// Compose the zip filename
		fnWOext := kopy.FileNameWOExt(filepath.Base(args[0])) // Returns a filename without an extension.
		zipFileName := fnWOext + ".zip"

		// To make directory path separator a universal, in Linux "/" and in Windows "\" to auto change
		// depends on the user's OS using the filepath.FromSlash organic Go's library.
		zipDest := filepath.FromSlash(path.Join(args[1], zipFileName))

		// List of Files to compressed.
		files := []string{src}

		os.MkdirAll(dst, os.ModePerm) // Create the root folder first
		// if err := kopy.ComFiles(zipDest, files, Sugar); err != nil {
		if err := kopy.ComFiles(zipDest, files); err != nil {
			fmt.Println(err)
			Sugar.Sugar.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			return
		}

		msg = `Done compressing the file:`
		fmt.Println(msg, src)
		Sugar.Sugar.Infow(msg, "src", src, "dst", zipDest, "log_time", time.Now().Format(itrlog.LogTimeFormat))
	},
}

func init() {
	rootCmd.AddCommand(comfileCmd)
}