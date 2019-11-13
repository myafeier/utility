package utility

import (
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"testing"
)

func TestOriginGenerateBehavior_Generate(t *testing.T) {

	config:=new(Config)
	//config.AppName=new(SingleText)
	//config.AppName.Title="测试小程序"
	//config.AppName.FontSize=30
	//config.AppName.FontColor=color.Black
	config.Title=new(SingleText)
	config.Title.Title="测试标题标题标题"
	config.Title.FontSize=35
	config.Title.FontColor=color.RGBA{255,0,0,255}

	mainFile,err:=os.OpenFile("./test_main.jpg",os.O_RDONLY,os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	codeFile,err:=os.OpenFile("./test_micro_app_code.png",os.O_RDONLY,os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	bgFile,err:=os.OpenFile("./920x1280.png",os.O_RDONLY,os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	fontData,err:=ioutil.ReadFile("./zh.ttf")
	if err != nil {
		t.Error(err)
	}

	config.MainImage,_,err=image.Decode(mainFile)
	if err != nil {
		t.Error(err)
	}
	config.MicroAppCode,_,err=image.Decode(codeFile)
	if err != nil {
		t.Error(err)
	}

	config.Background,_,err=image.Decode(bgFile)
	if err != nil {
		t.Error(err)
	}
	config.FontData=fontData
	config.PosterWidth=920
	config.PosterHeight=1280

	//config.MainImage=nil
	behavior:=new(OriginGenerateBehavior)
	buf,err:=GeneratePoster(config,behavior)
	if err != nil {
		t.Error(err)
	}
	err=ioutil.WriteFile("./test.png",buf.Bytes(),os.ModePerm)

	if err != nil {
		t.Error(err)
	}

}
