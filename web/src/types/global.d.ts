export {};

declare global {
  interface Window {
    Go: any;
    extractPkg: (data: Uint8Array) => any;
    repackPkg: (files: any[]) => Uint8Array;
  }
}