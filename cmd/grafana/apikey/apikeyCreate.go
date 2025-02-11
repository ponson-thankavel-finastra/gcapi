package apikey

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jtyr/gcapi/cmd/common"
)

// NewCmdCreate returns a new cobra command.
func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create (STACK_SLUG|--grafana-api-url STRING) NAME ROLE",
		Aliases: []string{"add"},
		Short:   "Create API key",
		Long:    "Create Grafana API key and print it out.",
		Args:    checkCreateArgs,
		Run:     runCreate,
	}

	cmd.Flags().BoolP("raw", "r", false, "show raw API response")
	cmd.Flags().StringP("grafana-api-url", "g", "", "Grafana API URL (e.g. https://grafana.domain.com/api).")
	cmd.Flags().Uint64P("seconds-to-live", "s", 0, "lifespan of the API key in seconds")

	return cmd
}

// checkCreateArgs checks if the positional arguments have correct
// value. If no args are specified, it prints out the command usage.
func checkCreateArgs(cmd *cobra.Command, args []string) error {
	gauFlag, err := common.GetGrafanaAPIURL(cmd)
	if err != nil {
		log.Fatalf("failed to get Grafana API URL: %s", err)
	}

	argsLen := len(args)

	if argsLen == 0 {
		cmd.Usage()
		os.Exit(0)
	} else if gauFlag != "" {
		if argsLen < 2 {
			return errors.New("requires NAME and ROLE argument")
		}

		if err := ak.SetName(args[0]); err != nil {
			return err
		}

		if err := ak.SetRole(args[1]); err != nil {
			return err
		}

		if err := ak.SetBaseURL(gauFlag); err != nil {
			return err
		}

		if token, err := common.GetGrafanaToken(cmd); err == nil {
			ak.SetGrafanaToken(token)
		} else {
			return fmt.Errorf("failed to get Grafana authorization token: %s", err)
		}
	} else if argsLen < 3 {
		return errors.New("requires STACK_SLUG, NAME and ROLE argument")
	} else {
		if err := ak.SetStackSlug(args[0]); err != nil {
			return err
		}

		if err := ak.SetName(args[1]); err != nil {
			return err
		}

		if err := ak.SetRole(args[2]); err != nil {
			return err
		}

		if token, err := common.GetToken(cmd); err == nil {
			ak.SetToken(token)
		} else {
			return fmt.Errorf("failed to get authorization token: %s", err)
		}
	}

	stlFlag, err := cmd.Flags().GetUint64("seconds-to-live")
	if err != nil {
		log.Fatalf("failed to get seconds-to-live flag value: %s", err)
	}

	if stlFlag > 0 {
		if err := ak.SetSecondsToLive(stlFlag); err != nil {
			return err
		}
	}

	return nil
}

// runCreate runs the command's action.
func runCreate(cmd *cobra.Command, args []string) {
	key, raw, err := ak.Create()
	if err != nil {
		log.Fatalf("failed to create API key: %s", err)
	}

	rawFlag, err := cmd.Flags().GetBool("raw")
	if err != nil {
		log.Fatalf("failed to get raw flag value: %s", err)
	}

	if rawFlag {
		fmt.Println(raw)
	} else {
		fmt.Println(key)
	}
}
