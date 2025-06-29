package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
	"encoding/binary"
	"fmt"
	"strings"
)

func Rmgrp(parametros []string) {
	fmt.Println(">> Procesando Comando RMGRP")
	var name string
	temp2 := strings.TrimRight(parametros[1], " ")
	temp := strings.Split(temp2, "=")

	if len(temp) != 2 {
		fmt.Println("Error no se reconoce el valor del parametro", temp[0])
		return
	}

	if strings.ToLower(temp[0]) == "name" {
		name = temp[1]
	} else {
		fmt.Println("Error parametro no reconocido", temp[0])
		return
	}
	CurrentUSR := Structs.CurrentUSR
	if CurrentUSR.STATUS {
		if CurrentUSR.Nombre == "root" {
			disk := CurrentUSR.Id[0:1]
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
			exec := false
			index := -1
			for i := 0; i < 4; i++ {
				identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
				if identificador == CurrentUSR.Id {
					exec = true
					index = i
					break
				}
			}

			if exec {
				var SuperBlk Structs.SuBlock
				err = Herramientas.ReadObj(disco, &SuperBlk, int64(mbr.Partitions[index].Start))
				if err != nil {
					fmt.Println("Error la partición no tiene formato")
					return
				}
				var inodo Structs.Inode
				Herramientas.ReadObj(disco, &inodo, int64(SuperBlk.SU_str_inode+int32(binary.Size(Structs.Inode{}))))
				var contenido string
				var fileBlock Structs.Fileblock
				for _, item := range inodo.In_blk {
					if item != -1 {
						Herramientas.ReadObj(disco, &fileBlock, int64(SuperBlk.SU_str_blk+(item*int32(binary.Size(Structs.Fileblock{})))))
						contenido += string(fileBlock.B_CONT[:])
					}
				}

				lineID := strings.Split(contenido, "\n")
				mod := false
				for i, reg := range lineID[:len(lineID)-1] {
					datos := strings.Split(reg, ",")
					if len(datos) == 3 {
						if datos[2] == name {
							if datos[0] != "0" {
								mod = true
								datos[0] = "0"
								mod := datos[0] + "," + datos[1] + "," + datos[2]
								lineID[i] = mod
							} else {
								fmt.Println("Error el grupo ya ha sido eliminado")
							}
							break
						}
					}
				}
				if mod {
					for i, reg := range lineID[:len(lineID)-1] {
						datos := strings.Split(reg, ",")
						if len(datos) == 5 {
							if datos[2] == name {
								if datos[0] != "0" {
									datos[0] = "0"
									mod := datos[0] + "," + datos[1] + "," + datos[2] + "," + datos[3] + "," + datos[4]
									lineID[i] = mod
								}
							}
						}
					}
				}

				if mod {
					mod := ""
					for _, reg := range lineID {
						mod += reg + "\n"
					}

					inicio := 0
					var fin int
					if len(mod) > 64 {
						fin = 64
					} else {
						fin = len(mod)
					}

					for _, newItem := range inodo.In_blk {
						if newItem != -1 {
							data := mod[inicio:fin]
							var newFileBlock Structs.Fileblock
							copy(newFileBlock.B_CONT[:], []byte(data))
							Herramientas.WrObj(disco, newFileBlock, int64(SuperBlk.SU_str_blk+(newItem*int32(binary.Size(Structs.Fileblock{})))))
							inicio = fin
							calc := len(mod[fin:])
							if calc > 64 {
								fin += 64
							} else {
								fin += calc
							}
						}
					}
				}

			} else {
				fmt.Println("Error detectado verifique el ID de la partición")
			}

		} else {
			fmt.Println("Error el usuario no cumple con los permisos minimos")
		}
	} else {
		fmt.Println("Error la sesión no fue iniciada")
	}
}
