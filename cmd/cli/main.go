package cli

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCMD := &cobra.Command{
		Use:   "go-template",
		Short: "tmp",
	}

	rootCMD.AddCommand(&cobra.Command{
		Use: "hello",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})

	_ = rootCMD.Execute()
}
