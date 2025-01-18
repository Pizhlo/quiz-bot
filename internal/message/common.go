package message

const StartMessage = "Нажми кнопку, чтобы начать квиз:"
const HelpMessage = "Это квиз-бот. Для того, чтобы пройти квиз, нажми /start и действуй по инструкции"

const Question = "Вопрос %d / %d\n\n<b>%s</b>"

const LevelEnd = "Раунд закончен!\n\nПравильных ответов: %d/%d"

var ErrorMessageChannelMessage = "#ошибки\n⚠️В чате с пользователем произошла ошибка!\n\nКоманда:\n<code>%+v</code>\n\nОшибка:\n<code>%+v</code>"
var ErrorMessageUser = "😔Во время обработки произошла ошибка, но я уже сообщил разработчику! Повтори запрос позднее"
