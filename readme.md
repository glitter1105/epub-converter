# EPUB 繁体转简体工具

这是一个使用 Go 语言开发的工具，用于将 EPUB 电子书的内容从繁体中文转换为简体中文。

## 使用方法

1. **安装依赖:**

    ```bash
    go mod tidy
    ```

2. **编译:**

    ```bash
    go build
    ```

3. **运行:**

    ```bash
    ./epub-converter -i input.epub -o output.epub
    ```

    *   `-i`: 指定输入的 EPUB 文件路径。
    *   `-o`: 指定输出的 EPUB 文件路径。

## 原理

该工具的主要步骤如下：

1. 解压 EPUB 文件。
2. 遍历解压后的文件，找到 HTML、XHTML 和 OPF 文件。
3. 将 HTML 和 XHTML 文件中的文本内容进行繁体到简体的转换。
4. 将 OPF 文件中的标题、描述和作者等元数据进行繁体到简体的转换。
5. 将转换后的文件重新打包成 EPUB 文件。

## 注意

*   该工具目前仅支持转换 HTML、XHTML 和 OPF 文件中的文本内容。
*   该工具使用 `github.com/liuzl/gocc` 库进行繁简转换，该库使用 Unicode 数据库进行转换。
*   如果遇到转换错误，请检查 EPUB 文件的结构是否正确，或者尝试使用其他 EPUB 阅读器打开该文件。

## 贡献

欢迎提出问题和改进建议，也欢迎提交 Pull Request。

## 许可证

本项目使用 MIT 许可证。