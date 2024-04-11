package main

import (
	"log"
	"strconv"

	"github.com/mtslzr/pokeapi-go"
	"github.com/xuri/excelize/v2"
)

type Movement struct {
	moveName, movePow, moveAcc, moveClass, moveType, moveEffect string
}

func main() {

	moveFile, err := excelize.OpenFile("movements.xlsx")

	if err != nil {
		log.Println(err)
		panic(err)
	}

	allRows, err := moveFile.GetRows("movements")

	if err != nil {
		log.Println(err)
		panic(err)
	}

	var errors []error
	var moves []Movement

	for i, rows := range allRows {

		if i == 0 {
			continue
		}

		for _, move := range rows {

			m, err := pokeapi.Move(move)

			if err != nil {
				errors = append(errors, err)
				continue
			}

			var move Movement

			move.moveName = m.Name
			move.movePow = strconv.Itoa(m.Power)
			move.moveAcc = strconv.Itoa(m.Accuracy)
			move.moveType = m.Type.Name
			move.moveClass = m.DamageClass.Name

			moves = append(moves, move)
		}

	}

	//TODO: agregar logica para escribir los movimientos encontrados en un nuevo excel

}
