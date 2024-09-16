package gemini

import (
	"context"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Gemini struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func New(apiKey string) *Gemini {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}

	return &Gemini{
		client: client,
		model:  model,
	}
}

func (g *Gemini) Close() error {
	return g.client.Close()
}

func (g *Gemini) Question(
	ctx context.Context,
	question string,
) (*genai.GenerateContentResponse, error) {
	resp, err := g.model.GenerateContent(ctx, genai.Text(question))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *Gemini) Chat() func(context.Context, string) (*genai.GenerateContentResponse, error) {
	chat := g.model.StartChat()

	return func(ctx context.Context, input string) (*genai.GenerateContentResponse, error) {
		res, err := chat.SendMessage(ctx, genai.Text(input))
		if err != nil {
			return nil, err
		}
		chat.History = append(chat.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(input),
			},
			Role: "user",
		})
		chat.History = append(chat.History, res.Candidates[0].Content)
		return res, nil
	}
}
