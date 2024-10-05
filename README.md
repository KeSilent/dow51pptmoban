# dow51pptmoban

批量下载 51pptmoban 中的免费模板。

# 编译

运行以下命令，生成 Windows 64 位的可执行文件：

```bash
GOOS=windows GOARCH=amd64 go build -o dow51pptmoban.exe
```

运行以下命令，生成 Linux 64 位的可执行文件：

```bash
GOOS=linux GOARCH=amd64 go build -o dow51pptmoban
```

运行以下命令，生成 macOS 64 位的可执行文件：

```bash
GOOS=darwin GOARCH=amd64 go build -o dow51pptmoban
```

# 使用说明

创建一个url.txt，在当中填写一个需要下载模板列表地址，如：https://www.51pptmoban.com/shangwu/

,然后在运行可执行文件的时候，将url.txt地址作为参数传递进去，下载文件会在当前目录中创建的download文件夹中，具体使用方法如下：

## windows 使用说明

### 1. **通过命令行传递参数**

这是最常用的方式。你可以使用 Windows 自带的命令提示符（Command Prompt）或者 PowerShell 来执行并传入参数。

**步骤：**

1.  打开命令提示符或 PowerShell：

    - 按 `Win + R`，输入 `cmd` 或者 `powershell`，然后按回车。

2.  使用 `cd` 命令进入 `.exe` 文件所在的目录。例如：

    ```bash
    cd C:\path\to\your\directory
    ```

3.  执行带有参数的命令：

    ```bash
    dow51pptmoban.exe -urlfile="C:\path\to\your\url.txt"
    ```

### 2. **创建批处理文件 (.bat)**

如果你希望通过双击运行并传递参数，你可以创建一个批处理文件（`.bat`）来执行 `.exe` 文件并传递参数。

**步骤：**

1.  在 `.exe` 文件同一个目录下，新建一个文本文件（例如 `run.bat`），内容如下：

    ```bat
    @echo off
    dow51pptmoban.exe -urlfile="C:\path\to\your\url.txt"
    pause
    ```

2.  保存并关闭文件，将文件后缀改为 `.bat`。

3.  双击 `run.bat` 文件即可执行 `dow51pptmoban.exe`，并传入 `url.txt` 的路径作为参数。

### 3. **使用快捷方式**

如果你不想通过命令行或批处理文件来运行，也可以为 `.exe` 创建一个快捷方式，并在快捷方式中传递参数。

**步骤：**

1.  右键点击 `.exe` 文件，选择“创建快捷方式”。

2.  右键点击刚刚创建的快捷方式，选择“属性”。

3.  在“目标”字段的末尾添加参数。例如：

    ```text
    "C:\path\to\dow51pptmoban.exe" -urlfile="C:\path\to\your\url.txt"
    ```

4.  点击“确定”保存设置。

5.  双击快捷方式，程序就会运行，并带有你指定的参数。

## linux、macOS 使用说明
在命令行中直接执行 "dow51pptmoban  -urlfile="\yourpath\url.txt"