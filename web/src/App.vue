<template>
  <div class="terminal-container">
    <!-- 背景效果层 -->
    <div class="bg-grid"></div>
    <div class="bg-noise"></div>
    <div class="bg-glow"></div>
    
    <!-- 扫描线效果 -->
    <div class="scanlines"></div>
    
    <!-- 顶部状态栏 -->
    <div class="top-bar">
      <div class="bar-section">
        <span class="bar-label">SYS.STATUS</span>
        <span class="bar-value" :class="wasmReady ? 'status-online' : 'status-offline'">
          {{ wasmReady ? 'ONLINE' : 'OFFLINE' }}
        </span>
      </div>
      <div class="bar-divider"></div>
      <div class="bar-section">
        <span class="bar-label">WASM.CORE</span>
        <span class="bar-value">{{ wasmStatus }}</span>
      </div>
      <div class="bar-divider"></div>
      <div class="bar-section">
        <span class="bar-label">SECURE.MODE</span>
        <span class="bar-value status-online">ENABLED</span>
      </div>
    </div>

    <!-- 主标题区 -->
    <header class="terminal-header">
      <div class="header-border top"></div>
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="main-title">
            <span class="title-bracket">[</span>
            <span class="title-text">WE_PKG_CYBER</span>
            <span class="title-bracket">]</span>
          </h1>
          <p class="sub-title">ASSET.EXTRACTION.MODULE // v2.0.0</p>
        </div>
        <div class="header-icon right">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
        </div>
      </div>
      <div class="header-border bottom"></div>
    </header>

    <!-- 主操作区 -->
    <main class="terminal-main">
      <!-- 左侧：解包模块 -->
      <section class="module module-extract">
        <div class="module-header">
          <div class="header-index">01</div>
          <h2 class="header-title">EXTRACT</h2>
          <div class="header-line"></div>
          <span class="header-badge">DECODE</span>
        </div>
        
        <div 
          class="drop-target" 
          :class="{ 'active': isDraggingUnpack, 'has-files': extractedFiles.length > 0 }"
          @click="triggerUnpack"
          @dragover.prevent="isDraggingUnpack = true"
          @dragleave="isDraggingUnpack = false"
          @drop.prevent="handleDropUnpack"
        >
          <div class="target-frame">
            <div class="frame-corner tl"></div>
            <div class="frame-corner tr"></div>
            <div class="frame-corner bl"></div>
            <div class="frame-corner br"></div>
          </div>
          
          <div class="target-content">
            <div class="target-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3"/>
              </svg>
            </div>
            <div class="target-status">{{ unpackStatus }}</div>
            <div class="target-hint">
              <span class="hint-primary">DROP .PKG FILE</span>
              <span class="hint-secondary">or click to browse</span>
            </div>
          </div>
          
          <input type="file" ref="unpackInput" @change="handleUnpack" accept=".pkg" />
        </div>

        <!-- 资产清单 -->
        <div class="asset-manifest" v-if="extractedFiles.length > 0">
          <div class="manifest-header">
            <span class="manifest-label">MANIFEST</span>
            <span class="manifest-count">{{ extractedFiles.length }} ITEMS</span>
          </div>
          <div class="manifest-grid">
            <div class="manifest-item">
              <div class="item-icon audio">♫</div>
              <div class="item-data">
                <span class="item-value">{{ assetStats.audio }}</span>
                <span class="item-label">AUDIO</span>
              </div>
            </div>
            <div class="manifest-item">
              <div class="item-icon visual">◈</div>
              <div class="item-data">
                <span class="item-value">{{ assetStats.visual }}</span>
                <span class="item-label">VISUAL</span>
              </div>
            </div>
            <div class="manifest-item">
              <div class="item-icon config">⚙</div>
              <div class="item-data">
                <span class="item-value">{{ assetStats.config }}</span>
                <span class="item-label">CONFIG</span>
              </div>
            </div>
            <div class="manifest-item">
              <div class="item-icon other">◉</div>
              <div class="item-data">
                <span class="item-value">{{ assetStats.other }}</span>
                <span class="item-label">OTHER</span>
              </div>
            </div>
          </div>
        </div>

        <button 
          class="action-btn btn-extract" 
          :disabled="extractedFiles.length === 0 || isProcessing"
          @click="downloadZip"
        >
          <span class="btn-glow"></span>
          <span class="btn-content">
            <span class="btn-icon">▼</span>
            <span class="btn-text">EXPORT ZIP</span>
          </span>
        </button>
      </section>

      <!-- 中央连接区 -->
      <div class="connector">
        <div class="connector-line"></div>
        <div class="connector-node">
          <div class="node-ring"></div>
          <div class="node-core"></div>
        </div>
        <div class="connector-line"></div>
      </div>

      <!-- 右侧：打包模块 -->
      <section class="module module-repack">
        <div class="module-header">
          <div class="header-index">02</div>
          <h2 class="header-title">REPACK</h2>
          <div class="header-line"></div>
          <span class="header-badge">ENCODE</span>
        </div>
        
        <div 
          class="drop-target" 
          :class="{ 'active': isDraggingRepack, 'ready': repackDataReady }"
          @click="triggerRepack"
          @dragover.prevent="isDraggingRepack = true"
          @dragleave="isDraggingRepack = false"
          @drop.prevent="handleDropRepack"
        >
          <div class="target-frame">
            <div class="frame-corner tl"></div>
            <div class="frame-corner tr"></div>
            <div class="frame-corner bl"></div>
            <div class="frame-corner br"></div>
          </div>
          
          <div class="target-content">
            <div class="target-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M17 8l-5-5-5 5M12 3v12"/>
              </svg>
            </div>
            <div class="target-status">{{ repackStatus }}</div>
            <div class="target-hint">
              <span class="hint-primary">DROP MOD.ZIP</span>
              <span class="hint-secondary">modified assets archive</span>
            </div>
          </div>
          
          <input type="file" ref="repackInput" @change="handleRepack" accept=".zip" />
        </div>

        <!-- 注入状态 -->
        <div class="injection-status" v-if="repackDataReady">
          <div class="status-header">
            <span class="status-indicator"></span>
            <span class="status-label">INJECTION READY</span>
          </div>
          <div class="status-details">
            <div class="detail-line">
              <span class="detail-key">FILES:</span>
              <span class="detail-value">{{ filesForWasm.length }}</span>
            </div>
            <div class="detail-line">
              <span class="detail-key">ALIGNMENT:</span>
              <span class="detail-value">BYTE-PERFECT</span>
            </div>
          </div>
        </div>

        <button 
          class="action-btn btn-repack" 
          :disabled="!repackDataReady || isProcessing"
          @click="buildPkg"
        >
          <span class="btn-glow"></span>
          <span class="btn-content">
            <span class="btn-icon">▲</span>
            <span class="btn-text">BUILD PKG</span>
          </span>
        </button>
      </section>
    </main>

    <!-- 底部信息栏 -->
    <footer class="terminal-footer">
      <div class="footer-section">
        <span class="footer-label">ENCRYPTION</span>
        <span class="footer-value">NONE // LOCAL ONLY</span>
      </div>
      <div class="footer-divider">◆</div>
      <div class="footer-section">
        <span class="footer-label">DATA.PRIVACY</span>
        <span class="footer-value">100% CLIENT-SIDE</span>
      </div>
      <div class="footer-divider">◆</div>
      <div class="footer-section author">
        <span class="footer-label">AUTHOR</span>
        <a href="https://github.com/fuzqing/we-pkg-cyber" target="_blank" rel="noopener" class="footer-link">
          <svg class="github-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
          <span class="footer-value">fuzqing</span>
        </a>
      </div>
    </footer>

    <!-- 处理遮罩 -->
    <div class="processing-overlay" v-if="isProcessing">
      <div class="overlay-content">
        <div class="spinner">
          <div class="spinner-ring"></div>
          <div class="spinner-ring"></div>
          <div class="spinner-ring"></div>
        </div>
        <div class="overlay-text">PROCESSING DATA STREAM...</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import JSZip from 'jszip';

