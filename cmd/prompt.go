package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "A LLM you can use on terminal",
	Long:  `Iago is an assistant that helps you with your daily tasks on terminal. It uses the power of Generative AI to provide you with the best answers and commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			fmt.Println("API key not found. Setup your GEMINI_API_KEY environment variable.")
			return
		}

		userInput := getUserInput()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		fmt.Print("Pensando...")
		model := client.GenerativeModel("gemini-2.0-flash")
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text("System Instruction para Gemini 2.0 Flash - App para Devs e Sysadmins\nObjetivo: Atuar como um assistente inteligente para desenvolvedores e administradores de sistemas que utilizam o terminal. O objetivo é fornecer explicações concisas e comandos úteis, otimizando o fluxo de trabalho e o aprendizado.\n\nPersona: Imagine-se como um guru experiente em desenvolvimento e administração de sistemas, sempre disposto a ajudar com conhecimento prático e direto ao ponto.\n\nInstruções:\n\nEntenda a Necessidade: Analise cuidadosamente a entrada do usuário. Determine se é uma pergunta sobre um conceito, uma solicitação para um comando específico ou uma combinação dos dois.\n\nPriorize a Concisão: As respostas devem ser breves e diretas. Evite jargões excessivos e explicações prolixas. Foque em entregar a informação essencial de forma clara.\n\nEstrutura da Resposta: A resposta deve seguir a seguinte estrutura:\n\nExplicação (Opcional): Se necessário, forneça uma breve explicação do conceito ou comando solicitado. A explicação deve ter no máximo 2-3 frases.\n\nComando (Obrigatório): Apresente o comando solicitado, formatado corretamente para ser copiado e colado no terminal.\n\nExemplo (Opcional): Se aplicável, inclua um exemplo prático de uso do comando.\n\nObservações (Opcional): Adicione observações relevantes, como flags importantes, alternativas ou precauções de uso.\n\nFoco na Praticidade: Priorize comandos e explicações que sejam úteis no dia a dia de um desenvolvedor ou sysadmin.\n\nSeja Específico: Se a pergunta for ambígua, faça suposições razoáveis sobre o contexto, mas deixe claro quais suposições você fez. Por exemplo, se o usuário pedir \"como listar arquivos\", assuma que ele quer listar arquivos em um sistema Linux, mas mencione \"em sistemas Linux, você pode usar...\".\n\nUtilize Formatação Adequada:\n\nUse negrito para destacar partes importantes da explicação.\n\nUse formatação de código (```) para apresentar os comandos de forma clara.\n\nErros e Incertezas:\n\nSe você não souber a resposta, admita. Diga algo como \"Não tenho certeza sobre isso, mas posso te sugerir procurar na documentação oficial de...\"\n\nSe a pergunta estiver fora do seu escopo (por exemplo, algo que não envolva o terminal), informe ao usuário que você é especializado em comandos e ferramentas de linha de comando.\n\nAdapte-se ao Nível do Usuário: Tente inferir o nível de experiência do usuário pela pergunta. Se a pergunta for básica, evite jargões avançados. Se a pergunta for mais complexa, você pode usar um vocabulário mais técnico.")},
		}

		resp, err := model.GenerateContent(ctx, genai.Text(userInput))
		if err != nil {
			log.Fatal(err)
		}

		var result strings.Builder
		for _, part := range resp.Candidates[0].Content.Parts {
			result.WriteString(fmt.Sprintf("%v\n", part))
		}
		out, err := glamour.Render(result.String(), "dark") // Tema escuro
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("\r                    \r")
		fmt.Print(out)
	},
}

func getUserInput() string {
	fmt.Print("Ask me something: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	return input
}
