package main

import (
	"os"
	"strings"
)

func checkArgs() {
	args := os.Args

	if len(args) == 1 {
		return
	}

	args = args[1:]

	if args[0] == "create_user" {
		argCreateUser(args)
	}

	if args[0] == "delete_user" {
		argsDeleteUser(args)
	}

	os.Exit(0)
}

func argCreateUser(args []string) {
	if len(args) < 4 {
		Error.Println("to create a user, provide the following in order: usertype, username, password")
	}

	userType := strings.ToUpper(args[1])
	username := args[2]
	password := args[3]

	if userType != string(FACULTY) && userType != string(STAFF) {
		Error.Fatalln("provide a valid user type, faculty or staff")
	}

	authCreateUser(username, password, UserType(userType))
	Info.Printf("successfully created user %s with type %s\n", username, userType)
}

func argsDeleteUser(args []string) {
	if len(args) < 2 {
		Error.Fatalln("to delete a user, provide the following: username")
	}

	username := args[1]
	authDeleteUser(username)
	Info.Printf("successfully deleted user %s\n", username)
}