interface WasmFile {
  name: string;
  data: Uint8Array;
}

declare global {
  interface Window {
    extractPkg: (data: Uint8Array) => any;
    repackPkg: (files: WasmFile[]) => Uint8Array;
    Go: any;
  }
}

const wasmReady = ref(false);
const wasmStatus = ref('INITIALIZING WASM NEURAL NET...');
const isProcessing = ref(false);

const unpackInput = ref<HTMLInputElement | null>(null);
const unpackStatus = ref('STANDBY');
const extractedFiles = ref<any[]>([]);

const repackInput = ref<HTMLInputElement | null>(null);
const repackStatus = ref('STANDBY');
const repackDataReady = ref(false);
let filesForWasm: WasmFile[] = [];

// ================= 资产分类计算属性 =================
const assetStats = computed(() => {
  const stats = { audio: 0, visual: 0, config: 0, other: 0 };
  extractedFiles.value.forEach(file => {
    const ext = file.name.split('.').pop()?.toLowerCase();
    if (['mp3', 'ogg', 'wav'].includes(ext)) stats.audio++;
    else if (['png', 'jpg', 'jpeg', 'mp4', 'webm', 'gif', 'tex'].includes(ext)) stats.visual++;
    else if (['json', 'txt', 'xml'].includes(ext)) stats.config++;
    else stats.other++;
  });
  return stats;
});

