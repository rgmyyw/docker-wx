package Xml

import "encoding/xml"

type VideoMsg struct {
	XMLName xml.Name `xml:"msg"`
	Video   struct {
		Aeskey         string `xml:"aeskey,attr"`
		Cdnthumbaeskey string `xml:"cdnthumbaeskey,attr"`
		Cdnvideourl    string `xml:"cdnvideourl,attr"`
		Length         uint32 `xml:"length,attr"`
		Playlength     uint32 `xml:"playlength,attr"`
		Cdnthumblength uint32 `xml:"cdnthumblength,attr"`
		Cdnthumbwidth  uint32 `xml:"cdnthumbwidth,attr"`
		Cdnthumbheight uint32 `xml:"cdnthumbheight,attr"`
		Cdnthumburl    string `xml:"cdnthumburl,attr"`
		Md5            string `xml:"md5,attr"`
		Newmd5         string `xml:"newmd5,attr"`
		Isad           uint32 `xml:"isad,attr"`
	} `xml:"videomsg"`
}
