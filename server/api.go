package main

import (
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	soundcloudapi "github.com/zackradisic/soundcloud-api"
)

type TrackInfo struct{
	ID int64 `json:"id"`
	Title string `json:"name"`
	Duration int64 `json:"value"`
	ArtworkURL string `json:"ArtworkURL"`
}

type UploaderInfo struct{
	ID int64 `json:"id"`
	Name string `json:"name"`
	Tracks []TrackInfo `json:"children"`
	AvatarURL string `json:"AvatarURL"`
}

type UserLikes struct{
	Name string `json:"name"`
	Likes []UploaderInfo `json:"children"`
	AvatarURL string `json:"AvatarURL"`
}

func FetchLikes(profileURL string) (UserLikes, error){

	clientID, err := getClientID()

	if err != nil {
		fmt.Println(err)
		return UserLikes{}, errors.New("Failed to connect to Soundcloud.")
	}

	//init souncloud api
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{ClientID: clientID})
	if err != nil {
		fmt.Println(err)
		return UserLikes{}, errors.New("Failed to connect to Soundcloud.")
	}

	//getting my likes
	likesPaginated, err := sc.GetLikes(soundcloudapi.GetLikesOptions{
		ProfileURL: profileURL,
		Limit: 50,
		Offset: "",
		Type: "all",
	})
	if err != nil {
		fmt.Println(err)
		return UserLikes{}, errors.New("Failed to get likes from Soundcloud.")
	}

	likes, err := likesPaginated.GetLikes()
	if err != nil {
		fmt.Println(err)
		return UserLikes{}, errors.New("Failed to get likes from Soundcloud.")
	}

	data := make(map[string]UploaderInfo)
	for i := 0; i < len(likes); i++ {
		if (likes[i].Track.Title == ""){
			continue	
		}
		
		var info TrackInfo
		if (likes[i].Track.ArtworkURL != ""){
			info = TrackInfo{
				ID: likes[i].Track.ID,
				Title: likes[i].Track.Title, 
				ArtworkURL: strings.ReplaceAll(likes[i].Track.ArtworkURL, "large", "t500x500"),
				Duration: likes[i].Track.FullDurationMS,
			}
		}else if(likes[i].Track.User.AvatarURL != ""){
			info = TrackInfo{
				ID: likes[i].Track.ID,
				Title: likes[i].Track.Title, 
				ArtworkURL: strings.ReplaceAll(likes[i].Track.User.AvatarURL, "large", "t500x500"),
				Duration: likes[i].Track.FullDurationMS,
			}
		}else{
			info = TrackInfo{
				ID: likes[i].Track.ID,
				Title: likes[i].Track.Title, 
				ArtworkURL: "https://d21buns5ku92am.cloudfront.net/26628/images/419679-1x1_SoundCloudLogo_cloudmark-f5912b-large-1645807040.jpg",
				Duration: likes[i].Track.FullDurationMS,
			}
		}

		uploader, ok := data[likes[i].Track.User.Username];

		if(ok){
			uploader.Tracks = append(uploader.Tracks, info)
			data[likes[i].Track.User.Username] = uploader
		}else{
			data[likes[i].Track.User.Username] = UploaderInfo{
				ID: likes[i].Track.User.ID,
				Name: likes[i].Track.User.Username, 
				AvatarURL: strings.ReplaceAll(likes[i].Track.User.AvatarURL, "large", "t500x500"),
				Tracks: []TrackInfo {info},
			}
		}	
	}

	user, err := sc.GetUser(soundcloudapi.GetUserOptions{ProfileURL: profileURL})
	if err != nil {
		fmt.Println(err)
		return UserLikes{}, errors.New("Failed to get likes from Soundcloud.")
	}

	userlikes := UserLikes{
		Likes: slices.Collect(maps.Values(data)), 
		Name: user.Username, 
		AvatarURL: strings.ReplaceAll(user.AvatarURL, "large", "t500x500"),
	}

	return userlikes, nil
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
