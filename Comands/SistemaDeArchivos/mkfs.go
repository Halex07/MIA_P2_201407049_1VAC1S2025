package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func Mkfs(parametros []string) {
	fmt.Println(">> Procesando comando MKFS")
	var id string
	fs := "2fs"
	paramOk := true

	for _, raw := range parametros[1:] {
		temp2 := strings.TrimRight(raw, " ")
		temp := strings.Split(temp2, "=")
		if len(temp) != 2 {
			fmt.Println("Error valor no reconocido del parametro ", temp[0])
			paramOk = false
			break
		}

		if strings.ToLower(temp[0]) == "id" {
			id = strings.ToUpper(temp[1])

		} else if strings.ToLower(temp[0]) == "fs" {
			if strings.ToLower(temp[1]) == "3fs" {
				fs = "3fs"
			} else if strings.ToLower(temp[1]) != "2fs" {
				fmt.Println("Error verifique los valores para Ext2 o  Ext3: ", temp[1])
				paramOk = false
				break
			}

		} else if strings.ToLower(temp[0]) == "type" {
			if strings.ToLower(temp[1]) != "full" {
				fmt.Println("Error valor del -type no reconocido")
				paramOk = false
				break
			}

		} else {
			fmt.Println("Error parametro: ", temp[0], "no reconocido")
			paramOk = false
			break
		}
	}

	if paramOk {
		if id != "" {
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

			formatear := true
			for i := 0; i < 4; i++ {
				identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
				if identificador == id {
					formatear = false
					confirmar := true
					var newSuperBlk Structs.SuBlock
					Herramientas.ReadObj(disco, &newSuperBlk, int64(mbr.Partitions[i].Start))
					if newSuperBlk.SU_FileSys_ty != 0 {
						fmt.Printf("¿Esta seguro de proceder con el formateo de la particion %s? (y/n): ", Structs.GETNOM(string(mbr.Partitions[i].Name[:])))
						var respuesta string
						fmt.Scanln(&respuesta)
						respuesta = strings.ToLower(respuesta)
						// Validar la respuesta
						if respuesta == "y" || respuesta == "si" {
							fmt.Println("Ejectuando espere.... ")
						} else {
							confirmar = false
							fmt.Println("la acción ha sido cancelada no se ha completado el formateo de la partición", Structs.GETNOM(string(mbr.Partitions[i].Name[:])))
						}
					}
					if confirmar {
						numerator := int(mbr.Partitions[i].Size) - binary.Size(Structs.SuBlock{})
						denominator := 4 + binary.Size(Structs.Inode{}) + 3*binary.Size(Structs.Fileblock{})
						if fs == "3fs" {
							denominator += binary.Size(Structs.Journaling{})
						}
						n := int32(numerator / denominator)
						newSuperBlk.SU_CountBLK = int32(3 * n)
						newSuperBlk.SU_Free_Blk = int32(3 * n)
						newSuperBlk.SU_CountInodo = n
						newSuperBlk.SU_Free_Inodo = n
						newSuperBlk.SU_inode_size = int32(binary.Size(Structs.Inode{}))
						newSuperBlk.SU_blk_size = int32(binary.Size(Structs.Fileblock{}))
						ahora := time.Now()
						copy(newSuperBlk.SU_Date_mon[:], ahora.Format("14/06/2026 12:00"))
						copy(newSuperBlk.SU_Date_um[:], ahora.Format("14/06/2026 12:00"))
						newSuperBlk.SU_mont_sys += 1
						newSuperBlk.SU_SysF = 0xEF53

						if fs == "2fs" {
							createEXT2(n, mbr.Partitions[i], newSuperBlk, ahora.Format("14/06/2026 12:00"), disco)
						} else {
							createEXT3(n, mbr.Partitions[i], newSuperBlk, ahora.Format("14/06/2026 12:00"), disco)
						}
						fmt.Println("La partición", id, " ha sido formateada exitosamente")
						if Structs.CurrentUSR.STATUS {
							var new Structs.USRINF
							Structs.CurrentUSR = new
						}
					}
					break
				}
			}

			if formatear {
				fmt.Println("Error no es posible formatear la partición ", id)
				fmt.Println("Error ID no identificado")
			}
		} else {
			fmt.Println("Error no se encuentra el Id de la particón")
		}
	}
}

