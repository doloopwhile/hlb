package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/hlblib"
	"github.com/mpppk/hlb/etc"
	"github.com/skratchdot/open-golang/open"
	"fmt"
	"strconv"
)

// browseissuesCmd represents the browseissues command
var browseissuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "browse issues",
	Long: `browse issues`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many issue IDs")
		}

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)

		var url string
		if len(args) == 0 {
			u, err := base.Client.GetIssues().GetIssuesURL(base.Remote.Owner, base.Remote.RepoName)
			etc.PanicIfErrorExist(err)
			url = u
		}else {
			id, err := strconv.Atoi(args[0])
			etc.PanicIfErrorExist(err)

			u, err := base.Client.GetIssues().GetURL(base.Remote.Owner, base.Remote.RepoName, id)
			etc.PanicIfErrorExist(err)
			url = u
		}

		if urlFlag {
			fmt.Println(url)
		} else {
			open.Run(url)
		}
	},
}

func init() {
	browseCmd.AddCommand(browseissuesCmd)
}
