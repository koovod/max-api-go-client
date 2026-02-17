package maxapi

import "encoding/json"

// BotInfo describes the current bot.
type BotInfo struct {
	UserID           int64        `json:"user_id"`
	FirstName        string       `json:"first_name"`
	LastName         string       `json:"last_name,omitempty"`
	Username         string       `json:"username,omitempty"`
	IsBot            bool         `json:"is_bot"`
	LastActivityTime int64        `json:"last_activity_time,omitempty"`
	Name             string       `json:"name,omitempty"`
	Description      string       `json:"description,omitempty"`
	AvatarURL        string       `json:"avatar_url,omitempty"`
	FullAvatarURL    string       `json:"full_avatar_url,omitempty"`
	Commands         []BotCommand `json:"commands,omitempty"`
}

// BotCommand describes a single supported slash command.
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description,omitempty"`
}

// Image describes a small media descriptor returned in chat info.
type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// Chat describes a MAX chat (group or dialog).
type Chat struct {
	ChatID            int64             `json:"chat_id"`
	Type              string            `json:"type"`
	Status            string            `json:"status"`
	Title             string            `json:"title,omitempty"`
	Icon              *Image            `json:"icon,omitempty"`
	LastEventTime     int64             `json:"last_event_time,omitempty"`
	Participants      []ChatParticipant `json:"participants,omitempty"`
	ParticipantsCount int               `json:"participants_count,omitempty"`
	OwnerID           *int64            `json:"owner_id,omitempty"`
	IsPublic          bool              `json:"is_public"`
	Link              string            `json:"link,omitempty"`
	Description       string            `json:"description,omitempty"`
	DialogWithUser    *UserWithPhoto    `json:"dialog_with_user,omitempty"`
	ChatMessageID     string            `json:"chat_message_id,omitempty"`
	PinnedMessage     *Message          `json:"pinned_message,omitempty"`
}

// ChatParticipant provides a compact representation of chat users in list responses.
type ChatParticipant struct {
	UserID         int64  `json:"user_id"`
	LastAccessTime int64  `json:"last_access_time,omitempty"`
	IsOwner        bool   `json:"is_owner"`
	IsAdmin        bool   `json:"is_admin"`
	JoinTime       int64  `json:"join_time,omitempty"`
	Alias          string `json:"alias,omitempty"`
}

// User represents either a bot or a human user.
type User struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	IsBot     bool   `json:"is_bot"`
}

// UserWithPhoto extends User with avatar URLs.
type UserWithPhoto struct {
	User
	AvatarURL     string `json:"avatar_url,omitempty"`
	FullAvatarURL string `json:"full_avatar_url,omitempty"`
}

// ChatMember describes permissions of a user within a chat.
type ChatMember struct {
	User
	LastActivityTime int64                 `json:"last_activity_time,omitempty"`
	LastAccessTime   int64                 `json:"last_access_time,omitempty"`
	IsOwner          bool                  `json:"is_owner"`
	IsAdmin          bool                  `json:"is_admin"`
	JoinTime         int64                 `json:"join_time,omitempty"`
	Permissions      []ChatAdminPermission `json:"permissions,omitempty"`
	Alias            string                `json:"alias,omitempty"`
}

// ChatAdminPermission enumerates admin rights.
type ChatAdminPermission string

// ChatsPage contains a page of chats with pagination marker.
type ChatsPage struct {
	Chats  []Chat `json:"chats"`
	Marker *int64 `json:"marker,omitempty"`
}

// ChatMembersPage contains chat members with pagination marker.
type ChatMembersPage struct {
	Members []ChatMember `json:"members"`
	Marker  *int64       `json:"marker,omitempty"`
}

// ChatActionRequest triggers specific behaviors (typing_on, mark_seen, etc).
type ChatActionRequest struct {
	Action string `json:"action"`
}

// PinMessageRequest pins a message inside a chat.
type PinMessageRequest struct {
	MessageID string `json:"message_id"`
	Notify    *bool  `json:"notify,omitempty"`
}

// SetAdminsRequest configures chat admins.
type SetAdminsRequest struct {
	Admins []ChatAdminAssignment `json:"admins"`
}

// ChatAdminAssignment pairs a user ID with permissions.
type ChatAdminAssignment struct {
	UserID      int64                 `json:"user_id"`
	Permissions []ChatAdminPermission `json:"permissions,omitempty"`
	Alias       string                `json:"alias,omitempty"`
}

