package main

import (
	"archive/zip"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/glitter1105/epub-converter/converter"
	"github.com/glitter1105/epub-converter/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputEpub := flag.String("i", "", "输入的 EPUB 文件路径")
	outputEpub := flag.String("o", "", "输出的 EPUB 文件路径")
	flag.Parse()

	if *inputEpub == "" || *outputEpub == "" {
		fmt.Println("请使用 -i 指定输入 EPUB 文件，使用 -o 指定输出 EPUB 文件")
		flag.Usage()
		return
	}

	err := convertEpub(*inputEpub, *outputEpub)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	fmt.Println("转换成功！")
}

func convertEpub(inputEpub, outputEpub string) error {
	// 打开 EPUB 文件
	r, err := zip.OpenReader(inputEpub)
	if err != nil {
		return fmt.Errorf("打开 EPUB 文件失败: %v", err)
	}
	defer r.Close()

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "epub-converter")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			// 记录日志，但不返回错误，避免覆盖原有错误
			log.Printf("删除临时目录失败: %v", err)
		}
	}()

	// 解压 EPUB 文件到临时目录
	err = utils.Unzip(inputEpub, tmpDir)
	if err != nil {
		return fmt.Errorf("解压 EPUB 文件失败: %v", err)
	}

	// 遍历临时目录中的文件
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 转换 HTML 和 XHTML 文件
		if strings.HasSuffix(info.Name(), ".html") || strings.HasSuffix(info.Name(), ".xhtml") {
			fmt.Printf("正在转换: %s\n", path)
			err = convertHTMLFile(path)
			if err != nil {
				return fmt.Errorf("转换 HTML 文件失败: %v", err)
			}
		}

		// 转换 OPF 文件
		if strings.HasSuffix(info.Name(), ".opf") {
			fmt.Printf("正在转换: %s\n", path)
			err = convertOPFFile(path)
			if err != nil {
				return fmt.Errorf("转换 OPF 文件失败: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 将临时目录压缩成新的 EPUB 文件
	err = utils.Zip(tmpDir, outputEpub)
	if err != nil {
		return fmt.Errorf("压缩 EPUB 文件失败: %v", err)
	}

	return nil
}

// 转换 HTML 文件
func convertHTMLFile(path string) error {
	// 读取文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 转换文件内容
	convertedData, err := converter.ConvertString(string(data))
	if err != nil {
		return err
	}

	// 将转换后的内容写回文件
	return os.WriteFile(path, []byte(convertedData), 0644)
}

// 转换 OPF 文件
func convertOPFFile(path string) error {
	// 读取文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 解析 XML
	var opf OPF
	err = xml.Unmarshal(data, &opf)
	if err != nil {
		return err
	}

	// 转换 title 和 description
	opf.Metadata.Title, err = converter.ConvertString(opf.Metadata.Title)
	if err != nil {
		return err
	}
	opf.Metadata.Description, err = converter.ConvertString(opf.Metadata.Description)
	if err != nil {
		return err
	}

	// 转换 creator
	for i, creator := range opf.Metadata.Creator {
		opf.Metadata.Creator[i].Text, err = converter.ConvertString(creator.Text)
		if err != nil {
			return err
		}
	}

	// 转换 manifest 中的 title
	for i, item := range opf.Manifest.Item {
		if item.Properties == "title" {
			opf.Manifest.Item[i].Title, err = converter.ConvertString(item.Title)
			if err != nil {
				return err
			}
		}
	}

	// 将转换后的 OPF 结构体重新编码为 XML
	convertedData, err := xml.MarshalIndent(opf, "  ", "    ")
	if err != nil {
		return err
	}

	// 在 XML 头部添加声明
	finalData := []byte(xml.Header + string(convertedData))

	// 将转换后的内容写回文件
	return os.WriteFile(path, finalData, 0644)
}

// OPF 文件结构体
type OPF struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
}

type Metadata struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Creator     []Creator `xml:"creator"`
}

type Creator struct {
	Text string `xml:",chardata"`
	Role string `xml:"role,attr"`
}

type Manifest struct {
	Item []Item `xml:"item"`
}

type Item struct {
	ID         string `xml:"id,attr"`
	Href       string `xml:"href,attr"`
	MediaType  string `xml:"media-type,attr"`
	Properties string `xml:"properties,attr"`
	Title      string `xml:"title,attr"`
}

type Spine struct {
	Itemref []Itemref `xml:"itemref"`
}

type Itemref struct {
	IDref string `xml:"idref,attr"`
}
