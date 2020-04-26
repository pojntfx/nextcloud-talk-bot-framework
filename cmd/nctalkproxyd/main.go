package main

import (
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
	configFileDefault = "/etc/" + keyPrefix + "yaml"
	configFileKey     = keyPrefix + "configFile"
	addrLocalKey      = keyPrefix + "addrLocal"
	addrRemoteKey     = keyPrefix + "addrRemote"
	usernameKey       = keyPrefix + "username"
	passwordKey       = keyPrefix + "password"
	dbpathKey         = keyPrefix + "dbpath"
)

var rootCmd = &cobra.Command{
	Use:   "nctalkproxyd",
	Short: "nctalkproxyd is a Nextcloud Talk API gRPC proxy daemon.",
	Long: `nctalkproxyd is a Nextcloud Talk API gRPC proxy daemon.

Find more information at:
https://pojntfx.github.io/nextcloud-talk-bot-framework/`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("nctalkproxyd")
		// convert Environment parameter names with camel case syntax
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !(viper.GetString(configFileKey) == configFileDefault) {
			viper.SetConfigFile(viper.GetString(configFileKey))
			log.Info("nctalkproxyd: using custom config file", rz.String("configFile", viper.GetString(configFileKey)))
			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		} else {
			log.Info("nctalkproxyd: using default config file", rz.String("configFile:", configFileDefault))
			viper.SetConfigFile(configFileDefault)
			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		}


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
					rz.String("user",  viper.GetString(usernameKey)))

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

func init() {
	var (
		configFileFlag string
		addrLocalFlag  string
		addrRemoteFlag string
		usernameFlag   string
		passwordFlag   string
		dbpathFlag     string
	)

	rootCmd.PersistentFlags().StringVarP(&configFileFlag, configFileKey, "c", configFileDefault, cmd.NcTalkProxyConfigurationFile)
	rootCmd.PersistentFlags().StringVarP(&addrLocalFlag, addrLocalKey, "l", cmd.NcTalkProxydDefaultAddrLocal, "NcTalkProxyd address.")
	rootCmd.PersistentFlags().StringVarP(&addrRemoteFlag, addrRemoteKey, "r", "https://mynextcloud.com", "Nextcloud bot address.")
	rootCmd.PersistentFlags().StringVarP(&usernameFlag, usernameKey, "u", "botusername", "Nextcloud bot account username.")
	rootCmd.PersistentFlags().StringVarP(&passwordFlag, passwordKey, "p", "botpassword", "Nextcloud bot account password.")
	rootCmd.PersistentFlags().StringVarP(&dbpathFlag, dbpathKey, "d", "/var/lib/nctalkproxyd", "Database path.")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(cmd.CouldNotBindFlagsErrorMessage, rz.Err(err))
	}

	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(cmd.CouldNotStartRootCommandErrorMessage, rz.Err(err))
	}
}
