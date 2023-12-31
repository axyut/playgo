package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

var screen = bufio.NewWriter(os.Stdout)

func hideCursor() {
	fmt.Fprint(screen, "\033[?25l")
}

func showCursor() {
	fmt.Fprint(screen, "\033[?25h")
}
func moveCursor(pos [2]int) {
	fmt.Fprintf(screen, "\033[%d;%dH", pos[1], pos[0])
}

func clear() {
	fmt.Fprint(screen, "\033[2J")
}

func draw(str string) {
	fmt.Fprint(screen, str)
}

// write all data in buffer to terminal
func render() {
	screen.Flush()
}

func termSize() (width int, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	return width, height
}

func border() {
	maxX, maxY := termSize()
	// ---- top ---- bottom
	for i := 1; i <= maxX; i++ {
		fmt.Fprintf(screen, "\033[1;%dH", i)
		fmt.Fprintf(screen, "_")
		fmt.Fprintf(screen, "\033[%d;%dH", maxY, i)
		fmt.Fprintf(screen, "_")
	}
	// | left | right |
	for i := 1; i <= maxY; i++ {
		fmt.Fprintf(screen, "\033[%d;1H", i)
		fmt.Fprintf(screen, "|")
		fmt.Fprintf(screen, "\033[%d;%dH", i, maxX)
		fmt.Fprintf(screen, "|")
	}
}

func seprator() {
	maxX, maxY := termSize()
	for i := 1; i <= maxY; i++ {
		fmt.Fprintf(screen, "\033[%d;%dH", i, (maxX/2)-2)
		fmt.Fprintf(screen, "|")
	}
}

func (p *Player) display() {
	clear()
	hideCursor()
	// border()
	seprator()
	maxX, maxY := termSize()

	// Playlist
	// TODO: will make playlist auto scroll one step when 1 songs finishes
	// so that prev 5 and next 5 will always be in place then section NEXT SONGS can be removed
	moveCursor(pos{3, 1})
	fmt.Fprintf(screen, "PLAYLIST (%d songs)", len(playlist))
	for i, v := range playlist {
		stripped := stripString(v)
		if i > maxY/(3) {
			moveCursor(pos{2, i + 3})
			fmt.Fprintf(screen, "  %d more Songs...", len(playlist)-i)
			moveCursor(pos{2, i + 4})
			fmt.Fprintf(screen, "%d. %s", len(playlist), stripString(playlist[len(playlist)-1]))
			break
		}
		moveCursor(pos{2, i + 3})
		fmt.Fprintf(screen, "%d. %s", i+1, stripped)
	}

	// Settings
	intH := int(float32(maxY) / 1.25)
	moveCursor(pos{2, intH - 1})
	fmt.Fprintf(screen, "SETTINGS")
	moveCursor(pos{3, intH})
	fmt.Fprintf(screen, "Shuffle: %t", UserSetting.Shuffle)
	moveCursor(pos{3, intH + 1})
	fmt.Fprintf(screen, "Repeat Song: %t", UserSetting.RepeatSong)
	moveCursor((pos{3, intH + 2}))
	fmt.Fprintf(screen, "Repeat playlist: %t", UserSetting.RepeatPlaylist)

	// Now Playing
	moveCursor(pos{maxX / 2, 1})
	fmt.Fprintf(screen, "NOW PLAYING")
	moveCursor(pos{maxX / 2, 3})
	fmt.Fprintf(screen, "%s", stripString(p.File.Name()))
	moveCursor(pos{maxX / 2, 4})
	fmt.Fprintf(screen, "%d:00 -------------------- 3:14s", timer)
	// song info
	// seek info

	// Next / Prev Song
	moveCursor(pos{maxX / 2, maxY / 4})
	fmt.Fprintf(screen, "UPCOMING SONGS")
	for i, _ := range playlist {
		if i == songs.currentSong {
			for j := 1; j <= 5; j++ {
				next := songs.currentSong + j
				if next >= len(playlist) {
					for next >= len(playlist) {
						next = next - len(playlist)
					}
				}
				moveCursor(pos{maxX / 2, (maxY / 4) + j})
				fmt.Fprintf(screen, "%s", stripString(playlist[next]))
			}
		}
		continue
	}

	// Notification
	moveCursor(pos{maxX / 2, int(float32(maxY)/1.25) - 1})
	fmt.Fprintln(screen, "NOTIFICATIONS")
	for i, v := range notifications {
		if i > 4 {
			break
		}
		moveCursor(pos{maxX / 2, int(float32(maxY)/1.25) + i})
		fmt.Fprintf(screen, " %s", stripString(v))
	}

	render()
}

func displayStats() {
	clear()
	showCursor()

	// moveCursor(pos{1, 1})
	// fmt.Fprintf(screen, "Stats :")
	moveCursor(pos{2, 2})
	fmt.Fprintf(screen, "Played         : %d song(s).", len(playedList)+(len(playlist)*completedPlaylist))
	moveCursor(pos{2, 3})
	fmt.Fprintf(screen, "Played list    : %d time(s).", completedPlaylist)
	moveCursor(pos{2, 4})
	fmt.Fprintf(screen, "Minutes played : 21 minute(s)")

	render()
	os.Exit(0)
}
