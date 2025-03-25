package cmd

import (
	"context"
	"log"
	"starPivot/internal/service/chat"
)

func main() {
	ctx := context.Background()

	log.Println("create message ")
	messages := chat.CreateMessageFromTemplate()
	log.Println("message: ", messages)

	log.Println("create llm")
	chatModel := chat.CreateOPenAIChatModel(ctx)
	log.Println("create success")

	log.Printf("===llm generate===\n")
	result := chat.Generate(ctx, chatModel, messages)
	log.Printf("result: %+v\n\n", result)

	log.Printf("===llm stream generate===\n")
	streamResult := chat.Stream(ctx, chatModel, messages)
	chat.ReportStream(streamResult)

}
