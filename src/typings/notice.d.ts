declare namespace Notice{

    interface NoticeInfo extends Common.InfoBase{
        title:string
        content:string
        displayType:number
        oneRead:number
        url:string
        isLogin:number
        noticeType: number  // 1公告 2站内信
        status: number      // 1启用 2停用
        userId: number
        targetUserIds: string
    }


}