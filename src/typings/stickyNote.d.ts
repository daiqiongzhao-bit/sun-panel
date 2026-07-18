declare namespace StickyNote {
  interface Info extends Common.InfoBase {
    userId: number
    content: string
    color: string
    posX: number
    posY: number
    width: number
    height: number
    zIndex: number
    status: number
  }
}