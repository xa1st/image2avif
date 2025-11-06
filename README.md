# 图片转AVIF工具 (image2avif) 🖼️

![Go Language](https://img.shields.io/badge/language-Go-blue?style=flat-square&logo=go)
![Apache 2.0 License](https://img.shields.io/badge/Apache%202.0-Source-green)
[![GitHub stars](https://img.shields.io/github/stars/xa1st/image2avif.svg?label=Stars&style=flat-square)](https://github.com/xa1st/image2avif)
[![GitHub forks](https://img.shields.io/github/forks/xa1st/image2avif.svg?label=Fork&style=flat-square)](https://github.com/xa1st/image2avif)
[![GitHub issues](https://img.shields.io/github/issues/xa1st/image2avif.svg?label=Issue&style=flat-square)](https://github.com/xa1st/image2avif/issues)
![](https://changkun.de/urlstat?mode=github&repo=xa1st/image2avif)
[![license](https://img.shields.io/badge/license-Apache%202.0-blue.svg?style=flat-square)](https://github.com/xa1st/image2avif/blob/master/LICENSE)

一款基于**Go语言**开发的高效命令行工具，专注于将主流图片格式快速转换为AVIF格式，兼顾压缩效率与图像质量，支持批量处理与并发转换。

## ✨ 核心特性

| 特性                | 说明                                      |
| ----------------- | --------------------------------------- |
| 🎨 **多格式支持**      | 完美兼容 PNG、BMP、JPG/JPEG、WebP、GIF 等主流图像格式 |
| 🚀 **并发转换**      | 自动利用CPU多核性能，并行处理多个文件，大幅提升转换效率    |
| 🔧 **质量可调**      | 支持 1-100 级质量参数调节，平衡文件大小与图像质量      |
| 📦 **批量处理**      | 支持通配符匹配（如 `*.png`），一键转换多个多个文件        |
| ⚡ **强制覆盖**      | 可选强制覆盖已存在的AVIF文件，灵活处理重复转换场景      |
| 🖥️ **跨平台兼容**    | 支持 Linux、macOS、Windows 系统，无需图形界面        |

## 🚀 快速开始

### 🔍 前提条件

- 无需安装Go环境，直接使用预编译可执行文件
- 操作系统：Linux（任意终端）、macOS（Terminal/iTerm2）、Windows（PowerShell/CMD）

### 🛠️ 安装与验证

1. **下载可执行文件**

   从项目发布页下载对应系统的预编译版本，保存到本地目录。

2. **配置环境变量（可选）**

   将工具所在目录添加到系统PATH，实现全局调用。

3. **验证安装**

   ```bash
   image2avif -v
   ```

   若显示版本信息（如 `图片转AVIF工具: 1.2.0`），则安装成功。

## 📖 使用指南

### 基本语法

```bash
image2avif [选项] <文件路径...>
```

### 示例用法

| 场景                | 命令示例                                  | 说明                     |
| ----------------- | ------------------------------------- | ---------------------- |
| 转换单个文件           | `image2avif photo.jpg`                | 在同目录生成 `photo.avif` |
| 指定质量转换           | `image2avif -q 90 picture.png`        | 以质量90转换PNG文件        |
| 批量转换多种格式         | `image2avif *.jpg *.png`              | 转换当前目录所有JPG和PNG文件 |
| 强制覆盖已存在文件        | `image2avif -f -q 70 oldimage.webp`   | 强制覆盖现有AVIF文件       |
| 转换特定目录文件         | `image2avif ~/Pictures/*.gif`         | 转换 Pictures 目录下所有GIF |

### 🧰 命令行参数说明

| 参数         | 简写   | 类型      | 说明                                  |
| ---------- | ---- | ------- | ----------------------------------- |
| `--quality` | `-q` | 数字（1-100） | 可选：设置AVIF图像质量（默认80，数值越高质量越好）   |
| `--force`   | `-f` | 开关      | 可选：强制覆盖已存在的AVIF输出文件              |
| `--help`    | `-h` | 开关      | 可选：显示完整帮助信息                      |
| `--version` | `-v` | 开关      | 可选：显示当前工具版本（如 v1.2.0）           |

## 🛠️ 编译指南（开发者）

如需自行编译源码，需先安装依赖库：

### 环境准备

| 操作系统   | 依赖安装命令                                  |
| ------ | --------------------------------------- |
| macOS  | `brew install libavif aom`              |
| Debian/Ubuntu | `apt-get install libavif-dev libaom-dev` |
| Windows | 使用 [MSYS2](https://www.msys2.org/) 执行：<br>`pacman -S mingw-w64-x86_64-libavif mingw-w64-x86_64-aom` |

### 编译步骤

```bash
# 克隆仓库（假设）
git clone https://github.com/xa1st/image2avif.git
cd image2avif

# 编译可执行文件
go build -o image2avif
```

## ⚠️ 常见问题

1. **Windows系统提示缺少 libaom.dll**
   - 下载 [libaom.dll](https://github.com/xa1st/image2avif/raw/refs/heads/main/dll/libaom.dll)
   - 放置于工具同目录或系统 `System32` 目录

2. **转换失败提示不支持的格式**
   - 检查文件扩展名是否正确（如 `.jpeg` 而非 `.jpe`）
   - 确认文件是否为工具支持的格式（PNG/BMP/JPG/WebP/GIF）

3. **大文件转换耗时过长**
   - 可降低质量参数（如 `-q 60`）
   - 减少并发处理的文件数量

## 🧩 技术栈说明

| 功能模块       | 依赖库           | 作用说明                          |
| ---------- | ------------- | ----------------------------- |
| AVIF编码    | `github.com/Kagami/go-avif` | 核心AVIF格式编码实现               |
| 多格式解码    | 标准库`image`系列 + `webp` | 支持PNG/JPG/GIF等格式的图像解码      |
| 命令行参数解析  | 标准库`flag`     | 处理用户输入的命令行选项与参数        |
| 并发控制      | 标准库`sync`     | 基于CPU核心数限制并发goroutine数量   |
| 文件路径处理    | 标准库`filepath` | 处理通配符匹配、路径解析与输出文件生成  |

## 📄 许可证

本项目基于 **Apache License 2.0** 开源，详见 [LICENSE](LICENSE) 文件。
