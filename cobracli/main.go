package cobracli

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rliebz/tusk/config"
	"github.com/rliebz/tusk/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Metadata is the metadata object configured by the root command.
var Metadata = &config.Metadata{}
var cfgFile string

var verbosity = struct {
	silent  bool
	quiet   bool
	verbose bool
}{}

var rootCmd = &cobra.Command{
	Use:   "tusk",
	Short: "A task runner built with simplicity in mind",
	Long:  "Here is a longer description.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help() // nolint: errcheck,gas
	},
}

// CreateCommand creates a Command.
func CreateCommand() *cobra.Command {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&Metadata.PrintVersion, "version", "V", false, "show the version and exit")
	rootCmd.PersistentFlags().BoolVarP(&verbosity.silent, "silent", "s", false, "print no output")
	rootCmd.PersistentFlags().BoolVarP(&verbosity.quiet, "quiet", "q", false, "print only command output and application errors")
	rootCmd.PersistentFlags().BoolVarP(&verbosity.verbose, "verbose", "v", false, "print verbose output")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "file", "f", "", "config file")

	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version")) // nolint: errcheck,gas
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))   // nolint: errcheck,gas
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))     // nolint: errcheck,gas
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")) // nolint: errcheck,gas

	return rootCmd
}

func initConfig() {
	if err := findConfigFile(Metadata); err != nil {
		ui.Error(err)
		os.Exit(1)
	}

	if Metadata.PrintVersion {
		// TODO: Version
		ui.Println("version")
		os.Exit(0)
	}
}

func findConfigFile(metadata *config.Metadata) error {
	var err error
	if cfgFile != "" {
		metadata.CfgText, err = ioutil.ReadFile(cfgFile)
		if err != nil {
			return err
		}
	} else {
		var found bool
		cfgFile, found, err = config.SearchForFile()
		if err != nil {
			return err
		}

		if found {
			metadata.CfgText, err = ioutil.ReadFile(cfgFile)
			if err != nil {
				return err
			}
		}
	}

	metadata.Directory = filepath.Dir(cfgFile)
	metadata.PrintVersion = viper.GetBool("version")

	if verbosity.silent {
		metadata.Verbosity = ui.VerbosityLevelSilent
	} else if verbosity.quiet {
		metadata.Verbosity = ui.VerbosityLevelQuiet
	} else if verbosity.verbose {
		metadata.Verbosity = ui.VerbosityLevelVerbose
	} else {
		metadata.Verbosity = ui.VerbosityLevelNormal
	}

	return nil
}
