package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Kagami/go-avif"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp" // 官方纯 Go WebP 解码库
)

const version = "1.3.0" // 应用版本号

var (
	qualite     int  // 图像质量
	forceFlag   bool // 新增：强制覆盖标志
	help        bool // 帮助标志
	versionFlag bool // 版本标志
)

// 支持的图像格式
var supportedFormats = map[string]bool{
	".png":  true,
	".bmp":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
	".gif":  true,
}

// 初始化命令行标志
func init() {
	flag.IntVar(&qualite, "q", 80, "设置AVIF图像质量 (1-100)")
	flag.BoolVar(&forceFlag, "f", false, "强制覆盖已存在的AVIF文件")
	flag.BoolVar(&help, "h", false, "显示帮助令牌")
	flag.BoolVar(&versionFlag, "v", false, "显示版本信息")
}

func main() {
	// 解析命令行标志
	flag.Parse()
	// 显示版本信息并退出
	if versionFlag {
		fmt.Println("图片转AVIF工具:", version)
		return
	}
	// 显示帮助信息并退出
	if help || flag.NArg() == 0 {
		showHelp()
		return
	}
	// 验证质量参数
	if qualite < 1 || qualite > 100 {
		fmt.Println("无效的图像质量参数")
		return
	}
	// 获取输入文件列表
	inputFiles := flag.Args()
	// 展开通配符
	expandedFiles, err := expandWildcards(inputFiles)
	if err != nil {
		fmt.Printf("解析文件列表时出错: %v\n", err)
		return
	}
	if len(expandedFiles) == 0 {
		fmt.Println("未找到匹配的输入文件。")
		return
	}
	// 处理每个输入文件
	processFiles(expandedFiles, qualite)
}

// 处理输入的文件列表
func processFiles(files []string, quality int) {
	// 启用协程处理文件
	var wg sync.WaitGroup
	// 获取CPU核心数
	var maxGoroutines = runtime.NumCPU()
	// 限制并发数
	var sem = make(chan bool, maxGoroutines)
	// 总结果数
	resultCount := 0
	// 成功转换数
	successCount := 0
	// 失败转换数
	failCount := 0
	// 显示开始信息
	fmt.Printf("开始转换 %d 个文件,质量: %d\n", len(files), quality)
	// 记录开始转换时间
	startTime := time.Now()
	// 遍历每个文件
	for _, file := range files {
		// 检查文件是否存在
		if !fileExists(file) {
			fmt.Printf("文件 %s 不存在\n", file)
			failCount++
			continue
		}
		// 检查文件格式是否支持
		ext := ext(file)
		if !supportedFormats[ext] {
			fmt.Printf("不支持的文件格式: %s\n", file)
			failCount++
			continue
		}
		// 检查转换过的AVIF文件是否已存在
		outputPath := getOutputPath(file)
		if fileExists(outputPath) {
			if forceFlag {
				fmt.Printf("AVIF文件已存在,强制覆盖: %s\n", outputPath)
				// 继续执行转换流程
			} else {
				fmt.Printf("AVIF文件已存在,跳过: %s\n", outputPath)
				failCount++
				continue
			}
		}
		// 限制并发数
		sem <- true
		wg.Add(1)
		resultCount++
		go func(inputPath, outputPath string, quality int) {
			// 结束时释放信号量和等待组
			defer func() {
				<-sem
				wg.Done()
			}()
			err := convertToAVIF(inputPath, outputPath, quality)
			if err != nil {
				fmt.Printf("转换失败: %s\n", err)
				failCount++
			} else {
				fmt.Printf("转换成功: %s\n", outputPath)
				successCount++
			}
		}(file, outputPath, quality)
	}
	// 等待所有任务完成
	wg.Wait()
	// 记录结束时间
	endTime := time.Now()
	// 显示总结
	fmt.Println()
	fmt.Printf("转换完成:本次转换共处理 %d 个文件,成功 %d 个,失败 %d 个,总耗时 %s 。\n", resultCount, successCount, failCount, endTime.Sub(startTime).String())
}

// expandWildcards 函数用于展开文件路径中的通配符模式,返回匹配的所有文件路径
func expandWildcards(patterns []string) ([]string, error) {
	var files []string
	seen := make(map[string]bool)

	for _, pattern := range patterns {
		// 检查模式是否包含通配符
		if strings.ContainsAny(pattern, "*?[]") {
			// 使用 filepath.Glob 查找匹配的文件
			matches, err := filepath.Glob(pattern)
			if err != nil {
				// 忽略无效的通配符,但打印警告
				fmt.Printf("警告: 无效的文件模式 '%s': %v\n", pattern, err)
				continue
			}
			// 遍历匹配结果,去重后添加到文件列表
			for _, match := range matches {
				if !seen[match] {
					files = append(files, match)
					seen[match] = true
				}
			}
		} else {
			// 如果不是通配符,直接添加到列表
			if !seen[pattern] {
				files = append(files, pattern)
				seen[pattern] = true
			}
		}
	}
	return files, nil
}

