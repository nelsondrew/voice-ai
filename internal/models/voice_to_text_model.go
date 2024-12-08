package models

type VoiceToTextModel struct {
	AudioFilePath string `json:"audio_file_path"`
	Text          string `json:"text"`
}
