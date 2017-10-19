// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.
package main

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/spf13/cobra"
)

var messageExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data from Mattermost in a format suitable for import into a third-party application",
}

var messageExportActianceCmd = &cobra.Command{
	Use:     "actiance",
	Short:   "Export Actiance XML",
	Long:    "Export data in the Actiance XML format.",
	Example: "  export actiance",
	RunE:    messageExportActianceCmdF,
}

func init() {
	messageExportCmd.AddCommand(
		messageExportActianceCmd,
	)
}

func messageExportActianceCmdF(cmd *cobra.Command, args []string) error {
	a, err := initDBCommandContextCobra(cmd)
	if err != nil {
		return err
	}

	if !*a.Config().MessageExportSettings.EnableExport {
		CommandPrintErrorln("ERROR: The message export feature is not enabled.")
		return nil
	}

	// TODO: check licensing? This code reports that the server is not licensed, when it is
	/*
		if !utils.IsLicensed() || !*utils.License().Features.MessageExport {
			fmt.Println("IsLicensed: " + strconv.FormatBool(utils.IsLicensed()))
			fmt.Println("License Features: " + utils.License().ToJson())
			CommandPrintErrorln("ERROR: The server is not licensed to use the message export feature.")
			return nil
		}
	*/

	if messageExportI := a.MessageExport; messageExportI != nil {
		job, err := messageExportI.StartSynchronizeJob(true)
		if err != nil || job.Status == model.JOB_STATUS_ERROR || job.Status == model.JOB_STATUS_CANCELED {
			CommandPrintErrorln("ERROR: Message export job failed. Please check the server logs")
		} else {
			CommandPrettyPrintln("SUCCESS: Message export job complete")
		}
	}

	return nil
}
