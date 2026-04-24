<template>
	<svg
		class="gantt-row-bars"
		:width="totalWidth"
		height="72"
		xmlns="http://www.w3.org/2000/svg"
		role="img"
		:aria-label="$t('project.gantt.taskBarsForRow', { rowId })"
		:data-row-id="rowId"
	>
		<GanttBarPrimitive
			v-for="bar in bars"
			:key="bar.id"
			:model="bar"
			:timeline-start="dateFromDate"
			:timeline-end="dateToDate"
			:on-update="(id, start, end) => emit('updateTask', id, start, end)"
		>
			<!-- Gradient definitions for partial-date bars -->
			<defs v-if="bar.meta?.dateType === 'startOnly' || bar.meta?.dateType === 'endOnly'">
				<linearGradient
					:id="`gradient-${bar.id}`"
					x1="0"
					y1="0"
					x2="1"
					y2="0"
				>
					<stop
						v-if="bar.meta?.dateType === 'endOnly'"
						offset="0%"
						:stop-color="getBarFill(bar)"
						stop-opacity="0"
					/>
					<stop
						v-if="bar.meta?.dateType === 'endOnly'"
						offset="40%"
						:stop-color="getBarFill(bar)"
						stop-opacity="1"
					/>
					<stop
						v-if="bar.meta?.dateType === 'startOnly'"
						offset="60%"
						:stop-color="getBarFill(bar)"
						stop-opacity="1"
					/>
					<stop
						v-if="bar.meta?.dateType === 'startOnly'"
						offset="100%"
						:stop-color="getBarFill(bar)"
						stop-opacity="0"
					/>
				</linearGradient>
			</defs>

			<!-- Main bar (regular task) -->
			<rect
				v-if="!bar.meta?.isParent"
				:x="getBarX(bar)"
				:y="4"
				:width="getBarWidth(bar)"
				:height="52"
				:rx="4"
				:fill="getBarFillAttr(bar)"
				:opacity="bar.meta?.isDone ? 0.5 : 1"
				:stroke="getBarStroke(bar)"
				:stroke-width="getBarStrokeWidth(bar)"
				:stroke-dasharray="isDateless(bar) ? '5,5' : 'none'"
				class="gantt-bar"
				role="button"
				:aria-label="getBarAriaLabel(bar)"
				:aria-pressed="isRowFocused"
				@pointerdown="handleBarPointerDown(bar, $event)"
			/>

			<!-- Progress bar overlay -->
			<rect
				v-if="!bar.meta?.isParent && getBarPercentDone(bar) > 0"
				:x="getBarX(bar)"
				:y="4"
				:width="getBarWidth(bar) * getBarPercentDone(bar)"
				:height="52"
				:rx="4"
				fill="#ff8c00"
				opacity="1"
				pointer-events="none"
			/>

			<!-- Parent summary bar (full height with diamond endpoints) -->
			<g
				v-if="bar.meta?.isParent"
				class="gantt-bar gantt-parent-bar"
				role="button"
				:aria-label="getBarAriaLabel(bar)"
				:aria-pressed="isRowFocused"
				@pointerdown="handleBarPointerDown(bar, $event)"
			>
				<rect
					:x="getBarX(bar)"
					:y="4"
					:width="getBarWidth(bar)"
					:height="52"
					:rx="4"
					:fill="getBarFillAttr(bar)"
					:opacity="bar.meta?.isDone ? 0.5 : 1"
					:stroke="getBarStroke(bar)"
					:stroke-width="getBarStrokeWidth(bar)"
					:stroke-dasharray="bar.meta?.hasDerivedDates ? '4,2' : 'none'"
				/>
				<!-- Left diamond -->
				<polygon
					:points="getLeftDiamondPoints(bar)"
					:fill="getParentDiamondFill(bar)"
					:stroke="getBarFillAttr(bar)"
					stroke-width="1"
					:opacity="bar.meta?.isDone ? 0.5 : 1"
				/>
				<!-- Right diamond -->
				<polygon
					:points="getRightDiamondPoints(bar)"
					:fill="getParentDiamondFill(bar)"
					:stroke="getBarFillAttr(bar)"
					stroke-width="1"
					:opacity="bar.meta?.isDone ? 0.5 : 1"
				/>
			</g>

			<!-- Left resize handle (hidden for endOnly bars) -->
			<rect
				v-if="bar.meta?.dateType !== 'endOnly'"
				:x="getBarX(bar) - RESIZE_HANDLE_OFFSET"
				:y="4"
				:width="6"
				:height="52"
				:rx="3"
				fill="var(--white)"
				stroke="var(--primary)"
				stroke-width="1"
				class="gantt-resize-handle gantt-resize-left"
				role="button"
				:aria-label="$t('project.gantt.resizeStartDate', { task: bar.meta?.label || bar.id })"
				@pointerdown="startResize(bar, 'start', $event)"
			/>

			<!-- Right resize handle (hidden for startOnly bars) -->
			<rect
				v-if="bar.meta?.dateType !== 'startOnly'"
				:x="getBarX(bar) + getBarWidth(bar) - RESIZE_HANDLE_OFFSET"
				:y="4"
				:width="6"
				:height="52"
				:rx="3"
				fill="var(--white)"
				stroke="var(--primary)"
				stroke-width="1"
				class="gantt-resize-handle gantt-resize-right"
				role="button"
				:aria-label="$t('project.gantt.resizeEndDate', { task: bar.meta?.label || bar.id })"
				@pointerdown="startResize(bar, 'end', $event)"
			/>

			<!-- Task label and labels with clipping -->
			<foreignObject
				:x="getBarX(bar)"
				:y="4"
				:width="getBarWidth(bar)"
				:height="52"
				aria-hidden="true"
			>
				<div
					xmlns="http://www.w3.org/1999/xhtml"
					class="gantt-bar-content"
					:style="{
						color: getBarTextColor(bar),
						textDecoration: bar.meta?.isDone ? 'line-through' : 'none',
					}"
				>
					<span class="gantt-bar-title-row">
						<AssigneeList
							v-if="getTaskAssignees(bar).length > 0"
							:assignees="getTaskAssignees(bar)"
							:avatar-size="20"
							class="gantt-bar-assignees"
						/>
						<span class="gantt-bar-title">{{ bar.meta?.label || bar.id }}</span>
					</span>
					<span class="gantt-bar-bottom-row">
						<span
							v-if="getTaskLabels(bar).length > 0"
							class="gantt-bar-labels"
						>
							<span
								v-for="label in getTaskLabels(bar)"
								:key="label.id"
								class="gantt-bar-label-badge"
								:style="{
									backgroundColor: label.hexColor,
									color: label.textColor,
								}"
							>{{ label.title }}</span>
						</span>
						<small
							v-if="bar.meta?.bucketName"
							class="gantt-bar-bucket"
						>{{ bar.meta.bucketName }}</small>
					</span>
				</div>
			</foreignObject>

			<!-- Hover zone to create a new milestone subtask (under markers) -->
			<rect
				:x="getBarX(bar)"
				:y="56"
				:width="getBarWidth(bar)"
				:height="16"
				fill="transparent"
				class="gantt-due-hover-zone"
				@pointermove="onDueHoverMove(bar, $event)"
				@pointerleave="onDueHoverLeave(bar)"
				@click.stop="onDueHoverClick(bar, $event)"
			/>

			<!-- Milestone markers (one per subtask with endDate) -->
			<g
				v-for="milestone in (bar.meta?.milestones ?? [])"
				:key="milestone.taskId"
				class="gantt-milestone-group"
			>
				<line
					:x1="getMilestoneCx(bar, milestone)"
					:x2="getMilestoneCx(bar, milestone)"
					:y1="DUE_STEM_TOP"
					:y2="DUE_STEM_BOTTOM"
					stroke="var(--grey-500)"
					stroke-width="2"
					stroke-linecap="round"
					pointer-events="none"
				/>
				<text
					v-if="collapsedIds.has(Number(bar.id))"
					:x="getMilestoneCx(bar, milestone) - DUE_DIAMOND_SIZE - 4"
					:y="DUE_DIAMOND_CY"
					text-anchor="end"
					dominant-baseline="middle"
					class="gantt-milestone-label"
				>{{ milestone.title }}</text>
				<polygon
					:points="getMilestonePoints(bar, milestone)"
					:fill="milestone.done ? '#ff8c00' : 'var(--danger)'"
					stroke="var(--white)"
					stroke-width="1.5"
					class="gantt-due-diamond"
					role="button"
					:aria-label="milestone.title"
					@pointerdown.stop="startMilestoneDrag(bar, milestone, $event)"
				>
					<title>{{ milestone.title }}</title>
				</polygon>
			</g>

			<!-- Ghost stem + diamond on hover -->
			<g
				v-if="hoveredDueDay?.barId === bar.id"
				pointer-events="none"
				opacity="0.5"
			>
				<line
					:x1="getGhostDiamondCenterX(bar) ?? 0"
					:x2="getGhostDiamondCenterX(bar) ?? 0"
					:y1="DUE_STEM_TOP"
					:y2="DUE_STEM_BOTTOM"
					stroke="var(--grey-500)"
					stroke-width="2"
					stroke-linecap="round"
				/>
				<polygon
					:points="getGhostDiamondPoints(bar)"
					fill="var(--danger)"
					stroke="var(--white)"
					stroke-width="1.5"
				/>
			</g>
		</GanttBarPrimitive>

		<!-- Collapse/expand chevron for each parent task in this row -->
		<g
			v-for="pbar in parentBars"
			:key="'chev-' + pbar.id"
			class="gantt-collapse-toggle"
			:transform="`translate(${Math.max(0, getBarX(pbar) - 14)}, 24)`"
			role="button"
			:aria-label="collapsedIds.has(Number(pbar.id))
				? $t('project.gantt.expandGroup', { task: pbar.meta?.label || '' })
				: $t('project.gantt.collapseGroup', { task: pbar.meta?.label || '' })"
			tabindex="0"
			@pointerdown.stop="emit('toggleCollapse', Number(pbar.id))"
			@keydown.enter.stop="emit('toggleCollapse', Number(pbar.id))"
		>
			<rect
				x="-2"
				y="-2"
				width="14"
				height="14"
				fill="transparent"
			/>
			<polygon
				v-if="collapsedIds.has(Number(pbar.id))"
				points="2,0 10,5 2,10"
				fill="var(--grey-500)"
			/>
			<polygon
				v-else
				points="0,2 10,2 5,10"
				fill="var(--grey-500)"
			/>
		</g>
	</svg>
