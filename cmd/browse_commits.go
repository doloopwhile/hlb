package cmd

import (
	"fmt"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var browsecommitsCmd = &cobra.Command{
	Use:   "commits",
	Short: "browse commits",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("warning: `browse commits` does not accept any args. They are ignored.")
		}

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		url, err := sw.GetCommitsURL()
		etc.PanicIfErrorExist(err)
		open.Run(url)
	},
}

func init() {
	browseCmd.AddCommand(browsecommitsCmd)
}