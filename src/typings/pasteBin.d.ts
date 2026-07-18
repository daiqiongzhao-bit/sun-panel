declare namespace PasteBin {
  interface Info extends Common.InfoBase {
    userId: number
    type: number      // 1文本 2文件
    title: string
    content: string
    fileName: string
    fileSize: number
    code: string
    expireAt: string
    accessCnt: number
    status: number
  }
}