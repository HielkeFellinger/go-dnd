package game_engine

import (
	"fmt"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type chatCommands struct {
}

func (c *chatCommands) handleRollChatCommand(message EventMessage, pool CampaignPool, clearedBody string) error {
	// Attempt to parse roll
	rawValue := c.cleanRollString(clearedBody)

	// Check string
	if _, err := regexp.Match(`(^[\dd\-+]+$)`, []byte(rawValue)); err != nil {
		errorMessage := getChatMessage("/roll command only allows dice, digits, plus or minus")
		errorMessage.Destinations = append(errorMessage.Destinations, message.Source)
		return err
	}

	// Break up operations
	fifoQueue := c.buildFifoCalcQueue(rawValue)

	// Run through operations (could have multiple dice rolls and operations)
	rollTotal := 0
	lastOperand := "+" // Sum should start with 0 (rollTotal) +
	sumHistoryPath := ""
	for _, value := range fifoQueue {

		// Update last operant
		if value == "+" || value == "-" {
			sumHistoryPath += " " + value + " "
			lastOperand = value
			continue
		}

		// Calculate dice rolls if encountered
		calculatedValue := 0
		if slices.Contains([]byte(value), 'd') {
			// Do not calculate broken dices strings
			if strings.Count(value, "d") > 1 {
				sumHistoryPath += " [ 0 ] "
				continue
			}

			// Parse dice roll formula and make them roll
			nrOfDice, diceEyes := c.translateDieRollStringToRollsAndNrOfEyes(value)
			sumHistoryPath, calculatedValue = c.handleRollsOfDice(nrOfDice, diceEyes, sumHistoryPath, calculatedValue)
		} else if n, err := strconv.Atoi(value); err == nil {
			// Handle a "just" an int
			calculatedValue = n
			sumHistoryPath += " " + value + " "
		}

		// Update rollTotal
		if lastOperand == "-" {
			rollTotal -= calculatedValue
		} else {
			// Else should always be plus
			rollTotal += calculatedValue
		}
	}
	if sumHistoryPath == "" {
		sumHistoryPath = "0"
	}

	rollResult := fmt.Sprintf(" Rolled [%s] for a total of: '%d'. \nCalculation Path: '%s'",
		rawValue, rollTotal, sumHistoryPath)
	chatMessage := getChatMessage(rollResult)
	chatMessage.Type = TypeChatCommandRoll
	chatMessage.Source = message.Source
	pool.TransmitEventMessage(chatMessage)

	return nil
}

func (c *chatCommands) translateDieRollStringToRollsAndNrOfEyes(value string) (int, int) {
	parts := strings.Split(value, "d")
	nrOfDice := 1
	diceEyes := 0

	if len(parts) == 2 {
		if nr, err := strconv.Atoi(parts[0]); err == nil {
			nrOfDice = nr
		}
		if nr, err := strconv.Atoi(parts[1]); err == nil {
			diceEyes = nr
		}
	}
	return nrOfDice, diceEyes
}

func (c *chatCommands) handleRollsOfDice(nrOfDice int, diceEyes int, sumHistoryPath string, calculatedValue int) (string, int) {
	// zero times n or n times 0 will always be zero
	if diceEyes == 0 || nrOfDice == 0 {
		sumHistoryPath += " [ 0 ] "
		return sumHistoryPath, calculatedValue
	}

	// Check if result total needs to be negative and always calc with positive numbers
	negate := false
	if nrOfDice < 0 {
		negate = true
		nrOfDice = 0 - nrOfDice
	}
	if diceEyes < 0 {
		negate = !negate // Allow for the possibility to cancel out
		diceEyes = 0 - diceEyes
	}

	// Simulate different rolls
	sumHistoryPath += " [ "
	for throw := 0; throw < nrOfDice; throw++ {
		throwValue := 1
		if diceEyes-1 != 0 {
			throwValue = rand.Intn(diceEyes) + 1
		}
		if negate {
			sumHistoryPath += strconv.Itoa(0-throwValue) + " "
		} else {
			sumHistoryPath += strconv.Itoa(throwValue) + " "
		}
		calculatedValue = calculatedValue + throwValue
	}
	sumHistoryPath += "] "

	// Negate the total if needed (nrOfDice or diceEyes) was negative
	if negate {
		calculatedValue = 0 - calculatedValue
	}
	return sumHistoryPath, calculatedValue
}

func (c *chatCommands) buildFifoCalcQueue(rawValue string) []string {
	fifoQueue := make([]string, 0)
	lastOperationIndex := -1
	lastDIndex := -1

	// Loop through the roll string and cut it into usable calculable parts
	for index, value := range rawValue {
		if value == 'd' {
			lastDIndex = index
		} else if value == '+' {
			fifoQueue = append(fifoQueue, rawValue[lastOperationIndex+1:index])
			fifoQueue = append(fifoQueue, string(rawValue[index]))
			lastOperationIndex = index
		} else if value == '-' {
			// Check; Minus could be part of a negative int value and is not always indicative of a new operation
			partOfDice := lastDIndex != -1 && lastDIndex > lastOperationIndex && lastDIndex == index-1
			partOfInt := lastOperationIndex == index-1
			if !partOfDice && !partOfInt {
				fifoQueue = append(fifoQueue, rawValue[lastOperationIndex+1:index])
				fifoQueue = append(fifoQueue, string(rawValue[index]))
				lastOperationIndex = index
			}
		}
		if index == len(rawValue)-1 {
			fifoQueue = append(fifoQueue, rawValue[lastOperationIndex+1:index+1])
		}
	}
	return fifoQueue
}

func (c *chatCommands) cleanRollString(clearedBody string) string {
	rawValue := strings.ReplaceAll(clearedBody, "/roll", "")
	rawValue = strings.ReplaceAll(rawValue, " ", "")
	rawValue = strings.ToLower(rawValue)

	// Small operant trimming
	rawValue = strings.ReplaceAll(rawValue, "++", "+")  // ++ should result in +
	rawValue = strings.ReplaceAll(rawValue, "++", "+")  // ++ should result in + (second loop)
	rawValue = strings.ReplaceAll(rawValue, "+-+", "+") // +-+ should result in +
	rawValue = strings.ReplaceAll(rawValue, "-+-", "-") // -+- should result in -
	rawValue = strings.ReplaceAll(rawValue, "--", "+")  // -- should cancel out
	rawValue = strings.Trim(rawValue, "+")
	rawValue = strings.TrimRight(rawValue, "-")
	return rawValue
}
