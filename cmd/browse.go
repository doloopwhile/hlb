package cmd

import (
	"fmt"

	"os"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var urlFlag bool

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse repo",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Printf(`Unknown subcommand "%v" for "browse"`, args[0])
			fmt.Println("\nDid you mean this?")

			suggestedCmdNames := cmd.SuggestionsFor(args[0])
			for _, s := range suggestedCmdNames {
				fmt.Println("\t", s)
			}
			os.Exit(1)
		}

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		url, err := base.Client.GetRepositories().GetURL(base.Remote.Owner, base.Remote.RepoName)

		etc.PanicIfErrorExist(err)

		if urlFlag {
			fmt.Println(url)
		} else {
			open.Run(url)
		}
	},
}

func init() {
	browseCmd.PersistentFlags().BoolVarP(&urlFlag, "url", "u", false,
		"outputs the URL rather than opening the browser")
	RootCmd.AddCommand(browseCmd)
}