onMounted(async () => {
  try {
    const go = new window.Go();
    // WASM 文件与 index.html 平级，基于当前页面 URL 计算路径
    const wasmUrl = new URL('main.wasm', window.location.href).href;
    const result = await WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject);
    go.run(result.instance);
    wasmReady.value = true;
    wasmStatus.value = 'ONLINE // FULL ACCESS GRANTED';
  } catch (e) {
    wasmStatus.value = 'FATAL ERROR // WASM INJECTION FAILED';
    console.error(e);
  }
});

const triggerUnpack = () => { if (wasmReady.value && !isProcessing.value) unpackInput.value?.click(); };

const handleUnpack = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;
  processPkgFile(file);
  target.value = ''; 
};

const downloadZip = async () => {
  isProcessing.value = true;
  unpackStatus.value = 'COMPRESSING PAYLOAD...';
  
  try {
    const zip = new JSZip();
    extractedFiles.value.forEach(f => zip.file(f.name, f.data));
    const blob = await zip.generateAsync({ type: 'blob' });
    
    triggerDownload(blob, 'scene_unpacked.zip');
    unpackStatus.value = 'PAYLOAD DELIVERED_';
  } catch (err) {
    console.error(err);
    unpackStatus.value = 'ERROR // ZIP FAILED';
  } finally {
    isProcessing.value = false;
  }
};

const isDraggingUnpack = ref(false);
const isDraggingRepack = ref(false);

const handleDropUnpack = (event: DragEvent) => {
  isDraggingUnpack.value = false;
  const files = event.dataTransfer?.files;
  if (files && files.length > 0) {
    processPkgFile(files[0]);
  }
};

const handleDropRepack = (event: DragEvent) => {
  isDraggingRepack.value = false;
  const files = event.dataTransfer?.files;
  if (files && files.length > 0) {
    processZipFile(files[0]);
  }
};

const processPkgFile = (file: File) => {
  unpackStatus.value = 'BYPASSING SECURITY...';
  const reader = new FileReader();
  
  reader.onload = (e) => {
    try {
      const uint8Array = new Uint8Array(e.target?.result as ArrayBuffer);
      const res = window.extractPkg(uint8Array);
      if (res.error) throw new Error(res.error);
      
      extractedFiles.value = res.files;
      unpackStatus.value = 'EXTRACTION COMPLETE';
    } catch (err: any) {
      alert("解析失败: " + err.message);
      unpackStatus.value = 'EXTRACTION FAILED';
    }
  };
  reader.readAsArrayBuffer(file);
};

