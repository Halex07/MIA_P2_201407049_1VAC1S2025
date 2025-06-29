package HerramientasInodos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func LookInodo(idInodo int32, path string, SuperBlk Structs.SuBlock, file *os.File) int32 {
	PasPath := strings.Split(path, "/")
	tempPath := PasPath[1:]
	var Inode0 Structs.Inode
	Herramientas.ReadObj(file, &Inode0, int64(SuperBlk.SU_str_inode+(idInodo*int32(binary.Size(Structs.Inode{})))))
	var FLDBLK Structs.FLDBLK
	for i := 0; i < 12; i++ {
		idBlk := Inode0.In_blk[i]
		if idBlk != -1 {
			Herramientas.ReadObj(file, &FLDBLK, int64(SuperBlk.SU_str_blk+(idBlk*int32(binary.Size(Structs.FLDBLK{})))))
			for j := 2; j < 4; j++ {
				apuntador := FLDBLK.B_CONT[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GETB_NOM(string(FLDBLK.B_CONT[j].B_name[:]))
					if tempPath[0] == pathActual {
						if len(tempPath) > 1 {
							return LookInodoRec(apuntador, tempPath[1:], SuperBlk.SU_str_inode, SuperBlk.SU_str_blk, file)
						} else {
							return apuntador
						}
					}
				}
			}
		}
	}
	return idInodo
}

func LookInodoRec(idInodo int32, path []string, iStart int32, bStart int32, file *os.File) int32 {
	var inodo Structs.Inode
	Herramientas.ReadObj(file, &inodo, int64(iStart+(idInodo*int32(binary.Size(Structs.Inode{})))))
	var FLDBLK Structs.FLDBLK
	for i := 0; i < 12; i++ {
		idBlk := inodo.In_blk[i]
		if idBlk != -1 {
			Herramientas.ReadObj(file, &FLDBLK, int64(bStart+(idBlk*int32(binary.Size(Structs.FLDBLK{})))))
			for j := 2; j < 4; j++ {
				apuntador := FLDBLK.B_CONT[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GETB_NOM(string(FLDBLK.B_CONT[j].B_name[:]))
					if path[0] == pathActual {
						if len(path) > 1 {
							return LookInodoRec(apuntador, path[1:], iStart, bStart, file)
						} else {
							return apuntador
						}
					}
				}
			}
		}
	}
	return -1
}

