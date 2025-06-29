package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	HerramientasInodos "MIA_P1_201407049/InodoTools"
	"MIA_P1_201407049/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Cat(parametros []string) {
	fmt.Println(">> Ejecutando CAT")
	var file []string
	if !Structs.CurrentUSR.STATUS {
		fmt.Println("Error la sesión no se encuentra iniciada")
		return
	}

	for _, raw := range parametros[1:] {
		temp2 := strings.TrimRight(raw, " ")
		temp := strings.Split(temp2, "=")
		if len(temp) != 2 {
			fmt.Println("Error valor del parametro ", temp[0], " no identificado")
			return
		}

		if strings.ToLower(temp[0]) == "file" {
			temp1 := strings.ReplaceAll(temp[1], "\"", "")
			file = append(file, temp1)
		} else {
			comando := strings.Split(strings.ToLower(temp[0]), "file")
			if comando[0] == "file" {
				_, IDError := strconv.Atoi(comando[1])
				if IDError != nil {
					fmt.Println("Error imposible obtener numero de fichero.")
					return
				}
				temp1 := strings.ReplaceAll(temp[1], "\"", "")
				file = append(file, temp1)
			} else {
				fmt.Println("Error parametro ", temp[0], " no reconocido")
				return
			}
		}
	}

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
		var contenido string
		var fileBlock Structs.Fileblock

		err := Herramientas.ReadObj(disco, &SuperBlk, int64(mbr.Partitions[part].Start))
		if err != nil {
			fmt.Println("Error la partición no tiene formato")
		}

		for _, item := range file {
			idInodo := HerramientasInodos.LookInodo(0, item, SuperBlk, disco)
			var inodo Structs.Inode
			if idInodo > 0 {
				Herramientas.ReadObj(disco, &inodo, int64(SuperBlk.SU_str_inode+(idInodo*int32(binary.Size(Structs.Inode{})))))
				for _, idBlock := range inodo.In_blk {
					if idBlock != -1 {
						Herramientas.ReadObj(disco, &fileBlock, int64(SuperBlk.SU_str_blk+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
						contenido += string(fileBlock.B_CONT[:])
					}
				}
				contenido += "\n"
			} else {
				fmt.Println("Error el archivo no ha sido localizado ", item)
				return
			}
			fmt.Println(contenido)
		}
	}
}