</template>

<script setup lang="ts">
import {computed, ref} from 'vue'
import dayjs from 'dayjs'
import {useI18n} from 'vue-i18n'

import type {GanttBarModel} from '@/composables/useGanttBar'
import {getTextColor, LIGHT} from '@/helpers/color/getTextColor'
import {MILLISECONDS_A_DAY} from '@/constants/date'
import {roundToNaturalDayBoundary} from '@/helpers/time/roundToNaturalDayBoundary'

import GanttBarPrimitive from './primitives/GanttBarPrimitive.vue'
import AssigneeList from '@/components/tasks/partials/AssigneeList.vue'
import type {IUser} from '@/modelTypes/IUser'

const props = defineProps<{
	bars: GanttBarModel[]
	totalWidth: number
	dateFromDate: Date
	dateToDate: Date
	dayWidthPixels: number
	isDragging: boolean
	isResizing: boolean
	isDraggingMilestone: boolean
	dragState: {
		barId: string
		startX: number
		originalStart: Date
		originalEnd: Date
		currentDays: number
		edge?: 'start' | 'end'
	} | null
	milestoneDragState: {
		barId: string
		taskId: number
		startX: number
		originalDate: Date
		currentDays: number
	} | null
	focusedRow: string | null
	focusedCell: number | null
	rowId: string
	collapsedIds: Set<number>
}>()

