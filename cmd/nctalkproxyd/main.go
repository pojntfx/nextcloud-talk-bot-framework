package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/pojntfx/nextcloud-talk-bot-framework/cmd"
	"github.com/pojntfx/nextcloud-talk-bot-framework/pkg/clients"
	nextcloudTalk "github.com/pojntfx/nextcloud-talk-bot-framework/pkg/protos/generated"
	"github.com/pojntfx/nextcloud-talk-bot-framework/pkg/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/bloom42/libs/rz-go"
	"gitlab.com/bloom42/libs/rz-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	keyPrefix         = "nctalkproxyd."
	configFileDefault = "nctalkproxyd" // viper will resolve supported format extension
	configFileKey     = keyPrefix + "configFile"
	addrLocalKey      = keyPrefix + "addrLocal"
	addrRemoteKey     = keyPrefix + "addrRemote"
	usernameKey       = keyPrefix + "username"
	passwordKey       = keyPrefix + "password"
	dbpathKey         = keyPrefix + "dbpath"
)

var (
	// commandline flags
	configFile     string
	addrLocalFlag  string
	addrRemoteFlag string
	usernameFlag   string
	passwordFlag   string
	dbpathFlag     string

	rootCmd = &cobra.Command{
		Use:   "nctalkproxyd",
		Short: "nctalkproxyd is a Nextcloud Talk API gRPC proxy daemon.",
		Long: `nctalkproxyd is a Nextcloud Talk API gRPC proxy daemon.

Find more information at:
https://pojntfx.github.io/nextcloud-talk-bot-framework/`,
		Version: "0.2",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// convert Environment parameter names with camel case syntax
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
		},

		// our main function, running the proxy routines
		RunE: func(cmd *cobra.Command, args []string) error {
			listener, err := net.Listen("tcp", viper.GetString(addrLocalKey))
			if err != nil {
				return err
			} else {
				log.Info("nctalkproxyd: listener established", rz.String("addrLocal", viper.GetString(addrLocalKey)))
			}

			server := grpc.NewServer()
			reflection.Register(server)

			chatChan := make(chan clients.Chat)
			chatRequestChan := make(chan bool)
			chatChans := []chan clients.Chat{}
			chatResponseChan := make(chan chan clients.Chat)
			statusChan, svcStatusChan := make(chan string), make(chan string)

			nextcloudTalkClient := clients.NewNextcloudTalk(
				viper.GetString(addrRemoteKey),
				viper.GetString(usernameKey),
				viper.GetString(passwordKey),
				viper.GetString(dbpathKey),
				chatChan,
				statusChan,
			)
			log.Info("nctalkproxyd: Bot connection to NextcloudTalk configured",
				rz.String("addrRemote", viper.GetString(addrRemoteKey)),
				rz.String("user", viper.GetString(usernameKey)))

			writeChan := func(token, message string) error {
				log.Info("writing chat from client to Nextcloud Talk",
					rz.String("token", token), rz.String("message", message))

				return nextcloudTalkClient.WriteChat(token, message)
			}

			defer nextcloudTalkClient.Close()
			if err := nextcloudTalkClient.Open(); err != nil {
				log.Fatal("could not open Nextcloud Talk client", rz.Err(err))
			}

			go func() {
				for {
					if err := nextcloudTalkClient.ReadRooms(); err != nil {
						log.Info("could not read rooms, retrying", rz.Err(err))
					}
				}
			}()

			go func() {
				for {
					if err := nextcloudTalkClient.ReadChats(); err != nil {
						log.Info("could not read chats, retrying", rz.Err(err))
					}
				}
			}()

			go func() {
				for status := range statusChan {
					log.Info("received Nextcloud client status", rz.String("status", status))
				}
			}()

			go func() {
				for status := range svcStatusChan {
					log.Info("received service status", rz.String("status", status))
				}
			}()

			go func() {
				for range chatRequestChan {
					log.Info("new client connected to service")

					chatChan := make(chan clients.Chat)

					chatChans = append(chatChans, chatChan)

					chatResponseChan <- chatChan
				}
			}()

			go func() {
				for chat := range chatChan {
					log.Info("writing chat from Nextcloud Talk to clients", rz.Any("chat", chat))

					for _, chatChan := range chatChans {
						chatChan <- chat
					}
				}
			}()

			nextcloudTalkService := services.NewNextcloudTalk(chatRequestChan, chatResponseChan, svcStatusChan, writeChan)

			nextcloudTalk.RegisterNextcloudTalkServer(server, nextcloudTalkService)

			log.Info("starting server")

			return server.Serve(listener)
		},
	}
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "nctalkproxyd is a Nextcloud Talk API gRPC proxy daemon.",
	Long:  `nctalkproxyd is implementet in the go language`,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, configFileKey, "c", configFileDefault, cmd.NcTalkProxyConfigurationFile)
	rootCmd.PersistentFlags().StringVarP(&addrLocalFlag, addrLocalKey, "l", cmd.NcTalkProxydDefaultAddrLocal, "NcTalkProxyd socket.")
	rootCmd.PersistentFlags().StringVarP(&addrRemoteFlag, addrRemoteKey, "r", "https://mynetxcloud.com", "Nextcloud bot URL.")
	rootCmd.PersistentFlags().StringVarP(&usernameFlag, usernameKey, "u", "botusername", "Nextcloud bot account username.")
	rootCmd.PersistentFlags().StringVarP(&passwordFlag, passwordKey, "p", "botpassword", "Nextcloud bot account password.")
	rootCmd.PersistentFlags().StringVarP(&dbpathFlag, dbpathKey, "d", "/var/lib/nctalkproxyd", "Database path.")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "Felix Pojntinger")
	viper.SetDefault("license", "AGPLv3")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(cmd.CouldNotBindFlagsErrorMessage, rz.Err(err))
	}

	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	// Search given config file in standard directoies
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("/etc/nctalkproxyd/")
	viper.AddConfigPath(".")

	// handle config file
	if len(viper.GetString(configFileKey)) > 0 {
		// Use config file from the flag.
		fmt.Println("Using config file:", configFile)
		viper.SetConfigName(configFile)

		// Searches for config file in given paths and read it
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				if !(viper.GetString(configFileKey) == configFileDefault) {
					log.Fatal("Failed to read given configuration file!", rz.Err(err))
				} else {
					log.Info("nctalkproxyd: default config file not found", rz.Err(err))
					log.Info("nctalkproxyd: using build-in defaults")
				}
			} else {
				// Config file was found but another error was produced
				log.Info("nctalkproxyd: error reading config file", rz.String("configFile", viper.GetString(configFile)))
			}
		}
	} else {
		log.Info("nctalkproxyd: no config file selected, using buildin flags ...")
	}

	// handle environment variables: if set, they take precedence
	viper.AutomaticEnv()

	// array with valid environment variables
	env_vars := []string{
		"nctalkproxyd_addrRemote",
		"nctalkproxyd_addrLocal",
		"nctalkproxyd_username",
		"nctalkproxyd_password",
		"nctalkproxyd_dbpath",
	}

	for i := 0; i < len(env_vars); i++ {
		if viper.Get(env_vars[i]) != nil {
			if env_vars[i] != "nctalkproxyd_password" {
				log.Debug("nctalkproxyd: environment value", rz.String(env_vars[i], viper.GetString(env_vars[i])))
			} else {
				log.Debug("nctalkproxyd: environment value", rz.String(env_vars[i], "*******"))
			}
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(cmd.CouldNotStartRootCommandErrorMessage, rz.Err(err))
	}
}
