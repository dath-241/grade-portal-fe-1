// cmd/api/other.go
package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"Grade_Portal_TelegramBot/internal/handlers"
	"Grade_Portal_TelegramBot/config"
)

func Start() {
	// Dùng hàm LoadConfig() định nghĩa trong file config.go - package config.
	cfg := config.LoadConfig()
	// Khởi tạo bot
	bot, err := tgbotapi.NewBotAPI(cfg.BOT_TOKEN)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err) // In lỗi chi tiết
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	// Cấu hình để nhận cập nhật
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	// Nhận cập nhật từ bot
	fmt.Println("Starting bot...")
	updates := bot.GetUpdatesChan(u)
	fmt.Println("Listening for updates...")

	// Vòng lặp để xử lý các bản cập nhật
	for update := range updates {
		if update.Message != nil {
			handlers.HandleUpdate(bot, update)
		}
	}
}