const emit = defineEmits<{
	(e: 'barPointerDown', bar: GanttBarModel, event: PointerEvent): void
	(e: 'startResize', bar: GanttBarModel, edge: 'start' | 'end', event: PointerEvent): void
	(e: 'startMilestoneDrag', bar: GanttBarModel, milestone: {taskId: number, date: Date, title: string, done: boolean, dateField: 'dueDate' | 'endDate'}, event: PointerEvent): void
	(e: 'createMilestone', bar: GanttBarModel, dayIndex: number, clientX: number, clientY: number): void
	(e: 'updateTask', id: string, newStart: Date, newEnd: Date): void
	(e: 'toggleCollapse', taskId: number): void
}>()

const {t} = useI18n({useScope: 'global'})

const RESIZE_HANDLE_OFFSET = 3

function addDays(dateOrValue: Date | string | number, days: number): Date {
	const date = new Date(dateOrValue)
	const newDate = new Date(date)
	newDate.setDate(newDate.getDate() + days)
	return newDate
}

const isRowFocused = computed(() => props.focusedRow === props.rowId)

const parentBars = computed(() => props.bars.filter(b => b.meta?.isParent))

function computeBarX(startDate: Date) {
	const daysDiff = dayjs(startDate).diff(dayjs(props.dateFromDate), 'day')
	const x = daysDiff * props.dayWidthPixels
	return x
}

