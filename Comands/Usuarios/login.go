package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Login(parametros []string) {
	fmt.Println("Login")
	var (
		user    string
		pass    string
		id      string
		paramOk bool = true
	)
	for _, raw := range parametros[1:] {
		temp2 := strings.TrimRight(raw, " ")
		temp := strings.Split(temp2, "=")
		if len(temp) != 2 {
			fmt.Println("Error valor de parametro", temp[0], "no reconocido")
			paramOk = false
			break
		}
		if strings.ToLower(temp[0]) == "ID" {
			id = strings.ToUpper(temp[1])

		} else if strings.ToLower(temp[0]) == "Usuario" {
			user = temp[1]

		} else if strings.ToLower(temp[0]) == "Password" {
			pass = temp[1]

		} else {
			fmt.Println("Error parametro ", temp[0], "no reconocido")
			paramOk = false
			break
		}
	}

	if paramOk {
		if id != "" && user != "" && pass != "" {
			if Structs.CurrentUSR.STATUS {
				fmt.Println("Error ya existe una sesión en curso, intente mas tarde")
				return
			}
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

			logear := false
			index := -1
			for i := 0; i < 4; i++ {
				identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
				if identificador == id {
					logear = true
					index = i
					break
				}
			}

			if logear {
				var SuperBlk Structs.SuBlock
				err := Herramientas.ReadObj(disco, &SuperBlk, int64(mbr.Partitions[index].Start))
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

				linea := strings.Split(contenido, "\n")

				loginFail := true
				for _, reg := range linea {
					usuario := strings.Split(reg, ",")

					if len(usuario) == 5 {
						if usuario[0] != "0" {
							if usuario[3] == user {
								if usuario[4] == pass {
									loginFail = false
									Structs.CurrentUSR.Id = id
									buscarGRPID(linea, usuario[2])
									USRID(usuario[0])
									Structs.CurrentUSR.Nombre = user
									Structs.CurrentUSR.STATUS = true
									fmt.Println("Bienvenido ", user)
								} else {
									loginFail = false
									fmt.Println("Precaucón la contraseña no es correcta")
								}
								break
							}
						}
					}
				}

				if loginFail {
					fmt.Println("Error usuario desconocido")
				}
			} else {
				fmt.Println("Error ID no reconocido")
			}
		} else {
			fmt.Println("Error no cumple con los parametros minimos")
		}
	}
}

func buscarGRPID(lineID []string, grupo string) {
	for _, reg := range lineID[:len(lineID)-1] {
		datos := strings.Split(reg, ",")
		if len(datos) == 3 {
			if datos[2] == grupo {
				id, IDError := strconv.Atoi(datos[0])
				if IDError != nil {
					fmt.Println("Error desconocido valide el ID")
					return
				}
				Structs.CurrentUSR.GRPID = int32(id)
				return
			}
		}
	}
}

func USRID(id string) {
	idU, IDError := strconv.Atoi(id)
	if IDError != nil {
		fmt.Println("Error desconocido valide el ID")
		return
	}
	Structs.CurrentUSR.USRID = int32(idU)
}
