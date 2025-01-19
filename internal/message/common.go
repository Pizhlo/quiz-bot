package message

const StartMessage = "Нажми кнопку, чтобы начать квиз:"
const HelpMessage = "Это квиз-бот. Для того, чтобы пройти квиз, нажми /start и действуй по инструкции"

const Question = "Вопрос %d / %d\n\n<b>%s</b>"

var ErrorMessageChannelMessage = "#ошибки\n⚠️В чате с пользователем произошла ошибка!\n\nКоманда:\n<code>%+v</code>\n\nОшибка:\n<code>%+v</code>"
var ErrorMessageUser = "😔Во время обработки произошла ошибка, но я уже сообщил разработчику! Повтори запрос позднее"

var ResultMessage = "<b>Викторина закончена!</b>🎉\n\nТвои результаты:\n\n%s"

var Result = "1️⃣Первый раунд - %d / %d\n2️⃣Второй раунд - %d / %d\n3️⃣Третий раунд - %d / %d\n\n⏱Время: %s"

var ChannelResultMessage = "Пользователь @%s закончил викторину!🎉\n\nРезультаты:\n\n%s"
