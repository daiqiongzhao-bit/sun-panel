<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { NButton, NInput, NPopconfirm, NSwitch, useMessage } from 'naive-ui'
import { SvgIcon } from '@/components/common'
import { t } from '@/locales'
import { getStickyNoteList, createStickyNote, updateStickyNote, deleteStickyNote } from '@/api/panel/stickyNote'

const ms = useMessage()

const PRESET_COLORS = [
  '#fff3bf',
  '#ffe0e6',
  '#d0f0fd',
  '#d3f9d8',
  '#e8daef',
  '#fce4c4',
  '#f5f5f5',
  '#ffffff',
  'transparent',
]

const DEFAULT_WIDTH = 200
const DEFAULT_HEIGHT = 150
const MIN_WIDTH = 150
const MIN_HEIGHT = 100

interface NoteData {
  id: number
  content: string
  color: string
  posX: number
  posY: number
  width: number
  height: number
  zIndex: number
  status: number
}

const notes = ref<NoteData[]>([])
let maxZIndex = 1

// 小绿点拖拽
const btnPos = ref({ x: 24, y: 120 })
const isDraggingBtn = ref(false)
const btnDragStart = ref({ x: 0, y: 0, elX: 24, elY: 120 })

function initBtnPos() {
  try {
    const saved = localStorage.getItem('sun-panel-sticky-btn-pos')
    if (saved) {
      const pos = JSON.parse(saved)
      if (pos.x !== undefined && pos.y !== undefined) {
        btnPos.value = { x: pos.x, y: pos.y }
        btnDragStart.value = { x: 0, y: 0, elX: pos.x, elY: pos.y }
      }
    }
  } catch { /* ignore */ }
}
function saveBtnPos() {
  localStorage.setItem('sun-panel-sticky-btn-pos', JSON.stringify(btnPos.value))
}
initBtnPos()

function onBtnDragStart(e: MouseEvent | TouchEvent) {
  isDraggingBtn.value = true
  const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX
  const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY
  btnDragStart.value = { x: clientX, y: clientY, elX: btnPos.value.x, elY: btnPos.value.y }
  document.addEventListener('mousemove', onBtnDragMove)
  document.addEventListener('mouseup', onBtnDragEnd)
  document.addEventListener('touchmove', onBtnDragMove, { passive: false })
  document.addEventListener('touchend', onBtnDragEnd)
}
function onBtnDragMove(e: MouseEvent | TouchEvent) {
  if (!isDraggingBtn.value) return
  e.preventDefault()
  const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX
  const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY
  const dx = btnDragStart.value.x - clientX
  const dy = btnDragStart.value.y - clientY
  btnPos.value = {
    x: Math.max(4, Math.min(window.innerWidth - 52, btnDragStart.value.elX + dx)),
    y: Math.max(4, Math.min(window.innerHeight - 52, btnDragStart.value.elY + dy)),
  }
}
function onBtnDragEnd() {
  isDraggingBtn.value = false
  document.removeEventListener('mousemove', onBtnDragMove)
  document.removeEventListener('mouseup', onBtnDragEnd)
  document.removeEventListener('touchmove', onBtnDragMove)
  document.removeEventListener('touchend', onBtnDragEnd)
  saveBtnPos()
}

// debounce map for content auto-save
const debounceTimers: Record<number, ReturnType<typeof setTimeout>> = {}

// --- Drag state ---
interface DragState {
  noteId: number | null
  type: 'move' | 'resize'
  startX: number
  startY: number
  startPosX: number
  startPosY: number
  startWidth: number
  startHeight: number
}

const dragState = ref<DragState>({
  noteId: null,
  type: 'move',
  startX: 0,
  startY: 0,
  startPosX: 0,
  startPosY: 0,
  startWidth: 0,
  startHeight: 0,
})

function loadNotes() {
  getStickyNoteList<Common.ListResponse<NoteData[]>>().then(({ code, data }) => {
    if (code === 0) {
      notes.value = data.list || []
      let max = 1
      for (const n of notes.value) {
        if (n.zIndex > max)
          max = n.zIndex
      }
      maxZIndex = max
    }
  })
}

function handleCreate() {
  // Place new note at a slightly random offset
  const offsetIdx = notes.value.filter(n => n.status === 1).length
  const newNote: NoteData = {
    id: 0,
    content: '',
    color: PRESET_COLORS[0],
    posX: 50 + (offsetIdx % 5) * 30,
    posY: 50 + (offsetIdx % 5) * 30,
    width: DEFAULT_WIDTH,
    height: DEFAULT_HEIGHT,
    zIndex: ++maxZIndex,
    status: 1,
  }
  createStickyNote<Common.Response<number>>(newNote).then(({ code, data }) => {
    if (code === 0) {
      newNote.id = data
      notes.value.push({ ...newNote })
    }
  })
}