function getDaysDifference(startDate: Date, endDate: Date): number {
	return Math.ceil(
		(roundToNaturalDayBoundary(endDate).getTime() - roundToNaturalDayBoundary(startDate, true).getTime()) /
MILLISECONDS_A_DAY,
	)
}

function computeBarWidth(bar: GanttBarModel) {
	const diff = getDaysDifference(bar.start, bar.end)
	const width = diff * props.dayWidthPixels
	return width
}

const originalStartX = computed(() => props.dragState?.originalStart
	? computeBarX(props.dragState.originalStart)
	: 0)

const getBarX = computed(() => (bar: GanttBarModel) => {
	if (props.isDragging && props.dragState?.barId === bar.id) {
		const offset = props.dragState.currentDays * props.dayWidthPixels
		return originalStartX.value + offset
	}

	if (props.isResizing && props.dragState?.barId === bar.id && props.dragState.edge === 'start') {
		const newStart = addDays(props.dragState.originalStart, props.dragState.currentDays)
		return computeBarX(newStart)
	}
	return computeBarX(bar.start)
})

const getBarWidth = computed(() => (bar: GanttBarModel) => {
	if (props.isResizing && props.dragState?.barId === bar.id) {
		if (props.dragState.edge === 'start') {
			const newStart = addDays(props.dragState.originalStart, props.dragState.currentDays)
			return Math.max(0, getDaysDifference(newStart, props.dragState.originalEnd) * props.dayWidthPixels)
		} else {
			const newEnd = addDays(props.dragState.originalEnd, props.dragState.currentDays)
			return Math.max(0, getDaysDifference(props.dragState.originalStart, newEnd) * props.dayWidthPixels)
		}
	}
	return computeBarWidth(bar)
})

// Diamond endpoint helpers for parent summary bars
const DIAMOND_SIZE = 5

function getLeftDiamondPoints(bar: GanttBarModel): string {
	const x = getBarX.value(bar) - DIAMOND_SIZE
	const cy = 30 // vertical center of the bar
	return `${x},${cy} ${x + DIAMOND_SIZE},${cy - DIAMOND_SIZE} ${x + DIAMOND_SIZE * 2},${cy} ${x + DIAMOND_SIZE},${cy + DIAMOND_SIZE}`
}

function getRightDiamondPoints(bar: GanttBarModel): string {
	const x = getBarX.value(bar) + getBarWidth.value(bar) + DIAMOND_SIZE
	const cy = 30
	return `${x - DIAMOND_SIZE * 2},${cy} ${x - DIAMOND_SIZE},${cy - DIAMOND_SIZE} ${x},${cy} ${x - DIAMOND_SIZE},${cy + DIAMOND_SIZE}`
}

function getParentDiamondFill(bar: GanttBarModel): string {
	// Use a darker shade for contrast on the full-height bar
	if (bar.meta?.color) {
		return 'var(--white)'
	}
	return 'var(--white)'
}

function isPartialDate(bar: GanttBarModel) {
	return bar.meta?.dateType === 'startOnly' || bar.meta?.dateType === 'endOnly'
}

function isDateless(bar: GanttBarModel) {
	return !bar.meta?.hasActualDates && !isPartialDate(bar)
}

function getBarFill(bar: GanttBarModel) {
	// Partial dates still have "actual" dates on one side — use the task color
	if (isPartialDate(bar)) {
		if (bar.meta?.color) {
			return bar.meta.color
		}
		return 'var(--primary)'
	}

	if (bar.meta?.hasActualDates) {
		if (bar.meta?.color) {
			return bar.meta.color
		}
		return 'var(--primary)'
	}

	return 'var(--grey-100)'
}

function getBarFillAttr(bar: GanttBarModel): string {
	if (isPartialDate(bar)) {
		return `url(#gradient-${bar.id})`
	}
	return getBarFill(bar)
}

function getBarStroke(bar: GanttBarModel) {
	if (isDateless(bar)) {
		return 'var(--grey-300)' // Gray for dashed border
	}
	return 'none'
}