const processZipFile = async (file: File) => {
  repackStatus.value = 'ANALYZING ARCHIVE...';
  isProcessing.value = true;
  filesForWasm = [];
  
  try {
    const zip = new JSZip();
    const loadedZip = await zip.loadAsync(file);
    const promises: Promise<void>[] = [];

    loadedZip.forEach((relativePath, zipEntry) => {
      if (!zipEntry.dir) {
        promises.push(
          zipEntry.async('uint8array').then((data) => {
            filesForWasm.push({ name: relativePath, data: data });
          })
        );
      }
    });

    await Promise.all(promises);
    repackDataReady.value = true;
    repackStatus.value = 'ARCHIVE LOADED';
  } catch (err) {
    console.error(err);
    repackStatus.value = 'LOAD FAILED';
  } finally {
    isProcessing.value = false;
  }
};

const triggerRepack = () => { if (wasmReady.value && !isProcessing.value) repackInput.value?.click(); };

const handleRepack = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;
  await processZipFile(file);
  target.value = '';
};

const buildPkg = () => {
  isProcessing.value = true;
  repackStatus.value = 'COMPILING CORE BINARY...';
  
  try {
    const pkgBytes = window.repackPkg(filesForWasm);
    const blob = new Blob([pkgBytes.buffer.slice(pkgBytes.byteOffset, pkgBytes.byteOffset + pkgBytes.byteLength) as ArrayBuffer], { type: 'application/octet-stream' });
    
    triggerDownload(blob, 'scene.pkg');
    repackStatus.value = 'MOD INJECTION SUCCESSFUL_';
    repackDataReady.value = false; 
  } catch (err) {
    console.error(err);
    repackStatus.value = 'ERROR // MATRIX COLLAPSE';
  } finally {
    isProcessing.value = false;
  }
};

const triggerDownload = (blob: Blob, filename: string) => {
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
};
</script>

<style>
/* ============================================
   WE_PKG_CYBER - Industrial Terminal Theme
   ============================================ */

@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600;700&family=Chakra+Petch:wght@400;500;600;700&display=swap');

:root {
  /* Core Colors */
  --bg-primary: #0a0a0f;
  --bg-secondary: #111118;
  --bg-tertiary: #1a1a24;
  
  /* Accent Colors */
  --accent-cyan: #00d4ff;
  --accent-cyan-dim: rgba(0, 212, 255, 0.15);
  --accent-cyan-glow: rgba(0, 212, 255, 0.4);
  --accent-amber: #ff9500;
  --accent-amber-dim: rgba(255, 149, 0, 0.15);
  --accent-amber-glow: rgba(255, 149, 0, 0.4);
  --accent-red: #ff3366;
  --accent-green: #00ff88;
  
  /* Text Colors */
  --text-primary: #e8e8f0;
  --text-secondary: #8b8b9e;
  --text-muted: #5a5a6e;
  
  /* Border Colors */
  --border-dim: rgba(255, 255, 255, 0.08);
  --border-active: rgba(0, 212, 255, 0.3);
  
  /* Fonts */
  --font-mono: 'JetBrains Mono', monospace;
  --font-display: 'Chakra Petch', sans-serif;
  
  /* Spacing */
  --space-xs: 0.5rem;
  --space-sm: 0.75rem;
  --space-md: 1rem;
  --space-lg: 1.5rem;
  --space-xl: 2rem;
  --space-2xl: 3rem;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--font-mono);
  background: var(--bg-primary);
  color: var(--text-primary);
  min-height: 100vh;
  overflow-x: hidden;
}

/* ============================================
   Background Effects
   ============================================ */

