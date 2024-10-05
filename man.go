package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 基础URL
const baseURL = "https://www.51pptmoban.com"

func main() {
	// 使用命令行参数传递 url.txt 的文件路径
	urlFilePath := flag.String("urlfile", "url.txt", "url.txt 文件路径")
	flag.Parse()

	// 从指定的文件中读取 mainURL
	mainURL, err := readURLFromFile(*urlFilePath)
	if err != nil {
		log.Fatalf("读取 %s 文件失败: %v", *urlFilePath, err)
	}

	// 创建下载文件夹
	downloadFolder := "download"
	err = os.MkdirAll(downloadFolder, os.ModePerm)
	if err != nil {
		log.Fatalf("创建下载文件夹失败: %v", err)
	}

	// 请求主页面，获取所有 "pptlist" 中的超链接
	res, err := http.Get(mainURL)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	defer res.Body.Close()

	// 检查状态码
	if res.StatusCode != http.StatusOK {
		log.Fatalf("请求失败，状态码: %d", res.StatusCode)
	}

	// 使用 goquery 解析主页面 HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("解析主页面 HTML 失败: %v", err)
	}

	// 查找 "pptlist" 中的所有链接
	doc.Find("div.pptlist a").Each(func(i int, s *goquery.Selection) {
		// 获取每个链接的 href
		link, exists := s.Attr("href")
		if exists {
			// 如果链接没有前缀，追加 "https://www.51pptmoban.com"
			fullLink := link
			if !strings.HasPrefix(link, "http") {
				fullLink = baseURL + link
			}

			// 打印完整的链接
			fmt.Printf("请求链接: %s\n", fullLink)

			// 请求该链接的页面，获取 "ppt_xz" 中 "down" 的超链接
			getDownLink(fullLink, downloadFolder, strconv.Itoa(i+1))
		}
	})
}

// 从文件中读取 URL
func readURLFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", fmt.Errorf("文件为空")
}

// 请求子页面，解析并获取 "ppt_xz" 中的 "down" 链接
func getDownLink(subPageURL, downloadFolder, index string) {
	// 请求子页面
	res, err := http.Get(subPageURL)
	if err != nil {
		log.Printf("请求子页面失败: %v", err)
		return
	}
	defer res.Body.Close()

	// 检查状态码
	if res.StatusCode != http.StatusOK {
		log.Printf("请求子页面失败，状态码: %d", res.StatusCode)
		return
	}

	// 解析子页面 HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("解析子页面 HTML 失败: %v", err)
		return
	}
	fileBaseName := ""
	// 查找div.wz下的所有a标签
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		// 确保处理的是第二个a标签
		fileBaseName, _ = decodeGB2312ToUTF8(s.Text())
	})
	//如果fileBaseName不为空，则通过逗号分割，获取第一个内容重新赋值fileBaseName
	if fileBaseName != "" {
		fileBaseName = strings.Split(fileBaseName, ",")[0]
	}

	// 查找页面中的 "ppt_xz" 中 class 为 "down" 的链接
	doc.Find("div.ppt_xz a.down").Each(func(i int, s *goquery.Selection) {
		downloadLink, exists := s.Attr("href")
		if exists {
			// 如果链接没有前缀，追加 baseURL
			fullDownloadLink := downloadLink
			if !strings.HasPrefix(downloadLink, "http") {
				fullDownloadLink = baseURL + downloadLink
			}

			fmt.Printf("找到的下载链接: %s\n", fullDownloadLink)

			// 请求 "down" 链接的页面，查找最终下载链接并下载文件
			getFinalDownloadLink(fullDownloadLink, downloadFolder, index, fileBaseName)
		}
	})
}

func decodeGB2312ToUTF8(input string) (string, error) {
	gb2312Decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Reader, _, _ := transform.String(gb2312Decoder, input)
	return utf8Reader, nil
}

// 请求下载页面，解析并获取 "down" 中的 "tjd0" 链接
func getFinalDownloadLink(downPageURL, downloadFolder, index, fileBaseName string) {
	// 请求下载页面
	res, err := http.Get(downPageURL)
	if err != nil {
		log.Printf("请求下载页面失败: %v", err)
		return
	}
	defer res.Body.Close()

	// 检查状态码
	if res.StatusCode != http.StatusOK {
		log.Printf("请求下载页面失败，状态码: %d", res.StatusCode)
		return
	}

	// 解析下载页面 HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("解析下载页面 HTML 失败: %v", err)
		return
	}

	// 查找 "div.wz" 中的 a 标签的值，作为文件名
	doc.Find("div.wz a").Each(func(i int, s *goquery.Selection) {
		fileName := fileBaseName + ".zip"

		// 查找最终的下载链接
		doc.Find("div.down a.tjd0").Each(func(i int, s *goquery.Selection) {
			finalLink, exists := s.Attr("href")
			if exists {
				// 将链接中的 ".." 替换为 "https://www.51pptmoban.com/e/DownSys"
				finalDownloadLink := strings.Replace(finalLink, "..", "https://www.51pptmoban.com/e/DownSys", 1)

				fmt.Printf("最终下载链接: %s\n", finalDownloadLink)

				// 下载文件到下载文件夹中，文件名为拼音标题
				downloadFile(finalDownloadLink, fileName, downloadFolder)
			}
		})
	})
}

// 下载文件到 download 文件夹中
func downloadFile(fileURL, fileName, downloadFolder string) {
	// 请求下载链接
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Printf("下载文件失败: %v", err)
		return
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("请求下载链接失败，状态码: %d", resp.StatusCode)
		return
	}

	// 创建文件路径
	filePath := filepath.Join(downloadFolder, fileName)

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("创建文件失败: %v", err)
		return
	}
	defer file.Close()

	// 将文件内容写入本地文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("写入文件失败: %v", err)
		return
	}

	fmt.Printf("文件 %s 下载成功！\n", fileName)
}
