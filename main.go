package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/mtslzr/pokeapi-go"
	"github.com/xuri/excelize/v2"
	"google.golang.org/api/option"
)

type Movement struct {
	moveName, movePow, moveAcc, moveClass, moveType, moveEffect string
}

//var translator *translate.TranslationClient

// Dont forget to close the file
var movesFile *excelize.File
var newFile *excelize.File

func init() {

	// t, err := getTranslator()

	// if err != nil {
	// 	log.Println(err)
	// }

	// translator = t

	mf, err := excelize.OpenFile("movements.xlsx")

	if err != nil {
		log.Println(err)
		panic(err)
	}

	movesFile = mf

}

func main() {

	//defer translator.Close()
	defer movesFile.Close()

	allRows, err := movesFile.GetRows("movements")

	if err != nil {
		log.Println(err)
		panic(err)
	}

	var readingErrors []error

	//TODO: agregar errores de traduccion para separarlos del resto
	//var translatingErrors []error

	var moves []Movement

	for i, rows := range allRows {

		if i == 0 {
			continue
		}

		for _, move := range rows {

			m, err := pokeapi.Move(move)

			if err != nil {
				log.Printf("Err in move %v: %v", move, err)
				readingErrors = append(readingErrors, err)
				continue
			}

			var move Movement

			//TODO: agregar traduccion a los campos correspondientes (traducir a español y traducir Pow/Acc en formato D10)
			move.moveName = m.Names[5].Name
			move.movePow = strconv.Itoa(m.Power)
			move.moveAcc = strconv.Itoa(m.Accuracy)
			move.moveType = m.Type.Name
			move.moveClass = m.DamageClass.Name
			move.moveEffect = m.EffectEntries[0].Effect

			moves = append(moves, move)
		}

	}

	if len(readingErrors) > 0 {
		log.Println(readingErrors)
	}

	var excelColumns []string = []string{"A", "B", "C", "D", "E", "F"}
	var writingErrors []error

	newFile = excelize.NewFile()

	var appendWritingErrors = func(err error) {

		if err != nil {
			writingErrors = append(writingErrors, err)
		}

	}

	for i, move := range moves {

		for j, columnLetter := range excelColumns {

			coords := fmt.Sprintf("%v%v", columnLetter, i+1)

			switch {

			case j == 0:
				appendWritingErrors(newFile.SetCellValue("Sheet1", coords, move.moveName))

			case j == 1:
				appendWritingErrors(newFile.SetCellValue("Sheet1", coords, move.movePow))

			case j == 2:
				appendWritingErrors(newFile.SetCellValue("Sheet1", coords, move.moveAcc))

			case j == 3:
				appendWritingErrors(newFile.SetCellValue("Sheet1", coords, move.moveType))

			case j == 4:
				appendWritingErrors(newFile.SetCellValue("Sheet1", coords, move.moveClass))

			case j == 5:
				appendWritingErrors(newFile.SetCellValue("Sheet1", coords, move.moveEffect))
			}

		}

	}

	if len(writingErrors) > 0 {
		log.Printf("%v errors writing", writingErrors)
		log.Println(writingErrors)
	}

	if len(readingErrors) > 0 {
		log.Printf("%v errors reading", readingErrors)
		log.Println(readingErrors)
	}

	err = newFile.SaveAs("movements_data.xlsx")

	if err != nil {
		panic(err)
	}

	//TODO: testear

}

// Dont forget to close the translator at the end with defer translator.Close()
func getTranslator() (*translate.TranslationClient, error) {

	ctx := context.Background()

	//TODO: habilitar billing en google cloud para habilitar el translator
	c, err := translate.NewTranslationClient(ctx, option.WithCredentialsFile("service-account.json"))

	if err != nil {
		log.Printf("err creating new translator: %v", err)
		return nil, err
	}

	return c, nil

}

func translateToSpanish(ctx context.Context, translator *translate.TranslationClient) (string, error) {

	req := &translatepb.TranslateTextRequest{
		Contents:           []string{"hi"},
		SourceLanguageCode: "en",
		TargetLanguageCode: "spa",
	}
	resp, err := translator.TranslateText(ctx, req)

	if err != nil {
		log.Printf("err translating to spanish: %v", err)
		return "", err
	}
	// TODO: Use resp.
	return resp.String(), nil

}
