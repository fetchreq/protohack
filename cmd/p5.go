/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/spf13/cobra"
)




// p3Cmd represents the p3 command
var p5Cmd = &cobra.Command{
	Use:   "p5",
	Short: "Mob in the Middle",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Set up logger	
		jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		myslog := slog.New(jsonHandler)
		myslog.Info("Starting Mob in the Middle")

		listener, err := net.Listen("tcp", "0.0.0.0:10000")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer listener.Close()


		// start the chat room
		go chatroomControllerMitM(myslog)

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}

			// Send welcome message
			msg := []byte(WelcomeMessageTemplate)
			msg = append(msg, '\n')
			conn.Write(msg)

			// Add user to the the chat room 
			go ChatRoomConnection(conn, myslog)

		}

	},
}

func init() {
	rootCmd.AddCommand(p5Cmd)
}


type SubscriptionEvil struct {
	ID uuid.UUID
	Name string
	Channel chan string
}

// The main controller for the chatroom
func chatroomControllerMitM(logger *slog.Logger) {
	cliensts := make(map[Subscription]bool)
	for {
		select {
		case sub := <- subscribe: 
			logger.Info(fmt.Sprintf("[new user mitm] %s", sub.Name))	

			conn, err := net.Dial("tcp", "2a03:b0c0:1:d0::116a:8001")
			if err != nil {
				logger.Error(fmt.Sprintf("Unable to connect to downstream for user: %s", sub.Name))
			}

			// Send an event to the new user to show the names of people in the chat
			userListMsg := fmt.Sprintf(UserListTemplate, GetListOfActiveNames(cliensts))
			sub.Channel <- userListMsg

			// Sned an event to annouse the new user to everyone else
			newUserMsg := fmt.Sprintf(UserEnterTemplate, sub.Name)
			for user := range cliensts {
				user.Channel <- newUserMsg
			}
			// Add the new user
			cliensts[sub] = true

		case event := <- publish:

			userMessage := fmt.Sprintf(NewMessageTemplate, event.FromUser, event.Message)
	
			// Sends a message to all other users in the chat room
			for user := range cliensts {

				if event.UserId != user.ID {
					user.Channel <- userMessage
				}
			}
		case unsub := <- unsubscribe: 
			delete(cliensts, unsub)	
			logger.Info(fmt.Sprintf("[user left mitm] %s", unsub.Name))	

			// Announce a user has left
			userLeftMsg := fmt.Sprintf(UserLeaveTemplate, unsub.Name)
			for user := range cliensts {
				user.Channel <- userLeftMsg
			}

		}
			
		
	}
}

