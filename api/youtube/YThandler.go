package YT

import (
	"encoding/json"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type idInfo struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoID"`
}
type item struct {
	Id idInfo `json:"id"`
}
type video struct {
	Items []item `json:"items"`
}

const searchURL = "https://www.googleapis.com/youtube/v3/search"

var YtToken = os.Args[2]

const videoURL = "https://www.youtube.com/watch?v="

func (v *video) getYTStruct(userSearch string) error {
	client := http.Client{}
	request, err := createRequest(userSearch)
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}
	return nil
}

func GetLink(userSearch string) string {

	return videoURL + GetVideoID(userSearch)
}

func GetVideoID(userSearch string) string {
	v := video{}
	err := v.getYTStruct(userSearch)
	if err != nil {
		panic(err)
	}
	return v.Items[0].Id.VideoID
}

func GetAudio(userSearch string) io.ReadCloser {
	client := youtube.Client{}

	v, err := client.GetVideo(GetVideoID(userSearch))
	if err != nil {
		panic(err)
	}
	formats := v.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(v, &formats[0])
	if err != nil {
		panic(err)
	}
	return stream
}

func DownloadAudio(s *io.ReadCloser) {
	var (
		dir           string = "songs"
		lastFileIndex int    = 0
	)
	f, _ := os.ReadDir(dir)
	if len(f) != 0 {
		lastFileIndex, _ = strconv.Atoi(string(f[len(f)-1].Name()[4]))
	}
	file, err := os.Create(filepath.Join(dir, fmt.Sprintf("song%s.mp3", strconv.Itoa(lastFileIndex+1))))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, *s)
	if err != nil {
		panic(err)
	}
}

func createRequest(userSearch string) (*http.Request, error) {
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("part", "id")
	query.Add("q", userSearch)
	query.Add("key", YtToken)
	req.URL.RawQuery = query.Encode()

	return req, nil
}
