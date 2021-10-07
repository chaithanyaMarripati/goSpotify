//here we get the token
//with the token, we need to call the spotify for name and currently listening songs

package spotify

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/chaithanyaMarripati/goSpotify/config"
)

type meSpotify struct {
	DisplayName string `json:"display_name"`
}
type currentSong struct {
	Item struct {
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
	} `json:"item"`
}

type Spotify interface {
	CurrentSong(accessToken string) (string, error)
	Profile(accessToken string) (string, error)
}

type HttpSpotify struct{}

func (s *HttpSpotify) CurrentSong(accessToken string) (string, error) {
	getUserSpotify := config.EnvVariables.GetMeSpotify
	client := &http.Client{}
	req, _ := http.NewRequest("GET", getUserSpotify, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	responsePayload := &meSpotify{}
	json.Unmarshal(body, responsePayload)
	return responsePayload.DisplayName, nil
}

func (s *HttpSpotify) Profile(accessToken string) (string, error) {
	getCurrentSpotify := config.EnvVariables.CurrentlyPlaying
	client := &http.Client{}
	req, _ := http.NewRequest("GET", getCurrentSpotify, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", nil
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var responsePayload currentSong
	json.Unmarshal(body, &responsePayload)
	return responsePayload.Item.Album.Name, nil
}