.terminal-container {
  min-height: 100vh;
  position: relative;
  display: flex;
  flex-direction: column;
  padding: var(--space-lg);
}

.bg-grid {
  position: fixed;
  inset: 0;
  background-image: 
    linear-gradient(rgba(0, 212, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 212, 255, 0.03) 1px, transparent 1px);
  background-size: 40px 40px;
  pointer-events: none;
  z-index: 0;
}

.bg-noise {
  position: fixed;
  inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
  opacity: 0.03;
  pointer-events: none;
  z-index: 1;
}

.bg-glow {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80vw;
  height: 80vh;
  background: radial-gradient(ellipse, rgba(0, 212, 255, 0.05) 0%, transparent 70%);
  pointer-events: none;
  z-index: 0;
}

.scanlines {
  position: fixed;
  inset: 0;
  background: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0, 0, 0, 0.15) 2px,
    rgba(0, 0, 0, 0.15) 4px
  );
  pointer-events: none;
  z-index: 1000;
  opacity: 0.4;
}

/* ============================================
   Top Bar
   ============================================ */

.top-bar {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: var(--space-md);
  padding: var(--space-sm) var(--space-md);
  background: var(--bg-secondary);
  border: 1px solid var(--border-dim);
  margin-bottom: var(--space-xl);
}

.bar-section {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
}

.bar-label {
  font-size: 0.65rem;
  font-weight: 600;
  color: var(--text-muted);
  letter-spacing: 0.1em;
}

.bar-value {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-secondary);
  font-family: var(--font-display);
}

.bar-divider {
  width: 1px;
  height: 16px;
  background: var(--border-dim);
}

.status-online {
  color: var(--accent-green);
  text-shadow: 0 0 10px rgba(0, 255, 136, 0.4);
}

.status-offline {
  color: var(--accent-red);
  text-shadow: 0 0 10px rgba(255, 51, 102, 0.4);
}

/* ============================================
   Header
   ============================================ */

.terminal-header {
  position: relative;
  z-index: 10;
  margin-bottom: var(--space-2xl);
}

.header-border {
  height: 2px;
  background: linear-gradient(90deg, 
    transparent 0%, 
    var(--accent-cyan) 20%, 
    var(--accent-cyan) 80%, 
    transparent 100%
  );
  opacity: 0.5;
}

.header-border.top {
  margin-bottom: var(--space-md);
}

.header-border.bottom {
  margin-top: var(--space-md);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-lg);
  padding: var(--space-md) 0;
}

.header-icon {
  width: 48px;
  height: 48px;
  color: var(--accent-cyan);
  opacity: 0.6;
}

.header-icon svg {
  width: 100%;
  height: 100%;
}

