package utils

import (
	"github.com/disintegration/imaging"
	"testing"
)

func TestExampleStream(t *testing.T) {
	reader := ExampleReadFrameAsJpeg("./public/bear.mp4", 5)
	img, err := imaging.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}
	err = imaging.Save(img, "./public/out1.jpeg")
	if err != nil {
		t.Fatal(err)
	}
}
//func TestSaveFrame(t *testing.T) {
//	pCodecCtxOrig := pFormatContext.Streams()[i].Codec()
//	pCodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(pCodecCtxOrig.GetCodecId()))
//	pCodecCtx := pCodec.AvcodecAllocContext3()
//
//	pFrameRGB := avutil.AvFrameAlloc()
//
//	SaveFrame(pFrameRGB, pCodecCtx.Width(), pCodecCtx.Height(), frameNumber)
//}
