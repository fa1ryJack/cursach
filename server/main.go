package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)


func GetClientID() (string){
	searching := true
	myURL :="https://w.soundcloud.com/player/?url=https://api.soundcloud.com/tracks/"
	trackID := 1
	for searching{
		myURL = myURL + strconv.Itoa(trackID)
		resp, err := http.Get(myURL)
		if err != nil {
			fmt.Println(err, "Failed to fetch SoundCloud Client ID")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("", errors.Wrap(err, "Failed to read body while fetching SoundCloud Client ID"))
		}

		bodyString := string(body)

		if (strings.Contains(bodyString,"script crossorigin src")){
			searching = false
		}else{
			trackID++
			continue
		}

		// The link to the JS file with the client ID looks like this:
		// <script crossorigin src="https://widget.sndcdn.com/widget-9-9e56aaba7dc4.js"></script>
		split := strings.Split(bodyString, `<script crossorigin src="`)
		urls := []string{}

		// Extract all the URLS that match our pattern
		for _, raw := range split {
			u := strings.Replace(raw, `"></script>`, "", 1)
			u = strings.Split(u, "\n")[0]

			if string([]rune(u)[0:25]) == `https://widget.sndcdn.com` {
				urls = append(urls, u)
			}
		}

		// It seems like our desired URL is always imported last,
		// so we use urls[len(urls) - 1]
		resp, err = http.Get(urls[len(urls)-1])
		if err != nil {
			fmt.Println("", errors.Wrap(err, "Failed to fetch SoundCloud Client ID"))
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("", errors.Wrap(err, "Failed to read body while fetching SoundCloud Client ID"))
		}

		bodyString = string(body)

		// Extract the client ID
		if strings.Contains(bodyString, `client_id:u?"`) {
			clientID := strings.Split(bodyString, `client_id:u?"`)[1]
			clientID = strings.Split(clientID, `"`)[0]
			fmt.Println(clientID, trackID)
			return clientID
		}else{
			continue
		}
	}
	return "failed"
}

