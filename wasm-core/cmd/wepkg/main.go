//go:build js && wasm

package main

import (
	"bytes"
	"encoding/binary"
	"sort"
	"strings"
	"syscall/js"
)

// extractPkg: 解包 (PKG 字节数组 -> JS 对象数组)
func extractPkg(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.Null()
	}
	jsData := args[0]
	length := jsData.Get("byteLength").Int()
	buf := make([]byte, length)
	js.CopyBytesToGo(buf, jsData)
	reader := bytes.NewReader(buf)

	// -- 读取头 --
	var strLen uint32
	if err := binary.Read(reader, binary.LittleEndian, &strLen); err != nil {
		return map[string]any{"error": "Failed to read header"}
	}
	verBuf := make([]byte, strLen)
	reader.Read(verBuf)
	version := string(verBuf)
	if len(version) < 4 || version[:4] != "PKGV" {
		return map[string]any{"error": "Not a valid PKG file (Got: " + version + ")"}
	}

	var fileCount uint32
	binary.Read(reader, binary.LittleEndian, &fileCount)

	// -- 读取索引 --
	type Entry struct {
		Name   string
		Offset uint32
		Size   uint32
	}
	entries := make([]Entry, fileCount)
	for i := uint32(0); i < fileCount; i++ {
		binary.Read(reader, binary.LittleEndian, &strLen)
		nameBuf := make([]byte, strLen)
		reader.Read(nameBuf)
		entries[i].Name = string(nameBuf)
		binary.Read(reader, binary.LittleEndian, &entries[i].Offset)
		binary.Read(reader, binary.LittleEndian, &entries[i].Size)
	}
	dsPtr, _ := reader.Seek(0, 1)

	// -- 准备返回的 JS 对象 --
	jsResults := js.Global().Get("Array").New()
	for i, entry := range entries {
		fileData := make([]byte, entry.Size)
		reader.Seek(dsPtr+int64(entry.Offset), 0)
		reader.Read(fileData)

		jsUint8Array := js.Global().Get("Uint8Array").New(entry.Size)
		js.CopyBytesToJS(jsUint8Array, fileData)

		fileObj := js.Global().Get("Object").New()
		fileObj.Set("name", entry.Name)
		fileObj.Set("size", entry.Size)
		fileObj.Set("data", jsUint8Array)
		jsResults.SetIndex(i, fileObj)
	}
	return map[string]any{"success": true, "version": version, "files": jsResults}
}

// repackPkg: 完美复刻 CLI + 智能路径降维 + 防御性清洗
func repackPkg(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.Null()
	}

	jsFiles := args[0]
	fileCount := jsFiles.Get("length").Int()

	if fileCount == 0 {
		return js.Null()
	}

	type FileNode struct {
		Path string
		Data []byte
		Size uint32
	}
	var files []FileNode

	// 1. 提取 JS 数据并清洗系统垃圾文件
	for i := 0; i < fileCount; i++ {
		jsFile := jsFiles.Index(i)
		path := jsFile.Get("name").String()
		path = strings.ReplaceAll(path, "\\", "/")

		// 过滤掉 MacOS 和 Windows 的压缩包隐藏垃圾
		if strings.Contains(path, "__MACOSX") || strings.HasSuffix(path, ".DS_Store") || strings.HasSuffix(path, "Thumbs.db") {
			continue
		}

		jsData := jsFile.Get("data")
		size := jsData.Get("byteLength").Int()

		data := make([]byte, size)
		js.CopyBytesToGo(data, jsData)

		files = append(files, FileNode{
			Path: path,
			Data: data,
			Size: uint32(size),
		})
	}

	if len(files) == 0 {
		return js.Null()
	}

	// 2. 【终极防御】智能剥离多余的外层文件夹！
	// 应对用户直接右键打包文件夹的情况，自动削掉前缀
	for {
		parts := strings.Split(files[0].Path, "/")
		// 如果最前面的路径已经没有 '/' 了（比如 project.json），说明已经到达真正根目录
		if len(parts) <= 1 {
			break
		}

		rootCandidate := parts[0] + "/"
		isCommonRoot := true

		// 检查是否所有文件都包在这个多余的文件夹里
		for _, f := range files {
			if !strings.HasPrefix(f.Path, rootCandidate) {
				isCommonRoot = false
				break
			}
		}

		// 如果是，强制把所有文件的这个前缀给裁掉！
		if isCommonRoot {
			for i := range files {
				files[i].Path = strings.TrimPrefix(files[i].Path, rootCandidate)
			}
		} else {
			break // 目录结构正常，停止剥离
		}
	}

	// 3. 强行按字母排序 (消除 JSZip 异步并发造成的乱序，防止引擎二分查找崩溃)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	// 4. 构建二进制流，完美复刻 CLI 的双重循环无缝对齐法
	pkgBuffer := new(bytes.Buffer)

	version := "PKGV0022"
	binary.Write(pkgBuffer, binary.LittleEndian, uint32(len(version)))
	pkgBuffer.WriteString(version)
	binary.Write(pkgBuffer, binary.LittleEndian, uint32(len(files)))

	var currentOffset uint32 = 0
	// 索引区
	for _, f := range files {
		binary.Write(pkgBuffer, binary.LittleEndian, uint32(len(f.Path)))
		pkgBuffer.WriteString(f.Path)
		binary.Write(pkgBuffer, binary.LittleEndian, currentOffset)
		binary.Write(pkgBuffer, binary.LittleEndian, f.Size)
		currentOffset += f.Size
	}

	// 数据区
	for _, f := range files {
		pkgBuffer.Write(f.Data)
	}

	// 5. 传回 JS
	finalBytes := pkgBuffer.Bytes()
	jsUint8Array := js.Global().Get("Uint8Array").New(len(finalBytes))
	js.CopyBytesToJS(jsUint8Array, finalBytes)

	return jsUint8Array
}

func main() {
	c := make(chan struct{})
	js.Global().Set("extractPkg", js.FuncOf(extractPkg))
	js.Global().Set("repackPkg", js.FuncOf(repackPkg))
	println("Wasm Engine (WE Pkg Tool) Ready!")
	<-c
}
