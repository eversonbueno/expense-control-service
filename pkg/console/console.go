package console

import (
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var commands = map[string]*cobra.Command{}

func PrintInfo(format string, a ...interface{}) {
	color.Cyan(format, a...)
}

func PrintSuccess(format string, a ...interface{}) {
	color.Green(format, a...)
}

func PrintError(format string, a ...interface{}) {
	color.Red(format, a...)
}

func AddCommand(name, description string, runFunc func(cmd *cobra.Command, args []string)) {
	commands[name] = &cobra.Command{
		Use:   name,
		Short: description,
		Run:   runFunc,
	}
}

func Help(cmd *cobra.Command, args []string) {
	PrintSuccess("Comandos disponíveis:")
	for name, command := range commands {
		if name == "help" {
			continue
		}
		PrintSuccess("- %s: %s\n", name, command.Short)
	}
}

func GetCommand(command string) (*cobra.Command, bool) {
	cmd, exists := commands[command]
	return cmd, exists
}

func ParseOptions(args []string) map[string]string {
	options := make(map[string]string)
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
			key := parts[0]
			value := ""
			if len(parts) > 1 {
				value = parts[1]
			} else {
				value = "true"
			}
			options[key] = value
		}
	}
	return options
}