// AddMembersRequest adds members by ID.
type AddMembersRequest struct {
	UserIDs []int64 `json:"user_ids"`
}

// Message represents a message payload from the API.
type Message struct {
	MessageID string         `json:"message_id"`
	Sender    *User          `json:"sender,omitempty"`
	Recipient *Recipient     `json:"recipient,omitempty"`
	Timestamp int64          `json:"timestamp,omitempty"`
	Link      *LinkedMessage `json:"link,omitempty"`
	Body      *MessageBody   `json:"body,omitempty"`
	Stat      *MessageStat   `json:"stat,omitempty"`
	URL       string         `json:"url,omitempty"`
}

// Recipient contains either chat or user reference.
type Recipient struct {
	Type   string `json:"type"`
	ChatID int64  `json:"chat_id,omitempty"`
	UserID int64  `json:"user_id,omitempty"`
}

// LinkedMessage references a replied/forwarded message.
type LinkedMessage struct {
	MessageID string       `json:"message_id"`
	Body      *MessageBody `json:"body,omitempty"`
	Sender    *User        `json:"sender,omitempty"`
}

// MessageBody describes the textual content and attachments.
type MessageBody struct {
	MidBody     bool           `json:"mid_body,omitempty"`
	Text        string         `json:"text,omitempty"`
	Link        *MessageLink   `json:"link,omitempty"`
	Format      *MessageFormat `json:"format,omitempty"`
	Attachments []Attachment   `json:"attachments,omitempty"`
	Notify      *bool          `json:"notify,omitempty"`
}

// NewMessageBody is used when sending or editing messages.
type NewMessageBody struct {
	Text        string         `json:"text,omitempty"`
	Link        *MessageLink   `json:"link,omitempty"`
	Format      *MessageFormat `json:"format,omitempty"`
	Attachments []Attachment   `json:"attachments,omitempty"`
	Notify      *bool          `json:"notify,omitempty"`
}

// Attachment is a generic attachment wrapper.
type Attachment struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// MessageLink describes inline buttons linking to messages.
type MessageLink struct {
	Text string `json:"text,omitempty"`
	URL  string `json:"url,omitempty"`
}

// MessageFormat contains formatting directives.
type MessageFormat struct {
	ParseMode string `json:"parse_mode,omitempty"`
}

// MessageStat contains aggregate metrics for a message.
type MessageStat struct {
	Views    int `json:"views,omitempty"`
	Forwards int `json:"forwards,omitempty"`
}

// MessageList is used by GET /messages.
type MessageList struct {
	Messages []Message `json:"messages"`
}

// Subscription represents a webhook subscription.
type Subscription struct {
	URL         string   `json:"url"`
	UpdateTypes []string `json:"update_types"`
	Secret      string   `json:"secret,omitempty"`
	Ports       []int    `json:"ports,omitempty"`
	CreatedAt   int64    `json:"created_at,omitempty"`
}

// SubscriptionsResponse wraps the subscriptions list.
type SubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

// UpdatePage contains updates with pagination.
type UpdatePage struct {
	Updates []Update `json:"updates"`
	Marker  *int64   `json:"marker,omitempty"`
}

// Update describes platform events (message_created, message_callback, etc).
type Update struct {
	UpdateID   int64           `json:"update_id"`
	UpdateType string          `json:"update_type"`
	Timestamp  int64           `json:"timestamp"`
	Message    *Message        `json:"message,omitempty"`
	Callback   *Callback       `json:"callback,omitempty"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}

// Callback is delivered when a button is pressed.
type Callback struct {
	CallbackID string          `json:"callback_id"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}

// UploadResponse describes resumable / multipart upload instructions.
type UploadResponse struct {
	URL   string `json:"url"`
	Token string `json:"token,omitempty"`
}

// SuccessResponse is returned by a variety of write endpoints.
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// VideoMetadata describes /videos/{videoToken} responses.
type VideoMetadata struct {
	Token     string           `json:"token"`
	URLs      VideoURLs        `json:"urls"`
	Thumbnail *PhotoAttachment `json:"thumbnail,omitempty"`
	Width     int              `json:"width,omitempty"`
	Height    int              `json:"height,omitempty"`
	Duration  int              `json:"duration,omitempty"`
}

// VideoURLs may contain multiple quality options.
type VideoURLs struct {
	HLS string `json:"hls,omitempty"`
	MP4 string `json:"mp4,omitempty"`
}

// PhotoAttachment describes an image payload (used for thumbnails).
type PhotoAttachment struct {
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Token  string `json:"token,omitempty"`
}
