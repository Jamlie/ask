package gemini

type partRequest struct {
	Text string `json:"text,omitempty"`
}

func (pr partRequest) String() string {
	return pr.Text
}

type contentRequest struct {
	Parts []partRequest `json:"parts,omitempty"`
}

type Request struct {
	Contents []contentRequest `json:"contents,omitempty"`
}

type partResponse struct {
	Text string `json:"text,omitempty"`
}

func (pr partResponse) String() string {
	return pr.Text
}

type contentResponse struct {
	Parts []partResponse `json:"parts,omitempty"`
	Role  string         `json:"role,omitempty"`
}

type safetyRating struct {
	Category    string `json:"category,omitempty"`
	Probability string `json:"probability,omitempty"`
}

type candidate struct {
	Content       contentResponse `json:"content,omitempty"`
	FinishReason  string          `json:"finishReason,omitempty"`
	Index         int             `json:"index,omitempty"`
	SafetyRatings []safetyRating  `json:"safetyRatings,omitempty"`
}

type Response struct {
	Candidates    []candidate   `json:"candidates"`
	UsageMetadata usageMetadata `json:"usageMetadata"`
}

type usageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount,omitempty"`
	CandidatesTokenCount int `json:"candidatesTokenCount,omitempty"`
	TotalTokenCount      int `json:"totalTokenCount,omitempty"`
}
