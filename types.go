package telegram

type User struct {
	ID                      int    `json:"id"`
	ISBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Username                string `json:"username"`
	LanguageCode            string `json:"language_code"`
	IsPremium               bool   `json:"is_premium"`
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessages  bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
	CanManageTopics       bool `json:"can_manage_topics"`
}

type Location struct {
	Longitude            float64
	Latitude             float64
	HorizontalAccuracy   float64
	LivePeriod           int
	Heading              int
	ProximityAlertRadius int
}

type ChatLocation struct {
	Location *Location `json:"location"`
	Address  string    `json:"address"`
}

type Chat struct {
	ID                                 int              `json:"id"`
	Type                               string           `json:"type"` // Type of chat, can be either “private”, “group”, “supergroup” or “channel”
	Title                              string           `json:"title"`
	Username                           string           `json:"username"`
	FirstName                          string           `json:"first_name"`
	LastName                           string           `json:"last_name"`
	IsForum                            bool             `json:"is_forum"`
	Photo                              *ChatPhoto       `json:"chat_photo"`
	ActiveUsernames                    []string         `json:"active_usernames"`
	EmojisStatusCustomEmojiID          string           `json:"emoji_status_custom_emoji_id"`
	Bio                                string           `json:"bio"`
	HasPrivateForwards                 bool             `json:"has_private_forwards"`
	HasRestrictedVoiceAndVideoMessages bool             `json:"has_restricted_voice_and_video_messages"`
	JoinToSendMessages                 bool             `json:"join_to_send_messages"`
	JoinByRequest                      bool             `json:"join_by_request"`
	Description                        string           `json:"description"`
	InviteLink                         string           `json:"invite_link"`
	PinnedMessage                      *Message         `json:"pinned_message"`
	Permissions                        *ChatPermissions `json:"permissions"`
	SlowModeDelay                      int              `json:"slow_mode_delay"`
	MessageAutoDeleteTime              int              `json:"message_auto_delete_time"`
	HasProtectedContent                bool             `json:"has_protected_content"`
	StickerSetName                     string           `json:"sticker_set_name"`
	CanSetStickerSet                   bool             `json:"can_set_sticker_set"`
	LinkedChatID                       int              `json:"linked_chat_id"`
	Location                           *ChatLocation    `json:"location"`
}

type Message struct {
	MessageID       int             `json:"message_id"`
	Text            string          `json:"text"`
	From            *User           `json:"from"`
	Chat            *Chat           `json:"chat"`
	MessageEntities []MessageEntity `json:"entities"`
}

type CallbackQuery struct {
	ID              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageID string   `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
	GameShortName   string   `json:"game_short_name"`
}

type InlineQuery struct {
	ID       string    `json:"id"`
	From     *User     `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType string    `json:"chat_type"`
	Location *Location `json:"location"`
}

// Update represents data given from "getUpdates" query.
// Doc https://core.telegram.org/bots/api#getting-updates TODO fill the structure
type Update struct {
	UpdateID          int            `json:"update_id"`
	Message           *Message       `json:"message"`
	EditedMessage     *Message       `json:"edited_message"`
	ChannelPost       *Message       `json:"channel_post"`
	InlineQuery       *InlineQuery   `json:"inline_query"`
	EditedChannelPost *Message       `json:"edited_channel_post"`
	CallbackQuery     *CallbackQuery `json:"callback_query"`
}

type UpdateResponse struct {
	OK          bool     `json:"ok"`
	Description string   `json:"description"`
	Result      []Update `json:"result"`
}

type Response struct {
	OK          bool     `json:"ok"`
	ErrorCode   int      `json:"error_code"`
	Description string   `json:"description"`
	Result      *Message `json:"result"`
}

func (u *Update) hasMessageText() bool {
	return u.Message != nil && u.Message.Text != ""
}
