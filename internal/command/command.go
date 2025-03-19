package command

import (
	commandHanler "awesomeProject/internal/command/handler"
	"github.com/spf13/cobra"
)

type Command struct {
	exampleH *commandHanler.ExampleHandler
}

// NewCommand .
func NewCommand(
	exampleHandler *commandHanler.ExampleHandler,
) *Command {
	return &Command{
		exampleH: exampleHandler,
	}
}

func Register(rootCmd *cobra.Command, newCmd func() (*Command, func(), error)) {
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "example",
			Short: "example command",
			Run: func(cmd *cobra.Command, args []string) {
				command, cleanup, err := newCmd()
				if err != nil {
					panic(err)
				}
				defer cleanup()

				command.exampleH.Hello(cmd, args)
			},
		},
	)
}