function getBarStrokeWidth(bar: GanttBarModel) {
	if (isDateless(bar)) {
		return '2'
	}
	return '0'
}

function getBarTextColor(bar: GanttBarModel) {
	if (isDateless(bar)) {
		return 'var(--grey-800)'
	}

	if (bar.meta?.color) {
		return getTextColor(bar.meta.color)
	}

	return LIGHT
}

function getBarAriaLabel(bar: GanttBarModel): string {
	const task = bar.meta?.label || bar.id
	const startDate = bar.start.toLocaleDateString()
	const endDate = bar.end.toLocaleDateString()

	let dateType: string
	if (bar.meta?.dateType === 'startOnly') {
		dateType = t('project.gantt.partialDatesStart')
	} else if (bar.meta?.dateType === 'endOnly') {
		dateType = t('project.gantt.partialDatesEnd')
	} else if (bar.meta?.hasActualDates) {
		dateType = t('project.gantt.scheduledDates')
	} else {
		dateType = t('project.gantt.estimatedDates')
	}

	return t('project.gantt.taskBarLabel', {task, startDate, endDate, dateType})
}

function handleBarPointerDown(bar: GanttBarModel, event: PointerEvent) {
	emit('barPointerDown', bar, event)
}

function startResize(bar: GanttBarModel, edge: 'start' | 'end', event: PointerEvent) {
	emit('startResize', bar, edge, event)
}

const DUE_DIAMOND_SIZE = 5
const DUE_STEM_TOP = 56
const DUE_STEM_BOTTOM = 61
const DUE_DIAMOND_CY = 66

function startMilestoneDrag(bar: GanttBarModel, milestone: {taskId: number, date: Date, title: string, done: boolean, dateField: 'dueDate' | 'endDate'}, event: PointerEvent) {
	emit('startMilestoneDrag', bar, milestone, event)
}

function getMilestoneCx(bar: GanttBarModel, milestone: {taskId: number, date: Date}): number {
	let x = computeBarX(milestone.date) + props.dayWidthPixels / 2
	if (props.isDraggingMilestone && props.milestoneDragState?.barId === bar.id && props.milestoneDragState.taskId === milestone.taskId) {
		x += props.milestoneDragState.currentDays * props.dayWidthPixels
	}
	return x
}

function getMilestonePoints(bar: GanttBarModel, milestone: {taskId: number, date: Date}): string {
	const cx = getMilestoneCx(bar, milestone)
	const s = DUE_DIAMOND_SIZE
	return `${cx},${DUE_DIAMOND_CY - s} ${cx + s},${DUE_DIAMOND_CY} ${cx},${DUE_DIAMOND_CY + s} ${cx - s},${DUE_DIAMOND_CY}`
}

const hoveredDueDay = ref<{barId: string, day: number} | null>(null)

function computeDayFromEvent(event: PointerEvent | MouseEvent): number | null {
	const svg = (event.currentTarget as Element)?.closest('svg')
	if (!svg) return null
	const rect = svg.getBoundingClientRect()
	const x = event.clientX - rect.left
	return Math.floor(x / props.dayWidthPixels)
}

function onDueHoverMove(bar: GanttBarModel, event: PointerEvent) {
	const day = computeDayFromEvent(event)
	if (day === null) return
	if (hoveredDueDay.value?.barId === bar.id && hoveredDueDay.value.day === day) return
	hoveredDueDay.value = {barId: bar.id, day}
}

function onDueHoverLeave(bar: GanttBarModel) {
	if (hoveredDueDay.value?.barId === bar.id) {
		hoveredDueDay.value = null
	}
}

function onDueHoverClick(bar: GanttBarModel, event: MouseEvent) {
	const day = computeDayFromEvent(event)
	if (day === null) return
	hoveredDueDay.value = null
	emit('createMilestone', bar, day, event.clientX, event.clientY)
}

function getGhostDiamondCenterX(bar: GanttBarModel): number | null {
	if (!hoveredDueDay.value || hoveredDueDay.value.barId !== bar.id) return null
	return hoveredDueDay.value.day * props.dayWidthPixels + props.dayWidthPixels / 2
}

function getGhostDiamondPoints(bar: GanttBarModel): string {
	const cx = getGhostDiamondCenterX(bar)
	if (cx === null) return ''
	const s = DUE_DIAMOND_SIZE
	return `${cx},${DUE_DIAMOND_CY - s} ${cx + s},${DUE_DIAMOND_CY} ${cx},${DUE_DIAMOND_CY + s} ${cx - s},${DUE_DIAMOND_CY}`
}

