package cmd

import (
	"libify/handlers"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a list of borrowed books",
	Long: ` The usage:

    libify show --password 123 --username almas`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials := handlers.LibraryCredentials{Password: password, Username: username}
		handlers.HandleShowBooks(&credentials)
	},
}
var (
	password string
	username string
)

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	showCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Profile password")
	showCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username")
}
