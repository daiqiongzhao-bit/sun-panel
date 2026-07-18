declare namespace LoginLog {
  interface Info extends Common.InfoBase {
    userId: number
    username: string
    ip: string
    location: string
    userAgent: string
    status: number  // 1成功 2失败
    remark: string
  }
}