func CRTFOLDER(idInode int32, carpeta string, initSuperBlk int64, disco *os.File) int32 {
	var SuperBlk Structs.SuBlock
	Herramientas.ReadObj(disco, &SuperBlk, initSuperBlk)
	var inodo Structs.Inode
	Herramientas.ReadObj(disco, &inodo, int64(SuperBlk.SU_str_inode+(idInode*int32(binary.Size(Structs.Inode{})))))
	fmt.Println("Un momemnto por favor... Creando Carpeta ", carpeta)

	for i := 0; i < 12; i++ {
		idBlk := inodo.In_blk[i]
		if idBlk != -1 {
			var FLDBLK Structs.FLDBLK
			Herramientas.ReadObj(disco, &FLDBLK, int64(SuperBlk.SU_str_blk+(idBlk*int32(binary.Size(Structs.FLDBLK{})))))
			for j := 2; j < 4; j++ {
				apuntador := FLDBLK.B_CONT[j].B_inodo
				if apuntador == -1 {
					copy(FLDBLK.B_CONT[j].B_name[:], carpeta)
					ino := SuperBlk.SU_fst_ino
					FLDBLK.B_CONT[j].B_inodo = ino
					Herramientas.WrObj(disco, FLDBLK, int64(SuperBlk.SU_str_blk+(idBlk*int32(binary.Size(Structs.FLDBLK{})))))
					var newInodo Structs.Inode
					newInodo.In_usrid = Structs.CurrentUSR.USRID
					newInodo.In_grpid = Structs.CurrentUSR.GRPID
					newInodo.In_size = 0
					ahora := time.Now()
					date := ahora.Format("12/06/2025 15:04")
					copy(newInodo.In_readDate[:], date)
					copy(newInodo.In_crtdate[:], date)
					copy(newInodo.In_modDate[:], date)
					copy(newInodo.In_typ[:], "0")
					copy(newInodo.In_modDate[:], "664")
					for i := int32(0); i < 15; i++ {
						newInodo.In_blk[i] = -1
					}
					block := SuperBlk.SU_fst_blk
					newInodo.In_blk[0] = block
					Herramientas.WrObj(disco, newInodo, int64(SuperBlk.SU_str_inode+(ino*int32(binary.Size(Structs.Inode{})))))
					var newFLDBLK Structs.FLDBLK
					newFLDBLK.B_CONT[0].B_inodo = ino
					copy(newFLDBLK.B_CONT[0].B_name[:], ".")
					newFLDBLK.B_CONT[1].B_inodo = FLDBLK.B_CONT[0].B_inodo
					copy(newFLDBLK.B_CONT[1].B_name[:], "..")
					newFLDBLK.B_CONT[2].B_inodo = -1
					newFLDBLK.B_CONT[3].B_inodo = -1
					Herramientas.WrObj(disco, newFLDBLK, int64(SuperBlk.SU_str_blk+(block*int32(binary.Size(Structs.FLDBLK{})))))
					SuperBlk.SU_Free_Inodo -= 1
					SuperBlk.SU_Free_Blk -= 1
					SuperBlk.SU_fst_blk += 1
					SuperBlk.SU_fst_ino += 1
					Herramientas.WrObj(disco, SuperBlk, initSuperBlk)
					Herramientas.WrObj(disco, byte(1), int64(SuperBlk.SU_btp_str_blk+block))
					Herramientas.WrObj(disco, byte(1), int64(SuperBlk.SU_btp_str_ino+ino))
					return ino
				}
			}
		} else {
			block := SuperBlk.SU_fst_blk
			inodo.In_blk[i] = block
			Herramientas.WrObj(disco, &inodo, int64(SuperBlk.SU_str_inode+(idInode*int32(binary.Size(Structs.Inode{})))))
			var FLDBLK Structs.FLDBLK
			bloque := inodo.In_blk[0]
			Herramientas.ReadObj(disco, &FLDBLK, int64(SuperBlk.SU_str_blk+(bloque*int32(binary.Size(Structs.FLDBLK{})))))
			var newFLDBLK1 Structs.FLDBLK
			newFLDBLK1.B_CONT[0].B_inodo = FLDBLK.B_CONT[0].B_inodo
			copy(newFLDBLK1.B_CONT[0].B_name[:], ".")
			newFLDBLK1.B_CONT[1].B_inodo = FLDBLK.B_CONT[1].B_inodo
			copy(newFLDBLK1.B_CONT[1].B_name[:], "..")
			ino := SuperBlk.SU_fst_ino
			newFLDBLK1.B_CONT[2].B_inodo = ino
			copy(newFLDBLK1.B_CONT[2].B_name[:], carpeta)
			newFLDBLK1.B_CONT[3].B_inodo = -1
			Herramientas.WrObj(disco, newFLDBLK1, int64(SuperBlk.SU_str_blk+(block*int32(binary.Size(Structs.FLDBLK{})))))
			var newInodo Structs.Inode
			newInodo.In_usrid = Structs.CurrentUSR.USRID
			newInodo.In_grpid = Structs.CurrentUSR.GRPID
			newInodo.In_size = 0
			ahora := time.Now()
			date := ahora.Format("12/06/2025 15:04")
			copy(newInodo.In_readDate[:], date)
			copy(newInodo.In_crtdate[:], date)
			copy(newInodo.In_modDate[:], date)
			copy(newInodo.In_typ[:], "0")
			copy(newInodo.In_modDate[:], "664")
			for i := int32(0); i < 15; i++ {
				newInodo.In_blk[i] = -1
			}
			block2 := SuperBlk.SU_fst_blk + 1
			newInodo.In_blk[0] = block2
			Herramientas.WrObj(disco, newInodo, int64(SuperBlk.SU_str_inode+(ino*int32(binary.Size(Structs.Inode{})))))
			var newFLDBLK2 Structs.FLDBLK
			newFLDBLK2.B_CONT[0].B_inodo = ino
			copy(newFLDBLK2.B_CONT[0].B_name[:], ".")
			newFLDBLK2.B_CONT[1].B_inodo = newFLDBLK1.B_CONT[0].B_inodo
			copy(newFLDBLK2.B_CONT[1].B_name[:], "..")
			newFLDBLK2.B_CONT[2].B_inodo = -1
			newFLDBLK2.B_CONT[3].B_inodo = -1
			Herramientas.WrObj(disco, newFLDBLK2, int64(SuperBlk.SU_str_blk+(block2*int32(binary.Size(Structs.FLDBLK{})))))
			SuperBlk.SU_Free_Inodo -= 1
			SuperBlk.SU_Free_Blk -= 2
			SuperBlk.SU_fst_blk += 2
			SuperBlk.SU_fst_ino += 1
			Herramientas.WrObj(disco, SuperBlk, initSuperBlk)
			Herramientas.WrObj(disco, byte(1), int64(SuperBlk.SU_btp_str_blk+block))
			Herramientas.WrObj(disco, byte(1), int64(SuperBlk.SU_btp_str_blk+block2))
			Herramientas.WrObj(disco, byte(1), int64(SuperBlk.SU_btp_str_ino+ino))
			return ino
		}
	}
	return 0
}
