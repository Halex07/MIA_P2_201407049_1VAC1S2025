package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Rep(parametros []string) {
	fmt.Println(">> Generando Reporte")
	var (
		name    string
		path    string
		id      string
		paramOK bool = true
	)
	for _, raw := range parametros[1:] {
		temp2 := strings.TrimRight(raw, " ")
		temp := strings.Split(temp2, "=")
		if len(temp) != 2 {
			fmt.Println("Error valor no reconocido del parametro", temp[0])
			paramOK = false
			break
		}

		if strings.ToLower(temp[0]) == "name" {
			name = strings.ToLower(temp[1])
		} else if strings.ToLower(temp[0]) == "path" {
			name = strings.ReplaceAll(temp[1], "\"", "")
			path = name
		} else if strings.ToLower(temp[0]) == "id" {
			id = strings.ToUpper(temp[1])
		} else if strings.ToLower(temp[0]) == "ruta" {
		} else {
			fmt.Println("Error parametro no reconocido", temp[0])
			paramOK = false
			break
		}
	}

	if paramOK {
		if name != "" && id != "" && path != "" {
			switch name {
			case "mbr":
				fmt.Println("reporte MBR")
				mbr(path, id)
			case "disk":
				fmt.Println("reporte DISK")
				disk(path, id)
			case "inode":
				fmt.Println("reporte INODE")
			case "journaling":
				fmt.Println("reporte JOURNALING")
				journal(path, id)
			case "block":
				fmt.Println("reporte BLOCK")
			case "bm_inode":
				fmt.Println("reporte BM INODE")
				bm_inode(path, id)
			case "bm_block":
				fmt.Println("reporte BM BLOCK")
				bm_block(path, id)
			case "tree":
				fmt.Println("reporte TREE")
				tree(path, id)
			case "sb":
				fmt.Println("reporte SuperBlk")
				sb(path, id)
			case "file":
				fmt.Println("reporte File")
			case "ls":
				fmt.Println("reporte ls")
			default:
				fmt.Println("Error reporte", name, " no reconocido")
			}
		} else {
			fmt.Println("Error verifique los parametros")
		}
	}
}

func mbr(path string, id string) {
	disk := id[0:1]
	temp := strings.Split(path, "/")
	nombre := strings.Split(temp[len(temp)-1], ".")[0]
	folder := "./MIA/P1/"
	ext := ".dsk"
	dirDisk := folder + disk + ext

	file, err := Herramientas.OpenFile(dirDisk)
	if err != nil {
		return
	}
	var mbr Structs.MBR
	if err := Herramientas.ReadObj(file, &mbr, 0); err != nil {
		return
	}
	defer file.Close()
	reportar := false
	for i := 0; i < 4; i++ {
		identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			break
		}
	}
	if reportar {
		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='SlateBlue' COLSPAN=\"2\"> Reporte MBR </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_tamano </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", mbr.MBRSZ)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#AFA1D1'> mbr_fecha_creacion </td> \n  <td bgcolor='#AFA1D1'> %s </td> \n </tr> \n", string(mbr.DATECRT[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_disk_signature </td> \n  <td bgcolor='Azure'> %d </td> \n </tr>  \n", mbr.Id)
		cad += Structs.ReportGrap(mbr, file)
		cad += "</table> > ]\n}"
		folder = filepath.Dir(path)
		rutaReporte := "." + folder + "/" + nombre + ".dot"

		Herramientas.ReportGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("Error ID no reconocido ")
	}
}

func disk(path string, id string) {
	disk := id[0:1]
	temp := strings.Split(path, "/")
	nombre := strings.Split(temp[len(temp)-1], ".")[0]
	folder := "./MIA/P1/"
	ext := ".dsk"
	dirDisk := folder + disk + ext
	file, err := Herramientas.OpenFile(dirDisk)
	if err != nil {
		return
	}
	var TempMBR Structs.MBR
	if err := Herramientas.ReadObj(file, &TempMBR, 0); err != nil {
		return
	}
	defer file.Close()
	cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n<tr> \n"
	cad += " <td bgcolor='SlateBlue'  ROWSPAN='3'> MBR </td>\n"
	cad += Structs.ReporDkGrap(TempMBR, file)
	cad += "\n</table> > ]\n}"
	folder = filepath.Dir(path)
	rutaReporte := "." + folder + "/" + nombre + ".dot"

	Herramientas.ReportGraphizMBR(rutaReporte, cad, nombre)
}

