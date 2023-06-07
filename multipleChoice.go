package goprettycli

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

type MultipleChoiceStruct struct {
	// choice is the text that will be displayed
	Choice string

	// id is the value that will be returned
	// its also used to run the function, make sure its unique
	Id string
}

/*
**MULTIPLE CHOICE**

lets the user select multiple choices with the spacebar
confirm with enter

Returns a slice of the selected choices
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

	// forceResult
	// if true, the user has to select at least one choice, otherwise it won't return
	false,

)

fmt.Println(selected) // ["1", "3"]
*/
func MultipleChoice(choices []MultipleChoiceStruct, forceResult bool) []string {
	// selected will hold the position of the cursor
	selected := 0
	// the active map will hold the selected choices
	// these are the choices that will be returned
	// they get selected by pressing space
	active := map[string]bool{}

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
			// make selected choice bold and add '»'
			// this is always the first choice
			fmt.Print("\033[1m", "» - "+choice.Choice, "\033[0m", "\n")
		} else {
			// non-selected choices
			fmt.Println("  - " + choice.Choice)
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

		// if enter is pressed, return the selected choices
		case key == keyboard.KeyEnter:

			if !forceResult {
				// if forceResult is false, return the selected choices, even if none are selected
				var activeChoices []string
				for id := range active {
					activeChoices = append(activeChoices, id)
				}
				return activeChoices
			} else {
				// if forceResult is true, return the selected choices only if at least one is selected
				if len(active) > 0 {
					var activeChoices []string
					for id := range active {
						activeChoices = append(activeChoices, id)
					}
					return activeChoices
				}
			}

		// if space is pressed, toggle the selected choice
		case key == keyboard.KeySpace:
			if active[choices[selected].Id] {
				// if its selected, delete it from the active map
				delete(active, choices[selected].Id)
			} else {
				// if its not selected, add it to the active map
				active[choices[selected].Id] = true
			}

		// if upArrow is pressed, move the cursor
		case key == keyboard.KeyArrowUp:
			// if the cursor is at the top, don't move it
			// if its lower 0, it would be over the first choice
			if selected > 0 {
				selected--
			}

		// if downArrow is pressed, move the cursor
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

		// print the choices
		for i, choice := range choices {
			if i == selected {
				// selected choice means whe moved the cursor to it, via up/down arrow
				// so we make selected choice bold and add '»'
				if active[choice.Id] {
					// if its also active, add a checkmark
					fmt.Print("\033[1m", "» ✓ "+choice.Choice, "\033[0m", "\n")
				} else {
					// if its not active, add a simple '-'
					fmt.Print("\033[1m", "» - "+choice.Choice, "\033[0m", "\n")
				}
			} else if active[choice.Id] {
				// if its active, add a checkmark
				fmt.Print("\033[1m", "  ✓ "+choice.Choice, "\033[0m", "\n")
			} else {
				// non-selected choices
				fmt.Println("  - " + choice.Choice)
			}
		}
	}
}
