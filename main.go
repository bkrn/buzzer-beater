//Package main launches buzzer-beater for door monitoring goodness
package main

import (
	"bkrn/buzzer-beater/doorserver"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	var command string
	flag.StringVar(&command, "command", "runserver", "Select command")
	flag.Parse()

	if command == "runserver" {
		control := doorserver.NewControl()
		http.ListenAndServe(":8080", control)
	} else if command == "adduser" {
		control := doorserver.NewControl()
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Name: ")
		name, _ := reader.ReadString('\n')
		fmt.Print("Enter phone: ")
		phone, _ := reader.ReadString('\n')
		fmt.Print("Enter email: ")
		email, _ := reader.ReadString('\n')
		//reader := bufio.NewReader(os.Stdin)
		var pw string
		var confirmpw string
		for pw == "" || pw != confirmpw {
			fmt.Print("Enter password: ")
			pw, _ = reader.ReadString('\n')
			fmt.Print("Confirm password: ")
			confirmpw, _ = reader.ReadString('\n')
		}
		usr := &doorserver.DoorUser{
			Name:     strings.TrimSpace(name),
			Phone:    strings.TrimSpace(phone),
			Email:    strings.TrimSpace(email),
			Password: strings.TrimSpace(pw),
		}
		err := usr.HashPass()
		if err != nil {
			log.Fatal(err)
		}
		usrdat, _ := json.Marshal(usr)
		nusr, err := control.Collections["users"].Post(usrdat)
		fmt.Println(nusr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(nusr)
	}

}
