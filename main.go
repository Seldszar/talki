package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"fyne.io/systray"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v2"

	"github.com/seldszar/talki/autorun"
	"github.com/seldszar/talki/collection"
	"github.com/seldszar/talki/discord"
)

type H = map[string]any

type Speaker struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
	Speaking    bool   `json:"speaking"`
	Deaf        bool   `json:"deaf"`
	Mute        bool   `json:"mute"`
}

var (
	channelID = ""

	voiceStates = collection.NewMap[string, discord.VoiceStateData]()
	speaking    = collection.NewMap[string, bool]()

	upgrader = websocket.Upgrader{
		ReadBufferSize: 0,

		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients = collection.NewSet[*websocket.Conn]()

	//go:embed all:public/dist
	publicFS embed.FS

	//go:embed icon.ico
	iconBytes []byte
)

func openURL(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).
			Start()

	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).
			Start()

	case "darwin":
		return exec.Command("open", url).
			Start()
	}

	return fmt.Errorf("cannot open url %s on this platform", url)
}

func manageEvents(client *discord.Client, cmd, channelID string) {
	if channelID == "" {
		return
	}

	client.Request(cmd, "VOICE_STATE_CREATE", discord.VoiceStateArgs{
		ChannelID: channelID,
	})

	client.Request(cmd, "VOICE_STATE_DELETE", discord.VoiceStateArgs{
		ChannelID: channelID,
	})

	client.Request(cmd, "VOICE_STATE_UPDATE", discord.VoiceStateArgs{
		ChannelID: channelID,
	})

	client.Request(cmd, "SPEAKING_START", discord.SpeakingArgs{
		ChannelID: channelID,
	})

	client.Request(cmd, "SPEAKING_STOP", discord.SpeakingArgs{
		ChannelID: channelID,
	})
}

func getSpeakers() []*Speaker {
	speakers := make([]*Speaker, 0)

	voiceStates.Each(func(k string, v discord.VoiceStateData) bool {
		speakers = append(speakers, &Speaker{
			ID:          v.User.ID,
			Name:        v.User.Username,
			DisplayName: v.Nick,
			AvatarURL:   fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.webp", v.User.ID, v.User.Avatar),
			Deaf:        v.VoiceState.Deaf || v.VoiceState.SelfDeaf,
			Mute:        v.VoiceState.Mute || v.VoiceState.SelfMute,
			Speaking:    speaking.Get(k),
		})

		return false
	})

	return speakers
}

func updateSpeakers(fn func()) {
	fn()

	clients.Each(func(c *websocket.Conn) bool {
		c.WriteJSON(H{
			"speakers": getSpeakers(),
		})

		return false
	})
}

func connectDiscord() error {
	client, err := discord.NewClient()

	if err != nil {
		return err
	}

	slog.Info("Connected to Discord client")

	for {
		var res discord.Response

		if err := client.Read(&res); err != nil {
			return err
		}

		switch res.Cmd {
		case "DISPATCH":
			switch res.Event {
			case "SPEAKING_START":
				var data discord.SpeakingData

				if err := res.UnmarshalData(&data); err != nil {
					return err
				}

				updateSpeakers(func() {
					speaking.Set(data.UserID, true)
				})

			case "SPEAKING_STOP":
				var data discord.SpeakingData

				if err := res.UnmarshalData(&data); err != nil {
					return err
				}

				updateSpeakers(func() {
					speaking.Set(data.UserID, false)
				})

			case "VOICE_STATE_CREATE", "VOICE_STATE_UPDATE":
				var data discord.VoiceStateData

				if err := res.UnmarshalData(&data); err != nil {
					return err
				}

				updateSpeakers(func() {
					voiceStates.Set(data.User.ID, data)
				})

			case "VOICE_STATE_DELETE":
				var data discord.VoiceStateData

				if err := res.UnmarshalData(&data); err != nil {
					return err
				}

				updateSpeakers(func() {
					voiceStates.Delete(data.User.ID)
				})

			case "VOICE_CHANNEL_SELECT":
				client.Request("GET_SELECTED_VOICE_CHANNEL", "", nil)
			}

		case "GET_SELECTED_VOICE_CHANNEL":
			var data discord.GetChannelData

			if err := res.UnmarshalData(&data); err != nil {
				return err
			}

			manageEvents(client, "UNSUBSCRIBE", channelID)

			updateSpeakers(func() {
				channelID = data.ID

				voiceStates.Clear()
				speaking.Clear()

				for _, v := range data.VoiceStates {
					voiceStates.Set(v.User.ID, v)
				}
			})

			manageEvents(client, "SUBSCRIBE", channelID)
		}
	}
}

func loopDiscord() error {
	for {
		err := connectDiscord()

		if err != nil {
			slog.Error(
				"An error occured with Discord client",
				"error", err.Error(),
			)
		}

		slog.Info("Disconnected from Discord, reconnecting in 3 seconds...")
		time.Sleep(3 * time.Second)
	}
}

func loopTrayIcon(ctx *cli.Context) error {
	ex, err := os.Executable()

	if err != nil {
		return err
	}

	ar := &autorun.AutoRun{
		Name:        "Talki",
		DisplayName: "Talki: Custom Discord Widget",
		Executable:  ex,
	}

	systray.Run(
		func() {
			systray.SetTitle("Talki: Custom Discord Widget")
			systray.SetTooltip("Talki: Custom Discord Widget")
			systray.SetIcon(iconBytes)

			openItem := systray.AddMenuItem("Open Widget Page", "")
			autorunItem := systray.AddMenuItemCheckbox("Start at Launch", "", ar.IsEnabled())
			closeItem := systray.AddMenuItem("Close", "")

			for {
				select {
				case <-openItem.ClickedCh:
					openURL(fmt.Sprintf("http://localhost:%d", ctx.Int("port")))

				case <-autorunItem.ClickedCh:
					if ar.IsEnabled() {
						if err := ar.Disable(); err == nil {
							autorunItem.Uncheck()
						}

						break
					}

					if err := ar.Enable(); err == nil {
						autorunItem.Check()
					}

				case <-closeItem.ClickedCh:
					systray.Quit()
				}
			}
		},
		func() {
			os.Exit(0)
		},
	)

	return nil
}

func getPublicFS(root string) (http.FileSystem, error) {
	if root == "" {
		fs, err := fs.Sub(publicFS, "public/dist")

		if err != nil {
			return nil, err
		}

		return http.FS(fs), nil
	}

	return http.Dir(root), nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Usage: "Port the server is listening",
				Value: 9786,
			},
			&cli.StringFlag{
				Name:  "public",
				Usage: "Optional path to the folder serving web assets",
			},
		},
		Action: func(ctx *cli.Context) error {
			fs, err := getPublicFS(ctx.String("public"))

			if err != nil {
				return err
			}

			http.Handle("/", http.FileServer(fs))
			http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
				conn, err := upgrader.Upgrade(w, r, nil)

				if err != nil {
					return
				}

				ctx := r.Context()

				defer clients.Delete(conn)
				defer conn.Close()

				clients.Add(conn)

				conn.WriteJSON(H{
					"speakers": getSpeakers(),
				})

				<-ctx.Done()
			})

			go loopDiscord()
			go loopTrayIcon(ctx)

			slog.Info(
				"Server is ready",
				"address", fmt.Sprintf("http://localhost:%d", ctx.Int("port")),
			)

			if err := http.ListenAndServe(fmt.Sprintf(":%d", ctx.Int("port")), nil); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error(
			"An error occured with the application",
			"error", err.Error(),
		)
	}
}
