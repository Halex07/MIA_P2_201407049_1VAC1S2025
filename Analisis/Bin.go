package Herramientas

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CrearDisco(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error creando disco path:  ", err)
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		newFile, err := os.Create(path)
		if err != nil {
			fmt.Println("Error creando disco: ", err)
			return err
		}
		defer newFile.Close()
	}
	return nil
}

func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile ==", err)
		return nil, err
	}
	return file, nil
}

func WrObj(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WrObj==", err)
		return err
	}
	return nil
}

func ReadObj(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObj==", err)
		return err
	}
	return nil
}

func ReportGraphizMBR(ruta string, contenido string, nombre string) {
	// Escribimos el archivo .dot
	err := os.WriteFile(ruta, []byte(contenido), 0644)
	if err != nil {
		fmt.Println("Error al escribir archivo .dot:", err)
		return
	}
	fmt.Println("Archivo .dot generado en:", ruta)

	// Ruta al archivo PNG de salida
	rutaPNG := filepath.Join(filepath.Dir(ruta), nombre+".png")

	// Ejecutamos el comando dot -Tpng input.dot -o output.png
	cmd := exec.Command("dot", "-Tpng", ruta, "-o", rutaPNG)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar imagen PNG:", err)
		return
	}

	fmt.Println("Imagen PNG generada en:", rutaPNG)
}

func Reporte(path string, contenido string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error creando el reporte path: ", err)
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return err
	}
	defer file.Close()
	_, err = file.WriteString(contenido)
	if err != nil {
		fmt.Println("Error en la escritura del archivo: ", err)
		return err
	}

	return err
}
func DelRaw1(size int32) []byte {
	datos := make([]byte, size)
	return datos
}

func WriRaw1(size int32) string {
	cad := strings.Repeat("L", int(size))
	return cad
}
