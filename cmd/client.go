/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runClient()
	},
}
var useWebsockets =false
func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().BoolVarP(&useWebsockets,"websocket", "", false, "is use websocket request")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


type Message struct {
	Id      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func runClient()  {

		if useWebsockets {
			ws, err := websocket.Dial("ws://localhost:8080/", "", "http://localhost:8080")
			for {
				var m Message
				err = websocket.JSON.Receive(ws, &m)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(err)
				}
				log.Printf("Received: %+v", m)
			}
		} else {
			log.Println("Sending request...")
			req, err := http.NewRequest("GET", "http://localhost:8080", nil)
			if err != nil {
				log.Fatal(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				log.Fatalf("Status code is not OK: %v (%s)", resp.StatusCode, resp.Status)
			}

			dec := json.NewDecoder(resp.Body)
			for {
				var m Message
				err := dec.Decode(&m)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(err)
				}
				log.Printf("Got response: %+v", m)
			}
		}
	log.Println("Server finished request...")
}