package game_engine

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func (e *baseEventMessageHandler) handleChatEventMessage(message EventMessage, pool CampaignPool) error {
	log.Printf("- Chat Data Events Type: '%d' Message: '%s'", message.Type, message.Id)

	// @todo Clean up mess
	if message.Type >= TypeChatBroadcast && message.Type <= TypeChatWhisper {

		justPassTroughMessage := true
		if message.Type == TypeChatBroadcast {
			// Check for commands

			// Undo escaping
			clearedBody := html.UnescapeString(message.Body)

			if strings.HasPrefix(clearedBody, "/roll") {

				// Attempt to parse roll
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

				// Check string
				if _, err := regexp.Match(`(^[\dd\-+]+$)`, []byte(rawValue)); err != nil {
					errorMessage := getChatMessage("/roll command only allows dice, digits, plus or minus")
					errorMessage.Destinations = append(errorMessage.Destinations, message.Source)
					return err
				}

				// Break up operations
				fifoQue := make([]string, 0)
				lastOperationIndex := -1
				lastDIndex := -1
				for index, value := range rawValue {
					if value == 'd' {
						lastDIndex = index
					} else if value == '+' {
						fifoQue = append(fifoQue, rawValue[lastOperationIndex+1:index])
						fifoQue = append(fifoQue, string(rawValue[index]))
						lastOperationIndex = index
					} else if value == '-' {
						// Minus cloud be part of the int value and not an operant
						partOfDice := lastDIndex != -1 && lastDIndex > lastOperationIndex && lastDIndex == index-1
						partOfInt := lastOperationIndex == index-1

						if !partOfDice && !partOfInt {
							fifoQue = append(fifoQue, rawValue[lastOperationIndex+1:index])
							fifoQue = append(fifoQue, string(rawValue[index]))
							lastOperationIndex = index
						}
					}
					if index == len(rawValue)-1 {
						fifoQue = append(fifoQue, rawValue[lastOperationIndex+1:index+1])
					}
				}

				// Run through operations
				sum := 0
				lastOperand := "+" // Sum should start with 0 (sum) +
				sumPath := ""
				for _, value := range fifoQue {

					if value == "+" || value == "-" {
						sumPath += " " + value + " "
						lastOperand = value
						continue
					}

					calculatedValue := 0
					if slices.Contains([]byte(value), 'd') {

						sumPath += " [ "
						if strings.Count(value, "d") > 1 {
							sumPath += "0 ] "
							continue // Do not calculate broken dices
						}

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

						if diceEyes != 0 && nrOfDice != 0 {
							negate := false
							if nrOfDice < 0 {
								negate = true
								nrOfDice = 0 - nrOfDice
							}
							if diceEyes < 0 {
								negate = !negate // Allow for the possibility to cancel out
								diceEyes = 0 - diceEyes
							}

							for throw := 0; throw < nrOfDice; throw++ {
								throwValue := 1
								if diceEyes-1 != 0 {
									if diceEyes == 2 {
										throwValue = rand.Intn(diceEyes) + 1
									} else {
										throwValue = rand.Intn(diceEyes-1) + 1
									}
								}
								if negate {
									sumPath += strconv.Itoa(0-throwValue) + " "
								} else {
									sumPath += strconv.Itoa(throwValue) + " "
								}

								calculatedValue = calculatedValue + throwValue
							}
							if negate {
								calculatedValue = 0 - calculatedValue
							}
						} else {
							sumPath += "0 "
						}
						sumPath += "] "
					} else {
						if n, err := strconv.Atoi(value); err == nil {
							calculatedValue = n
						}
						sumPath += " " + strconv.Itoa(calculatedValue) + " "
					}

					if lastOperand == "-" {
						sum = sum - calculatedValue
					} else if lastOperand == "+" {
						sum = sum + calculatedValue
					}
				}
				if sumPath == "" {
					sumPath = "0"
				}

				rollResult := fmt.Sprintf("Rolled [%s] for a total of: '%d'. \nSolution: '%s'", rawValue, sum, sumPath)
				chatMessage := getChatMessage(rollResult)
				chatMessage.Source = message.Source
				pool.TransmitEventMessage(chatMessage)
				justPassTroughMessage = false
			}

		}

		if justPassTroughMessage {
			// Just pass message trough
			pool.TransmitEventMessage(message)
		}
	}

	return nil
}

func getChatMessage(message string) EventMessage {
	chatMessage := NewEventMessage()
	chatMessage.Source = ServerUser
	chatMessage.Body = message
	chatMessage.Type = TypeChatBroadcast
	return chatMessage
}
