package bot

import (
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"io"
	"musicbot/api/youtube"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var playing = false

func RunBot(disToken string) {
	dg, err := discordgo.New("Bot " + disToken)
	if err != nil {
		fmt.Println("can't create discord session", err)
		return
	}
	dg.AddHandler(onMessage)
	err = dg.Open()
	if err != nil {
		fmt.Println("can't connect to the bot", err)
		return
	}
	fmt.Println("bot started!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	err = dg.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!link") {
		go doLink(m, s)
	} else if strings.EqualFold(m.Content, "!play") {
		if !playing {
			playing = true
			go play(m, s, playing)
		} else {
			s.ChannelMessageSend(m.ChannelID, "player is already active!")
		}
	} else if strings.Contains(m.Content, "!add") {
		go add(m, s)
	}

}

func doLink(m *discordgo.MessageCreate, s *discordgo.Session) {
	search := strings.Replace(m.Content, "!link ", "", 1)
	_, err := s.ChannelMessageSend(m.ChannelID, YT.GetLink(search))
	if err != nil {
		fmt.Println("message didn't send", err)
	}
}
func add(m *discordgo.MessageCreate, s *discordgo.Session) {
	search := strings.Replace(m.Content, "!add ", "", 1)
	stream := YT.GetAudio(search)
	YT.DownloadAudio(&stream)
	_, err := s.ChannelMessageSend(m.ChannelID, "Done!")
	if err != nil {
		fmt.Println("message didn't send", err)
	}
}

func play(m *discordgo.MessageCreate, s *discordgo.Session, playing bool) {
	dir := "songs"
	dgv, err := s.ChannelVoiceJoin(m.GuildID, m.ChannelID, false, false)
	if err != nil {
		panic(err)
	}
	iSong := 1
	for !isEmpty(dir) {
		dgvoice.PlayAudioFile(dgv, filepath.Join(dir, fmt.Sprintf("song%s.mp3", strconv.Itoa(iSong))), make(chan bool))
		os.Remove(filepath.Join(dir, fmt.Sprintf("song%s.mp3", strconv.Itoa(iSong))))
		iSong += 1
		fmt.Println(iSong)
	}
	playing = false
}

func isEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.ReadDir(1)
	if err == io.EOF {
		return true
	}
	return false
}
