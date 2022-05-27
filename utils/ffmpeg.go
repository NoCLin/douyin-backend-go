package utils

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
)

func ExampleReadFrameAsJpeg(inFileName string, frameNum int) (io.Reader, int64) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}

	return buf,int64(buf.Len())
}



//func SaveFrame(frame *avutil.Frame, width, height, frameNumber int) {
//	// Open file
//	fileName := fmt.Sprintf("frame%d.ppm", frameNumber)
//	file, err := os.Create(fileName)
//	if err != nil {
//		log.Println("Error Reading")
//	}
//	defer file.Close()
//
//	// Write header
//	header := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
//	file.Write([]byte(header))
//
//	// Write pixel data
//	for y := 0; y < height; y++ {
//		data0 := avutil.Data(frame)[0]
//		buf := make([]byte, width*3)
//		startPos := uintptr(unsafe.Pointer(data0)) + uintptr(y)*uintptr(avutil.Linesize(frame)[0])
//		for i := 0; i < width*3; i++ {
//			element := *(*uint8)(unsafe.Pointer(startPos + uintptr(i)))
//			buf[i] = element
//		}
//		file.Write(buf)
//	}
//}
