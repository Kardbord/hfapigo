package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kardbord/hfapigo/v3"
)

const HuggingFaceTokenEnv = "HUGGING_FACE_TOKEN"

func init() {
	key := os.Getenv(HuggingFaceTokenEnv)
	if key != "" {
		hfapigo.SetAPIKey(key)
	}
}

func init() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		quit()
	}()
}

const model = "facebook/blenderbot-400M-distill"

// Deprecated: HF's conversational endpoint seems to be under construction
// and slated to be either updated or replaced.
// TODO: Update or remove conversational support once it becomes
// clear what its replacement is.
func main() {
	fmt.Println("Enter your messages below. Hit enter to send. Use Ctrl+c or Ctrl+d to quit.")

	var input string
	var pastInputs, pastResps []string
	in := bufio.NewReader(os.Stdin)
	for {
		input = prompt(in)
		resp, err := sendRequest(input, pastResps, pastInputs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pastInputs = resp.Conversation.PastUserInputs
		pastResps = resp.Conversation.GeneratedResponses
		fmt.Printf("\nbot>%s\n", resp.GeneratedText)
	}
}

func quit() {
	fmt.Println("\nbot> Goodbye")
	os.Exit(0)
}

func sendRequest(input string, pastInputs, pastResps []string) (*hfapigo.ConversationalResponse, error) {
	type ChanRv struct {
		resp *hfapigo.ConversationalResponse
		err  error
	}
	ch := make(chan ChanRv)

	go func() {
		cresps, err := hfapigo.SendConversationalRequest(model, &hfapigo.ConversationalRequest{
			Inputs: hfapigo.ConverstationalInputs{
				Text:               input,
				GeneratedResponses: pastInputs,
				PastUserInputs:     pastResps,
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})

		ch <- ChanRv{cresps, err}
	}()

	for {
		select {
		case chrv := <-ch:
			return chrv.resp, chrv.err
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func prompt(in *bufio.Reader) string {
	fmt.Print("you> ")
	input, err := in.ReadString('\n')
	if err == io.EOF {
		quit()
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return input
}
