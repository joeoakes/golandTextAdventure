package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Game struct {
	XMLName xml.Name `xml:"game"`
	Rooms   []Room   `xml:"rooms>room"`
	Items   []Item   `xml:"items>item"`
	Events  []Event  `xml:"events>event"`
}

type Event struct {
	ID          string   `xml:"id,attr"`
	Description string   `xml:"description"`
	Choices     []Choice `xml:"choices>choice"`
}

type Choice struct {
	ID          string `xml:"id,attr"`
	Description string `xml:"description"`
	Outcome     string `xml:"outcome"`
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

// Character represents a game character.
type Character struct {
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Health      int      `xml:"health"`
	Inventory   []string `xml:"inventory>item"`
}

// Characters represents a collection of characters.
type Characters struct {
	Characters []Character `xml:"character"`
}

func checkItemsInRoom(roomID string, items []Item) {
	fmt.Println("Items in the room:")
	for _, item := range items {
		if item.Location == roomID {
			fmt.Println(item.Name + ": " + item.Description)
		}
	}
}

func printAvailableExits(exits []Exit) {
	fmt.Print("Exits: ")
	for i, exit := range exits {
		fmt.Print(exit.Direction)
		if i < len(exits)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()
}

func findCurrentRoom(roomID string, rooms []Room) *Room {
	for _, room := range rooms {
		if room.ID == roomID {
			return &room
		}
	}
	return nil
}

func main() {

	// Read and parse the game data from XML
	data, err := os.ReadFile("space.xml")
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
		currentRoomData := findCurrentRoom(currentRoom, game.Rooms)
		if currentRoomData == nil {
			fmt.Println("Error: Current room not found.")
			os.Exit(1)
		}
		// Print room description
		fmt.Println("\n" + currentRoomData.Name)
		fmt.Println(currentRoomData.Description)
		// Print available exits
		printAvailableExits(currentRoomData.Exits)
		// Check for items in the room
		checkItemsInRoom(currentRoom, game.Items)
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
