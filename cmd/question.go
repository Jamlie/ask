/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Jamlie/ask/internal/gemini"
	"github.com/Jamlie/ask/internal/logger"
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
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
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

		if len(args) != 1 {
			fmt.Println("gemini only receives one argument")
		}

		askQuestion(cmd, geminiAI, args[0])
	},
}

func init() {
	rootCmd.AddCommand(geminiCmd)

	geminiCmd.PersistentFlags().BoolP("no-csv", "n", false, "Use to stop saving questions in a csv file")
	geminiCmd.PersistentFlags().BoolP("chat", "c", false, "Used to start a chat with Gemini")
}

func askQuestion(cmd *cobra.Command, geminiAI *gemini.Gemini, question string) {
	if len(question) == 0 {
		fmt.Println(logger.Warn.String("Entered an empty input"))
		return
	}
	res, err := geminiAI.Question(context.Background(), question)
	if err != nil {
		log.Fatal(err)
	}

	answer := fmt.Sprint(res.Candidates[0].Content.Parts[0])

	unsave, _ := cmd.PersistentFlags().GetBool("no-csv")
	if !unsave {
		if err := saveToCSV([]string{question, answer}); err != nil {
			log.Println("Unable to save record into CSV")
		}
	}

	fmt.Println(answer)
}

func startChat(geminiAI *gemini.Gemini) {
	chatBot := geminiAI.Chat()
	rl, err := readline.New("> ")
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	interrupted := 0

	ctx := context.Background()

	for {
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
			log.Fatal(err)
		}

		interrupted--

		if len(input) == 0 {
			fmt.Println(logger.Warn.String("Entered an empty input"))
			continue
		}

		res, err := chatBot(ctx, input)
		if err != nil {
			log.Fatal(err)
		}

		answer := fmt.Sprint(res.Candidates[0].Content.Parts[0])
		fmt.Print(answer)
	}
}

func initViper(homeDir string) {
	viper.SetConfigName(".ask")
	viper.SetConfigType("toml")
	viper.AddConfigPath(homeDir)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(logger.Error.String("Cannot read .ask.toml"))
		os.Exit(0)
	}
}

func saveToCSV(record []string) error {
	csvFile, err := os.OpenFile("gemini.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	isNew, err := isNewFile(csvFile)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	if isNew {
		header := []string{"Question", "Answer"}
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	return writer.Write(record)
}

func isNewFile(f *os.File) (bool, error) {
	fileInfo, err := f.Stat()
	if err != nil {
		return false, err
	}

	return fileInfo.Size() == 0, nil
}
