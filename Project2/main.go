package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main(){
	ctx,cancel := context.WithCancel(context.Background())
	srvGroup,errCtx := errgroup.WithContext(ctx)
	// 创建一个 http Server
	srv := &http.Server{Addr:":8080"}
	srvGroup.Go(func()error{
		return httpServer(srv,"hello world")
	})
	srvGroup.Go(func()error{
		<-errCtx.Done()
		fmt.Println("http server stop")
		return srv.Shutdown(errCtx)
	})
	go listenKillSign(cancel)

	err := srvGroup.Wait()
	if err != nil {
		fmt.Printf("App start fail -> %v", err)
	} else {
		fmt.Printf("App shutdown success")
	}

}
func httpServer(server *http.Server,msg string)error{
	http.HandleFunc("/",func(writer http.ResponseWriter,request *http.Request){
		fmt.Fprintln(writer,msg)
	})
	return server.ListenAndServe()
}
func listenKillSign(cancel context.CancelFunc){
	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	cancel()
}
