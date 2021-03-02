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
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runServer()  {
	http.HandleFunc("/", Handle)

	log.Println("Serving...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}




// Heavily based on Kubernetes' (https://github.com/GoogleCloudPlatform/kubernetes) detection code.
var connectionUpgradeRegex = regexp.MustCompile("(^|.*,\\s*)upgrade($|\\s*,)")

func isWebsocketRequest(req *http.Request) bool {
	return connectionUpgradeRegex.MatchString(strings.ToLower(req.Header.Get("Connection"))) && strings.ToLower(req.Header.Get("Upgrade")) == "websocket"
}

func Handle(w http.ResponseWriter, r *http.Request) {
	// Handle websockets if specified.
	if isWebsocketRequest(r) {
		websocket.Handler(HandleWebSockets).ServeHTTP(w, r)
	} else {
		HandleHttp(w, r)
	}
	log.Println("Finished sending response...")
}

func HandleWebSockets(ws *websocket.Conn) {
	for i := 0; i < 5; i++ {
		log.Printf("Sending some data: %d", i)
		m := Message{
			Id:      i,
			Message: fmt.Sprintf("Sending you \"%d\"", i),
		}
		err := websocket.JSON.Send(ws, &m)
		if err != nil {
			log.Printf("Client stopped listening...")
			return
		}

		// Artificially induce a 1s pause
		time.Sleep(time.Second)
	}
}

func HandleHttp(w http.ResponseWriter, r *http.Request) {
	//cn, ok := w.(http.CloseNotifier)
	//if !ok {
	//	http.NotFound(w, r)
	//	return
	//}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Send the initial headers saying we're gonna stream the response.
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	enc := json.NewEncoder(w)

	for i := 0; i < 5; i++ {
		select {
		case <-r.Context().Done():
			log.Println("Client stopped listening")
			return
		default:
			// Artificially wait a second between reponses.
			time.Sleep(time.Second)

			log.Printf("Sending some data: %d", i)
			m := Message{
				Id:      i,
				Message: fmt.Sprintf("Sending you \"%d\"", i),
			}

			// Send some data.
			err := enc.Encode(m)
			if err != nil {
				log.Fatal(err)
			}
			flusher.Flush()
		}
	}
}