package Comandos

import (
	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Fdisk(parametros []string) {
	fmt.Println(">> Ejecutando Comando FDISK")

	var (
		size       int
		letter     string
		name       string
		unit       int    = 1024
		typee      string = "P"
		fit        string = "W"
		add        int
		opc        int
		paramOk    bool = true
		sizeInit   bool = false
		sizeValErr string
	)
	for _, raw := range parametros[1:] {
		temp2 := strings.TrimRight(raw, " ")
		temp := strings.Split(temp2, "=")
		if len(temp) != 2 {
			fmt.Println("Error valor de parametro", temp[0], " no identificado")
			paramOk = false
			break
		}
		if strings.ToLower(temp[0]) == "size" {
			sizeInit = true
			var err error
			size, err = strconv.Atoi(temp[1])
			if err != nil {
				sizeValErr = temp[1]
			}
		} else if strings.ToLower(temp[0]) == "driveletter" {
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
		} else if strings.ToLower(temp[0]) == "unit" {
			if strings.ToLower(temp[1]) == "b" {
				unit = 1
			} else if strings.ToLower(temp[1]) == "m" {
				unit = 1024 * 1024
			} else if strings.ToLower(temp[1]) != "k" {
				fmt.Println("Error unit los valores aceptados son: b, k, m. verificar: ", temp[1])
				paramOk = false
				break
			}

		} else if strings.ToLower(temp[0]) == "type" {
			if strings.ToLower(temp[1]) == "e" {
				typee = "E"
			} else if strings.ToLower(temp[1]) == "l" {
				typee = "L"
			} else if strings.ToLower(temp[1]) != "p" {
				fmt.Println("Error type los valores aceptados: e, l, p. verificar: ", temp[1])
				paramOk = false
				break
			}
		} else if strings.ToLower(temp[0]) == "fit" {
			if strings.ToLower(temp[1]) == "bf" {
				fit = "B"
			} else if strings.ToLower(temp[1]) == "ff" {
				fit = "F"
			} else if strings.ToLower(temp[1]) != "wf" {
				fmt.Println("Error en fit los valores aceptados son: BF, FF o WF. verificar: ", temp[1])
				paramOk = false
				break
			}
		} else if strings.ToLower(temp[0]) == "delete" {
			if strings.ToLower(temp[1]) == "full" {
				if opc == 0 {
					opc = 2
				}
			} else {
				fmt.Println("Error valor del comando delete no fue reconocido")
				paramOk = false
				break
			}
		} else if strings.ToLower(temp[0]) == "add" {
			var err error
			add, err = strconv.Atoi(temp[1])
			if err != nil {
				fmt.Println("Error el valor  \"add\" debe ser numerico, verificar:  ", temp[1])
				paramOk = false
				break
			} else {
				if opc == 0 {
					opc = 1
				}
			}
		} else {
			fmt.Println("Error parametro", temp[0], " no reconocido")
			paramOk = false
			break
		}
	}
	if opc == 0 && paramOk {
		if sizeInit {
			if sizeValErr == "" {
				if size <= 0 {
					fmt.Println("Error size debe ser un numero mayo a 0", size)
					paramOk = false
				}
			} else {
				fmt.Println("Error: size debe debe ser numerico verificar: ", sizeValErr)
				paramOk = false
			}
		} else {
			fmt.Println("Error no se ingreso el parametro size")
			paramOk = false
		}
	}
	if paramOk {
		if letter != "" && name != "" {
			filepath := "./MIA/P1/" + letter + ".dsk"
			disco, err := Herramientas.OpenFile(filepath)
			if err != nil {
				fmt.Println("Error imposible leer el disco")
				return
			}
			var mbr Structs.MBR
			if err := Herramientas.ReadObj(disco, &mbr, 0); err != nil {
				return
			}
			if opc == 0 {
				isPartExtend := false
				isName := true
				if typee == "E" {
					for i := 0; i < 4; i++ {
						tipo := string(mbr.Partitions[i].Type[:])
						if tipo != "E" {
							isPartExtend = true
						} else {
							isPartExtend = false
							isName = false
							fmt.Println("Error la partición extendida ya existe")
							fmt.Println("Error imposible crear la partición: ", name)
							break
						}
					}
				}
				if isName {
					for i := 0; i < 4; i++ {
						nombre := Structs.GETNOM(string(mbr.Partitions[i].Name[:]))
						if nombre == name {
							isName = false
							fmt.Println("Error partición ya registrada ", name)
							fmt.Println("Error imposible crear la partición: ", name)
							break
						}
					}
				}

				if isName {
					var partExtendida Structs.Partition
					if string(mbr.Partitions[0].Type[:]) == "E" {
						partExtendida = mbr.Partitions[0]
					} else if string(mbr.Partitions[1].Type[:]) == "E" {
						partExtendida = mbr.Partitions[1]
					} else if string(mbr.Partitions[2].Type[:]) == "E" {
						partExtendida = mbr.Partitions[2]
					} else if string(mbr.Partitions[3].Type[:]) == "E" {
						partExtendida = mbr.Partitions[3]
					}
					if partExtendida.Size != 0 {
						var actual Structs.EBR
						if err := Herramientas.ReadObj(disco, &actual, int64(partExtendida.Start)); err != nil {
							return
						}
						if Structs.GETNOM(string(actual.Name[:])) == name {
							isName = false
						} else {
							for actual.Next != -1 {
								if err := Herramientas.ReadObj(disco, &actual, int64(actual.Next)); err != nil {
									return
								}
								if Structs.GETNOM(string(actual.Name[:])) == name {
									isName = false
									break
								}
							}
						}
						if !isName {
							fmt.Println("Error partición: ", name, " ya registrada")
							fmt.Println("Error imposible crear la partición: ", name)
							return
						}
					}

				}
				sizeNewPart := size * unit
				guardar := false
				var newPart Structs.Partition
				if (typee == "P" || isPartExtend) && isName {
					sizeMBR := int32(binary.Size(mbr))
					mbr, newPart = fstAJ(mbr, typee, sizeMBR, int32(sizeNewPart), name, fit)
					guardar = newPart.Size != 0
					if guardar {
						if err := Herramientas.WrObj(disco, mbr, 0); err != nil {
							return
						}
						if isPartExtend {
							var ebr Structs.EBR
							ebr.Start = newPart.Start
							ebr.Next = -1
							if err := Herramientas.WrObj(disco, ebr, int64(ebr.Start)); err != nil {
								return
							}
						}
						var TempMBR2 Structs.MBR
						if err := Herramientas.ReadObj(disco, &TempMBR2, 0); err != nil {
							return
						}
						fmt.Println("Partición " + name + " creada de manera exitosa")
						Structs.PrintMBR(TempMBR2)
					} else {
						fmt.Println("Error imposible crear la partición: ", name)
					}
				} else if typee == "L" && isName {
					var partExtend Structs.Partition
					if string(mbr.Partitions[0].Type[:]) == "E" {
						partExtend = mbr.Partitions[0]
					} else if string(mbr.Partitions[1].Type[:]) == "E" {
						partExtend = mbr.Partitions[1]
					} else if string(mbr.Partitions[2].Type[:]) == "E" {
						partExtend = mbr.Partitions[2]
					} else if string(mbr.Partitions[3].Type[:]) == "E" {
						partExtend = mbr.Partitions[3]
					} else {
						fmt.Println("Error no existe la partición extendida")
					}
					if partExtend.Size != 0 {
						fstAJLogicas(disco, partExtend, int32(sizeNewPart), name, fit)
					}
				}
			} else if opc == 1 {
				add = add * unit
				if add < 0 {
					fmt.Println("Reducir espacio")
					reducir := true
					for i := 0; i < 4; i++ {
						nombre := Structs.GETNOM(string(mbr.Partitions[i].Name[:]))
						if nombre == name {
							reducir = false
							newSize := mbr.Partitions[i].Size + int32(add)
							if newSize > 0 {
								mbr.Partitions[i].Size += int32(add)
								if err := Herramientas.WrObj(disco, mbr, 0); err != nil {
									return
								}
								fmt.Println("ha sido creada la partición", name)
							} else {
								fmt.Println("Error partición a eliminar demasiado grande")
							}
							break
						}
					}
					if reducir {
						var partExtendida Structs.Partition
						if string(mbr.Partitions[0].Type[:]) == "E" {
							partExtendida = mbr.Partitions[0]
						} else if string(mbr.Partitions[1].Type[:]) == "E" {
							partExtendida = mbr.Partitions[1]
						} else if string(mbr.Partitions[2].Type[:]) == "E" {
							partExtendida = mbr.Partitions[2]
						} else if string(mbr.Partitions[3].Type[:]) == "E" {
							partExtendida = mbr.Partitions[3]
						}
						if partExtendida.Size != 0 {
							var actual Structs.EBR
							if err := Herramientas.ReadObj(disco, &actual, int64(partExtendida.Start)); err != nil {
								return
							}
							if Structs.GETNOM(string(actual.Name[:])) == name {
								reducir = false
							} else {
								for actual.Next != -1 {
									if err := Herramientas.ReadObj(disco, &actual, int64(actual.Next)); err != nil {
										return
									}
									if Structs.GETNOM(string(actual.Name[:])) == name {
										reducir = false
										break
									}
								}
							}
							if !reducir {
								actual.Size += int32(add)
								if actual.Size > 0 {
									if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
										return
									}
									fmt.Println("Se har reducido la partición ", name)
								} else {
									fmt.Println("Error partición a eliminar demasiado grande")
								}
							}
						}
					}

					if reducir {
						fmt.Println("Error partición a reducir no localizada")
					}
				} else if add > 0 {
					fmt.Println("aumentar espacio")
					evaluar := 0
					if Structs.GETNOM(string(mbr.Partitions[0].Name[:])) == name {
						if mbr.Partitions[1].Start == 0 {
							if mbr.Partitions[2].Start == 0 {
								if mbr.Partitions[3].Start == 0 {
									evaluar = int(mbr.MBRSZ - mbr.Partitions[0].GETEND())
								} else {
									evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[0].GETEND())
								}
							} else {
								evaluar = int(mbr.Partitions[2].Start - mbr.Partitions[0].GETEND())
							}
						} else {
							evaluar = int(mbr.Partitions[1].Start - mbr.Partitions[0].GETEND())
						}
						if evaluar > 0 && add <= evaluar {
							mbr.Partitions[0].Size += int32(add)
							if err := Herramientas.WrObj(disco, mbr, 0); err != nil {
								return
							}
							fmt.Println("Se ha aumentado el espacio de la partición ", name)
						} else {
							fmt.Println("Error imposible aumentar el espacio a la partición ", name)
						}
					} else if Structs.GETNOM(string(mbr.Partitions[1].Name[:])) == name {
						if mbr.Partitions[2].Start == 0 {
							if mbr.Partitions[3].Start == 0 {
								evaluar = int(mbr.MBRSZ - mbr.Partitions[1].GETEND())
							} else {
								evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[1].GETEND())
							}
						} else {
							evaluar = int(mbr.Partitions[2].Start - mbr.Partitions[1].GETEND())
						}
						if evaluar > 0 && add <= evaluar {
							mbr.Partitions[1].Size += int32(add)
							if err := Herramientas.WrObj(disco, mbr, 0); err != nil {
								return
							}
							fmt.Println("Espacio auemntado en la partición ", name, " exitosamente")
						} else {
							fmt.Println("Error imposible aumentar espacio esapacio a asignar demasidado grande para la partición ", name)
						}
					} else if Structs.GETNOM(string(mbr.Partitions[2].Name[:])) == name {
						if mbr.Partitions[3].Start == 0 {
							evaluar = int(mbr.MBRSZ - mbr.Partitions[2].GETEND())
						} else {
							evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[2].GETEND())
						}
						if evaluar > 0 && add <= evaluar {
							mbr.Partitions[2].Size += int32(add)
							if err := Herramientas.WrObj(disco, mbr, 0); err != nil {
								return
							}
							fmt.Println("Espacio auemntado en la partición ", name, " exitosamente")
						} else {
							fmt.Println("Error imposible aumentar espacio esapacio a asignar demasidado grande para la partición ", name)
						}
					} else if Structs.GETNOM(string(mbr.Partitions[3].Name[:])) == name {
						evaluar = int(mbr.MBRSZ - mbr.Partitions[3].GETEND())
						if evaluar > 0 && add <= evaluar {
							mbr.Partitions[3].Size += int32(add)
							if err := Herramientas.WrObj(disco, mbr, 0); err != nil {
								return
							}
							fmt.Println("Espacio auemntado en la partición ", name, " exitosamente")
						} else {
							fmt.Println("Error imposible aumentar espacio esapacio a asignar demasidado grande para la partición ", name)
						}
					} else {
						var partExtendida Structs.Partition
						if string(mbr.Partitions[0].Type[:]) == "E" {
							partExtendida = mbr.Partitions[0]
						} else if string(mbr.Partitions[1].Type[:]) == "E" {
							partExtendida = mbr.Partitions[1]
						} else if string(mbr.Partitions[2].Type[:]) == "E" {
							partExtendida = mbr.Partitions[2]
						} else if string(mbr.Partitions[3].Type[:]) == "E" {
							partExtendida = mbr.Partitions[3]
						}
						if partExtendida.Size != 0 {
							aumentar := false
							var actual Structs.EBR
							if err := Herramientas.ReadObj(disco, &actual, int64(partExtendida.Start)); err != nil {
								return
							}
							if Structs.GETNOM(string(actual.Name[:])) == name {
								aumentar = true
							} else {
								for actual.Next != -1 {
									if err := Herramientas.ReadObj(disco, &actual, int64(actual.Next)); err != nil {
										return
									}
									if Structs.GETNOM(string(actual.Name[:])) == name {
										aumentar = true
										break
									}
								}
							}
							if aumentar {
								if actual.Next != -1 {
									if add <= int(actual.Next)-int(actual.GETEND()) {
										actual.Size += int32(add)
										if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
											return
										}
										fmt.Println("Espacio auemntado en la partición ", name, " exitosamente")
									} else {
										fmt.Println("Error imposible aumentar espacio esapacio a asignar demasidado grande para la partición ", name)
									}
								} else {
									if add <= int(partExtendida.GETEND())-int(actual.GETEND()) {
										actual.Size += int32(add)
										if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
											return
										}
										fmt.Println("Espacio auemntado en la partición ", name, " exitosamente")
									} else {
										fmt.Println("Error imposible aumentar espacio esapacio a asignar demasidado grande para la partición ", name)
									}
								}
							} else {
								fmt.Println("Error particion a aumentar no registrada")
							}
						} else {
							fmt.Println("Error partición extendida no registrada")
						}
					}
				} else {
					fmt.Println("Error debe ingresar un valor mayo a 0 para aumentar o disminuir particiones")
				}
			} else if opc == 2 {
				fmt.Println("eliminar particion")
				del := true
				for i := 0; i < 4; i++ {
					nombre := Structs.GETNOM(string(mbr.Partitions[i].Name[:]))
					if nombre == name {
						var newPart Structs.Partition
						mbr.Partitions[i] = newPart
						if err := Herramientas.WrObj(disco, mbr, 0); err != nil {
							return
						}
						del = false
						fmt.Println("Se elimino la partición ", name, " Exitosamente")
						break
					}
				}
				if del {
					var partExtendida Structs.Partition
					if string(mbr.Partitions[0].Type[:]) == "E" {
						partExtendida = mbr.Partitions[0]
					} else if string(mbr.Partitions[1].Type[:]) == "E" {
						partExtendida = mbr.Partitions[1]
					} else if string(mbr.Partitions[2].Type[:]) == "E" {
						partExtendida = mbr.Partitions[2]
					} else if string(mbr.Partitions[3].Type[:]) == "E" {
						partExtendida = mbr.Partitions[3]
					}
					if partExtendida.Size != 0 {
						var actual Structs.EBR
						if err := Herramientas.ReadObj(disco, &actual, int64(partExtendida.Start)); err != nil {
							return
						}
						var anterior Structs.EBR
						var eliminar Structs.EBR
						if Structs.GETNOM(string(actual.Name[:])) == name {
							fmt.Println("Ingreso en el primer EBR")
							fmt.Println("Nombre primer EBR ", Structs.GETNOM(string(actual.Name[:])))
							del = false
						} else {
							for actual.Next != -1 {
								if err := Herramientas.ReadObj(disco, &anterior, int64(actual.Start)); err != nil {
									return
								}
								if err := Herramientas.ReadObj(disco, &actual, int64(actual.Next)); err != nil {
									return
								}
								if Structs.GETNOM(string(actual.Name[:])) == name {
									del = false
									break
								}
							}
						}
						if !del {
							sizeEBR := int32(binary.Size(actual))
							if actual.Next != -1 {
								if anterior.Size == 0 {
									actual.Size = 0
									actual.Name = eliminar.Name
									fmt.Println("Nombre modificado ", Structs.GETNOM(string(actual.Name[:])))
									if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
										return
									}
									if err := Herramientas.WrObj(disco, Herramientas.DelRaw1(actual.Size), int64(actual.Start+sizeEBR)); err != nil {
										return
									}
									fmt.Println("Se elimino la partición ", name, " Exitosamente")
								} else {
									anterior.Next = actual.Next
									if err := Herramientas.WrObj(disco, anterior, int64(anterior.Start)); err != nil {
										return
									}
									if err := Herramientas.WrObj(disco, Herramientas.DelRaw1(actual.Size+sizeEBR), int64(actual.Start)); err != nil {
										return
									}
									fmt.Println("Se elimino la partición ", name, " Exitosamente")
								}
							} else {
								if anterior.Size == 0 {
									actual.Size = 0
									actual.Name = eliminar.Name
									fmt.Println("Nombre cambiado ", Structs.GETNOM(string(actual.Name[:])))
									if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
										return
									}
									fmt.Println("Se elimino la partición ", name, " Exitosamente")
								} else {
									anterior.Next = -1
									if err := Herramientas.WrObj(disco, anterior, int64(anterior.Start)); err != nil {
										return
									}
									if err := Herramientas.WrObj(disco, Herramientas.DelRaw1(actual.Size+sizeEBR), int64(actual.Start)); err != nil {
										return
									}
									fmt.Println("Se elimino la partición ", name, " Exitosamente")
								}
							}
						}
					} else {
						fmt.Println("Error imposible localizar la partición lógica sin una particion extendida: ", name)
					}
				}
				if del {
					fmt.Println("Error no se localizo la partición que pretende eliminar.")
				}

			} else {
				fmt.Println("Error inesperado verificar.")
			}
			defer disco.Close()
			fmt.Println("======End FDISK======")
		} else {
			fmt.Println("Error no se detecto el parametro letter y/o name")
		}
	}
}

