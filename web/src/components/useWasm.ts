// src/composables/useWasm.ts (Repack 部分的完整循环)
import { ref, onMounted } from 'vue';
import JSZip from 'jszip';

// Wasm 接口声明
declare global {
  interface Window {
    extractPkg: (data: Uint8Array) => any;
    repackPkg: (files: { name: string, data: Uint8Array }[]) => Uint8Array;
  }
}

export function useWasm() {
  const wasmReady = ref(false);

  // 初始化（Vite 开发服务器会自动处理 public 的wasm）
  onMounted(async () => {
    const go = new (window as any).Go();
    try {
      // WASM 文件与 index.html 平级，基于当前页面 URL 计算路径
      const finalUrl = new URL('main.wasm', window.location.href).href;
      const result = await WebAssembly.instantiateStreaming(fetch(finalUrl), go.importObject);
      go.run(result.instance);
      wasmReady.value = true;
    } catch (e) {
      console.error("Wasm Loader error:", e);
    }
  });

  // Repack: 接受 ZIP 文件，在内存解压得到 Uint8Arrays，传给 Wasm，得到最终 PKG
  const doRepack = async (zipFile: File): Promise<Uint8Array | null> => {
    try {
      const zip = new JSZip();
      await zip.loadAsync(zipFile);

      const filesForWasm: { name: string, data: Uint8Array }[] = [];

      // 遍历 ZIP 结构，必须保留原相对路径
      const filePromises: Promise<any>[] = [];

      zip.forEach((path, entry) => {
        if (!entry.dir) {
          const promise = entry.async("uint8array").then(binaryData => {
            filesForWasm.push({ name: path, data: binaryData });
          });
          filePromises.push(promise);
        }
      });

      // 必须确保 ZIP 里的 96 个文件全部解压完成并塞进 filesForWasm
      await Promise.all(filePromises);

      console.log(`ZIP extracted in memory. Passing ${filesForWasm.length} files to Wasm Repacker.`);

      // 调用 Wasm 的打包函数
      // 这个时候 Go 内部通过两遍循环，修复偏移量计算 Bug
      const finalPkgBytes = window.repackPkg(filesForWasm);
      
      console.log(`Wasm Repack complete. PKG size: ${finalPkgBytes.length} bytes`);
      return finalPkgBytes;

    } catch (error) {
      console.error("Repack process failed:", error);
      return null;
    }
  };

  return { wasmReady, doRepack };
}