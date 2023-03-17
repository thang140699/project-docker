package main

import (
	utilities "WeddingBackEnd/ultilities"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

var (
	configPrefix string
	container    *Container

	mode         string
	configSource string
)

func main() {
	flag.Parse()
	fmt.Println("MD: ", mode)
	defer utilities.TimeTrack(time.Now(), fmt.Sprintf("Wedding API Service"))
	defer func() {
		fmt.Print("ef ")
		if e := recover(); e != nil {
			log.Panicln(e)
			main()
		}
	}()

	//load env
	var config Config
	err := utilities.LoadEnvFromFile(&config, configPrefix, configSource)
	if err != nil {
		log.Fatalln(err)
	}

	//load container
	container, err = NewContainer(config)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Server is running at : " + config.Binding)
	log.Fatalln(http.ListenAndServe(config.Binding, NewAPIv1(container)))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&configPrefix, "configPrefix", "wedding", "config prefix")
	flag.StringVar(&configSource, "configSource", "./controller/.env", "config source")
}