func fstAJ(mbr Structs.MBR, typee string, sizeMBR int32, sizeNewPart int32, name string, fit string) (Structs.MBR, Structs.Partition) {
	var newPart Structs.Partition
	var noPart Structs.Partition
	if mbr.Partitions[0].Size == 0 {
		newPart.SETINF(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if mbr.Partitions[1].Size == 0 {
			if mbr.Partitions[2].Size == 0 {
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MBRSZ-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart = noPart
						fmt.Println("Error Espacio insuficiente")
					}
				} else {
					if sizeNewPart <= mbr.Partitions[3].Start-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart.SETINF(typee, fit, mbr.Partitions[3].GETEND(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MBRSZ-newPart.Start {
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							mbr.Partitions[2].Correlativo = 3
						} else {
							newPart = noPart
							fmt.Println("Error Espacio insuficiente")
						}
					}
				}
			} else {
				if sizeNewPart <= mbr.Partitions[2].Start-sizeMBR {
					mbr.Partitions[0] = newPart
				} else {
					newPart.SETINF(typee, fit, mbr.Partitions[2].GETEND(), sizeNewPart, name, 4)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MBRSZ-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("Error Espacio insuficiente")
						}
					} else {
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							mbr.Partitions[1].Correlativo = 2
							mbr.Partitions[2].Correlativo = 3
						} else if sizeNewPart <= mbr.MBRSZ-mbr.Partitions[3].GETEND() {
							newPart.SETINF(typee, fit, mbr.Partitions[3].GETEND(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							mbr.Partitions[1].Correlativo = 2
							mbr.Partitions[2].Correlativo = 3
						} else {
							newPart = noPart
							fmt.Println("Error Espacio insuficiente")
						}
					}
				}
			}
		} else {
			if sizeNewPart <= mbr.Partitions[1].Start-sizeMBR {
				mbr.Partitions[0] = newPart
			} else {
				newPart.SETINF(typee, fit, mbr.Partitions[1].GETEND(), sizeNewPart, name, 3)
				if mbr.Partitions[2].Size == 0 {
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MBRSZ-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							newPart = noPart
							fmt.Println("Error Espacio insuficiente")
						}
					} else {
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							newPart.SETINF(typee, fit, mbr.Partitions[3].GETEND(), sizeNewPart, name, 4)
							if sizeNewPart <= mbr.MBRSZ-newPart.Start {
								mbr.Partitions[2] = mbr.Partitions[3]
								mbr.Partitions[3] = newPart
								mbr.Partitions[2].Correlativo = 3
							} else {
								newPart = noPart
								fmt.Println("Error Espacio insuficiente")
							}
						}
					}
				} else {
					if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
						mbr.Partitions[0] = mbr.Partitions[1]
						mbr.Partitions[1] = newPart
						mbr.Partitions[0].Correlativo = 1
						mbr.Partitions[1].Correlativo = 2
					} else if mbr.Partitions[3].Size == 0 {
						newPart.SETINF(typee, fit, mbr.Partitions[2].GETEND(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MBRSZ-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("Error Espacio insuficiente")
						}
					} else {
						newPart.SETINF(typee, fit, mbr.Partitions[2].GETEND(), sizeNewPart, name, 3)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							mbr.Partitions[0].Correlativo = 1
							mbr.Partitions[1].Correlativo = 2
						} else if sizeNewPart <= mbr.MBRSZ-mbr.Partitions[3].GETEND() {
							newPart.SETINF(typee, fit, mbr.Partitions[3].GETEND(), sizeNewPart, name, 4)
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							mbr.Partitions[0].Correlativo = 1
							mbr.Partitions[1].Correlativo = 2
							mbr.Partitions[2].Correlativo = 3
						} else {
							newPart = noPart
							fmt.Println("Error Espacio insuficiente")
						}
					}
				}
			}
		}
	} else if mbr.Partitions[1].Size == 0 {
		newPart.SETINF(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			mbr.Partitions[1].Correlativo = 2
		} else {
			newPart.SETINF(typee, fit, mbr.Partitions[0].GETEND(), sizeNewPart, name, 2)
			if mbr.Partitions[2].Size == 0 {
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MBRSZ-newPart.Start {
						mbr.Partitions[1] = newPart
					} else {
						newPart = noPart
						fmt.Println("Error Espacio insuficiente")
					}
				} else {
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[1] = newPart
					} else if sizeNewPart <= mbr.MBRSZ-mbr.Partitions[3].GETEND() {
						newPart.SETINF(typee, fit, mbr.Partitions[3].GETEND(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						mbr.Partitions[2].Correlativo = 3
					} else {
						newPart = noPart
						fmt.Println("Error Espacio insuficiente")
					}
				}
			} else {
				if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
					mbr.Partitions[1] = newPart
				} else {
					newPart.SETINF(typee, fit, mbr.Partitions[2].GETEND(), sizeNewPart, name, 3)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MBRSZ-newPart.Start {
							mbr.Partitions[3] = newPart
							mbr.Partitions[3].Correlativo = 4
						} else {
							newPart = noPart
							fmt.Println("Error Espacio insuficiente")
						}
					} else {
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							mbr.Partitions[1].Correlativo = 2
						} else if sizeNewPart <= mbr.MBRSZ-mbr.Partitions[3].GETEND() {
							newPart.SETINF(typee, fit, mbr.Partitions[3].GETEND(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							mbr.Partitions[1].Correlativo = 2
							mbr.Partitions[2].Correlativo = 3
						} else {
							newPart = noPart
							fmt.Println("Error Espacio insuficiente")
						}
					}
				}
			}
		}
	} else if mbr.Partitions[2].Size == 0 {
		newPart.SETINF(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			mbr.Partitions[2].Correlativo = 3
			mbr.Partitions[1].Correlativo = 2
		} else {
			newPart.SETINF(typee, fit, mbr.Partitions[0].GETEND(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				mbr.Partitions[2].Correlativo = 3
			} else {
				newPart.SETINF(typee, fit, mbr.Partitions[1].GETEND(), sizeNewPart, name, 3)
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MBRSZ-newPart.Start {
						mbr.Partitions[2] = newPart
					} else {
						newPart = noPart
						fmt.Println("Error Espacio insuficiente")
					}
				} else {
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[2] = newPart
					} else if sizeNewPart <= mbr.MBRSZ-mbr.Partitions[3].GETEND() {
						newPart.SETINF(typee, fit, mbr.Partitions[3].GETEND(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						mbr.Partitions[2].Correlativo = 3
					} else {
						newPart = noPart
						fmt.Println("Error Espacio insuficiente")
					}
				}
			}
		}
	} else if mbr.Partitions[3].Size == 0 {
		newPart.SETINF(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[3] = mbr.Partitions[2]
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			mbr.Partitions[3].Correlativo = 4
			mbr.Partitions[2].Correlativo = 3
			mbr.Partitions[1].Correlativo = 2
		} else {
			newPart.SETINF(typee, fit, mbr.Partitions[0].GETEND(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				mbr.Partitions[3].Correlativo = 4
				mbr.Partitions[2].Correlativo = 3
			} else if sizeNewPart <= mbr.Partitions[2].Start-mbr.Partitions[1].GETEND() {
				newPart.SETINF(typee, fit, mbr.Partitions[1].GETEND(), sizeNewPart, name, 3)
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = newPart
				mbr.Partitions[3].Correlativo = 4
			} else if sizeNewPart <= mbr.MBRSZ-mbr.Partitions[2].GETEND() {
				newPart.SETINF(typee, fit, mbr.Partitions[2].GETEND(), sizeNewPart, name, 4)
				mbr.Partitions[3] = newPart
			} else {
				newPart = noPart
				fmt.Println("Error Espacio insuficiente")
			}
		}
	} else {
		newPart = noPart
		fmt.Println("Error la partición ya no esta disponible")
	}

	return mbr, newPart
}

func fstAJLogicas(disco *os.File, partExtend Structs.Partition, sizeNewPart int32, name string, fit string) {
	save := true
	var actual Structs.EBR
	sizeEBR := int32(binary.Size(actual))
	if err := Herramientas.ReadObj(disco, &actual, int64(partExtend.Start)); err != nil {
		return
	}
	if actual.Size == 0 {
		if actual.Next == -1 {
			if sizeNewPart+sizeEBR <= partExtend.Size {
				actual.SETINF(fit, partExtend.Start, sizeNewPart, name, -1)
				if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				save = false
				fmt.Println("Se ha creado la partición ", name, " exitosamente")
			} else {
				fmt.Println("Error el espacio es insuficiente")
			}
		} else {
			disponible := actual.Next - partExtend.Start
			if sizeNewPart+sizeEBR <= disponible {
				actual.SETINF(fit, partExtend.Start, sizeNewPart, name, actual.Next)
				if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				save = false
				fmt.Println("Se ha creado la partición ", name, " exitosamente")
			} else {
				fmt.Println("Error Espacio insuficiente")
			}
		}
	}

	if save {
		for actual.Next != -1 {
			if sizeNewPart+sizeEBR <= actual.Next-actual.GETEND() {
				break
			}
			if err := Herramientas.ReadObj(disco, &actual, int64(actual.Next)); err != nil {
				return
			}

		}
		if actual.Next == -1 {
			if sizeNewPart+sizeEBR <= (partExtend.GETEND() - actual.GETEND()) {
				actual.Next = actual.GETEND()
				if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				newStart := actual.GETEND()
				actual.SETINF(fit, newStart, sizeNewPart, name, -1)
				if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				fmt.Println("Se ha creado la partición ", name, " exitosamente")
			} else {
				fmt.Println("Error Espacio insuficiente")
			}
		} else {
			if sizeNewPart+sizeEBR <= (actual.Next - actual.GETEND()) {
				siguiente := actual.Next
				actual.Next = actual.GETEND()
				if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				newStart := actual.GETEND()
				actual.SETINF(fit, newStart, sizeNewPart, name, siguiente)
				if err := Herramientas.WrObj(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				fmt.Println("Se ha creado la partición ", name, " exitosamente")
			} else {
				fmt.Println("Error Espacio insuficiente")
			}
		}
	}
}
