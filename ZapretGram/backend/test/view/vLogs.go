package view

import (
	ethModel "ZapretGram/backend/Core/ethernet/Model"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Цвета, да, раскрашиваем всякую хуйню
func colloredStrings(Type string, str string) string {
	//Цвета
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	//white := color.New(color.BgHiWhite).SprintFunc()

	//str = blue(str)

	//=========================================================
	switch strings.ToLower(Type) {
	case "status":
		//=====================================================
		switch strings.ToLower(str) {
		case "ok":
			str = "[" + green("Status") + "]: " + green("OK") + "\n"

		case "error", "err":
			str = "[" + red("Status") + "]: " + red("ERROR") + "\n"

		default:
			str = "[" + yellow("OTHER") + "]: " + str + "\n"
		}
		//=====================================================

	//=========================================================
	case "time":
		str = "[" + blue(Type) + "]: " + str + "\n"

	case "message", "msg":
		str = "[" + blue(Type) + "]: " + str + "\n"

	case "type":
		str = "[" + yellow("Type") + "]: " + str + "\n"

	case "nil":
		str = "[" + red("This nil") + "]: " + str + "\n"

	default:
		str = "[" + Type + "]: " + str + "\n"
	}
	//=========================================================

	return str
}

// time
func currtime() string {
	t := time.Now().Format("15:04:05")
	return fmt.Sprintf("%s", colloredStrings("Time", t))
}

// status
func status(Status string) string {
	var status string
	status = (colloredStrings("Status", Status))
	return status
}

func OtherdataTcp(message *string, data interface{}) {
	val := reflect.ValueOf(data)

	// Если не структура — выходим
	if val.Kind() != reflect.Struct {
		*message += fmt.Sprintf("[Data is not a struct]: %v\n", val.Interface())
		return
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Если поле — структура, рекурсивно
		if fieldValue.Kind() == reflect.Struct {
			// Проверяем, не нулевая ли структура
			zero := reflect.Zero(fieldValue.Type()).Interface()
			if reflect.DeepEqual(fieldValue.Interface(), zero) {
				*message += fmt.Sprintf("[%s]: <empty>", field.Name)
			} else {
				OtherdataTcp(message, fieldValue.Interface())
			}
			continue
		}

		// Преобразуем значение в строку
		valueStr := fmt.Sprintf("%v", fieldValue.Interface())

		// Если поле пустое — отмечаем <empty>
		if valueStr == "" {
			valueStr = "<empty>"
		}

		colored := colloredStrings(field.Name, valueStr)
		*message += colored
	}
}

// Общая функция для TCP структур
func respAndRequTcp(message *string, tcp interface{}) {
	val := reflect.ValueOf(tcp)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Основные поля
	if f := val.FieldByName("Status"); f.IsValid() {
		*message += status(f.String())
	}
	for _, fieldName := range []string{"DateTime", "CorrId"} {
		if f := val.FieldByName(fieldName); f.IsValid() {
			*message += colloredStrings(fieldName, f.String())
		}
	}

	// Остальные данные
	*message += "\n\n---Otherdata---\n"
	if f := val.FieldByName("Data"); f.IsValid() {
		OtherdataTcp(message, f.Interface())
	}
}

func LogTcp(Data ...interface{}) {
	Message := "\n======================================================\n"
	Message += currtime()

	for _, v := range Data {
		if v == nil {
			continue
		}

		Message += colloredStrings("Type", reflect.TypeOf(v).String())

		switch val := v.(type) {
		case string:
			Message += status(val)

		case *ethModel.RequestTcp, *ethModel.ResponseTcp:
			respAndRequTcp(&Message, val)

		default:
			fmt.Println("Неизвестный тип:", reflect.TypeOf(val))
		}
	}

	Message += "\n======================================================\n"
	fmt.Printf(Message)
}
