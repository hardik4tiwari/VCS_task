package cmd

import (
	"Backup/config"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	// "github.com/fatih/color"
	// "github.com/itrepablik/kopy"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/itrepablik/itrlog"
	// "go.uber.org/zap"

	homedir "github.com/mitchellh/go-homedir"
)

// DiskSN is to store the user's HD Serial No at first runtime.
var DiskSN string

// IsLicenseKeyValid default to 'true'
var IsLicenseKeyValid bool = true // Set license key valid

// IsUserNameEntered default to 'false'
var IsUserNameEntered bool = false

// IsPasswordEntered default to 'false'
var IsPasswordEntered bool = false

// IsLogCopiedFile default to 'true'
var IsLogCopiedFile bool = true

// IsAutoBK default to 'false'
var IsAutoBK bool = false

// IgnoreFileTypes collects the list of files to be ignored during the backup operation.
var IgnoreFileTypes interface{}

// IgnoreFT is a []strings type that collects the list of files to be ignored during the backup operation.
var IgnoreFT []string

// STCopyDIRF is the copy dir struct for data collection from the 'config.yaml' file.
type STCopyDIRF struct {
	src, dst      string
	runEvery      int
	interval      string
	retentionDays int
	intervalType  string
}

// STCopyDIRD is the copy dir struct for data collection from the 'config.yaml' file.
type STCopyDIRD struct {
	src, dst      string
	runEvery      int
	interval      string
	runAt         string
	intervalType  string
	retentionDays int
}

// STCopyMD is the copymd struct for data collection from the 'config.yaml' file.
type STCopyMD struct {
	src, dst       string
	runEvery       int
	interval       string
	runAt          string
	intervalType   string
	copyModNumDays string
}

// STCopyMDF is the copymd struct for data collection from the 'config.yaml' file.
type STCopyMDF struct {
	src, dst       string
	runEvery       int
	interval       string
	intervalType   string
	copyModNumDays string
}

// MapCopyDIRD is to store copydir_daily automated backup.
var MapCopyDIRD = make(map[int]STCopyDIRD)

// CURCopyDIRD is to store copydir_daily automated backup.
var CURCopyDIRD STCopyDIRD

// MapCopyDIRF is to store copydir_frequently automated backup.
var MapCopyDIRF = make(map[int]STCopyDIRF)

// CURCopyDIRF is to store copydir_frequently automated backup.
var CURCopyDIRF STCopyDIRF

// MapCopyMD is to store copymd automated backup.
var MapCopyMD = make(map[int]STCopyMD)

// CURCopyMD is to store copymd automated backup.
var CURCopyMD STCopyMD

// MapCopyMDF is to store copymd_frequently automated backup.
var MapCopyMDF = make(map[int]STCopyMDF)

// CURCopyMDF is to store copymd_frequently automated backup.
var CURCopyMDF STCopyMDF

// IsBKCompressed common automated backup variables
var IsBKCompressed bool = false

// BKSD array of strings for copydir collections.
var BKSD []string

// BKSF array of strings for copydir collections.
var BKSF []string

// BKMD array of strings for copymd collections.
var BKMD []string

// BKMDF array of strings for copymd collections.
var BKMDF []string

// MDays is the common variables to store number of days param to execute the copymd command.
var MDays int = 0

// IsBKItemsFound is to check if any backup items entered by the user from the 'config.yaml' file.
var IsBKItemsFound bool = false

// NumFilesCopied counts the number of files copied.
var NumFilesCopied int = 0

// NumFoldersCopied counts the number of folders copied.
var NumFoldersCopied int = 0

// MaxLogFileSizeInMB gets the max log file size value in megabytes.
var MaxLogFileSizeInMB int = 100 // mb

// MaxAgeLogInDays get the max age of a log files in days.
var MaxAgeLogInDays int = 0 // 0 days means, it won't delete older backup logs

// Sugar type is the *itrlog.ITRLogger initialization
var Sugar *itrlog.ITRLogger

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   config.AppName,
	Short: config.AppDesc,
	Long: config.AppDesc,
	Version: config.Appversion,
	Run: func(cmd *cobra.Command, args []string) {
	  // Do Stuff Here
	},
  }
  
  func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Fprintln(os.Stderr, err)
	  os.Exit(1)
	}
  }

 func init() {
	myFigure:= figure.NewFigure(config.AppName,"",true)
	myFigure.Print()
 }

func init() {
	cobra.OnInitialize(initConfig)
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	// Handle errors reading the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create the "config.yaml" asap.
			f, err := os.OpenFile("config.yaml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer f.Close()
		} else {
			// Config file was found but another error was produced
			log.Fatalf("fatal error config file: %v", err)
		}
	}

	// Set the default values, these will be fetch even though not found in "config.yaml" file.
	viper.SetDefault("license", "")                         // Set to blank value for the license
	viper.SetDefault("default.copy_mod_files_num_days", -1) // Set it to 1 day
	viper.SetDefault("logging.log_copied_file", true)       // Set to true to log each single file copied.
	viper.SetDefault("ignore.file_types", ".thumb, .db")    // Set the default common ignored file types.
	viper.SetDefault("ignore.folders", "")                  // Set the default ignored folder here, leave it blank.

	// Get the default value for the "max_log_file_size_in_mb" setting.
	maxLogFileSize := viper.Get("logging.max_log_file_size_in_mb")
	// MaxLogFileSizeInMB = maxLogFileSize.(int)
	// if _, ok := maxLogFileSize.(int); !ok {
	// 	MaxLogFileSizeInMB = 100 // default: mb
	// }
	if maxLogFileSize != nil {
		if value, ok := maxLogFileSize.(int); ok {
			MaxLogFileSizeInMB = value
		} else {
			// Handle the case where the value is not an integer
			fmt.Println("Value is not an integer")
		}
	} else {
		// Handle the case where the key is not found in the configuration
		fmt.Println("Key not found in configuration")
	}

	// Get the default value for the "max_age_in_days" setting.
	maxLogAge := viper.Get("logging.max_age_in_days")
	MaxAgeLogInDays = maxLogAge.(int)
	if _, ok := maxLogAge.(int); !ok {
		MaxAgeLogInDays = 0 // default: days
	}

	// Zap / Lamberjack Logger initialization
	Sugar = itrlog.InitLog(MaxLogFileSizeInMB, MaxAgeLogInDays, "logs", "gokopy_log_")

	// Check if need to log each copied file.
	isNeedToLogCopiedFile := viper.Get("logging.log_copied_file")
	bc, err := strconv.ParseBool(fmt.Sprintf("%v", isNeedToLogCopiedFile))
	if err != nil {
		bc = false
		Sugar.Sugar.Errorw("", "error:", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
	}

	IsLogCopiedFile = bc

	// Get the default value for the "copy_mod_files_num_days" setting.
	modDays := viper.Get("default.copy_mod_files_num_days")
	MDays = modDays.(int)
	if _, ok := modDays.(int); !ok {
		MDays = -1
	}

	viper.WatchConfig() // Tell the viper to watch any new changes to the config file.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gokopy" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gokopy")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}