function handleDelete(id: number) {
  deleteStickyNote({ id }).then(({ code }) => {
    if (code === 0) {
      notes.value = notes.value.filter(n => n.id !== id)
      ms.success('便签已删除')
    }
  })
}

function handleToggleStatus(note: NoteData) {
  const newStatus = note.status === 1 ? 0 : 1
  updateStickyNote({
    id: note.id,
    status: newStatus,
  }).then(({ code }) => {
    if (code === 0) {
      // 成功后更新本地状态
      const found = notes.value.find(n => n.id === note.id)
      if (found) {
        found.status = newStatus
      }
      ms.success(newStatus === 1 ? '便签已启用' : '便签已停用')
    } else {
      ms.error('操作失败')
    }
  }).catch(() => {
    ms.error('操作失败')
  })
}

function handleColorChange(note: NoteData, color: string) {
  note.color = color
  updateStickyNote({
    id: note.id,
    color,
  })
}

function bringToFront(note: NoteData) {
  note.zIndex = ++maxZIndex
}

function handleContentInput(note: NoteData, value: string) {
  note.content = value
  // debounce auto-save 500ms
  if (debounceTimers[note.id]) {
    clearTimeout(debounceTimers[note.id])
  }
  debounceTimers[note.id] = setTimeout(() => {
    updateStickyNote({
      id: note.id,
      content: note.content,
    })
  }, 500)
}

// --- Mouse drag (move) ---
function onNoteMouseDown(e: MouseEvent, note: NoteData) {
  // Only left button, ignore if target is textarea or resize handle
  if (e.button !== 0) return
  const target = e.target as HTMLElement
  if (target.tagName === 'TEXTAREA' || target.classList.contains('sticky-resize-handle') || target.closest('.sticky-color-dot') || target.closest('.sticky-delete-btn') || target.closest('.sticky-status-btn')) return

  bringToFront(note)
  dragState.value = {
    noteId: note.id,
    type: 'move',
    startX: e.clientX,
    startY: e.clientY,
    startPosX: note.posX,
    startPosY: note.posY,
    startWidth: note.width,
    startHeight: note.height,
  }
  e.preventDefault()
}

// --- Mouse drag (resize) ---
function onResizeMouseDown(e: MouseEvent, note: NoteData) {
  if (e.button !== 0) return
  e.stopPropagation()
  e.preventDefault()
  dragState.value = {
    noteId: note.id,
    type: 'resize',
    startX: e.clientX,
    startY: e.clientY,
    startPosX: note.posX,
    startPosY: note.posY,
    startWidth: note.width,
    startHeight: note.height,
  }
}

function onMouseMove(e: MouseEvent) {
  const ds = dragState.value
  if (!ds.noteId) return

  const note = notes.value.find(n => n.id === ds.noteId)
  if (!note) return

  if (ds.type === 'move') {
    const dx = e.clientX - ds.startX
    const dy = e.clientY - ds.startY
    note.posX = ds.startPosX + dx
    note.posY = ds.startPosY + dy
  }
  else if (ds.type === 'resize') {
    const dx = e.clientX - ds.startX
    const dy = e.clientY - ds.startY
    note.width = Math.max(MIN_WIDTH, ds.startWidth + dx)
    note.height = Math.max(MIN_HEIGHT, ds.startHeight + dy)
  }
}

function onMouseUp() {
  const ds = dragState.value
  if (!ds.noteId) return

  const note = notes.value.find(n => n.id === ds.noteId)
  if (note) {
    updateStickyNote({
      id: note.id,
      posX: Math.round(note.posX),
      posY: Math.round(note.posY),
      width: Math.round(note.width),
      height: Math.round(note.height),
      zIndex: note.zIndex,
    })
  }
  dragState.value.noteId = null
}

onMounted(() => {
  loadNotes()
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
})

onUnmounted(() => {
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', onMouseUp)
  // Clear all debounce timers
  for (const key in debounceTimers) {
    if (debounceTimers[key]) clearTimeout(debounceTimers[key])
  }
})
</script>

