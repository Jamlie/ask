/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/Jamlie/ask/internal/gemini"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// geminiCmd represents the question command
var geminiCmd = &cobra.Command{
	Use:   "gemini",
	Short: "Asks Gemini a question.",
	Long:  "Asks Gemini a question and gives the answer",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			slog.Error("question only receives one argument")
		}

		question := args[0]

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		viper.SetConfigName(".ask")
		viper.SetConfigType("toml")
		viper.AddConfigPath(homeDir)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
		}

		apiKey := viper.GetString("gemini_api")
		sendQuestion := gemini.Question(apiKey)
		res, err := sendQuestion(context.Background(), question)
		if err != nil {
			log.Fatal(err)
		}

		answer := res.Candidates[0].Content.Parts[0]

		unsave, _ := cmd.PersistentFlags().GetBool("no-csv")
		if !unsave {
			if err := saveToCSV([]string{question, answer.Text}); err != nil {
				log.Println("Unable to save record into CSV")
			}
		}

		fmt.Println(answer)
	},
}

func init() {
	rootCmd.AddCommand(geminiCmd)

	geminiCmd.PersistentFlags().BoolP("no-csv", "n", false, "Use to stop saving questions in a csv file")
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
