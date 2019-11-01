package util

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/draw"
	"github.com/nfnt/resize"
	"image/png"
	_ "image/jpeg"
	_ "image/gif"
	"unicode/utf8"
)

type Config struct {
	PosterWidth  int
	PosterHeight int
	Background   image.Image //背景图，可选，默认白背景
	MainImage    image.Image //主图，会等比例缩放为的高度500以内，自适应的图片
	MicroAppCode image.Image //小程序码，会缩放为430 * 430，800像素处居中显示
	Title        *SingleText //主标题 ，位置 主图下30像素，居中显示
	AppName      *SingleText // 小程序名称，750像素处居中显示
	FontData     []byte      //字体文件
}

type SingleText struct {
	Title     string
	FontSize  int
	FontColor color.Color
}

type GenerateBehavior interface {
	Generate(config *Config) (poster *bytes.Buffer, err error)
}

func GeneratePoster(config *Config, behavior GenerateBehavior) (poster *bytes.Buffer, err error) {
	return behavior.Generate(config)
}

type OriginGenerateBehavior struct {
	Board draw.Image
}


// 默认显示模式
//  BackGround    //背景图，可选，默认白背景,920* 1280
//	MainImage     //主图，会等比例缩放至 高度500，宽度920 以内，自适应的图片高度或长度，满屏显示尺寸： 920*500
//	MicroAppCode //小程序码，会缩放为430 * 430，800像素以下 居中显示
//	Title        //主标题 ，主图下30像素，居中显示
//	AppName      // 小程序名称，800像素以下居中显示
func (self *OriginGenerateBehavior) Generate(config *Config) (poster *bytes.Buffer, err error) {
	self.Board = image.NewNRGBA(image.Rect(0, 0, config.PosterWidth, config.PosterHeight))
	if config.Background == nil {
		draw.Draw(self.Board, image.Rect(0, 0, config.PosterWidth, config.PosterHeight), image.NewUniform(color.White), image.Point{0, 0}, draw.Over)
	} else {
		draw.Draw(self.Board, image.Rect(0, 0, config.PosterWidth, config.PosterHeight), config.Background, image.Point{0, 0}, draw.Over)
	}
	err=config.ScaleMainImage(config.PosterWidth,500)
	if err != nil {
		return
	}
	err=config.ScaleMicroAppCode(430,430)
	if err != nil {
		return
	}
	fmt.Println(config.MainImage.Bounds().Max)
	fmt.Println(config.MicroAppCode.Bounds().Max)

	fmt.Println(-(config.PosterWidth-config.MainImage.Bounds().Max.X)/2,-30)
	fmt.Println(-(config.PosterWidth-config.MicroAppCode.Bounds().Max.X)/2,-800)
	draw.Draw(self.Board,image.Rect(0,0,config.PosterWidth,config.PosterHeight),config.MainImage,image.Point{-(config.PosterWidth-config.MainImage.Bounds().Max.X)/2,-30},draw.Over)
	draw.Draw(self.Board,image.Rect(0,0,config.PosterWidth,config.PosterHeight),config.MicroAppCode,image.Point{-(config.PosterWidth-config.MicroAppCode.Bounds().Max.X)/2,-800},draw.Over)


	//写入文字

	font,err:=freetype.ParseFont(config.FontData)
	if err != nil {
		return
	}



	c:=freetype.NewContext()
	c.SetFont(font)
	c.SetDPI(72)
	c.SetClip(self.Board.Bounds())
	c.SetDst(self.Board)


	c.SetFontSize(float64(config.Title.FontSize))
	c.SetSrc(image.NewUniform(config.Title.FontColor))
	//article.PostTitle
	_,err=c.DrawString(config.Title.Title,freetype.Pt((self.Board.Bounds().Max.X-int(c.PointToFixed(float64(config.Title.FontSize))>>6)*utf8.RuneCountInString(config.Title.Title))/2, config.MainImage.Bounds().Max.Y+60+int(c.PointToFixed(float64(config.Title.FontSize))>>6)/2))
	if err != nil {
		return
	}

	c.SetFontSize(float64(config.AppName.FontSize))
	c.SetSrc(image.NewUniform(config.AppName.FontColor))
	//article.PostTitle
	_,err=c.DrawString(config.AppName.Title,freetype.Pt((self.Board.Bounds().Max.X-int(c.PointToFixed(float64(config.AppName.FontSize))>>6)*utf8.RuneCountInString(config.AppName.Title))/2, 750+int(c.PointToFixed(float64(config.Title.FontSize))>>6)/2))
	if err != nil {
		return
	}

	poster=new(bytes.Buffer)

	err=png.Encode(poster,self.Board)
	if err != nil {
		return
	}

	return
}




func (self *Config)ScaleMicroAppCode(width,height int)(err error){
	if self.MicroAppCode==nil{
		err=fmt.Errorf("小程序码不存在")
		return
	}
	if self.MicroAppCode.Bounds().Max.X<width||self.MicroAppCode.Bounds().Max.Y<height{
		err=fmt.Errorf("小程序码尺寸太小，会损失精度")
		return
	}

	if float32(self.MicroAppCode.Bounds().Max.X)/float32(self.MicroAppCode.Bounds().Max.Y)>=float32(width)/float32(height) {
		self.MicroAppCode=resize.Resize(uint(width),0,self.MicroAppCode,resize.Lanczos3)
	}else{
		self.MicroAppCode=resize.Resize(0,uint(height),self.MicroAppCode,resize.Lanczos3)
	}
	//t,err:=os.OpenFile("./m1.png",os.O_CREATE|os.O_RDWR,os.ModePerm)
	//err=png.Encode(t,self.MicroAppCode)
	return
}

func (self *Config)ScaleMainImage(width,height int)(err error){
	if self.MainImage==nil{
		err=fmt.Errorf("小程序码不存在")
		return
	}
	if self.MainImage.Bounds().Max.X<=width&&self.MainImage.Bounds().Max.Y<=height{
		//足够小了，没必要缩放
		return
	}

	if float32(self.MainImage.Bounds().Max.X)/float32(self.MainImage.Bounds().Max.Y)>=float32(width)/float32(height) {
		self.MainImage=resize.Resize(uint(width),0,self.MainImage,resize.Lanczos3)
	}else{
		self.MainImage=resize.Resize(0,uint(height),self.MainImage,resize.Lanczos3)
	}
	//t,err:=os.OpenFile("./m2.png",os.O_CREATE|os.O_RDWR,os.ModePerm)
	//err=png.Encode(t,self.MainImage)
	return
}