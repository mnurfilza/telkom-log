/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Logs struct {
	Message string `json:"message"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "telkom",
	Short: "A brief description of your application",
	Long:  `Telkom Log adalah cli untuk convert log dari /var/log linux`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		tipe, _ := cmd.Flags().GetString("type")
		output, _ := cmd.Flags().GetString("out")

		dat, err := os.Open(os.Args[1])
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		defer dat.Close()

		var response []string
		scanner := bufio.NewScanner(dat)

		for scanner.Scan() {
			response = append(response, scanner.Text())
		}

		dt := []byte(strings.Join(response, "\n"))
		var outputPath string

		if output != "" {
			outputPath = output
		}

		if len(os.Args) == 3 {
			log.Println("Insert Path Log")
			os.Exit(1)
		} else if len(os.Args) == 2 {
			convertToFileText(dt, outputPath)
		} else {
			//check if exist
			_, err := os.Stat(os.Args[1])
			if err != nil {
				if os.IsNotExist(err) {
					log.Println("FIle Does Not Exists")
					os.Exit(1)
				}
			}

			//check tipe to convert
			switch tipe {
			case "json":
				convertToFileJson(response, outputPath)
			case "text":
				convertToFileText(dt, outputPath)
			default:
				convertToFileText(dt, outputPath)
			}

		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	rootCmd.Flags().StringP("type", "t", "", "command for type convert")
	rootCmd.Flags().StringP("out", "o", "", "command for out to export file")

}

func convertToFileText(req []byte, output string) {

	//if output tidak kosong akan cetak ke path sesuai dengan inputan
	if output != "" {
		outputPath := cleanOutputPath(output)
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			// your file does not exist
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				log.Println(err)
				os.Exit(1)
			}
		}

		if err := os.WriteFile(output, req, 0644); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if _, err := os.Stat(dir + "/output"); os.IsNotExist(err) {
			if err := os.Mkdir(dir+"/output", os.ModePerm); err != nil {
				log.Println(err)
				os.Exit(1)
			}
		}

		fileName := "output_" + time.Now().Format("2006-01-02_15:04:05") + ".txt"
		if err := os.WriteFile(dir+"/output/"+fileName, req, 0644); err != nil {
			log.Println(err)
			os.Exit(1)
		}

	}
	return
}

func convertToFileJson(req []string, output string) {
	//convert to json
	var fs *os.File
	for _, v := range req {
		dat, err := json.Marshal(&Logs{Message: v})
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if output != "" {

			if _, err := os.Stat(output); os.IsNotExist(err) {
				// your file does not exist
				f, err := os.Create(output)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				fs = f
			}

			w := bufio.NewWriter(fs)
			w.Write(dat)
			w.Flush()

		} else {
			dir, err := os.Getwd()
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			if _, err := os.Stat(dir + "/output"); os.IsNotExist(err) {
				os.MkdirAll(dir+"/output", os.ModePerm)
				fileName := "output_" + time.Now().Format("2006-01-02_15:04:05") + ".json"
				f, err := os.Create(dir + "/output/" + fileName)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				fs = f
			}

			w := bufio.NewWriter(fs)
			w.Write(dat)
			w.Flush()

		}
	}

}

func cleanOutputPath(req string) string {
	var newString []string
	str := strings.Split(req, "/")

	for _, v := range str {
		if v == str[len(str)-1] {
			continue
		}

		newString = append(newString, v)
	}

	return strings.Join(newString, "/")

}
