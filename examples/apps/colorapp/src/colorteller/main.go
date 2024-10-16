package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
)

const defaultPort = "8080"
const defaultColor = "black"
const defaultStage = "default"

func getServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port != "" {
		return port
	}

	return defaultPort
}

func getColor() string {
	color := os.Getenv("COLOR")
	if color != "" {
		return color
	}

	return defaultColor
}

func getStage() string {
	stage := os.Getenv("STAGE")
	if stage != "" {
		return stage
	}

	return defaultStage
}

type colorHandler struct{}

func (h *colorHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("color requested, responding with", getColor())
	fmt.Fprint(writer, getColor())
}

type pingHandler struct{}

func (h *pingHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("ping requested, reponding with HTTP 200")
	writer.WriteHeader(http.StatusOK)
}

type retryTestHandler struct{}

func (h *retryTestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a := rand.Intn(100)
	if a < 50 {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "调用接口失败，这是随机返回失败")
		return
	}

	// ���用接口成功，这是��机返回成功
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "调用接口成功，这是随机返回成功")
}

func main() {
	log.Println("starting server, listening on port " + getServerPort())
	xraySegmentNamer := xray.NewFixedSegmentNamer(fmt.Sprintf("%s-colorteller-%s", getStage(), getColor()))
	http.Handle("/", xray.Handler(xraySegmentNamer, &colorHandler{}))
	http.Handle("/ping", xray.Handler(xraySegmentNamer, &pingHandler{}))
	http.Handle("/retryTest", xray.Handler(xraySegmentNamer, &retryTestHandler{}))
	http.ListenAndServe(":"+getServerPort(), nil)
}
