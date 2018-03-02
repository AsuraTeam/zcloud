package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/runtime/schema"
	restclient "k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	//"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	//"strings"
	"io"
	//"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Handler(r io.Reader,w io.Writer,containername string,podname string,namespace string) {
	//log.SetFlags(log.Llongfile)
	//config, err := clientcmd.BuildConfigFromFlags("", "./config")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	groupversion := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	config := restclient.Config{}
	config.Host = "http://10.16.55.6:8080"
	config.GroupVersion = &groupversion
	config.APIPath = "/api"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	restclient, err := rest.RESTClientFor(&config)
	if err != nil {
		log.Fatalln(err)
	}

	req := restclient.Post().
		Resource("pods").
		Name(podname).
		Namespace(namespace).
		SubResource("exec")
		//Param("container", containername).
		//Param("stdin", "true").
		//Param("stdout", "true").
		//Param("stderr", "true").
		//Param("command", "sh").Param("tty", "true")

	req.VersionedParams(
		&v1.PodExecOptions{
			Container: containername,
			Command:   []string{"sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		},
		scheme.ParameterCodec,
	)

	log.Println(req.URL().String())
	executor, err := remotecommand.NewSPDYExecutor(
		&config, http.MethodPost, req.URL(),
	)
	if err != nil {
		log.Println(err)
	}
	//strings.NewReader("touc /aa.txt")
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:              r,
		Stdout:             w,
		Stderr:             w,
		Tty:                true,
		//TerminalSizeQueue:  Terminal.s,
	})

	if err != nil {
		log.Println(err)
	}
}

func echoHandler(ws *websocket.Conn) {
	Handler(ws, ws,"auto-3","auto-3-7b45c8757c-bngtf","auto-3--dfsad")
}


func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}