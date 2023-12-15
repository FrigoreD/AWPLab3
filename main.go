package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

// Данные для подключения к БД
const (
	dbHost     = "localhost"
	dbUser     = "root"
	dbPassword = "root"
	dbName     = "quiz_web_app"
)

// Question Структура вопроса
type Question struct {
	ID        int
	Text      string
	Options   []string
	CorrectID int
}

// QuizApp Структура опроса
type QuizApp struct {
	questions []Question
	current   int
	score     int
}

// Точка входа
func main() {
	appInstance := app.New()
	mainWindow := appInstance.NewWindow("Quiz App")
	mainWindow.Resize(fyne.NewSize(600, 400))

	contentContainer := container.NewMax()
	var db *sql.DB
	var err error

	// Устанавливаем соединение с базой данных
	if db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)); err != nil {
		log.Fatal(err)
	}

	q := getQuestions(db)

	quizList := widget.NewList(
		func() int {
			return len(q.questions)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewLabel("Question"))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			text.SetText(q.questions[id].Text)
		})

	quizList.OnSelected = func(id widget.ListItemID) {
		if q.current >= len(q.questions) {
			showResult(contentContainer, q)
			return
		}

		question := q.questions[id]
		options := widget.NewRadioGroup(question.Options, func(selected string) {
			selectedID := getSelectedOptionID(selected, question.Options)
			if selectedID == question.CorrectID {
				q.score++
			}
			q.current++
			if q.current < len(q.questions) {
				quizList.Select(q.current)
			} else {
				showResult(contentContainer, q)
			}
		})

		contentContainer.Objects = []fyne.CanvasObject{
			container.NewVBox(
				widget.NewLabel(wrapText(question.Text, 50)),
				container.NewCenter(options),
			),
		}
	}

	splitView := container.NewHSplit(quizList, contentContainer)
	splitView.Offset = 0.3
	mainWindow.SetContent(splitView)
	mainWindow.ShowAndRun()
}

// Функция для получения вопросов из базы данных
func getQuestions(db *sql.DB) QuizApp {
	rows, err := db.Query("SELECT question, answer1, answer2, correct_answer FROM questions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var quiz QuizApp

	for rows.Next() {
		var question, answer1, answer2 string
		var correct_answer int
		err := rows.Scan(&question, &answer1, &answer2, &correct_answer)
		if err != nil {
			log.Fatal(err)
		}

		questionQuiz := Question{
			Text:      question,
			Options:   []string{answer1, answer2},
			CorrectID: correct_answer,
		}

		quiz.questions = append(quiz.questions, questionQuiz)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return quiz
}

// Показ результата
func showResult(content *fyne.Container, quiz QuizApp) {
	resultLabel := widget.NewLabel(fmt.Sprintf("Quiz completed! Your score: %d/%d", quiz.score, len(quiz.questions)))
	restartButton := widget.NewButton("Retake Quiz", func() {
		quiz.current = 0
		quiz.score = 0
	})

	content.Objects = []fyne.CanvasObject{
		container.NewVBox(
			resultLabel,
			container.NewCenter(restartButton),
		),
	}
}

// Обёртка
func wrapText(text string, lineLength int) string {
	return wrap(text, lineLength)
}

// Обёртка
func wrap(s string, lineLength int) string {
	words := strings.Fields(s)
	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+len(word) < lineLength {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	lines = append(lines, currentLine)
	return strings.Join(lines, "\n")
}

// Получение выбранного ответа
func getSelectedOptionID(selected string, options []string) int {
	for i, option := range options {
		if option == selected {
			return i
		}
	}
	return -1
}
