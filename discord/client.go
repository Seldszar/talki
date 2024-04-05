package discord

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Channel struct {
	ID   string `json:"id"`
	Type uint8  `json:"type"`
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	Bot           bool   `json:"bot"`
}

type VoiceState struct {
	Deaf     bool `json:"deaf"`
	Mute     bool `json:"mute"`
	SelfDeaf bool `json:"self_deaf"`
	SelfMute bool `json:"self_mute"`
	Suppress bool `json:"suppress"`
}

type AuthorizeArgs struct {
	ClientID string   `json:"client_id"`
	Scopes   []string `json:"scopes"`
	Prompt   string   `json:"prompt"`
}

type VoiceChannelSelectData struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

type VoiceStateData struct {
	Nick       string     `json:"nick"`
	Mute       bool       `json:"mute"`
	Volume     float64    `json:"volume"`
	VoiceState VoiceState `json:"voice_state"`
	User       User       `json:"user"`
}

type GetChannelData struct {
	Channel

	VoiceStates []VoiceStateData `json:"voice_states"`
}

type AuthenticateArgs struct {
	AccessToken string `json:"access_token"`
}

type SpeakingData struct {
	UserID string `json:"user_id"`
}

type SpeakingArgs struct {
	ChannelID string `json:"channel_id"`
}

type VoiceStateArgs struct {
	ChannelID string `json:"channel_id"`
}

type AuthorizeData struct {
	Code string `json:"code"`
}

type Request[T any] struct {
	Cmd   string `json:"cmd"`
	Nonce string `json:"nonce"`
	Evt   string `json:"evt"`
	Args  T      `json:"args"`
}

type Response struct {
	Cmd   string          `json:"cmd"`
	Nonce string          `json:"nonce"`
	Event string          `json:"evt"`
	Data  json.RawMessage `json:"data"`
}

func (r *Response) UnmarshalData(v any) error {
	return json.Unmarshal(r.Data, v)
}

func fetchAccessToken(code string) (string, error) {
	body, err := json.Marshal(map[string]string{
		"code": code,
	})

	if err != nil {
		return "", err
	}

	res, err := http.Post("https://streamkit.discord.com/overlay/token", "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", err
	}

	var data struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.AccessToken, nil
}

type Client struct {
	conn *websocket.Conn
}

func (c *Client) Request(cmd, evt string, args any) error {
	return c.conn.WriteJSON(Request[any]{cmd, uuid.NewString(), evt, args})
}

func (c *Client) Read(res *Response) error {
	if err := c.conn.ReadJSON(res); err != nil {
		return err
	}

	switch res.Cmd {
	case "AUTHENTICATE":
		c.Request("SUBSCRIBE", "VOICE_CHANNEL_SELECT", nil)
		c.Request("GET_SELECTED_VOICE_CHANNEL", "", nil)

	case "AUTHORIZE":
		var data AuthorizeData

		if err := res.UnmarshalData(&data); err != nil {
			return err
		}

		token, err := fetchAccessToken(data.Code)

		if err != nil {
			return err
		}

		c.Request("AUTHENTICATE", "", AuthenticateArgs{
			AccessToken: token,
		})

	case "DISPATCH":
		switch res.Event {
		case "READY":
			c.Request("AUTHORIZE", "", AuthorizeArgs{
				ClientID: "207646673902501888",
				Scopes:   []string{"rpc"},
				Prompt:   "none",
			})
		}
	}

	return nil
}

func NewClient() (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:6463?v=1&client_id=207646673902501888", http.Header{
		"Origin": []string{"https://streamkit.discord.com"},
	})

	if err != nil {
		return nil, err
	}

	return &Client{conn}, nil
}
