package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	HerramientasInodos "MIA_P1_201407049/InodoTools"
	"MIA_P1_201407049/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Mkfile(parametros []string) {
	fmt.Println(">> Ejecutando comando MKFILE")
	var (
		path string
		cont string
	)
	size := 0
	r := false
	if !Structs.CurrentUSR.STATUS {
		fmt.Println("Error la sesión no se encuentra iniciada")
		return
	}

	for _, raw := range parametros[1:] {
		temp2 := strings.TrimRight(raw, " ")
		temp := strings.Split(temp2, "=")
		if strings.ToLower(temp[0]) == "path" {
			if len(temp) != 2 {
				fmt.Println("Error valor no reconocido del parametro: ", temp[0])
				return
			}
			temp1 := strings.ReplaceAll(temp[1], "\"", "")
			path = temp1

		} else if strings.ToLower(temp[0]) == "size" {
			if len(temp) != 2 {
				fmt.Println("Error valor no reconocido del parametro:  ", temp[0])
				return
			}
			var err error
			size, err = strconv.Atoi(temp[1])
			if err != nil {
				fmt.Println("Error el tamaño solo admite valores enteros ", temp[1])
				return
			}
			if size < 0 {
				fmt.Println("Error el tamaño debe ser un valor positivo: ", temp[1])
				return
			}
		} else if strings.ToLower(temp[0]) == "Contenido" {
			if len(temp) != 2 {
				fmt.Println("Error valor no reconocido del parametro ", temp[0])
				return
			}
			temp1 := strings.ReplaceAll(temp[1], "\"", "")
			cont = temp1

			//R
		} else if strings.ToLower(temp[0]) == "r" {
			if len(temp) != 1 {
				fmt.Println("Error valor no reconocido del parametro ", temp[0])
				return
			}
			r = true

		} else {
			fmt.Println("Error parametro desconocido: ", temp[0])
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
			finRuta := len(stepPath) - 1
			idInicial := int32(0)
			idActual := int32(0)
			crear := -1
			for i, itemPath := range stepPath[1:finRuta] {
				idActual = HerramientasInodos.LookInodo(idInicial, "/"+itemPath, SuperBlk, disco)
				if idInicial != idActual {
					idInicial = idActual
				} else {
					crear = i + 1
					break
				}
			}
			if crear != -1 {
				if r {
					for _, item := range stepPath[crear:finRuta] {
						idInicial = HerramientasInodos.CRTFOLDER(idInicial, item, int64(mbr.Partitions[part].Start), disco)
						if idInicial == 0 {
							fmt.Println("Error imposible crear el directorio")
							return
						}
					}
				} else {
					fmt.Println("Error directorio ", stepPath[crear], " no registrado, se se permite crear sin directorio raiz")
					return
				}

			}
			idNuevo := HerramientasInodos.LookInodo(idInicial, "/"+stepPath[finRuta], SuperBlk, disco)
			if idNuevo == idInicial {
				fmt.Println("Error Archivo no registrado")
				if cont == "" {
					CreateFL(idInicial, stepPath[finRuta], size, "", int64(mbr.Partitions[part].Start), disco)
				}
			} else {
				fmt.Println("El archivo ya ha sido registrado")
			}

		}
	} else {
		fmt.Println("Error falta el parametro -path")
		fmt.Println("R ", r)
		fmt.Println("Cont ", cont)
	}
}