func createEXT2(n int32, particion Structs.Partition, newSuperBlk Structs.SuBlock, date string, file *os.File) {
	fmt.Println("====== CREAR EXT2 ======")
	fmt.Println("N: ", n)
	newSuperBlk.SU_FileSys_ty = 2
	newSuperBlk.SU_btp_str_ino = particion.Start + int32(binary.Size(Structs.SuBlock{}))
	newSuperBlk.SU_btp_str_blk = newSuperBlk.SU_btp_str_ino + n
	newSuperBlk.SU_str_inode = newSuperBlk.SU_btp_str_blk + 3*n
	newSuperBlk.SU_str_blk = newSuperBlk.SU_str_inode + n*int32(binary.Size(Structs.Inode{}))
	newSuperBlk.SU_Free_Inodo -= 2
	newSuperBlk.SU_Free_Blk -= 2
	newSuperBlk.SU_fst_ino = int32(2)
	newSuperBlk.SU_fst_blk = int32(2)
	bmInodeData := make([]byte, n)
	bmInodeErr := Herramientas.WrObj(file, bmInodeData, int64(newSuperBlk.SU_btp_str_ino))
	if bmInodeErr != nil {
		fmt.Println("Error: ", bmInodeErr)
		return
	}
	bmBlockData := make([]byte, 3*n)
	bmBlockErr := Herramientas.WrObj(file, bmBlockData, int64(newSuperBlk.SU_btp_str_blk))
	if bmBlockErr != nil {
		fmt.Println("Error: ", bmBlockErr)
		return
	}
	var newInode Structs.Inode
	for i := 0; i < 15; i++ {
		newInode.In_blk[i] = -1
	}

	for i := int32(0); i < n; i++ {
		err := Herramientas.WrObj(file, newInode, int64(newSuperBlk.SU_str_inode+i*int32(binary.Size(Structs.Inode{}))))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
	fileBlocks := make([]Structs.Fileblock, 3*n)
	fileBlocksErr := Herramientas.WrObj(file, fileBlocks, int64(newSuperBlk.SU_btp_str_blk))
	if fileBlocksErr != nil {
		fmt.Println("Error: ", fileBlocksErr)
		return
	}

	var Inode0 Structs.Inode
	Inode0.In_usrid = 1
	Inode0.In_grpid = 1
	Inode0.In_size = 0
	copy(Inode0.In_readDate[:], date)
	copy(Inode0.In_crtdate[:], date)
	copy(Inode0.In_modDate[:], date)
	copy(Inode0.In_typ[:], "0")
	copy(Inode0.In_usrpr[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.In_blk[i] = -1
	}
	Inode0.In_blk[0] = 0
	var FLDBLK0 Structs.FLDBLK
	FLDBLK0.B_CONT[0].B_inodo = 0
	copy(FLDBLK0.B_CONT[0].B_name[:], ".")
	FLDBLK0.B_CONT[1].B_inodo = 0
	copy(FLDBLK0.B_CONT[1].B_name[:], "..")
	FLDBLK0.B_CONT[2].B_inodo = 1
	copy(FLDBLK0.B_CONT[2].B_name[:], "users.txt")
	FLDBLK0.B_CONT[3].B_inodo = -1

	var Inode1 Structs.Inode
	Inode1.In_usrid = 1
	Inode1.In_grpid = 1
	Inode1.In_size = int32(binary.Size(Structs.FLDBLK{}))
	copy(Inode1.In_readDate[:], date)
	copy(Inode1.In_crtdate[:], date)
	copy(Inode1.In_modDate[:], date)
	copy(Inode1.In_typ[:], "1")
	copy(Inode0.In_usrpr[:], "664")
	for i := int32(0); i < 15; i++ {
		Inode1.In_blk[i] = -1
	}
	Inode1.In_blk[0] = 1
	data := "1,G,root\n1,U,root,root,1234\n"
	var fileBlock1 Structs.Fileblock
	copy(fileBlock1.B_CONT[:], []byte(data))
	Herramientas.WrObj(file, newSuperBlk, int64(particion.Start))
	Herramientas.WrObj(file, byte(1), int64(newSuperBlk.SU_btp_str_ino))
	Herramientas.WrObj(file, byte(1), int64(newSuperBlk.SU_btp_str_ino+1))
	Herramientas.WrObj(file, byte(1), int64(newSuperBlk.SU_btp_str_blk))
	Herramientas.WrObj(file, byte(1), int64(newSuperBlk.SU_btp_str_blk+1))
	Herramientas.WrObj(file, Inode0, int64(newSuperBlk.SU_str_inode))
	Herramientas.WrObj(file, Inode1, int64(newSuperBlk.SU_str_inode+int32(binary.Size(Structs.Inode{}))))
	Herramientas.WrObj(file, FLDBLK0, int64(newSuperBlk.SU_str_blk))
	Herramientas.WrObj(file, fileBlock1, int64(newSuperBlk.SU_str_blk+int32(binary.Size(Structs.Fileblock{}))))
}

func createEXT3(n int32, particion Structs.Partition, newSuperBlk Structs.SuBlock, date string, file *os.File) {
	fmt.Println("====== CREAR EXT3 ======")
	fmt.Println("N: ", n)
	fmt.Println("Journaling: ", binary.Size(Structs.Journaling{}))
	newSuperBlk.SU_FileSys_ty = 3
	newSuperBlk.SU_btp_str_ino = particion.Start + int32(binary.Size(Structs.SuBlock{})) + int32(binary.Size(Structs.Journaling{}))
	newSuperBlk.SU_btp_str_blk = newSuperBlk.SU_btp_str_ino + n
	newSuperBlk.SU_str_inode = newSuperBlk.SU_btp_str_blk + 3*n
	newSuperBlk.SU_str_blk = newSuperBlk.SU_str_inode + n*int32(binary.Size(Structs.Inode{}))
	newSuperBlk.SU_Free_Inodo -= 2
	newSuperBlk.SU_Free_Blk -= 2
	newSuperBlk.SU_fst_ino = int32(2)
	newSuperBlk.SU_fst_blk = int32(2)
	var newJrnal Structs.Journaling
	newJrnal.Ultimo = 0
	newJrnal.Size = int32(binary.Size(Structs.Journaling{}))
	dataJ := newJrnal.Contenido[0]
	copy(dataJ.Ope[:], "MKDIR")
	copy(dataJ.Path[:], "/")
	copy(dataJ.CONT[:], "-")
	copy(dataJ.Date[:], date)
	newJrnal.Contenido[0] = dataJ
	bmInodeData := make([]byte, n)
	bmInodeErr := Herramientas.WrObj(file, bmInodeData, int64(newSuperBlk.SU_btp_str_ino))
	if bmInodeErr != nil {
		fmt.Println("Error: ", bmInodeErr)
		return
	}
	bmBlockData := make([]byte, 3*n)
	bmBlockErr := Herramientas.WrObj(file, bmBlockData, int64(newSuperBlk.SU_btp_str_blk))
	if bmBlockErr != nil {
		fmt.Println("Error: ", bmBlockErr)
		return
	}
	var newInode Structs.Inode
	for i := 0; i < 15; i++ {
		newInode.In_blk[i] = -1
	}
	for i := int32(0); i < n; i++ {
		err := Herramientas.WrObj(file, newInode, int64(newSuperBlk.SU_str_inode+i*int32(binary.Size(Structs.Inode{}))))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
	fileBlocks := make([]Structs.Fileblock, 3*n)
	fileBlocksErr := Herramientas.WrObj(file, fileBlocks, int64(newSuperBlk.SU_btp_str_blk))
	if fileBlocksErr != nil {
		fmt.Println("Error: ", fileBlocksErr)
		return
	}
	var Inode0 Structs.Inode
	Inode0.In_usrid = 1
	Inode0.In_grpid = 1
	Inode0.In_size = 0
	copy(Inode0.In_readDate[:], date)
	copy(Inode0.In_crtdate[:], date)
	copy(Inode0.In_modDate[:], date)
	copy(Inode0.In_typ[:], "0")
	copy(Inode0.In_usrpr[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.In_blk[i] = -1
	}

	Inode0.In_blk[0] = 0

	var FLDBLK0 Structs.FLDBLK
	FLDBLK0.B_CONT[0].B_inodo = 0
	copy(FLDBLK0.B_CONT[0].B_name[:], ".")
	FLDBLK0.B_CONT[1].B_inodo = 0
	copy(FLDBLK0.B_CONT[1].B_name[:], "..")
	FLDBLK0.B_CONT[2].B_inodo = 1
	copy(FLDBLK0.B_CONT[2].B_name[:], "users.txt")
	FLDBLK0.B_CONT[3].B_inodo = -1
	var Inode1 Structs.Inode
	Inode1.In_usrid = 1
	Inode1.In_grpid = 1
	Inode1.In_size = int32(binary.Size(Structs.FLDBLK{}))
	copy(Inode1.In_readDate[:], date)
	copy(Inode1.In_crtdate[:], date)
	copy(Inode1.In_modDate[:], date)
	copy(Inode1.In_typ[:], "1")
	copy(Inode0.In_usrpr[:], "664")
	for i := int32(0); i < 15; i++ {
		Inode1.In_blk[i] = -1
	}
	Inode1.In_blk[0] = 1
	data := "1,G,root\n1,U,root,root,123\n"
	var fileBlock1 Structs.Fileblock
	copy(fileBlock1.B_CONT[:], []byte(data))
	Herramientas.WrObj(file, newSuperBlk, int64(particion.Start))
	Herramientas.WrObj(file, newJrnal, int64(particion.Start+int32(binary.Size(Structs.SuBlock{}))))
	Herramientas.WrObj(file, byte(1), int64(newSuperBlk.SU_btp_str_ino))
	Herramientas.WrObj(file, byte(1), int64(newSuperBlk.SU_btp_str_ino+1))
	Herramientas.WrObj(file, byte(1), int64(newSuperBlk.SU_btp_str_blk))
	Herramientas.WrObj(file, byte(1), int64(newSuperBlk.SU_btp_str_blk+1))
	Herramientas.WrObj(file, Inode0, int64(newSuperBlk.SU_str_inode))
	Herramientas.WrObj(file, Inode1, int64(newSuperBlk.SU_str_inode+int32(binary.Size(Structs.Inode{}))))
	Herramientas.WrObj(file, FLDBLK0, int64(newSuperBlk.SU_str_blk))
	Herramientas.WrObj(file, fileBlock1, int64(newSuperBlk.SU_str_blk+int32(binary.Size(Structs.Fileblock{}))))
}
