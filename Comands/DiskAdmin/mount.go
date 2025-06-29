package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Mount(parametros []string) {
	fmt.Println(">> Ejecutando Comando MOUNT")
	var (
		letter  string
		name    string
		paramOk bool = true
	)
	for _, raw := range parametros[1:] {
		temp2 := strings.TrimRight(raw, " ")
		temp := strings.Split(temp2, "=")

		if len(temp) != 2 {
			fmt.Println("Valor no reconocido del parametro ", temp[0])
			paramOk = false
			return
		}

		if strings.ToLower(temp[0]) == "driveletter" {
			letter = strings.ToUpper(temp[1])
			folder := "./MIA/P1/"
			ext := ".dsk"
			path := folder + string(letter) + ext
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("Error disco ", letter, " no registrado")
				paramOk = false
				break
			}
		} else if strings.ToLower(temp[0]) == "name" {
			name = strings.ReplaceAll(temp[1], "\"", "")
			name = strings.TrimSpace(name)
		} else {
			fmt.Println("Error parametro ", temp[0], "no reconocido")
			paramOk = false
			break
		}
	}

	if paramOk {
		if letter != "" && name != "" {
			filepath := "./MIA/P1/" + letter + ".dsk"
			disco, err := Herramientas.OpenFile(filepath)
			if err != nil {
				fmt.Println("Error imposible leer el disco, intente nuevamente")
				return
			}
			var mbr Structs.MBR
			if err := Herramientas.ReadObj(disco, &mbr, 0); err != nil {
				return
			}
			defer disco.Close()

			mount := true
			reportar := false
			for i := 0; i < 4; i++ {
				nombre := Structs.GETNOM(string(mbr.Partitions[i].Name[:]))
				if nombre == name {
					mount = false
					if string(mbr.Partitions[i].Status[:]) != "A" {
						if string(mbr.Partitions[i].Type[:]) != "E" {
							id := strings.ToUpper(letter) + strconv.Itoa(i+1) + "49"
							copy(mbr.Partitions[i].Status[:], "A")
							copy(mbr.Partitions[i].Id[:], id)
							if err := Herramientas.WrObj(disco, mbr, 0); err != nil {
								return
							}
							reportar = true
							fmt.Println("Partición ", name, " montada de manera exitosa")
						} else {
							fmt.Println("Error imposible montar la partición")
						}
					} else {
						fmt.Println("Error la paritición ya ha sido montada")
					}
					break
				}
			}

			if mount {
				fmt.Println("Error partición ", name, " imposible de montar")
				fmt.Println("Error partición no encontrada")
			}

			if reportar {
				fmt.Println("\nLista de particiones en el sistema\n ")
				for i := 0; i < 4; i++ {
					estado := string(mbr.Partitions[i].Status[:])
					if estado == "A" {
						fmt.Printf("Partition %d: name: %s, status: %s, id: %s, tipo: %s, start: %d, size: %d, fit: %s, correlativo: %d\n", i, string(mbr.Partitions[i].Name[:]), string(mbr.Partitions[i].Status[:]), string(mbr.Partitions[i].Id[:]), string(mbr.Partitions[i].Type[:]), mbr.Partitions[i].Start, mbr.Partitions[i].Size, string(mbr.Partitions[i].Fit[:]), mbr.Partitions[i].Correlativo)
					}
				}
			}
		} else {
			fmt.Println("Error no se encontro la letra de asignación ")
		}
	}
}
