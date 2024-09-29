/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jamlie/ask/internal/color"
	"github.com/Jamlie/ask/internal/gemini"
	"github.com/charmbracelet/glamour"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// geminiCmd represents the question command
var geminiCmd = &cobra.Command{
	Use:        "gemini",
	Short:      "Asks Gemini a question.",
	Long:       "Asks Gemini a question and gives the answer",
	SuggestFor: []string{"gemni", "gemin"},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		initViper(homeDir)

		apiKey := viper.GetString("gemini_api")

		geminiAI := gemini.New(apiKey)
		defer geminiAI.Close()

		chat, _ := cmd.PersistentFlags().GetBool("chat")
		if chat {
			startChat(geminiAI)
			return
		}

		doc, _ := cmd.PersistentFlags().GetBool("doc")
		msg, _ := cmd.PersistentFlags().GetString("msg")
		path, _ := cmd.PersistentFlags().GetString("path")

		if !doc && (len(msg) != 0 || len(path) != 0) {
			fmt.Fprintln(os.Stderr, color.Error.String("cannot use msg and path without defining doc"))
			return
		}

		if doc {
			documentRenderer(ctx, geminiAI, msg, path)
			return
		}

		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, color.Error.String("gemini only receives one argument"))
			return
		}

		askQuestion(geminiAI, args[0])
	},
}

func init() {
	rootCmd.AddCommand(geminiCmd)

	geminiCmd.PersistentFlags().BoolP("chat", "c", false, "Used to start a chat with Gemini")
	geminiCmd.PersistentFlags().BoolP("doc", "d", false, "Used to indicate that a document will be sent")
	geminiCmd.PersistentFlags().StringP("msg", "m", "", "Used to send a message to after doc is defined Gemini")
	geminiCmd.PersistentFlags().StringP("path", "p", "", "Used to specify path")
}

func documentRenderer(ctx context.Context, geminiAI *gemini.Gemini, msg, path string) {
	if path == "" {
		fmt.Fprintln(os.Stderr, color.Error.String("doc's path cannot be empty"))
		return
	}

	if msg == "" {
		fmt.Fprintln(os.Stderr, color.Error.String("doc's msg cannot be empty"))
		return
	}

	if !isValidPath(path) {
		fmt.Fprintf(os.Stderr, color.Error.String("\"%s\" is not valid path\n"), path)
		return
	}

	res, err := geminiAI.Document(ctx, path, msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.Error.String(err.Error()))
		return
	}

	answer := fmt.Sprint(res.Candidates[0].Content.Parts[0])

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	mdAnswer, err := r.Render(answer)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.Error.String(err.Error()))
		return
	}
	fmt.Print(string(mdAnswer))
}

func askQuestion(geminiAI *gemini.Gemini, q string) {
	question := strings.TrimSpace(q)
	if len(question) == 0 {
		fmt.Fprintln(os.Stderr, color.Warn.String("Entered an empty input"))
		return
	}
	res, err := geminiAI.Question(context.Background(), question)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.Error.String(err.Error()))
		return
	}

	answer := fmt.Sprint(res.Candidates[0].Content.Parts[0])

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	mdAnswer, err := r.Render(answer)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.Error.String(err.Error()))
		return
	}
	fmt.Print(string(mdAnswer))
}

func startChat(geminiAI *gemini.Gemini) {
	chatBot := geminiAI.Chat()
	rl, err := readline.New(color.Success.String("> "))
	if err != nil {
		fmt.Fprintln(os.Stderr, color.Error.String(err.Error()))
	}
	defer rl.Close()

	interrupted := 0

	ctx := context.Background()

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	var inputBuilder strings.Builder
	isMultiline := false

	const (
		multilineBeginning = "@\""
		multilineEnding    = "\"@"
	)

	for {
		if isMultiline {
			rl.SetPrompt(color.Placeholder.String("... "))
		} else {
			rl.SetPrompt(color.Success.String("> "))
		}

		input, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			if interrupted == 0 {
				fmt.Println("(To exit, press Ctrl+C again or Ctrl+D)")
				interrupted++
				continue
			}
			return
		}
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, color.Error.String(err.Error()))
			return
		}

		if interrupted != 0 {
			interrupted--
		}

		input = strings.TrimSpace(input)

		if len(input) == 0 && !isMultiline {
			fmt.Fprintln(os.Stderr, color.Warn.String("Entered an empty input"))
			continue
		}

		if strings.HasPrefix(input, multilineBeginning) {
			isMultiline = true
			inputBuilder.WriteString(strings.TrimPrefix(input, multilineBeginning))
			continue
		}

		if isMultiline && (strings.HasSuffix(input, multilineEnding) || input == multilineEnding) {
			inputBuilder.WriteString(strings.TrimSuffix(input, multilineEnding))
			isMultiline = false
		} else if isMultiline {
			inputBuilder.WriteString(input)
			continue
		} else {
			inputBuilder.WriteString(input)
		}

		fullInput := inputBuilder.String()
		fullInput = strings.TrimSpace(fullInput)

		if len(fullInput) == 0 {
			fmt.Fprintln(os.Stderr, color.Warn.String("Entered an empty input"))
			continue
		}

		inputBuilder.Reset()

		res, err := chatBot(ctx, fullInput)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.Error.String(err.Error()))
			return
		}

		answer := fmt.Sprint(res.Candidates[0].Content.Parts[0])

		mdAnswer, err := r.Render(answer)
		if err != nil {
			fmt.Fprintln(os.Stderr, color.Error.String(err.Error()))
			return
		}
		fmt.Print(string(mdAnswer))
	}
}

func initViper(homeDir string) {
	viper.SetConfigName(".ask")
	viper.SetConfigType("toml")
	viper.AddConfigPath(homeDir)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, color.Error.String("Cannot read .ask.toml"))
		return
	}
}

func isNewFile(f *os.File) (bool, error) {
	fileInfo, err := f.Stat()
	if err != nil {
		return false, err
	}

	return fileInfo.Size() == 0, nil
}

func isValidPath(path string) bool {
	cleanedPath := filepath.Clean(path)

	_, err := os.Stat(cleanedPath)
	return !os.IsNotExist(err)
}
