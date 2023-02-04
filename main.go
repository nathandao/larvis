package main

import (
	"fmt"
	"os"
)

func printUsage() {
	fmt.Println("Usage: larvis FIRST_HAND SECOND_HAND")
	fmt.Println()
}

// exitWithError prints the error and command usage before exiting with an error code
func exitWithError(err error) {
	fmt.Println()
	fmt.Println(err.Error())
	fmt.Println()
	printUsage()
	os.Exit(1)
}

// validateArguments checks if exactly 2 poker hands are provided.
func validateArguments(pokerHands []string) error {
	if len(pokerHands) < 2 {
		return fmt.Errorf("urg, %d poker hand was given. Please input 2 poker hands", len(pokerHands))
	}

	if len(pokerHands) > 2 {
		return fmt.Errorf("woaa, you input %d hands! I can only work with 2", len(pokerHands))
	}

	return nil
}

func main() {
	pokerHands := os.Args[1:]

	if err := validateArguments(pokerHands); err != nil {
		exitWithError(err)
	}

	var handsErr error
	for _, hand := range pokerHands {
		if err := validateHand(hand); err != nil {
			err = fmt.Errorf("- %s: %s", hand, err.Error())
			if handsErr == nil {
				handsErr = err
			} else {
				handsErr = fmt.Errorf("%s\n%s", handsErr.Error(), err.Error())
			}
		}
	}

	if handsErr != nil {
		exitWithError(fmt.Errorf("oops, looks like there are issues with your poker hands:\n%s", handsErr.Error()))
	}

	r := calculateResult(pokerHands[0], pokerHands[1])
	fmt.Println(r)
}
