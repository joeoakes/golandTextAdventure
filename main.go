package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Game struct {
	XMLName xml.Name `xml:"game"`
	Rooms   []Room   `xml:"rooms>room"`
	Items   []Item   `xml:"items>item"`
}

type Room struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Exits       []Exit `xml:"exits>exit"`
}

type Exit struct {
	Direction string `xml:"direction,attr"`
	Target    string `xml:"target,attr"`
}

type Item struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Location    string `xml:"location,attr"`
}

func main() {
	// Read and parse the game data from XML
	data, err := ioutil.ReadFile("game.xml")
	if err != nil {
		fmt.Println("Error reading game data:", err)
		return
	}

	var game Game
	err = xml.Unmarshal(data, &game)
	if err != nil {
		fmt.Println("Error parsing game data:", err)
		return
	}

	// Initialize game state
	currentRoom := "start"

	fmt.Println("Welcome to the Text Adventure Game!")

	for {
		// Find the current room
		var currentRoomData Room
		for _, room := range game.Rooms {
			if room.ID == currentRoom {
				currentRoomData = room
				break
			}
		}

		// Print room description
		fmt.Println("\n" + currentRoomData.Name)
		fmt.Println(currentRoomData.Description)

		// Print available exits
		fmt.Print("Exits: ")
		for i, exit := range currentRoomData.Exits {
			fmt.Print(exit.Direction)
			if i < len(currentRoomData.Exits)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println()

		// Check for items in the room
		fmt.Println("Items in the room:")
		for _, item := range game.Items {
			if item.Location == currentRoom {
				fmt.Println(item.Name + ": " + item.Description)
			}
		}

		// Prompt for user input
		fmt.Print("\nEnter a direction to move or 'quit' to exit: ")
		var userInput string
		fmt.Scanln(&userInput)

		// Handle user input
		if userInput == "quit" {
			fmt.Println("Thanks for playing!")
			os.Exit(0)
		}

		// Check if the chosen direction is a valid exit
		var newRoom string
		for _, exit := range currentRoomData.Exits {
			if userInput == exit.Direction {
				newRoom = exit.Target
				break
			}
		}

		if newRoom != "" {
			currentRoom = newRoom
		} else {
			fmt.Println("You can't go that way.")
		}
	}
}