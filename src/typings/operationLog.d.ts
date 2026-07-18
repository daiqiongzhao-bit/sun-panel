declare namespace OperationLog {
  interface Info extends Common.InfoBase {
    userId: number
    username: string
    module: string
    action: string
    method: string
    path: string
    ip: string
    userAgent: string
    requestBody: string
    responseCode: number
    duration: number
    remark: string
  }
}