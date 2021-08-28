package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/viper"

	"github.com/sl4mmer/crawlAndBench/pkg/common"
	 "github.com/sl4mmer/crawlAndBench/pkg/controller"
	"github.com/sl4mmer/crawlAndBench/pkg/rest"
)

func main()  {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	conf := common.Opts{}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigChan
		cancel()
	}()

	logic := controller.NewService(&conf)
	go logic.CacheClearanceLoop()

	handler := rest.Service{Queerer: logic}
	router := http.NewServeMux()
	router.HandleFunc("/sites", handler.Handle)
	srv := http.Server{
		Handler: router,
		Addr: ":8080",
	}
	go srv.ListenAndServe()
	<-ctx.Done()
}

