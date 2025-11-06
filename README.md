# 图片转AVIF工具 (Go语言版)

一个使用Go语言开发的简单易用的命令行工具，用于将常见图片格式转换为AVIF格式。

## 功能特点

- 支持多种图片格式转换为AVIF
- 可调节AVIF质量参数
- 批量转换支持
- 简单的命令行界面

## 支持的输入格式

- PNG (.png)
- BMP (.bmp)
- GIF (.png)
- JPG (.jpg, .jpeg)
- WebP (.webp)

## 安装说明

### 前提条件

- 无需安装Go语言环境（工具已预编译为可执行文件）

### 安装步骤

1. **下载dist下的可执行文件**

2. **将工具目录添加到系统PATH环境变量**
   - 右键点击"此电脑" → "属性" → "高级系统设置" → "环境变量"
   - 在"系统变量"中找到"Path" → 点击"编辑"
   - 点击"新建" → 输入工具所在目录的完整路径
   - 点击"确定"保存更改

3. **验证安装**
   打开新的命令提示符窗口，输入：
   ```bash
   image2avif --version
   ```
   如果显示版本信息，则安装成功。

## 使用方法

### 基本用法

```bash
image2avif <图片文件>
```

### 示例

1. **转换单个文件**
   ```bash
   image2avif image.jpg
   ```
   这将在同一目录下生成`image.avif`文件。

2. **转换多个文件质量**
   ```bash
   image2avif -q 90 photo.png
   ```
   使用质量参数（1-100，默认80）。

3. **批量转换**
   ```bash
   image2avif *.png *.jpg
   ```
   转换当前目录下所有PNG和JPG文件。

4. **拖放使用**
   将图片文件拖放到到`image2avif`文件上，自动转换为AVIF格式。

### 命令行选项

- `-q <数值>`: 设置AVIF质量 (1-100, 默认: 80)
- `-v`: 显示版本信息
- `-h`: 显示帮助信息

## 注意事项

1. 转换后的文件将保存在原文件的同一目录下，文件名相同，扩展名为`.avif`。

2. 如果目标文件已存在，转换将失败，不会覆盖现有文件。

3. 转换质量越高，生成的文件越大，但图像质量更好。

## 性能提示

- Go语言版本通常比Node.js版本运行更快，占用内存更少
- 对于大量文件转换，可以分批进行，避免占用过多系统资源
- 质量参数设置为80通常是一个很好的平衡点，提供良好的压缩率和图像质量

## 编译使用

要用到cgo，请确保已安装gcc和libavif库。

### MACOS

```bash
brew install libavif aom
```
### Debian/Ubuntu (使用 apt):

```bash
apt-get install libavif-dev libaom-dev
```

### Windows 请使用

[MSYS2](https://www.msys2.org/)

```bash
pacman -Syu
pacman -S mingw-w64-x86_64-libavif mingw-w64-x86_64-aom
go build
```

### 已知错误

   1. 找不到 libaom.dll，以下方法3选1即可

      请点击此处下载：[libaom.dll](https://github.com/xa1st/image2avif/raw/refs/heads/main/dll/libaom.dll)

      如果想自行打包编译，请访问：https://aomedia.googlesource.com/aom

      也可以在[msys2](https://www.msys2.org/)中使用pacman安装libaom-dev

## 许可证
APACHE2.0