package main

import (
	"dnd-app/dice"
	"fmt"
	"log"
	"strconv"
)

func getUserInput(prompt string) int {
	var input string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	parsed, err := strconv.Atoi(input)
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}
	return parsed
}

func main() {
	userID := getUserInput("Enter user ID: ")
	diceSides := getUserInput("Enter the number of sides on the dice (e.g., 6 for d6, 20 for d20): ")
	numberOfRolls := getUserInput("Enter the number of rolls: ")
	rngChoice := getUserInput("Choose RNG (1: Standard, 2: Karmic, 3: Pseudo): ")

	seedGen := &dice.SeedGenerator{}
	seed := seedGen.GenerateSeed(userID)
	fmt.Printf("Generated cryptographic seed for user %d: %d\n", userID, seed)

	var rng dice.RNG
	switch rngChoice {
	case 1:
		rng = dice.NewRandomNumberGenerator(seed)
	case 2:
		rng = dice.NewKarmicDiceRNG(dice.NewRandomNumberGenerator(seed))
	case 3:
		rng = dice.NewPseudoRandomRNG(seed)
	default:
		log.Fatalf("Invalid RNG choice")
	}

	diceRoller := dice.NewDiceRoller(rng)
	diceRoller.RollAndPrint(diceSides, numberOfRolls)
}
