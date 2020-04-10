package main

import (
	"net"
	"strings"
	"time"

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
	keyPrefix         = "nxtalkproxyd."
	configFileDefault = ""
	configFileKey     = keyPrefix + "configFile"
	laddrKey          = keyPrefix + "laddr"
	raddrKey          = keyPrefix + "raddr"
	usernameKey       = keyPrefix + "username"
	passwordKey       = keyPrefix + "password"
	dbpathKey         = keyPrefix + "dbpath"
)

var rootCmd = &cobra.Command{
	Use:   "nxtalkproxyd",
	Short: "nxtalkproxyd is a Nextcloud Talk API gRPC proxy daemon.",
	Long: `nxtalkproxyd is a Nextcloud Talk API gRPC proxy daemon.

Find more information at:
https://pojntfx.github.io/nextcloud-talk-bot-framework/`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("nxtalkproxyd")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !(viper.GetString(configFileKey) == configFileDefault) {
			viper.SetConfigFile(viper.GetString(configFileKey))

			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		}

		listener, err := net.Listen("tcp", viper.GetString(laddrKey))
		if err != nil {
			return err
		}

		server := grpc.NewServer()
		reflection.Register(server)

		chatChan := make(chan clients.Chat)
		chatRequestChan := make(chan bool)
		chatChans := []chan clients.Chat{}
		chatResponseChan := make(chan chan clients.Chat)
		writeChan := func(token, message string) error {
			log.Info("writing chat from client to Nextcloud Talk", rz.String("token", token), rz.String("message", message))

			return nil
		}

		go func() {
			for {
				chatChan <- clients.Chat{
					ID:               1,
					Token:            "testToken",
					ActorID:          "testActorID",
					ActorDisplayName: "testDisplayName",
					Message:          "testMessage",
				}

				time.Sleep(time.Second * 2)
			}
		}()

		go func() {
			for range chatRequestChan {
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

		nextcloudTalkService := services.NewNextcloudTalk(chatRequestChan, chatResponseChan, writeChan)

		nextcloudTalk.RegisterNextcloudTalkServer(server, nextcloudTalkService)

		log.Info("starting server")

		return server.Serve(listener)
	},
}

func init() {
	var (
		configFileFlag string
	)

	rootCmd.PersistentFlags().StringVarP(&configFileFlag, configFileKey, "f", configFileDefault, cmd.ConfigurationFileDocs)

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
