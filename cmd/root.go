/*
Copyright © 2024 achyut koirala <axyut.github.io>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/axyut/cold/internal/app"
	"github.com/axyut/cold/internal/config"
	"github.com/axyut/cold/internal/types"

	"github.com/spf13/cobra"
)

var Version = "dev-build"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cold",
	Short: "A CLI Music Player",
	Long: `A CLI Music Player that plays mp3 files from a directory, defaults to the current directory,
if not found any music files, plays from ~/Music/. It provides a simple interface to play, pause,
skip, and repeat songs.`,
	Version: getVersion(),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// If logging is enabled, logs will be output to debug.log.
		// if enableLogging {
		// 	f, err := tea.LogToFile("debug.log", "debug")
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	defer func() {
		// 		if err = f.Close(); err != nil {
		// 			log.Fatal(err)
		// 		}
		// 	}()
		// }

		setting := getTempSettings(cmd, args)
		config, err := config.Parse(setting)
		if err != nil {
			log.Fatal(err)
		}

		app.StartApp(&config)
	},
	Example: `cold # no commands defaults to config's start directory
cold . # if no audio files, defaults to ~/Music/
cold ~/Music -e a.mp3 -e b.mp3`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(updateCmd)

	rootCmd.PersistentFlags().StringArrayP("exclude", "e", []string{}, "File/s to ignore while playing files in directory")
	rootCmd.PersistentFlags().StringArrayP("include", "i", []string{}, "Include File/s to play with files in directory")
	rootCmd.PersistentFlags().StringArrayP("playonly", "p", []string{}, "Only File/s to play")
	rootCmd.PersistentFlags().StringP("renderer", "r", "", "Application Renderer [raw, tea]")
	rootCmd.PersistentFlags().Bool("icons", true, "Show icons [true/false]")
	rootCmd.PersistentFlags().Bool("hidden", false, "Play Hidden Files [true/false]")
	rootCmd.PersistentFlags().Bool("logging", false, "Enable logging player [true/false]")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getTempSettings(cmd *cobra.Command, args []string) *types.TempSetting {
	startDir := ""
	if len(args) > 0 {
		// fmt.Println("Start Dir: ", args[0])
		startDir = args[0]

	}

	enableLogging, err := cmd.Flags().GetBool("logging")
	if err != nil {
		log.Fatal(err)
	}

	renderer, err := cmd.Flags().GetString("renderer")
	// log.Println("Renderer: ", renderer)
	if err != nil {
		log.Fatal(err)
	}

	showIcons, err := cmd.Flags().GetBool("icons")
	if err != nil {
		log.Fatal(err)
	}
	showHidden, err := cmd.Flags().GetBool("hidden")
	if err != nil {
		log.Fatal(err)
	}

	exclude, err := cmd.Flags().GetStringArray("exclude")
	if err != nil {
		log.Fatal(err)
	}

	include, err := cmd.Flags().GetStringArray("include")
	if err != nil {
		log.Fatal(err)
	}

	playOnly, err := cmd.Flags().GetStringArray("playonly")
	if err != nil {
		log.Fatal(err)
	}
	return &types.TempSetting{
		StartDir:      startDir,
		EnableLogging: enableLogging,
		Renderer:      renderer,
		ShowIcons:     showIcons,
		ShowHidden:    showHidden,
		Exclude:       exclude,
		Include:       include,
		PlayOnly:      playOnly,
	}
}

func getVersion() string {
	var goVer string = runtime.Version()

	return fmt.Sprintf(`%s Built with %s

♪♫♫♫♪  ♫          ♪     ♪    ♪   ♪♫♫♪       ♪♫♫♫♪
♫   ♫  ♫         ♪ ♪     ♪  ♪   ♫          ♫     ♫
♫♫♫♪♪  ♫        ♪   ♪     ♪♪   ♫   ♪♫♫♫♪  ♫       ♫
♫      ♫       ♪♪♪♪♪♪♪    ♪     ♫   ♫  ♫   ♫     ♫
♫      ♫♪♪♪♪♪ ♪       ♪  ♪       ♪♫♫♪  ♪    ♪♫♫♫♪  .`, Version, goVer)
}
