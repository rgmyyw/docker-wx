package Xml

import "encoding/xml"

type ImgMsg struct {
	XMLName xml.Name `xml:"msg"`
	Img     struct {
		Aeskey         string `xml:"aeskey,attr"`
		Cdnthumbaeskey string `xml:"cdnthumbaeskey,attr"`
		Cdnthumburl    string `xml:"cdnthumburl,attr"`
		Cdnthumblength int32  `xml:"cdnthumblength,attr"`
		Cdnthumbheight int32  `xml:"cdnthumbheight,attr"`
		Cdnthumbwidth  int32  `xml:"cdnthumbwidth,attr"`
		Cdnmidheight   int32  `xml:"cdnmidheight,attr"`
		Cdnmidwidth    int32  `xml:"cdnmidwidth,attr"`
		Cdnhdheight    int32  `xml:"cdnhdheight,attr"`
		Cdnhdwidth     int32  `xml:"cdnhdwidth,attr"`
		Cdnmidimgurl   string `xml:"cdnmidimgurl,attr"`
		Cdnbigimgurl   string `xml:"cdnbigimgurl,attr"`
		Hdlength       int32  `xml:"hdlength,attr"`
		Length         int32  `xml:"length,attr"`
		Md5            string `xml:"md5,attr"`
		Encryver       int32  `xml:"encryver,attr"`
		Hevcmidsize    int32  `xml:"hevc_mid_size,attr"`
	} `xml:"img"`
}
