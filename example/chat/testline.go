package main
import (
	"fmt"
	"github.com/sbinet/liner"
	"regexp"
	"strings"
	"syscall"
)
/*
* CLI Logic
 */

var SupportedCliCommands = []string{
	"exit",
	"quit",
	"filters",
	"config",
	"deploy",
	"build",
	"help",
	"host",
	"hosts",
	"instance",
	"instances",
	"region",
	"regions",
	"registered",
	"reset",
	"service",
	"services",
	"version",
	"versions",
	"start",
	"stop",
	"restart",
	"register",
	"unregister",
	"log",
	"daemon",
	"daemon",
}

var serviceRegex = regexp.MustCompile("service ([^:]+):")

func tabCompleter(line string) []string {
	opts := make([]string, 0)

	if strings.HasPrefix(line, "reset") {
		filters := []string{
			"reset hosts",
			"reset instance",
			"reset region",
			"reset registered",
			"reset service",
			"reset version",
			"reset config",
		}

		for _, cmd := range filters {
			if strings.HasPrefix(cmd, line) {
				opts = append(opts, cmd)
			}
		}
	//} else if serviceRegex.MatchString(line) {
	//	cmds = make([]string, 0)
	//	matches := serviceRegex.FindAllStringSubmatch(line, -1)
	}

	return opts
}

func main() {
	term := liner.NewLiner()
    defer term.Close()

	fmt.Println("Skynet Interactive Shell")

	term.SetCompleter(tabCompleter)

	for {
		l, e := term.Prompt("> ")
		if e != nil {
			break
		}

		s := string(l)
		parts := strings.Split(s, " ")
        //fmt.Println("receive input:",s)

		validCommand := true

		switch parts[0] {
		case "exit", "quit":
			term.Close()
			syscall.Exit(0)
		case "help", "h":
            InteractiveShellHelp()

		case "services":
		case "regions":
		case "filters":
		default:
			validCommand = false
			fmt.Println("Unknown Command - type 'help' for a list of commands")
		}

		if validCommand {
			term.AppendHistory(s)
		}
	}
}

func confirm(term *liner.State, msg string) bool {
	confirm, _ := term.Prompt(msg + ", Are you sure? (Y/N) > ")
	if confirm == "Y" || confirm == "y" {
		return true
	}

	return false
}

func InteractiveShellHelp() {
	fmt.Print(`
  Commands:
  hosts: List all hosts available that meet the specified criteria
  instances: List all instances available that meet the specified criteria
  regions: List all regions available that meet the specified criteria
  services: List all services available that meet the specified criteria
  versions: List all services available that meet the specified criteria
  config: Set config file for Build/Deploy (defaults to ./build.cfg)
  log: Set log level of services that meet the specified criteria log <level>, options are DEBUG, TRACE, INFO, WARN, FATAL, PANIC
  daemon log: Set log level of daemons that meet the specified criteria daemon log <level>, options are DEBUG, TRACE, INFO, WARN, FATAL, PANIC
  daemon stop: Stop daemons that match the specified criteria

  Filters:
  filters - list current filters
  reset <filter> - reset all filters or specified filter
  region <region> - Add a region to filter, all commands will be scoped to these regions until reset
  service <service> - Add a service to filter, all commands will be scoped to these services until reset
  host <host> - Add host to filter, all commands will be scoped to these hosts until reset
  instance <uuid> - Add an instance to filter, all commands will be scoped to this instance until reset

`)
}

