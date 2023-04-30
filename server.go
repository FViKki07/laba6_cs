// авторы: Игорь Стребежев и Андрей Андропов
package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
)

var concurrentClientCount int32 = 0

var numSystem = 10

func main() {
	server, _ := net.Listen("tcp", "127.0.0.1:28563")
	defer server.Close()

	for {
		connection, err := server.Accept()
		// запускаем обработку запросов от подключившегося клиента в новом параллельном потоке
		if err == nil {
			go processClient(connection)
		}
	}
}

func processClient(socket net.Conn) {
	atomic.AddInt32(&concurrentClientCount, 1)
	defer atomic.AddInt32(&concurrentClientCount, -1)
	fmt.Printf("%d concurrent clients are connected\n", atomic.LoadInt32(&concurrentClientCount))

	//неполный запрос от клиентов (продолжение которого еще не доставилось по сети)
	dataForProcessing := ""
	data := make([]byte, 2048)
	defer socket.Close()

	for {
		count, err := socket.Read(data)
		if err == io.EOF || count == 0 {
			break // клиент закрыл соединение
		}
		fmt.Printf("Data received: " + string(data[:count]))
		request := dataForProcessing + string(data[:count])
		queries := strings.Split(request, "\r\n")
		if len(queries) > 1 {
			dataForProcessing = queries[len(queries)-1] // запоминаем последний неполный запрос
			queries = queries[:len(queries)-1]
		}
		for _, query := range queries {
			words := strings.Split(query, " ")
			a := words[0]
			b, _ := strconv.Atoi(words[1])

			if len(words) == 3 {
				numSystem, _ = strconv.Atoi(words[2])
			}

			var answer, _ = ConvertInt(a, b, numSystem)
			fmt.Printf("answer %s \n", answer)
			/*switch {
			  case words[1] == "+": answer = a + b
			  case words[1] == "*": answer = a * b
			}*/
			//v, _ := ConvertInt(bin, 2, 16)
			//var strAnswer = strconv.Itoa(answer) + "\r\n"
			var strAnswer = answer + "\r\n"
			socket.Write([]byte(strAnswer))
			fmt.Printf("Answer: " + strAnswer)
		}
	}
}

func ConvertInt(val string, base, toBase int) (string, error) {
	i, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, toBase), nil
}
