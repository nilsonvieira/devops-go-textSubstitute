package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	// Lista de nomes de arquivos a serem lidos
	fileNames := []string{"<FILE-A>", "FILE-B"}
	// Diretório onde os arquivos estão localizados
	directory := "<PATH>"

	// Tipo de substituição (Java17 ou Java8)
	javaVersion := "Java8"

	for _, fileName := range fileNames {
		filePath := fmt.Sprintf("%s/%s", directory, fileName)
		if err := processFile(filePath, javaVersion); err != nil {
			fmt.Printf("Erro ao processar o arquivo %s: %v\n", fileName, err)
		}
	}
}

func processFile(filePath, javaVersion string) error {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("não foi possível abrir o arquivo: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(filePath + ".tmp")
	if err != nil {
		return fmt.Errorf("não foi possível criar o arquivo temporário: %w", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, `env.JAVA_OPTIONS="`) {
			switch javaVersion {
			case "Java17":
				line = strings.Replace(line, `env.JAVA_OPTIONS="`, `env.TOOL_JAVA="JAVA_TOOL_OPTIONS"`+"\n"+`env.JAVA_OPTIONS="-Dspring.config.additional-location=file:/app/config/application.properties`, 1)
			case "Java8":
				line = strings.Replace(line, `env.JAVA_OPTIONS="`, `env.TOOL_JAVA="JAVA_OPTIONS"`+"\n"+`env.JAVA_OPTIONS="-Dspring.config.additional-location=file:/app/config/application.properties`, 1)
			}
		}
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("não foi possível escrever no arquivo temporário: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("erro ao ler o arquivo: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("erro ao salvar o arquivo temporário: %w", err)
	}

	if err := os.Rename(filePath+".tmp", filePath); err != nil {
		return fmt.Errorf("não foi possível renomear o arquivo temporário: %w", err)
	}

	return nil
}
