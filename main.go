package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/wendellast/Commit-IA/cmd/bot"
	"github.com/wendellast/Commit-IA/cmd/prompt"
	"github.com/wendellast/Commit-IA/cmd/typecmd"
)

func getProjectPath() (string, error) {
	projectPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("erro ao obter o diretório atual: %v", err)
	}
	return projectPath, nil
}

func runGitDiff(projectPath string) (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	cmd.Dir = projectPath

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("erro ao executar o comando git diff: %s", stderr.String())
	}

	return out.String(), nil
}

func executeCLI(outDiff string, url string) *cobra.Command {
	return &cobra.Command{
		Use:   "commitia",
		Short: "A CLI tool for handling commits",
		Run: func(cmd *cobra.Command, args []string) {
			description, err := cmd.Flags().GetString("description")

			if err != nil {
				fmt.Printf("Erro ao receber descrição do commit : %v\n", err)
				return
			}

			language, err := cmd.Flags().GetString("language")

			if err != nil {
				fmt.Printf("Erro ao receber idioma do commit : %v\n", err)
				return
			}

			if language == "" {
				language = "pt-br"
			}

			fmt.Println(language)

			commitMessage := prompt.CreateCommitMessage(outDiff, language, description)

			response, err := bot.SendMessageToBot(url, commitMessage)
			if err != nil {
				fmt.Printf("Erro ao enviar a mensagem: %v\n", err)
				return
			}

			err = typecmd.Typecmd(response.Response)

			if err != nil {
				fmt.Printf("Erro ao digitar o comando: %v\n", err)
				return
			}
		},
	}
}

func main() {
	projectPath, err := getProjectPath()
	if err != nil {
		fmt.Println(err)
		return
	}

	outDiff, err := runGitDiff(projectPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	llm := "https://hublast.com/gui-api/send-message-gui-commitia"

	rootCmd := executeCLI(outDiff, llm)
	rootCmd.Flags().StringP("description", "d", "", "Descrição básica do que fez no commit")
	rootCmd.Flags().StringP("language", "l", "", "Idioma ao qual deve ser escrito o commit")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Erro ao executar o comando CLI: %v\n", err)
		os.Exit(1)
	}
}
