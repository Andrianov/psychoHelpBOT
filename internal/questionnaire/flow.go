package questionnaire

import "dev/tgbot/cmd/service/internal/models"

var FlowSteps = []models.Step{
	{
		Name:     "👨 ФИО",
		Question: "Как к вам можно обращаться?",
	},
	{
		Name:     "🔞 Совершеннолетний?",
		Question: "Подтверждаете ли вы, что вам 18 лет и более?",
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
