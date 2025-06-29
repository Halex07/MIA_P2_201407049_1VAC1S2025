package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Mkgrp(parametros []string) {
	fmt.Println(">> Procesando comando MKGRP")
	var name string
	temp2 := strings.TrimRight(parametros[1], " ")
	temp := strings.Split(temp2, "=")
	if len(temp) != 2 {
		fmt.Println("Error valor no reconocido del parametro ", temp[0])
		return
	}

	if strings.ToLower(temp[0]) == "Nombre" {
		name = temp[1]
		if len(name) > 10 {
			fmt.Println("Error el nombre excede la cantidad maxima de 10 caracteres")
			return
		}
	} else {
		fmt.Println("Erorr parametro: ", temp[0], "no reconocido")
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

				var (
					inodo     Structs.Inode
					contenido string
					fileBlock Structs.Fileblock
					idFB      int32
				)
				Herramientas.ReadObj(disco, &inodo, int64(SuperBlk.SU_str_inode+int32(binary.Size(Structs.Inode{}))))
				for _, item := range inodo.In_blk {
					if item != -1 {
						Herramientas.ReadObj(disco, &fileBlock, int64(SuperBlk.SU_str_blk+(item*int32(binary.Size(Structs.Fileblock{})))))
						contenido += string(fileBlock.B_CONT[:])
						idFB = item
					}
				}
				lineID := strings.Split(contenido, "\n")
				for _, reg := range lineID[:len(lineID)-1] {
					datos := strings.Split(reg, ",")
					if len(datos) == 3 {
						if datos[2] == name {
							fmt.Println("Error el grupo ya se encuentra registrado ")
							return
						}
					}
				}

				id := -1
				var IDError error
				for i := len(lineID) - 2; i >= 0; i-- {
					reg := strings.Split(lineID[i], ",")
					if reg[1] == "G" {
						if reg[0] != "0" {
							id, IDError = strconv.Atoi(reg[0])
							if IDError != nil {
								fmt.Println("Error no es posible asignar ID al grupo")
								return
							}
							id++
							break
						}
					}
				}
				if id != -1 {
					CurrentCont := string(fileBlock.B_CONT[:])
					posNull := strings.IndexByte(CurrentCont, 0)
					data := fmt.Sprintf("%d,G,%s\n", id, name)
					if posNull != -1 {
						Freep := 64 - (posNull + len(data))
						if Freep > 0 {
							copy(fileBlock.B_CONT[posNull:], []byte(data))
							Herramientas.WrObj(disco, fileBlock, int64(SuperBlk.SU_str_blk+(idFB*int32(binary.Size(Structs.Fileblock{})))))
						} else {
							data1 := data[:len(data)+Freep]
							copy(fileBlock.B_CONT[posNull:], []byte(data1))
							Herramientas.WrObj(disco, fileBlock, int64(SuperBlk.SU_str_blk+(idFB*int32(binary.Size(Structs.Fileblock{})))))
							InfSave := true
							for i, item := range inodo.In_blk {
								if item == -1 {
									InfSave = false
									inodo.In_blk[i] = SuperBlk.SU_fst_blk
									SuperBlk.SU_Free_Blk -= 1
									SuperBlk.SU_fst_blk += 1
									data2 := data[len(data)+Freep:]
									var newFileBlock Structs.Fileblock
									copy(newFileBlock.B_CONT[:], []byte(data2))
									Herramientas.WrObj(disco, SuperBlk, int64(mbr.Partitions[index].Start))
									Herramientas.WrObj(disco, byte(1), int64(SuperBlk.SU_btp_str_blk+inodo.In_blk[i]))
									Herramientas.WrObj(disco, inodo, int64(SuperBlk.SU_str_inode+int32(binary.Size(Structs.Inode{}))))
									Herramientas.WrObj(disco, newFileBlock, int64(SuperBlk.SU_str_blk+(inodo.In_blk[i]*int32(binary.Size(Structs.Fileblock{})))))
									break
								}
							}

							if InfSave {
								fmt.Println("Error el espacio no es suficiente no se puede crear un nuevo registro")
							}
						}
					}
				}

			} else {
				fmt.Println("Error inseperado verifique el ID de la partición")
			}

		} else {
			fmt.Println("Error el usuario no cuenta con los permisos minimos")
		}
	} else {
		fmt.Println("Error la sesión no ha sido inicida ")
	}
}
