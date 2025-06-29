package Comandos

import (
	"fmt"
	"os"
	"strings"
)

func Rmdisk(parametros []string) {
	fmt.Println(">> Ejecutando comando RMDISK")
	temp2 := strings.TrimRight(parametros[1], " ")
	temp := strings.Split(temp2, "=")

	if len(temp) != 2 {
		fmt.Println("Error valor del parametro  ", temp[0], " no reconocido")
		return
	}

	if strings.ToLower(temp[0]) == "driveletter" {
		letter := strings.ToUpper(temp[1])
		folder := "./MIA/P1/"
		ext := ".dsk"
		path := folder + string(letter) + ext

		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			fmt.Println("Erro el disco ", letter, " no fue localizado")
			return
		}
		fmt.Printf("¿Precaución esta acción eliminara el disco %s  desea continuar? (y/n): ", letter)
		var respuesta string
		fmt.Scanln(&respuesta)

		respuesta = strings.ToLower(respuesta)

		if respuesta == "y" || respuesta == "si" {
			err2 := os.Remove(path)
			if err2 != nil {
				fmt.Println("Error imposible remover el disco")
				return
			}
			fmt.Println("Disco ", letter, "fue elimindado de manera exitosa: ")
		} else {
			fmt.Println("Precaución la operación ha sido cancelada")
		}

	} else {
		fmt.Println("Error paramatro ", temp[0], " no reconocido")
	}
}
