package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Mkdisk(parametros []string) {
	fmt.Println(">> Ejecutando Comando MKDISK")
	var (
		size     int
		fit      string = "F"
		unit     int    = 1048576
		paramOk  bool   = true
		sizeInit bool   = false
	)
	for _, raw := range parametros[1:] {
		temp2 := strings.TrimRight(raw, " ")
		temp := strings.Split(temp2, "=")
		if len(temp) != 2 {
			fmt.Println("Error valor de parametro ", temp[0], "no reconocido")
			paramOk = false
			break
		}
		if strings.ToLower(temp[0]) == "size" {
			sizeInit = true
			var err error
			size, err = strconv.Atoi(temp[1])
			if err != nil {
				fmt.Println("MError size debe tener valor numerico revisar: ", temp[1])
				paramOk = false
				break
			} else if size <= 0 {
				fmt.Println("Error size debe tener un valor positivo revisar:  ", temp[1])
				paramOk = false
				break
			}
		} else if strings.ToLower(temp[0]) == "fit" {
			if strings.ToLower(temp[1]) == "bf" {
				fit = "B"
			} else if strings.ToLower(temp[1]) == "wf" {
				fit = "W"
			} else if strings.ToLower(temp[1]) != "ff" {
				fmt.Println("Error -fit verifique los valores ingresados,  Revisar: ", temp[1])
				paramOk = false
				break
			}
		} else if strings.ToLower(temp[0]) == "unit" {
			if strings.ToLower(temp[1]) == "k" {
				unit = 1024
			} else if strings.ToLower(temp[1]) != "m" {
				fmt.Println("Error -unit verifique los valores ingresados Revisar: ", temp[1])
				paramOk = false
				break
			}
		} else {
			fmt.Println("Error parametro  ", temp[0], " no reconocido")
			paramOk = false
			break
		}
	}

	if paramOk {
		if sizeInit {
			tam := size * unit
			var path string
			var disco string
			folder := "./MIA/P1/"
			ext := ".dsk"
			for letra := 'A'; letra <= 'Z'; letra++ {
				path = folder + string(letra) + ext
				_, err := os.Stat(path)
				if os.IsNotExist(err) {
					disco = string(letra) + ext
					break
				}
			}
			err := Herramientas.CrearDisco(path)
			if err != nil {
				fmt.Println("Error:: ", err)
			}
			file, err := Herramientas.OpenFile(path)
			if err != nil {
				return
			}
			datos := make([]byte, tam)
			newErr := Herramientas.WrObj(file, datos, 0)
			if newErr != nil {
				fmt.Println("Error: ", newErr)
				return
			}
			ahora := time.Now()
			segundos := ahora.Second()
			minutos := ahora.Minute()
			cad := fmt.Sprintf("%02d%02d", segundos, minutos)
			idtemp, err := strconv.Atoi(cad)
			if err != nil {
				fmt.Println("Error imposible castear fecha en entero para id")
			}
			var newMBR Structs.MBR
			newMBR.MBRSZ = int32(tam)
			newMBR.Id = int32(idtemp)
			copy(newMBR.Fit[:], fit)
			copy(newMBR.DATECRT[:], ahora.Format("14/06/2024 12:00"))
			if err := Herramientas.WrObj(file, newMBR, 0); err != nil {
				return
			}
			defer file.Close()

			fmt.Println("\n Disco ", disco, " creado de forma exitosa")
			var TempMBR Structs.MBR
			if err := Herramientas.ReadObj(file, &TempMBR, 0); err != nil {
				return
			}
			Structs.PrintMBR(TempMBR)

			fmt.Println("\n======End MKDISK======")
		} else {
			fmt.Println("Error falta parametro size")
		}
	}
}
