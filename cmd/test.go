/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stress test called")
		url, _ := cmd.Flags().GetString("url")
		requests, _ := cmd.Flags().GetInt("requests")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		loadBalancer(concurrency, requests, url)

		cmd.Help()
	},
}

var requestsRealizados uint64 = 0
var statusCodeOk uint64 = 0
var statusCodeMovedTemporarily uint64 = 0
var statusCodeBadRequest uint64 = 0
var statusCodeNotFound uint64 = 0
var statusCodeInternalServerError uint64 = 0
var statusCode uint64 = 0

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringP("url", "u", "http://google.com", "URL do serviço a ser testado")
	testCmd.Flags().IntP("requests", "r", 1000, "Número total de requests")
	testCmd.Flags().IntP("concurrency", "c", 10, "Número de chamadas simultâneas")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func worker(workerId int, data <-chan string, wg *sync.WaitGroup) {
	for url := range data {
		fmt.Printf("worker %d received url %s\n", workerId, url)

		req, err := http.Get(url)
		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Printf("HTTP Response Status - workerId: %d, %s, %d \n", req.StatusCode, http.StatusText(req.StatusCode), workerId)

		switch req.StatusCode {
		case 200:
			atomic.AddUint64(&statusCodeOk, 1)
		case 302:
			atomic.AddUint64(&statusCodeMovedTemporarily, 1)
		case 400:
			atomic.AddUint64(&statusCodeBadRequest, 1)
		case 404:
			atomic.AddUint64(&statusCodeNotFound, 1)
		case 500:
			atomic.AddUint64(&statusCodeInternalServerError, 1)
		default:
			atomic.AddUint64(&statusCode, 1)
		}

		atomic.AddUint64(&requestsRealizados, 1)

		wg.Done()
	}
}

func publish(ch chan string, url string, requests int) {
	for i := 0; i < requests; i++ {
		ch <- url
	}
	close(ch)
}

func loadBalancer(concurrency, requests int, url string) {
	fmt.Println("url called: " + url)
	fmt.Println("requests called: " + fmt.Sprint(requests))
	fmt.Println("concurrency called: " + fmt.Sprint(concurrency))

	start := time.Now()
	data := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(requests)

	// inicializa os workers
	for i := 0; i < concurrency; i++ {
		go worker(i, data, &wg)
	}

	// publica a url para o worker processar
	publish(data, url, requests)

	wg.Wait()

	fmt.Println("Tempo total gasto na execução:", time.Since(start))
	fmt.Println("Quantidade total de requests realizados:", requestsRealizados)
	fmt.Println("Quantidade de requests com status HTTP 200:", statusCodeOk)
	fmt.Println("Distribuição de outros códigos de status HTTP (como 302, 400, 404, 500, etc.):", statusCodeMovedTemporarily,
		statusCodeBadRequest, statusCodeNotFound, statusCodeInternalServerError, statusCode)
}
