package controllers

type VoiceAssistantController struct {
	voiceToText *VoiceToTextController
	chatGPT     *ChatGPTController
}

func (v *VoiceAssistantController) ProcessVoiceInput(audioFilePath string) (string, error) {
	// Convert voice to text
	transcribedText, err := v.voiceToText.ConvertVoiceToText(audioFilePath)
	if err != nil {
		return "", err
	}

	// Process with ChatGPT
	response, err := v.chatGPT.ProcessConversation(transcribedText)
	if err != nil {
		return "", err
	}

	return response, nil
}