.main-title {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  font-family: var(--font-display);
  font-size: 2.5rem;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.title-bracket {
  color: var(--accent-cyan);
  opacity: 0.6;
  font-weight: 400;
}

.title-text {
  background: linear-gradient(180deg, var(--text-primary) 0%, var(--accent-cyan) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.sub-title {
  text-align: center;
  font-size: 0.75rem;
  color: var(--text-muted);
  letter-spacing: 0.3em;
  margin-top: var(--space-xs);
}

/* ============================================
   Main Layout
   ============================================ */

.terminal-main {
  position: relative;
  z-index: 10;
  flex: 1;
  display: flex;
  align-items: stretch;
  justify-content: center;
  gap: var(--space-xl);
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

/* ============================================
   Module Base
   ============================================ */

.module {
  flex: 1;
  max-width: 480px;
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
}

.module-header {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding-bottom: var(--space-sm);
  border-bottom: 1px solid var(--border-dim);
}

.header-index {
  font-family: var(--font-display);
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-muted);
  line-height: 1;
}

.module-extract .header-index {
  color: var(--accent-cyan);
}

.module-repack .header-index {
  color: var(--accent-amber);
}

.header-title {
  font-family: var(--font-display);
  font-size: 1.25rem;
  font-weight: 600;
  letter-spacing: 0.1em;
  color: var(--text-primary);
}

.header-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, var(--border-dim) 0%, transparent 100%);
}

.header-badge {
  font-size: 0.6rem;
  font-weight: 600;
  letter-spacing: 0.15em;
  padding: 4px 10px;
  border: 1px solid var(--border-dim);
  color: var(--text-secondary);
}

.module-extract .header-badge {
  border-color: var(--accent-cyan);
  color: var(--accent-cyan);
  background: var(--accent-cyan-dim);
}

.module-repack .header-badge {
  border-color: var(--accent-amber);
  color: var(--accent-amber);
  background: var(--accent-amber-dim);
}

/* ============================================
   Drop Target
   ============================================ */

.drop-target {
  position: relative;
  aspect-ratio: 4/3;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.drop-target input {
  display: none;
}

.target-frame {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.frame-corner {
  position: absolute;
  width: 20px;
  height: 20px;
  border: 2px solid var(--border-dim);
  transition: all 0.3s ease;
}

.frame-corner.tl { top: 0; left: 0; border-right: 0; border-bottom: 0; }
.frame-corner.tr { top: 0; right: 0; border-left: 0; border-bottom: 0; }
.frame-corner.bl { bottom: 0; left: 0; border-right: 0; border-top: 0; }
.frame-corner.br { bottom: 0; right: 0; border-left: 0; border-top: 0; }

.target-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-md);
  text-align: center;
  padding: var(--space-lg);
}

.target-icon {
  width: 64px;
  height: 64px;
  color: var(--text-muted);
  transition: all 0.3s ease;
}

.target-icon svg {
  width: 100%;
  height: 100%;
}

.target-status {
  font-family: var(--font-display);
  font-size: 0.875rem;
  font-weight: 600;
  letter-spacing: 0.1em;
  color: var(--text-secondary);
}

.target-hint {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.hint-primary {
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.15em;
  color: var(--text-muted);
}

.hint-secondary {
  font-size: 0.65rem;
  color: var(--text-muted);
  opacity: 0.7;
}

/* Drop Target States */
.module-extract .drop-target {
  background: linear-gradient(135deg, rgba(0, 212, 255, 0.03) 0%, transparent 50%);
  border: 1px solid var(--border-dim);
}

.module-extract .drop-target:hover,
.module-extract .drop-target.active {
  border-color: var(--accent-cyan);
  background: var(--accent-cyan-dim);
}

.module-extract .drop-target:hover .frame-corner,
.module-extract .drop-target.active .frame-corner {
  border-color: var(--accent-cyan);
  width: 30px;
  height: 30px;
}

.module-extract .drop-target:hover .target-icon,
.module-extract .drop-target.active .target-icon {
  color: var(--accent-cyan);
}

.module-extract .drop-target.has-files {
  border-color: var(--accent-green);
  background: rgba(0, 255, 136, 0.05);
}

.module-extract .drop-target.has-files .frame-corner {
  border-color: var(--accent-green);
}

.module-extract .drop-target.has-files .target-icon {
  color: var(--accent-green);
}

.module-repack .drop-target {
  background: linear-gradient(225deg, rgba(255, 149, 0, 0.03) 0%, transparent 50%);
  border: 1px solid var(--border-dim);
}

.module-repack .drop-target:hover,
.module-repack .drop-target.active {
  border-color: var(--accent-amber);
  background: var(--accent-amber-dim);
}

.module-repack .drop-target:hover .frame-corner,
.module-repack .drop-target.active .frame-corner {
  border-color: var(--accent-amber);
  width: 30px;
  height: 30px;
}

.module-repack .drop-target:hover .target-icon,
.module-repack .drop-target.active .target-icon {
  color: var(--accent-amber);
}

.module-repack .drop-target.ready {
  border-color: var(--accent-green);
  background: rgba(0, 255, 136, 0.05);
}

.module-repack .drop-target.ready .frame-corner {
  border-color: var(--accent-green);
}

.module-repack .drop-target.ready .target-icon {
  color: var(--accent-green);
}

/* ============================================
   Asset Manifest
   ============================================ */

.asset-manifest {
  background: var(--bg-secondary);
  border: 1px solid var(--border-dim);
  padding: var(--space-md);
}

.manifest-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-md);
  padding-bottom: var(--space-sm);
  border-bottom: 1px solid var(--border-dim);
}

