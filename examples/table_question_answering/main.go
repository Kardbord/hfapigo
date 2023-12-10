package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/Kardbord/hfapigo/v2"
)

const TableRows = 2

const query = "What is the population of Anytown?"

var table = func() map[string][]string {
	// Table looks like this:
	// +-----------------------+
	// |City      | Population |
	// |-----------------------|
	// |Anytown   |      12345 |
	// |Someplace |      7890  |
	// +-----------------------+
	return map[string][]string{
		"City":       {"Anytown", "Someplace"},
		"Population": {"12345", "7890"},
	}
}

func main() {
	outputTable()
	fmt.Println("\nSending Query:", query)
	sendRequest()
}

func outputTable() {
	fmt.Println("Data Table")
	fmt.Println("+---------------------+")
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', tabwriter.Debug)

	rowQueue := [][]string{}
	for hdr, col := range table() {
		fmt.Fprintf(w, "%s\t", hdr)
		rowQueue = append(rowQueue, col)
	}
	fmt.Fprintln(w)

	for _, col := range rowQueue {
		if len(col) != TableRows {
			panic("all columns must have the same number of rows")
		}
	}

	for i := 0; i < TableRows; i++ {
		for _, col := range rowQueue {
			fmt.Fprintf(w, "%s\t", col[i])
		}
		fmt.Fprintln(w)
	}

	w.Flush()
	fmt.Println("+---------------------+")
}

func sendRequest() {
	type ChanRv struct {
		resp *hfapigo.TableQuestionAnsweringResponse
		err  error
	}
	ch := make(chan ChanRv)

	go func() {
		tqaResp, err := hfapigo.SendTableQuestionAnsweringRequest(hfapigo.RecommendedTableQuestionAnsweringModel, &hfapigo.TableQuestionAnsweringRequest{
			Inputs: hfapigo.TableQuestionAnsweringInputs{
				Query: query,
				Table: table(),
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		ch <- ChanRv{tqaResp, err}
	}()

	for {
		select {
		case chrv := <-ch:
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}
			fmt.Println("\nAnswer:", chrv.resp.Answer)
			return

		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
