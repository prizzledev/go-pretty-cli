package goprettycli

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

type ChoiceStruct struct {
	// choice is the text that will be displayed
	Choice string

	// id is the value that will be returned
	// its also used to run the function, make sure its unique
	Id string
}

/*
**SINGLE CHOICE**

let the user select one choice with the arrow keys
confirm with enter

# Returns the id of the selected choice as a string

Example:
selected := Choice(

	[]ChoiceStruct{
		{
			choice: "Choice 1",
			id:     "1",
		},
		{
			choice: "Choice 2",
			id:     "2",
		},
		{
			choice: "Choice 3",
			id:     "3",
		},
	},

)

fmt.Println(selected) // ["1"]
*/
func Choice(choices []ChoiceStruct) string {
	// selected will hold the position of the cursor
	selected := 0

	// open keyboard
	// this is needed to read the keys
	// use defer to close it after the function is done
	// its from "github.com/eiannone/keyboard"
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	// initial print of choices
	// here the first choice is always selected
	for i, choice := range choices {
		if i == selected {
			// Make selected choice bold and add '»'
			// this is always the first choice
			fmt.Print("\033[1m", "» "+choice.Choice, "\033[0m", "\n")
		} else {
			// non-selected choices
			fmt.Println("- " + choice.Choice)
		}
	}

	for {
		// Read keyboard
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		// rn no special char is used
		_ = char

		// switch between the different keys
		switch {

		// if enter or space is pressed, return the id of the selected choice
		case key == keyboard.KeySpace || key == keyboard.KeyEnter:
			// indexing the choices with selected
			return choices[selected].Id

		// if arrow up is pressed, move the cursor up
		case key == keyboard.KeyArrowUp:
			// if the cursor is at the top, don't move it
			// if its lower 0, it would be over the first choice
			if selected > 0 {
				selected--
			}

		// if arrow down is pressed, move the cursor down
		case key == keyboard.KeyArrowDown:
			// if the cursor is at the bottom, don't move it
			// if its higher than the len(choices), it would be under the last choice
			if selected < len(choices)-1 {
				selected++
			}
		}

		// move the cursor up by the len(choices) lines
		// this is so we can overwrite the previous choices
		// so it looks like we're moving the cursor and not stacking up the choices
		fmt.Printf("\033[%dA", len(choices))

		// print the choices again
		// to overwrite the previous ones
		for i, choice := range choices {
			// not anymore the first choice
			// user probably moved the cursor
			if i == selected {
				// Make selected choice bold and add '»'
				fmt.Print("\033[1m", "» "+choice.Choice, "\033[0m", "\n")
			} else {
				// non-selected choices
				fmt.Println("- " + choice.Choice)
			}
		}
	}
}
