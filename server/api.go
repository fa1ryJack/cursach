package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	soundcloudapi "github.com/zackradisic/soundcloud-api"
)

type BasicInfo struct{
	ID int64
	Title string
	ArtworkURL string
}

func FetchMyLikes(profileURL string) ([]BasicInfo, error){

	clientID, err := getClientID()

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Failed to connect to Soundcloud.")
	}

	//init souncloud api
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{ClientID: clientID})
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Failed to connect to Soundcloud.")
	}

	//getting my likes
	likesPaginated, err := sc.GetLikes(soundcloudapi.GetLikesOptions{
		ProfileURL: profileURL,
		Limit: 1000,
		Offset: "",
		Type: "all",
	})
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Failed to get likes from Soundcloud.")
	}

	likes, err := likesPaginated.GetLikes()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Failed to get likes from Soundcloud.")
	}

	var data []BasicInfo
	for i := 0; i < len(likes); i++ {
		if (likes[i].Track.Title == ""){
			continue
			
		}
		var info BasicInfo
		if (likes[i].Track.ArtworkURL != ""){
			info = BasicInfo{
				ID: likes[i].Track.ID,
				Title: likes[i].Track.Title, 
				ArtworkURL: likes[i].Track.ArtworkURL}
		}else{
			info = BasicInfo{
				ID: likes[i].Track.ID,
				Title: likes[i].Track.Title, 
				ArtworkURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQHnPaUNBw_Kr6J7M77WWMbUoCDTq75SZXNDw&s"}
		}
		data = append(data, info)
	}


	return data, nil
}

func getClientID() (string, error){
	searching := true
	myURL :="https://w.soundcloud.com/player/?url=https://api.soundcloud.com/tracks/"
	trackID := 1
	for searching{
		myURL = myURL + strconv.Itoa(trackID)
		resp, err := http.Get(myURL)
		if err != nil {
			return "", errors.Wrap(err, "Failed to fetch SoundCloud.")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err, "Failed to read body while fetching SoundCloud Client ID")
			continue
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
			continue
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("", errors.Wrap(err, "Failed to read body while fetching SoundCloud Client ID"))
			continue
		}

		bodyString = string(body)

		// Extract the client ID
		if strings.Contains(bodyString, `client_id:u?"`) {
			clientID := strings.Split(bodyString, `client_id:u?"`)[1]
			clientID = strings.Split(clientID, `"`)[0]
			return clientID, nil
		}
	}
	return "", errors.New("Failed to find client_id")
}