package main

import (
	"fmt"

	"github.com/tidwall/gjson"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func firstFive(localFile string) {

	err := ffmpeg_go.Input(localFile,
		ffmpeg_go.KwArgs{"ss": "0:00"}).Output("/home/fryan/Videos/out.webm", ffmpeg_go.KwArgs{"c": "copy", "to": "2:00"}).OverWriteOutput().Run()

	if err != nil {
		fmt.Println("error ", err)
	}
}

func getVideoLength(localFile string) {
	a, _ := ffmpeg_go.Probe(localFile)

	duration := gjson.Get(a, "format.duration").Float()

	fmt.Printf("Duration: %v", duration)
}