func convertToAVIF(inputPath, outputPath string, quality int) error {
	// 打开文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("无法打开输入文件: %v", err)
	}
	defer inputFile.Close()

	// 读图
	img, err := decodeImage(inputFile, ext(inputPath))
	if err != nil {
		return fmt.Errorf("无法读取图片文件: %v", err)
	}
	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()
	// 编码为AVIF格式
	err = encodeAVIF(img, outputFile, quality)
	if err != nil {
		return fmt.Errorf("无法编码为AVIF格式: %v", err)
	}
	return nil
}

// decodeImage
// 根据文件扩展名解码对应的图片格式
func decodeImage(r io.Reader, ext string) (image.Image, error) {
	// 根据文件扩展名选择对应的解码器进行图片解码
	switch ext {
	case ".png":
		return png.Decode(r)
	case ".bmp":
		return bmp.Decode(r)
	case ".jpg", ".jpeg":
		return jpeg.Decode(r)
	case ".webp":
		// 使用官方纯 Go 的 WebP 解码器
		return webp.Decode(r)
	case ".gif":
		return gif.Decode(r)
	}
	// 如果扩展名不匹配任何支持的格式,返回错误信息
	return nil, fmt.Errorf("不支持的文件格式: %s", ext)
}

// encodeAVIF
// 将图像编码为AVIF格式并写入到指定的io.Writer中
func encodeAVIF(img image.Image, w io.Writer, quality int) error {
	// 注意：go-avif 库的 Quality 选项是与压缩比成反比的,数值越小质量越高,0是最高质量
	// 转换逻辑：用户质量 100 -> avif Quality 0 (最高质量)；用户质量 1 -> avif Quality 62 (最低质量)
	avifQuality := 63 - (quality*63)/100 // 转换为 0-63 范围,并反转
	if avifQuality < 0 {
		avifQuality = 0
	} else if avifQuality > 63 {
		avifQuality = 63
	}

	// 创建编码选项
	Options := avif.Options{
		Quality: avifQuality, // avif Quality 0 是最高质量
		Speed:   4,           // 默认速度 (0最慢/最高压缩,10最快/最低压缩)
	}
	// 使用指定选项进行编码
	return avif.Encode(w, img, &Options)
}

// ext 获取指定路径文件的扩展名并转换为小写
func ext(path string) string {
	return strings.ToLower(filepath.Ext(path))
}

// getOutputPath
// 根据输入文件路径生成对应的AVIF输出文件路径
func getOutputPath(inputPath string) string {
	// 提取输入文件的目录路径
	dir := filepath.Dir(inputPath)
	// 获取输入文件的文件名
	filename := filepath.Base(inputPath)
	// 获取文件扩展名
	ext := filepath.Ext(filename)
	// 去除文件名中的扩展名部分
	name := strings.TrimSuffix(filename, ext)
	// 组合并返回新的AVIF文件路径
	return filepath.Join(dir, name+".avif")
}

// fileExists
// 检查指定路径的文件或目录是否存在
func fileExists(path string) bool {
	// 使用os.Stat检查文件状态
	_, err := os.Stat(path)
	// 如果没有错误,文件存在
	return !os.IsNotExist(err)
}

// 显示帮助信息
func showHelp() {
	fmt.Println("图片转AVIF格式工具")
	fmt.Printf("当前版本:%s,作者:猫东东 <https://bsay.de>\n", version)
	fmt.Println("用于将任意的图片文件转换成avif,从而实现图片的压缩功能。")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println(" image2avif [选项] <文件...>")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println(" -q <数值>  设置AVIF质量 (1-100, 默认: 80)")
	fmt.Println(" -f      强制覆盖已存在的AVIF文件")
	fmt.Println(" -v      显示版本信息")
	fmt.Println(" -h      显示帮助信息")
	fmt.Println()
	fmt.Println("支持的格式:")
	fmt.Println(" .png, .bmp, .jpg, .jpeg, .webp, .gif")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println(" image2avif image.jpg")
	fmt.Println(" image2avif -q 90 photo.png")
	fmt.Println(" image2avif *.png *.jpg")
	fmt.Println()
}