func sb(path string, id string) {
	disk := id[0:1]
	temp := strings.Split(path, "/")
	nombre := strings.Split(temp[len(temp)-1], ".")[0]
	folder := "./MIA/P1/"
	ext := ".dsk"
	dirDisk := folder + disk + ext

	file, err := Herramientas.OpenFile(dirDisk)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	if err := Herramientas.ReadObj(file, &mbr, 0); err != nil {
		return
	}

	defer file.Close()

	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break
		}
	}

	if reportar {

		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='darkgreen' COLSPAN=\"2\"> <font color='white'> Reporte SuperBlk </font> </td> \n </tr> \n"
		cad += Structs.ReportSb(mbr.Partitions[part], file)
		cad += "</table> > ]\n}"
		folder = filepath.Dir(path)
		rutaReporte := "." + folder + "/" + nombre + ".dot"

		Herramientas.ReportGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("Erro ID no reconocido")
	}
}

func journal(path string, id string) {
	disk := id[0:1]
	temp := strings.Split(path, "/")
	nombre := strings.Split(temp[len(temp)-1], ".")[0]
	folder := "./MIA/P1/"
	ext := ".dsk"
	dirDisk := folder + disk + ext

	file, err := Herramientas.OpenFile(dirDisk)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	if err := Herramientas.ReadObj(file, &mbr, 0); err != nil {
		return
	}
	defer file.Close()
	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break
		}
	}
	if reportar {

		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += Structs.ReportJrnal(mbr.Partitions[part], file)
		cad += "</table> > ]\n}"
		folder = filepath.Dir(path)
		rutaReporte := "." + folder + "/" + nombre + ".dot"

		Herramientas.ReportGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("Error ID no reconocido")
	}
}

func bm_inode(path string, id string) {
	disk := id[0:1]
	temp := strings.Split(path, "/")
	nombre := strings.Split(temp[len(temp)-1], ".")[0]
	folder := "./MIA/P1/"
	ext := ".dsk"
	dirDisk := folder + disk + ext

	file, err := Herramientas.OpenFile(dirDisk)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	if err := Herramientas.ReadObj(file, &mbr, 0); err != nil {
		return
	}

	defer file.Close()

	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break
		}
	}

	if reportar {
		var SuperBlk Structs.SuBlock
		err := Herramientas.ReadObj(file, &SuperBlk, int64(mbr.Partitions[part].Start))
		if err != nil {
			fmt.Println("Error no se detecto ningun formato")
			return
		}

		cad := ""
		inicio := SuperBlk.SU_btp_str_ino
		fin := SuperBlk.SU_btp_str_blk
		count := 1
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			Herramientas.ReadObj(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += "0 "
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}
		folder = filepath.Dir(path)
		rutaReporte := "." + folder + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)
	} else {
		fmt.Println("Error ID no reconocido")
	}
}

func bm_block(path string, id string) {
	disk := id[0:1]
	temp := strings.Split(path, "/")
	nombre := strings.Split(temp[len(temp)-1], ".")[0]
	folder := "./MIA/P1/"
	ext := ".dsk"
	dirDisk := folder + disk + ext

	file, err := Herramientas.OpenFile(dirDisk)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	if err := Herramientas.ReadObj(file, &mbr, 0); err != nil {
		return
	}
	defer file.Close()
	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break
		}
	}

	if reportar {
		var SuperBlk Structs.SuBlock
		err := Herramientas.ReadObj(file, &SuperBlk, int64(mbr.Partitions[part].Start))
		if err != nil {
			fmt.Println("Error no se detecto ningun formato")
			return
		}

		cad := ""
		inicio := SuperBlk.SU_btp_str_blk
		fin := SuperBlk.SU_str_inode
		count := 1
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			Herramientas.ReadObj(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += "0 "
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}
		folder = filepath.Dir(path)
		rutaReporte := "." + folder + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)
	} else {
		fmt.Println("Erro ID desconocido")
	}
}

