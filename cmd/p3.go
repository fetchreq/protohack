/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

const (
	InvalidName = "Only alphanumeric characters are allowed in names"
	WelcomeMessageTemplate = "Welcome to budgetchat! What shall I call you?"
	UserListTemplate = "* This room contains %s"
	UserEnterTemplate = "* %s has entered the room"
	UserLeaveTemplate = "* %s has left the room"
	NewMessageTemplate = "[%s] %s"
)

type ClientStatus int

const (
	Joining ClientStatus = iota
	Joined 
)

var (
	// Send a new Subscription to add some one to the chatroom
	subscribe = make(chan Subscription)
	// Send a channel here to unsubscribe.
	unsubscribe = make(chan Subscription)
	// Send events here to publish them.
	publish = make(chan ClientEvent)
)

// p3Cmd represents the p3 command
var p3Cmd = &cobra.Command{
	Use:   "p3",
	Short: "Budget Chat",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Set up logger	
		jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		myslog := slog.New(jsonHandler)
		myslog.Info("Starting Budget Chat")

		listener, err := net.Listen("tcp", "0.0.0.0:10000")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer listener.Close()


		// start the chat room
		go chatroomController(myslog)

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
	rootCmd.AddCommand(p3Cmd)
}



// The main controller for the chatroom
func chatroomController(logger *slog.Logger) {
	cliensts := make(map[Subscription]bool)
	for {
		select {
		case sub := <- subscribe: 
			logger.Info(fmt.Sprintf("[new user] %s", sub.Name))	

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
			logger.Info(fmt.Sprintf("[user left] %s", unsub.Name))	

			// Announce a user has left
			userLeftMsg := fmt.Sprintf(UserLeaveTemplate, unsub.Name)
			for user := range cliensts {
				user.Channel <- userLeftMsg
			}

		}
			
		
	}
}

// Gets a comma seperated list of users who are in the chat room
func GetListOfActiveNames(cliensts map[Subscription]bool) string {
	var nameList strings.Builder
	for user := range cliensts {
		nameList.WriteString(fmt.Sprintf("%s, ", user.Name))
	}

	return strings.TrimSuffix(nameList.String(), ", ")
}

type Subscription struct {
	ID uuid.UUID
	Name string
	Channel chan string
}

type ClientEvent struct {
	UserId uuid.UUID
	FromUser string
	Message string
}


func ChatRoomConnection(conn net.Conn, logger *slog.Logger) {
	var user Subscription
	defer conn.Close()
	
	userChannel := make(chan string)
	scanner := bufio.NewScanner(conn)
	status := Joining

	for scanner.Scan() {
		buffer := scanner.Bytes()
		
		msg := string(buffer)
		if status == Joining {
			if IsNameNotAllowed(msg) {
				output := []byte(InvalidName)
				output = append(output, '\n')
				conn.Write(output)
				break
			}

			id := uuid.New()
			user = Subscription{ID: id, Name: msg, Channel: userChannel }

			// Add new user to chat room
			subscribe <- user

			// Start listening for events from the server 
			go UserEventPublisher(userChannel, conn)

			status = Joined

		} else if status == Joined {
			publish <- ClientEvent{UserId: user.ID, FromUser: user.Name, Message: msg}

		}
		
	}
	
	if status == Joined {
		// Publish a user has left only if they finished joining
		unsubscribe <- user
	}

}

// Checks if a name is not allowed, allowed names are only able to have alphanumeric characters
func IsNameNotAllowed(name string) bool {
	disallowed := false 
	for _, letter := range name {
		if !unicode.IsLetter(letter) && !unicode.IsDigit(letter) && letter != '\n' {
			disallowed = true
		}
	}
	return len(strings.TrimSpace(name)) == 0 || len(name) < 1 || disallowed
}

// Listens to events from the controller to be published to a user
func UserEventPublisher(channel chan string, conn net.Conn) {
	for {
		message := <- channel
		output := []byte(message)
		output = append(output, '\n')
		conn.Write(output)

	}
}
