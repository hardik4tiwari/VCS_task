package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/kopy"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var copylstCmd = &cobra.Command{
	Use:   "copylst",
	Short: "Copy the latest files from a specified directory based on the modified date and time",
	Long: `copylst command will copy the latest files including the sub-folders files based on the modified date and time from the specified folder.

Open the "config.yaml" configuration file, you can change the following default settings such as:

default:
	copy_mod_files_num_days: -7 # This must be a negative value interpreted as the previous days to start copying the files.
	
	ignore:
		file_type_or_folder_name: .db, folder_name # You can specify file extentions or folder name seperated with comma.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		NumFilesCopied = 0 // Reset this variable

		// Get the default value for the "copy_mod_files_num_days" setting.
		modDays := viper.Get("default.copy_mod_files_num_days")
		mDays := modDays.(int)
		if _, ok := modDays.(int); !ok {
			mDays = -1
		}

		// Get the list of ignored file types.
		IgnoreFileTypes = viper.Get("ignore.file_type_or_folder_name")
		IGFT := fmt.Sprint(IgnoreFileTypes)
		IgnoreFT = strings.Split(IGFT, ",")

		src := filepath.FromSlash(args[0])
		dst := filepath.FromSlash(args[1])

		msg := `Starts copying the latest files from:`
		fmt.Println(msg, src)
		Sugar.Sugar.Infow(msg, "src", src, "dst", dst, "log_time", time.Now().Format(itrlog.LogTimeFormat))

		// Starts copying the latest files from.
		// if err := kopy.WalkDIRModLatest(src, dst, mDays, IsLogCopiedFile, IgnoreFT, Sugar); err != nil {
		if err := kopy.WalkDIRModLatest(src, dst, mDays, IsLogCopiedFile, IgnoreFT); err != nil {
			fmt.Println(err)
			Sugar.Sugar.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			return
		}

		// Give some info back to the user's console and the logs as well.
		msg = `Successfully copied the latest files from:`
		fmt.Println(msg, src, " Number of Files Copied: ", NumFilesCopied)
		Sugar.Sugar.Infow(msg, "src", src, "dst", dst, "copied_files", NumFilesCopied, "log_time", time.Now().Format(itrlog.LogTimeFormat))
	},
}

func init() {
	rootCmd.AddCommand(copylstCmd)
}