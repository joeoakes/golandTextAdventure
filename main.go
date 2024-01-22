package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Game struct {
	XMLName    xml.Name    `xml:"game"`
	Rooms      []Room      `xml:"rooms>room"`
	Items      []Item      `xml:"items>item"`
	Events     []Event     `xml:"events>event"`
	Characters []Character `xml:"characters>character"`
	StartLevel string      `xml:"start-level"`
	Levels     []Level     `xml:"levels>level"`
}

type Level struct {
	ID          string `xml:"id,attr"`
	Description string `xml:"description"`
	Rooms       []Room `xml:"rooms>room"`
}

type Event struct {
	ID          string   `xml:"id,attr"`
	Description string   `xml:"description"`
	Choices     []Choice `xml:"choices>choice"`
}

type Choice struct {
	ID          string `xml:"id,attr"`
	Description string `xml:"description,attr"`
	Outcome     string `xml:"outcome"`
}

type Room struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name,attr"`
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

func getLevels(levels []Level) {
	// Iterate through levels and rooms
	for _, level := range levels {
		fmt.Println("\nLevel:", level.Description)
		for _, room := range level.Rooms {
			fmt.Printf("\nRoom: %s\nDescription: %s\nExits:\n", room.Name, room.Description)
			for _, exit := range room.Exits {
				fmt.Printf("- %s to %s\n", exit.Direction, exit.Target)
			}
		}
	}
}

func getCharacters(characters []Character) {
	fmt.Println("Game Characters:")
	for _, character := range characters {
		fmt.Printf("Name: %s\nDescription: %s\nHealth: %d\nInventory: %v\n\n",
			character.Name, character.Description, character.Health, character.Inventory)
	}
}

// Iterate through game events and choices
func getEvents(events []Event) {
	for _, event := range events {
		fmt.Println("\nEvent:", event.Description)
		for _, choice := range event.Choices {
			fmt.Printf("Choice %s: %s (Leads to outcome: %s)\n", choice.ID, choice.Description, choice.Outcome)
		}
	}
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
	data, err := os.ReadFile("water.xml")
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

	getCharacters(game.Characters)
	getEvents(game.Events)

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
		fmt.Print("\nEnter a direction to move or 'quit' or 'q' to exit: ")
		var userInput string
		fmt.Scanln(&userInput)
		// Handle user input
		if userInput == "quit" || userInput == "q" {
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