.manifest-label {
  font-size: 0.65rem;
  font-weight: 600;
  letter-spacing: 0.15em;
  color: var(--text-muted);
}

.manifest-count {
  font-size: 0.7rem;
  font-weight: 600;
  color: var(--accent-cyan);
  font-family: var(--font-display);
}

.manifest-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-sm);
}

.manifest-item {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-sm);
  background: var(--bg-tertiary);
  border: 1px solid var(--border-dim);
}

.item-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-dim);
}

.item-icon.audio { color: #00d4ff; }
.item-icon.visual { color: #ff6b9d; }
.item-icon.config { color: #a78bfa; }
.item-icon.other { color: #fbbf24; }

.item-data {
  display: flex;
  flex-direction: column;
}

.item-value {
  font-family: var(--font-display);
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.item-label {
  font-size: 0.6rem;
  letter-spacing: 0.1em;
  color: var(--text-muted);
  margin-top: 2px;
}

/* ============================================
   Injection Status
   ============================================ */

.injection-status {
  background: var(--bg-secondary);
  border: 1px solid var(--border-dim);
  padding: var(--space-md);
}

.status-header {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  margin-bottom: var(--space-md);
  padding-bottom: var(--space-sm);
  border-bottom: 1px solid var(--border-dim);
}

.status-indicator {
  width: 8px;
  height: 8px;
  background: var(--accent-green);
  border-radius: 50%;
  box-shadow: 0 0 10px var(--accent-green);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-label {
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.1em;
  color: var(--accent-green);
  font-family: var(--font-display);
}

.status-details {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.detail-line {
  display: flex;
  justify-content: space-between;
  font-size: 0.7rem;
}

.detail-key {
  color: var(--text-muted);
  letter-spacing: 0.05em;
}

.detail-value {
  color: var(--text-secondary);
  font-weight: 500;
  font-family: var(--font-display);
}

/* ============================================
   Action Buttons
   ============================================ */

.action-btn {
  position: relative;
  padding: var(--space-md) var(--space-lg);
  background: transparent;
  border: none;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s ease;
}

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-glow {
  position: absolute;
  inset: 0;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.btn-content {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-sm);
  font-family: var(--font-display);
  font-size: 0.875rem;
  font-weight: 600;
  letter-spacing: 0.1em;
}

.btn-icon {
  font-size: 1rem;
}

/* Extract Button */
.btn-extract {
  border: 1px solid var(--accent-cyan);
  color: var(--accent-cyan);
  background: var(--accent-cyan-dim);
}

.btn-extract .btn-glow {
  background: radial-gradient(ellipse at center, var(--accent-cyan-glow) 0%, transparent 70%);
}

.btn-extract:not(:disabled):hover {
  background: var(--accent-cyan);
  color: var(--bg-primary);
}

.btn-extract:not(:disabled):hover .btn-glow {
  opacity: 0.3;
}

/* Repack Button */
.btn-repack {
  border: 1px solid var(--accent-amber);
  color: var(--accent-amber);
  background: var(--accent-amber-dim);
}

.btn-repack .btn-glow {
  background: radial-gradient(ellipse at center, var(--accent-amber-glow) 0%, transparent 70%);
}

.btn-repack:not(:disabled):hover {
  background: var(--accent-amber);
  color: var(--bg-primary);
}

.btn-repack:not(:disabled):hover .btn-glow {
  opacity: 0.3;
}

/* ============================================
   Connector
   ============================================ */

.connector {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-sm);
  padding: 0 var(--space-md);
}

.connector-line {
  width: 2px;
  flex: 1;
  max-height: 100px;
  background: linear-gradient(180deg, 
    transparent 0%, 
    var(--border-dim) 20%, 
    var(--border-dim) 80%, 
    transparent 100%
  );
}

.connector-node {
  position: relative;
  width: 24px;
  height: 24px;
}

.node-ring {
  position: absolute;
  inset: 0;
  border: 2px solid var(--accent-cyan);
  border-radius: 50%;
  opacity: 0.3;
  animation: rotate 10s linear infinite;
}

.node-core {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 8px;
  height: 8px;
  background: var(--accent-cyan);
  border-radius: 50%;
  box-shadow: 0 0 15px var(--accent-cyan);
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* ============================================
   Footer
   ============================================ */

.terminal-footer {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-lg);
  padding: var(--space-md);
  margin-top: var(--space-2xl);
  border-top: 1px solid var(--border-dim);
}

.footer-section {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
}

.footer-label {
  font-size: 0.6rem;
  font-weight: 600;
  letter-spacing: 0.1em;
  color: var(--text-muted);
}

.footer-value {
  font-size: 0.7rem;
  font-weight: 500;
  color: var(--text-secondary);
  font-family: var(--font-display);
}

.footer-divider {
  color: var(--border-dim);
  font-size: 0.5rem;
}

.footer-section.author {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
}

.footer-link {
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
  transition: all 0.2s ease;
}

.footer-link:hover {
  opacity: 0.8;
}

.footer-link:hover .footer-value {
  color: var(--accent-cyan);
}

.github-icon {
  width: 14px;
  height: 14px;
  color: var(--text-secondary);
  transition: color 0.2s ease;
}

.footer-link:hover .github-icon {
  color: var(--accent-cyan);
}

/* ============================================
   Processing Overlay
   ============================================ */

.processing-overlay {
  position: fixed;
  inset: 0;
  background: rgba(10, 10, 15, 0.9);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.overlay-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-lg);
}

.spinner {
  position: relative;
  width: 80px;
  height: 80px;
}

.spinner-ring {
  position: absolute;
  inset: 0;
  border: 2px solid transparent;
  border-top-color: var(--accent-cyan);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.spinner-ring:nth-child(1) {
  animation-duration: 1s;
}

.spinner-ring:nth-child(2) {
  inset: 10px;
  border-top-color: var(--accent-amber);
  animation-duration: 1.5s;
  animation-direction: reverse;
}

.spinner-ring:nth-child(3) {
  inset: 20px;
  border-top-color: var(--accent-green);
  animation-duration: 2s;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.overlay-text {
  font-family: var(--font-display);
  font-size: 0.875rem;
  font-weight: 600;
  letter-spacing: 0.2em;
  color: var(--accent-cyan);
  animation: blink 1s ease-in-out infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* ============================================
   Responsive
   ============================================ */

@media (max-width: 1024px) {
  .terminal-main {
    flex-direction: column;
    align-items: center;
  }
  
  .module {
    max-width: 100%;
    width: 100%;
  }
  
  .connector {
    flex-direction: row;
    padding: var(--space-md) 0;
  }
  
  .connector-line {
    width: 100px;
    height: 2px;
    max-height: none;
  }
  
  .main-title {
    font-size: 1.75rem;
  }
  
  .header-icon {
    width: 32px;
    height: 32px;
  }
}

@media (max-width: 640px) {
  .terminal-container {
    padding: var(--space-md);
  }
  
  .top-bar {
    flex-wrap: wrap;
    gap: var(--space-sm);
  }
  
  .main-title {
    font-size: 1.25rem;
  }
  
  .header-content {
    gap: var(--space-sm);
  }
  
  .terminal-footer {
    flex-direction: column;
    gap: var(--space-sm);
  }
  
  .footer-divider {
    display: none;
  }
}
</style>