function getBarPercentDone(bar: GanttBarModel): number {
	const task = bar.meta?.task as {percentDone?: number} | undefined
	return task?.percentDone ?? 0
}

function getTaskLabels(bar: GanttBarModel): Array<{id: number, title: string, hexColor: string, textColor: string}> {
	const task = bar.meta?.task as {labels?: Array<{id: number, title: string, hexColor: string, textColor: string}>} | undefined
	return task?.labels ?? []
}

function getTaskAssignees(bar: GanttBarModel): IUser[] {
	const task = bar.meta?.task as {assignees?: IUser[]} | undefined
	return task?.assignees ?? []
}
</script>

<style scoped lang="scss">
.gantt-row-bars {
	position: absolute;
	inset-block-start: 0;
	inset-inline-start: 0;
	pointer-events: none;
	z-index: 4;

	.gantt-bar {
		cursor: grab;
		pointer-events: all;

		&:hover {
			opacity: 0.8;
		}

		&:active {
			cursor: grabbing;
		}
	}

	:deep(text) {
		pointer-events: none;
		user-select: none;
	}
}

.gantt-collapse-toggle {
	pointer-events: all;
	cursor: pointer;

	&:hover polygon {
		fill: var(--grey-700);
	}

	&:focus {
		outline: none;

		polygon {
			fill: var(--primary);
		}
	}
}

.gantt-bar-text {
	font-size: .85rem;
	pointer-events: none;
	user-select: none;
}

.gantt-bar-content {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: flex-end;
	gap: 2px;
	height: 100%;
	padding: 2px 8px;
	overflow: hidden;
	white-space: nowrap;
	font-size: .8rem;
	pointer-events: none;
	user-select: none;
}

.gantt-bar-title-row {
	display: flex;
	align-items: center;
	justify-content: flex-end;
	gap: 4px;
	max-width: 100%;
	overflow: hidden;
}

.gantt-bar-assignees {
	flex-shrink: 0;
	pointer-events: auto;
}

.gantt-bar-title {
	font-weight: bold;
	overflow: hidden;
	text-overflow: ellipsis;
	min-width: 0;
	flex-shrink: 1;
}

.gantt-bar-bottom-row {
	display: flex;
	align-items: center;
	justify-content: flex-end;
	gap: 4px;
	max-width: 100%;
	overflow: hidden;
}

.gantt-bar-labels {
	display: inline-flex;
	gap: 2px;
	flex-direction: row-reverse;
	flex-shrink: 1;
	overflow: hidden;
}

.gantt-bar-label-badge {
	border-radius: 3px;
	padding: 0 3px;
	font-size: .65rem;
	line-height: 1.4;
	flex-shrink: 0;
}

.gantt-bar-bucket {
	opacity: 0.7;
	font-weight: normal;
	font-size: .7rem;
	flex-shrink: 0;
}

.gantt-parent-bar {
	cursor: grab;
	pointer-events: all;
}

.gantt-milestone-label {
	font-size: 0.7rem;
	fill: var(--grey-700);
	pointer-events: none;
	user-select: none;
	opacity: 0;
	transition: opacity 0.15s ease;
}

.gantt-milestone-group:hover .gantt-milestone-label {
	opacity: 1;
}

.gantt-due-diamond {
	cursor: grab;
	pointer-events: all;

	&:hover {
		filter: brightness(1.1);
	}

	&:active {
		cursor: grabbing;
	}
}

.gantt-due-hover-zone {
	cursor: crosshair;
	pointer-events: all;
}

:deep(.gantt-resize-handle) {
	cursor: col-resize !important;
	opacity: 0;
	transition: opacity 0.2s ease;
	pointer-events: all; // Ensure they receive pointer events
}

// Show resize handles on bar hover
:deep(g:hover) .gantt-resize-handle {
	opacity: 0.8;

	&:hover {
		opacity: 1;
		cursor: inherit; // Use the specific cursor defined above
	}
}

// Focus styles for task bars
:deep(g[role="slider"]:focus) {
	outline: none; // Remove default browser outline
	
	.gantt-bar {
		stroke: var(--primary) !important;
		stroke-width: 3 !important;
	}
}
</style>
