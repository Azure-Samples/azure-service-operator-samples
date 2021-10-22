/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Azure/azure-votes-postgresql/pkg/azvotes"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azvotes",
	Short: "Azure Votes app for postgresql server",
	Long:  `Starts a web app that demos use of PostgreSQL Server.`,
	Run: func(cmd *cobra.Command, args []string) {

		user := os.Getenv("USERNAME")
		password := os.Getenv("PASSWORD")
		server := os.Getenv("SERVER")
		database := viper.GetString("DATABASE")
		port := viper.GetInt("PORT")

		// Confirm that all required fields are set
		if user == "" {
			log.Fatal("USERNAME environment variable must be set")
		}
		if password == "" {
			log.Fatal("PASSWORD environment variable must be set")
		}
		if server == "" {
			log.Fatal("SERVER environment variable must be set")
		}
		if database == "" {
			log.Fatal("DATABASE environment variable must be set")
		}
		if port == 0 {
			log.Fatal("PORT environment variable must be set")
		}

		dbClient, err := azvotes.NewClient(server, user, password, port, database)
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()
		err = dbClient.Ping(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("db is live")

		err = dbClient.Init(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("db is initialized")

		app := azvotes.NewVoteServer(dbClient)
		log.Fatal(app.Start())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	viper.AutomaticEnv()
	rootCmd.PersistentFlags().StringP("dbuser", "u", "", "database username")
	rootCmd.PersistentFlags().StringP("dbpass", "p", "", "database password")
	rootCmd.PersistentFlags().StringP("server", "s", "", "server name")
	rootCmd.PersistentFlags().StringP("database", "d", "", "database name")
	rootCmd.PersistentFlags().IntP("port", "t", 0, "server port")

	viper.BindPFlag("dbuser", rootCmd.PersistentFlags().Lookup("dbuser"))
	viper.BindPFlag("dbpass", rootCmd.PersistentFlags().Lookup("dbpass"))
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("port", 1433)

}
