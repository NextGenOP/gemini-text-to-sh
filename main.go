package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Access your API key as an environment variable
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	model.SetTemperature(0.65)

	// Get the input argument from the command line
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Please provide at least one argument")
	}
	inputText := args[0]
	executeOutput := false
	if len(args) > 1 && args[0] == "-x" {
		executeOutput = true
		inputText = args[1]
	}
	fullPrompt := fmt.Sprintf("Convert this text to a command that works in a BASH shell, return the command with comment referring to https://explainshell.com/explain?cmd=[URL ENCODED COMMAND] only . Text to convert: %s", inputText)
	// Generate content using the input text
	resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Candidates) < 1 {
		fmt.Println("No response generated")
	} else {
		c := resp.Candidates[0]
		if c.Content != nil {
			for _, part := range c.Content.Parts {
				if executeOutput {
					var execute string
					fmt.Println("Response:", part)
					fmt.Print("Do you want to execute it? (y/n): ")
					fmt.Scanln(&execute)

					// Execute the response if the user confirms
					if strings.ToLower(execute) == "y" {
						// Execute the response here

						cmd := strings.TrimSpace(fmt.Sprintf("%v", part))
						fmt.Println("Executing the response:", cmd)
						out, err := exec.CommandContext(ctx, "bash", "-c", cmd).Output()
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("Output:", string(out))
					} else {
						fmt.Println("Response not executed.")
					}
				} else {
					fmt.Println(part)
				}
			}
		} else {
			fmt.Println("<empty response from model>")
		}
	}
}
