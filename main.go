package main 

import (
	"fmt"
	"net/http"
	"strconv"
	"os"
	"bufio"
	"strings"
)

/**************************************************
purpose:

parameters:

notes:
**************************************************/
func main() {

	reader := bufio.NewReader(os.Stdin)

	var running []string

	for {
		fmt.Print(">> ")
		command, err := getNextCommand(reader)
		if err != nil {
			println("Error reading command")
			continue
		}

		if len(command) == 0 {
			continue
		}

		switch string(command[0]) {
		case "run":
			https, port, er := parseRunCommand(command)
			if (er != nil) {
				continue 
			}
			// i somehow need a handle to this 
			go runServer(https, port)
			running = append(running, port)

		case "running":
			for _, port := range running {
				fmt.Println(port)
			}

		case "exit":
			println("Dirty exiting")
			os.Exit(0)

		default:
			println("Invalid command: " + command[0])
		}

	}

}



/**************************************************
purpose:

parameters:

notes:
**************************************************/

func getNextCommand(reader *bufio.Reader) ([]string, error) {

	text, er := reader.ReadString('\n')
	if er != nil {
		return nil, er 
	}

	text = strings.Replace(text, "\n", "", -1)

	command := strings.Fields(text)

	return command, nil
}


/**************************************************
purpose:

parameters:

notes:
**************************************************/

func printRunUseage() {
	fmt.Println("run [http/https] [port, >1000, <10000]")
}

func parseRunCommand(args []string) (bool, string, error) {

	// dummy error for returning... nasty lmao
	_, err := strconv.Atoi("jdjdjd")

	if (len(args) != 3) {
		fmt.Println("not enough command line arguments")
		printRunUseage()
		return false, "", err
	}

	var https bool 

	if (args[1] == "https") {
		https = true 		
	} else if (args[1] == "http") {
		https = false
	} else {
		fmt.Println("enter either https or http")
		printRunUseage()
		return false, "", err
	}

	port, e := strconv.Atoi(args[2])
	if e != nil {
		fmt.Println("Enter an integer port number")
		printRunUseage()
		return false, "", err

	}
	if (port <= 1000 || port >= 10000) {
		fmt.Println("1000 >= port >= 10000")
		printRunUseage()
		return false, "", err
	}

	return https, args[2], nil
}

/**************************************************
purpose:

parameters:

notes:
**************************************************/


type handler struct {
	port string
}


func runServer(https bool, port string) {

	h := new(handler)
	h.port = port

	// create new server mux, instead of using the deafult net/http one
	// this is so you can run multiple at the same time
	server := http.NewServeMux()

	server.Handle("/", h)

	var err error

	fmt.Println("Listening...")

	if (https) {
		fmt.Println("running https server on port " + port)
		err = http.ListenAndServeTLS(":" + port, "https-server.crt", "https-server.key", server);
	} else {
		fmt.Println("running http server on port " + port)
		err = http.ListenAndServe(":" + port, server)
	}

	if (err != nil) {
		fmt.Println("Error starting server")
	}


}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		//println("Got / on port " + h.port)
		fmt.Fprintf(w, "/ on port " + h.port)
	case "/page1":
		//println("Got /page1 on port " + h.port)
		fmt.Fprintf(w, "page1 on port " + h.port)
	case "/page2":
		//println("Got /page2 on port " + h.port)
		fmt.Fprintf(w, "page2 on port " + h.port)
	default:
		//println("404, url " + r.URL.Path + ", port " + h.port)
		fmt.Fprintf(w, "404 on port " + h.port)
	}
}


