func tree(path string, id string) {
	disk := id[0:1]
	temp := strings.Split(path, "/")
	nombre := strings.Split(temp[len(temp)-1], ".")[0]
	folder := "./MIA/P1/"
	ext := ".dsk"
	dirDisk := folder + disk + ext
	file, err := Herramientas.OpenFile(dirDisk)
	if err != nil {
		return
	}

	var mbr Structs.MBR
	if err := Herramientas.ReadObj(file, &mbr, 0); err != nil {
		return
	}
	defer file.Close()
	reportar := false
	part := -1
	for i := 0; i < 4; i++ {
		identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			reportar = true
			part = i
			break
		}
	}

	if reportar {

		var SuperBlk Structs.SuBlock
		err := Herramientas.ReadObj(file, &SuperBlk, int64(mbr.Partitions[part].Start))
		if err != nil {
			fmt.Println("Error no se detecto ningun formato para la partición")
			return
		}

		var Inode0 Structs.Inode
		Herramientas.ReadObj(file, &Inode0, int64(SuperBlk.SU_str_inode))
		cad := "digraph { \n graph [pad=0.5, nodesep=0.5, ranksep=1] \n node [ shape=plaintext ] \n rankdir=LR \n"
		cad += "\n Inodo0 [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
		cad += "    <tr> <td bgcolor='skyblue' colspan=\"2\" port='P0'> Inodo 0 </td> </tr> \n"

		for i := 0; i < 12; i++ {
			cad += fmt.Sprintf("    <tr> <td> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode0.In_blk[i])
		}
		for i := 12; i < 15; i++ {
			cad += fmt.Sprintf("    <tr> <td bgcolor='pink'> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode0.In_blk[i])
		}
		cad += "   </table> \n  > \n ]; \n"

		for i := 0; i < 15; i++ {
			bloque := Inode0.In_blk[i]
			if bloque != -1 {
				cad += treeBlock(bloque, string(Inode0.In_typ[:]), 0, i+1, SuperBlk, file)
			}
		}
		cad += "\n}"
		folder = filepath.Dir(path)
		rutaReporte := "." + folder + "/" + nombre + ".dot"

		Herramientas.ReportGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("Error ID no reconocido")
	}
}
func treeBlock(idBloque int32, tipo string, idPadre int32, p int, SuperBlk Structs.SuBlock, file *os.File) string {
	cad := fmt.Sprintf("\n Bloque%d [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n", idBloque)

	if tipo == "0" {
		var FLDBLK Structs.FLDBLK
		Herramientas.ReadObj(file, &FLDBLK, int64(SuperBlk.SU_str_blk+(idBloque*int32(binary.Size(Structs.FLDBLK{})))))
		cad += fmt.Sprintf("    <tr> <td bgcolor='orchid' colspan=\"2\" port='P0'> Bloque %d </td> </tr> \n", idBloque)
		cad += fmt.Sprintf("    <tr> <td> . </td> <td port='P1'> %d </td> </tr> \n", FLDBLK.B_CONT[0].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> .. </td> <td port='P2'> %d </td> </tr> \n", FLDBLK.B_CONT[1].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> %s </td> <td port='P3'> %d </td> </tr> \n", Structs.GETB_NOM(string(FLDBLK.B_CONT[2].B_name[:])), FLDBLK.B_CONT[2].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> %s </td> <td port='P4'> %d </td> </tr> \n", Structs.GETB_NOM(string(FLDBLK.B_CONT[3].B_name[:])), FLDBLK.B_CONT[3].B_inodo)
		cad += "   </table> \n  > \n ]; \n"
		cad += fmt.Sprintf("\n Inodo%d:P%d -> Bloque%d:P0; \n", idPadre, p, idBloque)
		for i := 2; i < 4; i++ {
			inodo := FLDBLK.B_CONT[i].B_inodo
			if inodo != -1 {
				cad += treeInodo(inodo, idBloque, i+1, SuperBlk, file)
			}
		}
	} else {
		var fileBlock Structs.Fileblock
		Herramientas.ReadObj(file, &fileBlock, int64(SuperBlk.SU_str_blk+(idBloque*int32(binary.Size(Structs.Fileblock{})))))
		cad += fmt.Sprintf("    <tr> <td bgcolor='#ffff99' port='P0'> Bloque %d </td> </tr> \n", idBloque)
		cad += fmt.Sprintf("    <tr> <td> %s </td> </tr> \n", Structs.GETB_CONT(string(fileBlock.B_CONT[:])))
		cad += "   </table> \n  > \n ]; \n"
		cad += fmt.Sprintf("\n Inodo%d:P%d -> Bloque%d:P0; \n", idPadre, p, idBloque)
	}

	return cad
}

func treeInodo(idInodo int32, idPadre int32, p int, SuperBlk Structs.SuBlock, file *os.File) string {
	var Inode Structs.Inode
	Herramientas.ReadObj(file, &Inode, int64(SuperBlk.SU_str_inode+(idInodo*int32(binary.Size(Structs.Inode{})))))
	cad := fmt.Sprintf("\n Inodo%d [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n", idInodo)
	if string(Inode.In_typ[:]) == "0" {
		cad += fmt.Sprintf("    <tr> <td bgcolor='skyblue' colspan=\"2\" port='P0'> Inodo %d </td> </tr> \n", idInodo)
	} else {
		cad += fmt.Sprintf("    <tr> <td bgcolor='#7FC97F' colspan=\"2\" port='P0'> Inodo %d </td> </tr> \n", idInodo)
	}
	for i := 0; i < 12; i++ {
		cad += fmt.Sprintf("    <tr> <td> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode.In_blk[i])
	}
	for i := 12; i < 15; i++ {
		cad += fmt.Sprintf("    <tr> <td bgcolor='pink'> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode.In_blk[i])
	}
	cad += "   </table> \n  > \n ]; \n"
	cad += fmt.Sprintf("\n Bloque%d:P%d -> Inodo%d:P0; \n", idPadre, p, idInodo)

	for i := 0; i < 15; i++ {
		bloque := Inode.In_blk[i]
		if bloque != -1 {
			cad += treeBlock(bloque, string(Inode.In_typ[:]), idInodo, i+1, SuperBlk, file)
		}
	}

	return cad
}