<template>
  <div class="sticky-notes-container">
    <div
      v-for="note in notes"
      :key="note.id"
      class="sticky-note-card"
      :style="{
        left: `${note.posX}px`,
        top: `${note.posY}px`,
        width: `${note.width}px`,
        height: `${note.height}px`,
        zIndex: note.zIndex,
        backgroundColor: note.color,
      }"
      @mousedown="onNoteMouseDown($event, note)"
    >
      <!-- Top bar -->
      <div class="sticky-note-header">
        <!-- Color dots -->
        <div class="sticky-colors">
          <span
            v-for="color in PRESET_COLORS"
            :key="color"
            class="sticky-color-dot"
            :class="{ active: note.color === color }"
            :style="{ backgroundColor: color }"
            @mousedown.stop
            @click.stop="handleColorChange(note, color)"
          />
        </div>
        <div class="sticky-actions">
          <!-- Enable/Disable toggle -->
          <span class="sticky-btn sticky-status-btn" :title="note.status === 1 ? '已启用（点击停用）' : '已停用（点击启用）'" @mousedown.stop @click.stop="handleToggleStatus(note)">
            <SvgIcon :icon="note.status === 1 ? 'mdi:eye' : 'mdi:eye-off'" style="font-size: 13px;" />
          </span>
          <!-- Delete button - more visible -->
          <NPopconfirm @positive-click="handleDelete(note.id)">
            <template #trigger>
              <span class="sticky-btn sticky-delete-btn" @mousedown.stop @click.stop>
                <SvgIcon icon="mdi:close" style="font-size: 16px; color: #e74c3c;" />
              </span>
            </template>
            确定要删除这个便签吗？
          </NPopconfirm>
        </div>
      </div>
      <!-- Content area -->
      <div class="sticky-note-body">
        <NInput
          type="textarea"
          :autosize="true"
          :value="note.content"
          :bordered="false"
          :placeholder="note.status === 0 ? '（已停用）' : ''"
          :disabled="note.status === 0"
          class="sticky-textarea"
          @update:value="handleContentInput(note, $event)"
          @mousedown.stop
        />
      </div>
      <!-- Resize handle -->
      <div
        class="sticky-resize-handle"
        @mousedown.stop="onResizeMouseDown($event, note)"
      />
    </div>

    <!-- Floating add button - 可拖拽小绿点 -->
    <div
      class="sticky-add-btn"
      :style="{ left: btnPos.x + 'px', bottom: btnPos.y + 'px' }"
      @mousedown="onBtnDragStart"
      @touchstart="onBtnDragStart"
    >
      <NButton circle type="primary" size="tiny" color="#22c55e" @click="!isDraggingBtn && handleCreate()">
        <template #icon>
          <SvgIcon icon="mdi:note-plus-outline" :size="14" />
        </template>
      </NButton>
    </div>
  </div>
</template>

<style scoped>
.sticky-notes-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.sticky-note-card {
  position: absolute;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.12);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  pointer-events: auto;
  cursor: move;
  user-select: none;
  backdrop-filter: blur(2px);
  transition: box-shadow 0.2s;
}

.sticky-note-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.18);
}

.sticky-note-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 6px 2px 6px;
  flex-shrink: 0;
}

.sticky-colors {
  display: flex;
  gap: 4px;
  align-items: center;
}

.sticky-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.sticky-color-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  cursor: pointer;
  border: 1px solid rgba(0, 0, 0, 0.15);
  transition: transform 0.15s, box-shadow 0.15s;
}
.sticky-color-dot:nth-child(9) {
  border: 2px dashed rgba(0,0,0,0.2);
  background-image: repeating-linear-gradient(45deg, rgba(0,0,0,0.05) 0px, rgba(0,0,0,0.05) 2px, transparent 2px, transparent 6px) !important;
}

.sticky-color-dot:hover {
  transform: scale(1.3);
}

.sticky-color-dot.active {
  box-shadow: 0 0 0 2px rgba(0, 0, 0, 0.25);
  transform: scale(1.2);
}

.sticky-btn {
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  transition: background-color 0.15s, opacity 0.15s;
}

.sticky-status-btn:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

.sticky-delete-btn:hover {
  background-color: rgba(231, 76, 60, 0.15);
}

.sticky-note-body {
  flex: 1;
  overflow: hidden;
  padding: 0 6px 2px 6px;
}

.sticky-textarea {
  width: 100%;
  height: 100%;
  background-color: transparent !important;
}

.sticky-textarea :deep(.n-input__textarea-el) {
  background-color: transparent !important;
  box-shadow: none !important;
  resize: none;
}

.sticky-resize-handle {
  position: absolute;
  right: 0;
  bottom: 0;
  width: 14px;
  height: 14px;
  cursor: nwse-resize;
  opacity: 0.3;
  background: linear-gradient(135deg, transparent 50%, rgba(0, 0, 0, 0.3) 50%);
  border-radius: 0 0 8px 0;
}

.sticky-add-btn {
  position: fixed;
  pointer-events: auto;
  z-index: 9999;
  padding: 4px;
  cursor: grab;
  user-select: none;
}
.sticky-add-btn:active {
  cursor: grabbing;
}
</style>