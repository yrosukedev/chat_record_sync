package wecom_chat

const (
	WeComMessageTypeText  = "text"
	WeComMessageTypeImage = "image"
	WeComMessageTypeVoice = "voice"
	WeComMessageTypeVideo = "video"
)

type WeComChatRecord struct {
	Seq     uint64   `json:"seq,omitempty"`
	MsgID   string   `json:"msgid,omitempty"`
	Action  string   `json:"action,omitempty"`
	From    string   `json:"from,omitempty"`
	ToList  []string `json:"tolist,omitempty"`
	RoomID  string   `json:"roomid,omitempty"`
	MsgTime int64    `json:"msgtime,omitempty"`
	MsgType string   `json:"msgtype,omitempty"`

	Text  *TextMessage  `json:"text,omitempty"`
	Image *ImageMessage `json:"image,omitempty"`
	Voice *VoiceMessage `json:"voice,omitempty"`
	Video *VideoMessage `json:"video,omitempty"`
}

type TextMessage struct {
	Content string `json:"content,omitempty"`
}

type ImageMessage struct {
	SdkFileID string `json:"sdkfileid,omitempty"`
	Md5Sum    string `json:"md5sum,omitempty"`
	FileSize  uint32 `json:"filesize,omitempty"`
}

type VoiceMessage struct {
	SdkFileID  string `json:"sdkfileid,omitempty"`
	VoiceSize  uint32 `json:"voice_size,omitempty"`
	PlayLength uint32 `json:"play_length,omitempty"`
	Md5Sum     string `json:"md5sum,omitempty"`
}

type VideoMessage struct {
	SdkFileID  string `json:"sdkfileid,omitempty"`
	FileSize   uint32 `json:"filesize,omitempty"`
	PlayLength uint32 `json:"play_length,omitempty"`
	Md5Sum     string `json:"md5sum,omitempty"`
}