func CreateFL(idInodo int32, file string, size int, contenido string, initSuperBlk int64, disco *os.File) {
	var superB Structs.SuBlock
	Herramientas.ReadObj(disco, &superB, initSuperBlk)
	var inodoFile Structs.Inode
	Herramientas.ReadObj(disco, &inodoFile, int64(superB.SU_str_inode+(idInodo*int32(binary.Size(Structs.Inode{})))))
	for i := 0; i < 12; i++ {
		idBloque := inodoFile.In_blk[i]
		if idBloque != -1 {
			var FLDBLK Structs.FLDBLK
			Herramientas.ReadObj(disco, &FLDBLK, int64(superB.SU_str_blk+(idBloque*int32(binary.Size(Structs.FLDBLK{})))))
			for j := 2; j < 4; j++ {
				apuntador := FLDBLK.B_CONT[j].B_inodo
				if apuntador == -1 {
					copy(FLDBLK.B_CONT[j].B_name[:], file)
					ino := superB.SU_fst_ino
					FLDBLK.B_CONT[j].B_inodo = ino
					Herramientas.WrObj(disco, FLDBLK, int64(superB.SU_str_blk+(idBloque*int32(binary.Size(Structs.FLDBLK{})))))
					var newInodo Structs.Inode
					newInodo.In_usrid = Structs.CurrentUSR.USRID
					newInodo.In_grpid = Structs.CurrentUSR.GRPID
					newInodo.In_size = int32(size)
					ahora := time.Now()
					date := ahora.Format("14/06/2025 12:00")
					copy(newInodo.In_readDate[:], date)
					copy(newInodo.In_crtdate[:], date)
					copy(newInodo.In_modDate[:], date)
					copy(newInodo.In_typ[:], "1")
					copy(newInodo.In_modDate[:], "664")

					for i := int32(0); i < 15; i++ {
						newInodo.In_blk[i] = -1
					}
					guardarContenido := ""
					if contenido == "" {
						digit := 0
						for i := 0; i < size; i++ {
							if digit == 10 {
								digit = 0
							}
							guardarContenido += strconv.Itoa(digit)
							digit++
						}
					}
					fileblock := superB.SU_fst_blk
					inicio := 0
					fin := 0
					sizeContenido := len(guardarContenido)
					if sizeContenido < 64 {
						fin = len(guardarContenido)
					} else {
						fin = 64
					}
					for i := int32(0); i < 12; i++ {
						newInodo.In_blk[i] = fileblock
						data := guardarContenido[inicio:fin]
						var newFileBlock Structs.Fileblock
						copy(newFileBlock.B_CONT[:], []byte(data))
						Herramientas.WrObj(disco, newFileBlock, int64(superB.SU_str_blk+(fileblock*int32(binary.Size(Structs.Fileblock{})))))
						superB.SU_Free_Blk -= 1
						superB.SU_fst_blk += 1
						Herramientas.WrObj(disco, byte(1), int64(superB.SU_btp_str_blk+fileblock))
						calc := len(guardarContenido[fin:])
						if calc > 64 {
							inicio = fin
							fin += 64
						} else if calc > 0 {
							inicio = fin
							fin += calc
						} else {
							break
						}
						fileblock++
					}
					Herramientas.WrObj(disco, newInodo, int64(superB.SU_str_inode+(ino*int32(binary.Size(Structs.Inode{})))))
					superB.SU_Free_Inodo -= 1
					superB.SU_fst_ino += 1
					Herramientas.WrObj(disco, superB, initSuperBlk)
					Herramientas.WrObj(disco, byte(1), int64(superB.SU_btp_str_ino+ino))

					return
				}
			}
		} else {
			block := superB.SU_fst_blk
			inodoFile.In_blk[i] = block
			Herramientas.WrObj(disco, &inodoFile, int64(superB.SU_str_inode+(idInodo*int32(binary.Size(Structs.Inode{})))))
			var FLDBLK Structs.FLDBLK
			bloque := inodoFile.In_blk[0]
			Herramientas.ReadObj(disco, &FLDBLK, int64(superB.SU_str_blk+(bloque*int32(binary.Size(Structs.FLDBLK{})))))
			var newFLDBLK1 Structs.FLDBLK
			newFLDBLK1.B_CONT[0].B_inodo = FLDBLK.B_CONT[0].B_inodo
			copy(newFLDBLK1.B_CONT[0].B_name[:], ".")
			newFLDBLK1.B_CONT[1].B_inodo = FLDBLK.B_CONT[1].B_inodo
			copy(newFLDBLK1.B_CONT[1].B_name[:], "..")
			ino := superB.SU_fst_ino
			newFLDBLK1.B_CONT[2].B_inodo = ino
			copy(newFLDBLK1.B_CONT[2].B_name[:], file)
			newFLDBLK1.B_CONT[3].B_inodo = -1
			Herramientas.WrObj(disco, newFLDBLK1, int64(superB.SU_str_blk+(block*int32(binary.Size(Structs.FLDBLK{})))))
			Herramientas.WrObj(disco, byte(1), int64(superB.SU_btp_str_blk+block))
			superB.SU_fst_blk += 1
			superB.SU_Free_Blk -= 1
			var newInodo Structs.Inode
			newInodo.In_usrid = Structs.CurrentUSR.USRID
			newInodo.In_grpid = Structs.CurrentUSR.GRPID
			newInodo.In_size = int32(size)
			ahora := time.Now()
			date := ahora.Format("14/06/2025 12:00")
			copy(newInodo.In_readDate[:], date)
			copy(newInodo.In_crtdate[:], date)
			copy(newInodo.In_modDate[:], date)
			copy(newInodo.In_typ[:], "1")
			copy(newInodo.In_modDate[:], "664")

			for i := int32(0); i < 15; i++ {
				newInodo.In_blk[i] = -1
			}
			guardarContenido := ""
			if contenido == "" {
				digit := 0
				for i := 0; i < size; i++ {
					if digit == 10 {
						digit = 0
					}
					guardarContenido += strconv.Itoa(digit)
					digit++
				}
			}
			fileblock := superB.SU_fst_blk

			inicio := 0
			fin := 0
			sizeContenido := len(guardarContenido)
			if sizeContenido < 64 {
				fin = len(guardarContenido)
			} else {
				fin = 64
			}
			for i := int32(0); i < 12; i++ {
				newInodo.In_blk[i] = fileblock
				data := guardarContenido[inicio:fin]
				var newFileBlock Structs.Fileblock
				copy(newFileBlock.B_CONT[:], []byte(data))
				Herramientas.WrObj(disco, newFileBlock, int64(superB.SU_str_blk+(fileblock*int32(binary.Size(Structs.Fileblock{})))))
				superB.SU_Free_Blk -= 1
				superB.SU_fst_blk += 1
				Herramientas.WrObj(disco, byte(1), int64(superB.SU_btp_str_blk+fileblock))
				calc := len(guardarContenido[fin:])
				if calc > 64 {
					inicio = fin
					fin += 64
				} else if calc > 0 {
					inicio = fin
					fin += calc
				} else {
					break
				}
				fileblock++
			}

			Herramientas.WrObj(disco, newInodo, int64(superB.SU_str_inode+(ino*int32(binary.Size(Structs.Inode{})))))
			superB.SU_Free_Inodo -= 1
			superB.SU_fst_ino += 1
			Herramientas.WrObj(disco, superB, initSuperBlk)
			Herramientas.WrObj(disco, byte(1), int64(superB.SU_btp_str_ino+ino))

			return

		}
	}

}
