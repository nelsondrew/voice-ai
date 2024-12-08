package interfaces

type VoiceToTextInterface interface {
	ConvertVoiceToText(audioFilePath string) (string, error)
}
