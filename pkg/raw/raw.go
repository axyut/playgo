package raw

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/axyut/cold/internal/list"
	"github.com/axyut/cold/internal/types"
	"github.com/mattn/go-tty"
)

const REFRESH_RATE = time.Millisecond * 200

type TUI struct {
	Playlist *list.Playlist
	Notifs   *[]string
	Songs    *types.Activelist
	Setting  *types.Config
}

func NewUI(playlist *list.Playlist, notifs *[]string, setting *types.Config) *TUI {
	return &TUI{
		Playlist: playlist,
		Notifs:   notifs,
		Setting:  setting,
	}
}

var maxX, maxY = termSize()

func (tui TUI) Display() {
	go func() {
		for {
			music := tui.Setting.Music
			playlist := (*tui.Playlist).List
			notifs := *tui.Notifs
			currentSong := tui.Playlist.CurrentSong
			Shuffle := music.Shuffle
			RepeatSong := music.RepeatSong
			RepeatPlaylist := music.RepeatPlaylist

			clear()
			HideCursor()
			// border()
			seprator()

			currentlyPlaying(playlist, currentSong)
			displayPrevSongs(playlist, currentSong)
			displayNextSongs(playlist, currentSong)
			displaySettings(Shuffle, RepeatSong, RepeatPlaylist)
			displayNowPlaying(playlist, currentSong)
			displayNotifications(notifs)

			render()
			time.Sleep(REFRESH_RATE)
		}
	}()

}

func (ui TUI) HandleInterrupt(playedList []string, completedPlaylist int) {
	HideCursor()

	// handle CTRL C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			ui.DisplayStats(playedList, completedPlaylist)
		}
	}()
}

type ListenKeyAndAction struct {
	Key    rune
	Action func()
}

func (ui TUI) ListenForKey(listenKeys []ListenKeyAndAction) {
	go func() {
		tty, err := tty.Open()
		if err != nil {
			panic(err)
		}

		for {

			if char, err := tty.ReadRune(); err == nil {
				for _, keyAction := range listenKeys {
					switch char {
					case keyAction.Key:
						keyAction.Action()
					}
				}
			}
		}
	}()
}

func displayPrevSongs(playlist []types.Song, currentSong int) {
	// Playlist
	// TODO: will make playlist scrollable with cursor later on
	// will later implement this to show all playlist when certain key pressed.
	moveCursor(pos{3, 1})
	fmt.Fprintf(screen, "PLAYLIST (%d songs)", len(playlist))

	// 1. iterate playedList? start currentSong from top and move down as playedList increases upto ~5 prev songs
	// 2. iterate playlist? start at constant row and provide what would be last ~5 songs then gradually add playedList
	// prev songs, idk which option better, now implementing 2.
	totalPrevSongs := 3
	for i := 1; i <= totalPrevSongs; i++ {
		prev := currentSong - i
		if prev <= -1 {
			for prev <= -1 {
				prev = prev + len(playlist)
			}
		}
		moveCursor(pos{2, (maxY / 4) - i})
		color.Magenta(stripString(playlist[prev].Name))
	}
}

func currentlyPlaying(playlist []types.Song, currentSong int) {
	moveCursor(pos{2, maxY / 4})
	color.Reversed()
	color.Cyan(fmt.Sprintf("⏯️ %s", stripString(playlist[currentSong].Name)))
}

func displayNextSongs(playlist []types.Song, currentSong int) {
	totalNextSongs := 6
	for j := 1; j <= totalNextSongs; j++ {
		next := currentSong + j
		if next >= len(playlist) {
			for next >= len(playlist) {
				next = next - len(playlist)
			}
		}
		moveCursor(pos{2, (maxY / 4) + j})
		color.Blue(stripString(playlist[next].Name))
	}

}

func displaySettings(Shuffle, RepeatSong, RepeatPlaylist bool) {
	intH := int(float32(maxY) / 1.25)
	moveCursor(pos{2, intH - 1})
	fmt.Fprintf(screen, "SETTINGS")
	moveCursor(pos{3, intH})
	fmt.Fprintf(screen, "[r] Repeat Song: %t", RepeatSong)
	moveCursor((pos{3, intH + 1}))
	fmt.Fprintf(screen, "[t] Repeat playlist: %t", RepeatPlaylist)
	moveCursor(pos{3, intH + 2})
	fmt.Fprintf(screen, "[y] Shuffle: %t", Shuffle)
}

func displayNowPlaying(playlist []types.Song, currentSong int) {
	moveCursor(pos{maxX / 2, 1})
	fmt.Fprintf(screen, "NOW PLAYING")
	moveCursor(pos{maxX / 2, 3})
	fmt.Fprintf(screen, "%s", stripString(playlist[currentSong].Name))
	moveCursor(pos{maxX / 2, 4})
	fmt.Fprintf(screen, "0:00 -------------------- 3:14s")
}

func displayNotifications(notifs []string) {
	moveCursor(pos{maxX / 2, int(float32(maxY)/1.25) - 1})
	fmt.Fprintln(screen, "NOTIFICATIONS")
	for i, v := range notifs {
		if i > 4 {
			break
		}
		moveCursor(pos{maxX / 2, int(float32(maxY)/1.25) + i})
		fmt.Fprintf(screen, " %s", stripString(v))
	}
}

func (tui TUI) DisplayStats(playedList []string, completedPlaylist int) {
	clear()
	showCursor()

	moveCursor(pos{2, 2})
	fmt.Fprintf(screen, "Played         : %d song(s).", len(playedList)+(len(tui.Playlist.List)*completedPlaylist))
	moveCursor(pos{2, 3})
	fmt.Fprintf(screen, "Played list    : %d time(s).", completedPlaylist)
	moveCursor(pos{2, 4})
	fmt.Fprintf(screen, "Minutes played : 21 minute(s)")

	render()
	os.Exit(0)
}
