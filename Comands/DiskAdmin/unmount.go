package Comandos

import (
	"fmt"
	"strings"

	Herramientas "MIA_P1_201407049/Analisis"
	"MIA_P1_201407049/Structs"
)

func Unmount(parametros []string) {
	fmt.Println(">> Ejecutando Comando UNMOUNT")

	if len(parametros) < 2 {
		fmt.Println("Error: Parámetro insuficiente.")
		return
	}

	// Procesar parámetro id
	param := strings.TrimSpace(parametros[1])
	parts := strings.SplitN(param, "=", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "id") {
		fmt.Printf("Error: Parámetro inválido -> %s\n", param)
		return
	}
	id := strings.ToUpper(strings.TrimSpace(parts[1]))
	if len(id) < 1 {
		fmt.Println("Error: ID vacío.")
		return
	}

	// Ruta del disco
	diskLetter := id[0:1]
	diskPath := fmt.Sprintf("./MIA/P1/%s.dsk", diskLetter)
	disco, err := Herramientas.OpenFile(diskPath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el disco:", err)
		return
	}
	defer disco.Close()

	var mbr Structs.MBR
	if err := Herramientas.ReadObj(disco, &mbr, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	encontrada := false
	for i := 0; i < 4; i++ {
		identificador := Structs.GETID(string(mbr.Partitions[i].Id[:]))
		if identificador == id {
			encontrada = true
			name := Structs.GETNOM(string(mbr.Partitions[i].Name[:]))
			mbr.Partitions[i].Id = [16]byte{} // Limpiar ID
			copy(mbr.Partitions[i].Status[:], "I")
			if err := Herramientas.WrObj(disco, mbr, 0); err != nil {
				fmt.Println("Error al escribir cambios en el disco:", err)
				return
			}
			fmt.Println("Partición", name, "desmontada exitosamente.")
			break
		}
	}

	if !encontrada {
		fmt.Printf("UNMOUNT Error: No se encontró la partición con ID '%s'\n", id)
		return
	}

	fmt.Println("\nLista de particiones montadas:")
	for i, part := range mbr.Partitions {
		if string(part.Status[:]) == "A" {
			fmt.Printf("Partition %d: name: %s, status: %s, id: %s, tipo: %s, start: %d, size: %d, fit: %s, correlativo: %d\n",
				i,
				Structs.GETNOM(string(part.Name[:])),
				string(part.Status[:]),
				string(part.Id[:]),
				string(part.Type[:]),
				part.Start,
				part.Size,
				string(part.Fit[:]),
				part.Correlativo,
			)
		}
	}
}
