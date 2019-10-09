/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-votes-sql/pkg/azvotes"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azvotes",
	Short: "Azure Votes app for sql server",
	Long:  `Starts a web app that demos use of SQL Server.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("dbuser") == "" {
			return fmt.Errorf("no db user provided")
		}
		if viper.GetString("dbpass") == "" {
			return fmt.Errorf("no db password provided")
		}
		if viper.GetString("server") == "" {
			return fmt.Errorf("no server provided")
		}
		if viper.GetString("database") == "" {
			return fmt.Errorf("no db provided")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		user := viper.GetString("dbuser")
		pass := viper.GetString("dbpass")
		database := viper.GetString("database")
		server := viper.GetString("server")
		port := viper.GetInt("port")

		dbClient, err := azvotes.NewClient(server, user, pass, port, database)
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		err = dbClient.Ping(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("db is live")

		err = dbClient.Init()
		if err != nil {
			log.Println(err)
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
