package handlers

import (
	"Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/services"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config) {
	resp, err := services.GetStudentInfo(update.Message.Chat.ID, cfg)
	if err != nil {
		response := "Không tìm thấy thông tin đăng nhập. Hãy đăng nhập trước khi sử dụng dịch vụ: " + err.Error()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	response, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	msgText := fmt.Sprintf("```json\n%s\n```", response)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func HandleGrade(bot *tgbotapi.BotAPI, update tgbotapi.Update, semesterOrCourseID string, cfg *config.Config) {
	resp, err := services.GetGrades(update.Message.Chat.ID, semesterOrCourseID, cfg)
	var response string

	if err != nil {
		response = "Không thể lấy dữ liệu điểm: " + err.Error()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	result := map[string]interface{}{
		"course_id":   semesterOrCourseID,
		"course_name": resp.Name,
		"scores": map[string]interface{}{
			"BT":  resp.Score.BT,
			"TN":  resp.Score.TN,
			"BTL": resp.Score.BTL,
			"GK":  resp.Score.GK,
			"CK":  resp.Score.CK,
		},
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		response = "Lỗi khi tạo JSON: " + err.Error()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	msgText := fmt.Sprintf("```json\n%s\n```", string(jsonData))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
	// } else {
	// 	response = fmt.Sprintf("Kết quả điểm cho %s:\n________\n%s:\n", semesterOrCourseID, resp.Name)

	// 	if resp.Score.BT != nil {
	// 		response += fmt.Sprintf("  - BT: %.1f\n", *resp.Score.BT)
	// 	} else {
	// 		response += "  - BT: null\n"
	// 	}

	// 	if resp.Score.TN != nil {
	// 		response += fmt.Sprintf("  - TN: %.1f\n", *resp.Score.TN)
	// 	} else {
	// 		response += "  - TN: null\n"
	// 	}

	// 	if resp.Score.BTL != nil {
	// 		response += fmt.Sprintf("  - BTL: %.1f\n", *resp.Score.BTL)
	// 	} else {
	// 		response += "  - BTL: null\n"
	// 	}

	// 	if resp.Score.GK != nil {
	// 		response += fmt.Sprintf("  - Giữa kỳ: %.1f\n", *resp.Score.GK)
	// 	} else {
	// 		response += "  - GK: null\n"
	// 	}

	// 	if resp.Score.CK != nil {
	// 		response += fmt.Sprintf("  - CK: %.1f\n", *resp.Score.CK)
	// 	} else {
	// 		response += "  - CK: null\n"
	// 	}
	// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	// 	bot.Send(msg)
	// }
}

func HandleAllGrade(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config) {
	resp, err := services.GetAllGrades(update.Message.Chat.ID, cfg)
	var response interface{}
	if err != nil {
		response = map[string]string{
			"error": "Không thể lấy dữ liệu điểm: " + err.Error(),
		}
	} else {
		var grades []map[string]interface{}
		for _, grade := range resp.AllGrades {
			gradeData := map[string]interface{}{
				"course_id":   grade.Ms,
				"course_name": grade.Name,
				"scores": map[string]interface{}{
					"BT":  grade.Score.BT,
					"TN":  grade.Score.TN,
					"BTL": grade.Score.BTL,
					"GK":  grade.Score.GK,
					"CK":  grade.Score.CK,
				},
			}
			grades = append(grades, gradeData)
		}
		response = map[string]interface{}{
			"grades": grades,
		}
	}
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		// Xử lý lỗi khi mã hóa JSON
		fmt.Println("Lỗi khi mã hóa JSON:", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Không thể xử lý dữ liệu JSON.")
		bot.Send(msg)
		return
	}
	fmt.Println("Dữ liệu JSON trả về:", string(responseJSON))
	// Gửi phản hồi dưới dạng JSON thô
	msgText := fmt.Sprintf("```json\n%s\n```", string(responseJSON))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(msgText))
	msg.ParseMode = "MarkdownV2" // Nếu bạn muốn hiển thị trong markdown
	bot.Send(msg)
	// response = "Kết quả điểm:\n________\n"
	// for _, grade := range resp.AllGrades {
	// 	response += fmt.Sprintf("Mã Môn: %s\nMôn: %s\n", grade.Ms, grade.Name)

	// 	if grade.Score.BT != nil {
	// 		response += fmt.Sprintf("  - BT: %.1f\n", *grade.Score.BT)
	// 	} else {
	// 		response += "  - BT: null\n"
	// 	}

	// 	if grade.Score.TN != nil {
	// 		response += fmt.Sprintf("  - TN: %.1f\n", *grade.Score.TN)
	// 	} else {
	// 		response += "  - TN: null\n"
	// 	}

	// 	if grade.Score.BTL != nil {
	// 		response += fmt.Sprintf("  - BTL: %.1f\n", *grade.Score.BTL)
	// 	} else {
	// 		response += "  - BTL: null\n"
	// 	}

	// 	if grade.Score.GK != nil {
	// 		response += fmt.Sprintf("  - Giữa kỳ: %.1f\n", *grade.Score.GK)
	// 	} else {
	// 		response += "  - GK: null\n"
	// 	}

	// 	if grade.Score.CK != nil {
	// 		response += fmt.Sprintf("  - CK: %.1f\n", *grade.Score.CK)
	// 	} else {
	// 		response += "  - CK: null\n"
	// 	}
	// 	response += "________\n"
	// }
}

// msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
// bot.Send(msg)

// msgText := fmt.Sprintf("```json\n%s\n```", response)
// msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
// msg.ParseMode = "MarkdownV2"
// bot.Send(msg)
