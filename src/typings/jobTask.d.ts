declare namespace JobTask {
  interface Info extends Common.InfoBase {
    name: string
    jobType: number      // 1提醒 2自定义
    cronExpr: string
    content: string
    status: number        // 1运行中 2暂停
    lastRunAt: string | null
    nextRunAt: string | null
    creatorId: number
  }
}

declare namespace JobLog {
  interface Info extends Common.InfoBase {
    jobId: number
    jobName: string
    startTime: string
    duration: number
    status: number      // 1成功 2失败
    errorMsg: string
  }
}
