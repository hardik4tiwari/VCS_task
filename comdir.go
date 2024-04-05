package cmd

import (
	"bytes"
	"fmt"
	"path"

	// "Backup/config"
	"io"
	"os"
	"path/filepath"
	"time"

	// "github.com/fatih/color"
	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/kopy"
	"github.com/spf13/cobra"
)

var comdircmd = &cobra.Command{
	Use: "Comdir",
	Short: "Compress the directory or a folder",
	Long:  `ada`,
	// Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	// },
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src:= filepath.FromSlash(args[0])
		dst:= filepath.FromSlash(args[1])
		
		msg := `Start compressing the directory or a folder:`
		fmt.Println(msg, src)
		Sugar.Sugar.Infow(msg, "src", src, "dst", dst, "log_time", time.Now().Format(itrlog.LogTimeFormat))

		// Compose the zip filename
		fnWOext := kopy.FileNameWOExt(filepath.Base(args[0])) // Returns a filename without an extension.
		zipDir := fnWOext + kopy.ComFileFormat

		// To make directory path separator a universal, in Linux "/" and in Windows "\" to auto change
		// depends on the user's OS using the filepath.FromSlash organic Go's library.
		zipDest := filepath.FromSlash(path.Join(args[1], zipDir))

		// Start compressing the entire directory or a folder using the tar + gzip
		var buf bytes.Buffer
		if err := kopy.CompressDIR(src, &buf, IgnoreFT); err != nil {
			fmt.Println(err)
			Sugar.Sugar.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			return
		}

		// write the .tar.gzip
		os.MkdirAll(dst, os.ModePerm) // Create the root folder first
		fileToWrite, err := os.OpenFile(zipDest, os.O_CREATE|os.O_RDWR, os.FileMode(600))
		if err != nil {
			Sugar.Sugar.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			panic(err)
		}
		if _, err := io.Copy(fileToWrite, &buf); err != nil {
			Sugar.Sugar.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			panic(err)
		}
		defer fileToWrite.Close()

		msg = `Done compressing the directory or a folder:`
		fmt.Println(msg, src)
		Sugar.Sugar.Infow(msg, "dst", zipDest, "log_time", time.Now().Format(itrlog.LogTimeFormat))
	},
}

func init(){
	rootCmd.AddCommand(comdircmd)
}