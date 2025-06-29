package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	HerramientasInodos "MIA_P1_201407049/InodoTools"
	"MIA_P1_201407049/Structs"
	"fmt"
	"strings"
)

func Mkdir(parametros []string) {
	fmt.Println(">> Ejecuando Comando MKDIR")
	var (
		path string
		r    bool = false
	)
	if !Structs.CurrentUSR.STATUS {
		fmt.Println("Error no se ha iniciado sesión")
		return
	}

	for _, raw := range parametros[1:] {
		temp2 := strings.TrimRight(raw, " ")
		temp := strings.Split(temp2, "=")
		if strings.ToLower(temp[0]) == "path" {

			if len(temp) != 2 {
				fmt.Println("Erro valor del parametro ", temp[0], "no reconocido")
				return
			}
			temp1 := strings.ReplaceAll(temp[1], "\"", "")
			path = temp1

		} else if strings.ToLower(temp[0]) == "r" {
			if len(temp) != 1 {
				fmt.Println("Error valor no reconocido del parametro ", temp[0])
				return
			}
			r = true

		} else {
			fmt.Println("Error parametro ", temp[0], " no reconocido")
			return
		}
	}

	if path != "" {
		id := Structs.CurrentUSR.Id
		disk := id[0:1]
		folder := "./MIA/P1/"
		ext := ".dsk"
		dirDisk := folder + disk + ext
		disco, err := Herramientas.OpenFile(dirDisk)
		if err != nil {
			return
		}

		var mbr Structs.MBR
		if err := Herramientas.ReadObj(disco, &mbr, 0); err != nil {
			return
		}

		defer disco.Close()
		buscar := false
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				buscar = true
				part = i
				break
			}
		}

		if buscar {
			var SuperBlk Structs.SuBlock

			err := Herramientas.ReadObj(disco, &SuperBlk, int64(mbr.Partitions[part].Start))
			if err != nil {
				fmt.Println("Error la partición no tiene formato")
			}
			stepPath := strings.Split(path, "/")
			strID := int32(0)
			CurrentID := int32(0)
			crear := -1
			for i, itemPath := range stepPath[1:] {
				CurrentID = HerramientasInodos.LookInodo(strID, "/"+itemPath, SuperBlk, disco)
				if strID != CurrentID {
					strID = CurrentID
				} else {
					crear = i + 1
					break
				}
			}

			if crear != -1 {
				if crear == len(stepPath)-1 {
					HerramientasInodos.CRTFOLDER(strID, stepPath[crear], int64(mbr.Partitions[part].Start), disco)
				} else {
					if r {
						for _, item := range stepPath[crear:] {
							strID = HerramientasInodos.CRTFOLDER(strID, item, int64(mbr.Partitions[part].Start), disco)
							if strID == 0 {
								fmt.Println("Error imposible crear el directorio")
								return
							}
						}
					} else {
						fmt.Println("No tiene permisos suficientes para crear el directorio raíz")
					}
				}
			} else {
				fmt.Println("Error el directorio ya se encuentra registrado")
			}
		}
	} else {
		fmt.Println("Error falta el parametro path")
		fmt.Println("R ", r)
	}
}
