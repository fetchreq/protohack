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
	"time"
	"unicode"

	"github.com/spf13/cobra"
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


		channel := make(chan string)

		go chatroom(myslog)

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}

			go chatRoomHandler(conn, myslog, channel)

		}

	},
}

func init() {
	rootCmd.AddCommand(p3Cmd)
}

const helloMessage = "Welcome to budgetchat! What shall I call you?"
type ClientStatus  int

const (
	Joining ClientStatus = iota
	Joined 
)

type EventType  int

const (
	NewMessage EventType = iota
	NewUserBroadCast
	UserList
	UserLeave
)

var (
	// Send a channel here to get room events back.  It will send the entire
	// archive initially, and then new messages as they come in.
	subscribe = make(chan Subscription)
	// Send a channel here to unsubscribe.
	unsubscribe = make(chan Subscription)
	// Send events here to publish them.
	publish = make(chan Event)
)
func chatroom (logger *slog.Logger) {
	cliensts := make(map[Subscription]bool)
	for {
		select {
		case sub := <- subscribe: 
			logger.Info(fmt.Sprintf("[new user] %s", sub.name))	
			var nameList strings.Builder
			for user := range cliensts {
				nameList.WriteString(fmt.Sprintf("%s, ", user.name))
			}

			names := nameList.String()
			names = strings.TrimSuffix(names, ", ")
			
			sub.channel <- Event{eventType: UserList, message: names}
			for user := range cliensts {

				user.channel <- Event{eventType: NewUserBroadCast, message: sub.name}
			}
			cliensts[sub] = true

		case event := <- publish:
			for user := range cliensts {
				if event.userId != user.id {
					user.channel <- Event{eventType: NewMessage, message: event.message, name: event.name}
				}
			}
		case unsub := <- unsubscribe: 
			delete(cliensts, unsub)	
			logger.Info(fmt.Sprintf("[user left] %s", unsub.name))	
			for user := range cliensts {

				user.channel <- Event{eventType: UserLeave, message: unsub.name}
			}

		}
			
		
	}
}

type Subscription struct {
	id string
	name string
	channel chan Event 
}

type Event struct {
	userId string
	name string
	eventType EventType
	message string
}

func chatRoomHandler(conn net.Conn, logger *slog.Logger, channel chan string) {
	var user Subscription
	defer conn.Close()
	msg := []byte(helloMessage)
	msg = append(msg, '\n')
	conn.Write(msg)
	
	c := make(chan Event)
	scanner := bufio.NewScanner(conn)
	status := Joining
	for scanner.Scan() {
		buffer := scanner.Bytes()
		
		msg := string(buffer)
		if status == Joining {
			allowed := true
			for _, letter := range msg {
				if !unicode.IsLetter(letter) && !unicode.IsDigit(letter) && letter != '\n' {
					allowed = false	
				}
			}
			if len(strings.TrimSpace(msg)) == 0 || len(msg) < 1 || !allowed {
				break
			}

			user = Subscription{id: fmt.Sprintf("%s_%s_%s", time.Nanosecond.String(), conn.RemoteAddr(), conn.RemoteAddr().Network()), name: msg, channel: c }
			go messageHandler(c, conn)
			subscribe <- user
			status = Joined
		} else if status == Joined {
			publish <- Event{userId: user.id, name: user.name, message: msg}

		}
		
	}
	
	if status == Joined {
		unsubscribe <- user
	}

}
const (
	UserListTemplate = "* This room contains %s"
	UserEnterTemplate = "* %s has entered the room"
	UserLeaveTemplate = "* %s has left the room"
	NewMessageTemplate = "[%s] %s"
)

func messageHandler(c chan Event, conn net.Conn) {
	for {
		select {
		case event := <- c:
			var message string
			if event.eventType == UserList {
				message = fmt.Sprintf(UserListTemplate, event.message)

			} else if event.eventType == NewUserBroadCast {
				message = fmt.Sprintf(UserEnterTemplate, event.message)

			} else if event.eventType == UserLeave {
				message = fmt.Sprintf(UserLeaveTemplate, event.message)

			} else if event.eventType == NewMessage {
				message = fmt.Sprintf(NewMessageTemplate, event.name, event.message)

			}

			output := []byte(message)
			output = append(output, '\n')
			conn.Write(output)

		}
	}
}
