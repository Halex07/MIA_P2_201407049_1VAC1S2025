package Comandos

import (
	"MIA_P1_201407049/Structs"
	"fmt"
)

func Logout() {
	fmt.Println(">> Proceso en ejecución Saliendo")
	if !Structs.CurrentUSR.STATUS {
		fmt.Println("LOGOUT ERROR: No existe una sesion iniciada")
		return
	}
	Structs.CurrentUSR.STATUS = false
	fmt.Println("la sesión se ha cerrado de manera exitosa ", Structs.CurrentUSR.Nombre)
	Structs.CurrentUSR.Id = ""
	Structs.CurrentUSR.GRPID = 0
	Structs.CurrentUSR.USRID = 0
	Structs.CurrentUSR.Nombre = ""
}
