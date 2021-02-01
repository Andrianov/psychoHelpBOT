package questionnaire

import "github.com/Andrianov/psychoHelpBOT/internal/models"

var FlowSteps = []models.Step{
	{
		Name:     "👨 ФИО",
		Question: "Как к вам можно обращаться?",
	},
	{
		Name:     "🔞 Совершеннолетний?",
		Question: "Являетесь ли вы соверешеннолетним?",
		Options:  []string{"да", "нет"},
	},
	{
		Name:     "🧠 Состояние",
		Question: "Опишите свое состояние сейчас/последние несколько дней",
	},
	{
		Name:     "📞 Контакты",
		Question: "Контактный номер/никнейм телеграма, чтобы психолог мог с вами связаться",
	},
	{
		Name:     "💬 Способ работы",
		Question: "Укажите комфортный для вас способ работы",
		Options:  []string{"Skype", "Telegram", "WhatsApp"},
	},
}
