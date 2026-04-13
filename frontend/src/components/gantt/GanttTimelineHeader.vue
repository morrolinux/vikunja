<template>
	<div
		class="gantt-timeline"
		role="columnheader"
		:aria-label="$t('project.gantt.timelineHeader')"
	>
		<!-- Upper timeunit for months -->
		<div
			class="gantt-timeline-months"
			role="row"
			:aria-label="$t('project.gantt.monthsRow')"
		>
			<div
				v-for="monthGroup in monthGroups"
				:key="monthGroup.key"
				class="timeunit-month"
				:style="{ width: `${monthGroup.width}px` }"
				role="columnheader"
				:aria-label="$t('project.gantt.monthLabel', {month: monthGroup.label})"
			>
				{{ monthGroup.label }} {{ monthGroup.sponsorInfo }}
			</div>
		</div>
        
		<!-- Lower timeunit for days -->
		<div
			class="gantt-timeline-days"
			role="row"
			:aria-label="$t('project.gantt.daysRow')"
		>
			<div
				v-for="date in timelineData"
				:key="date.toISOString()"
				class="timeunit"
				:style="{ width: `${dayWidthPixels}px` }"
				role="columnheader"
				:aria-label="dateIsToday(date) 
					? $t('project.gantt.dayLabelToday', {
						date: date.toLocaleDateString(),
						weekday: weekDayFromDate(date)
					})
					: $t('project.gantt.dayLabel', {
						date: date.toLocaleDateString(),
						weekday: weekDayFromDate(date)
					})"
			>
				<div
					class="timeunit-wrapper"
					:class="{'today': dateIsToday(date), 'special-day': dateIsSpecial(date)}"
				>
					<span>{{ date.getDate() }}</span>
					<span class="weekday">
						{{ weekDayFromDate(date) }}
					</span>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import {useGlobalNow} from '@/composables/useGlobalNow'
import {useWeekDayFromDate} from '@/helpers/time/formatDate'
import dayjs from 'dayjs'

import type {ITask} from '@/modelTypes/ITask'

const props = defineProps<{
    timelineData: Date[]
    dayWidthPixels: number
    tasks?: Map<ITask['id'], ITask>
}>()

const weekDayFromDate = useWeekDayFromDate()
const { now: today } = useGlobalNow()

const dateIsToday = computed(() => {
	const todayStr = today.value.toDateString()
	return (date: Date) => date.toDateString() === todayStr
})

function dateIsSpecial(date: Date): boolean {
	return date.getDay() === 1
}

// TODO: do not hardcode publish dates, let the user decide
function countPublishDays(year: number, month: number): number {
	const firstDay = new Date(year, month, 1)
	const lastDay = new Date(year, month + 1, 0)
	let count = 0
	const current = new Date(firstDay)
	while (current <= lastDay) {
		if (current.getDay() === 1) {
			count++
		}
		current.setDate(current.getDate() + 1)
	}
	return count
}

function countSponsorsInMonth(year: number, month: number): number {
	if (!props.tasks) return 0
	let count = 0
	const isSponsorTask = (task: ITask) =>
		task.labels.some((label) => label.description?.toLowerCase().includes('sponsor'))

	props.tasks.forEach(task => {
		if (task.endDate) {
			const taskEndMonth = task.endDate.getMonth()
			const taskEndDay = task.endDate.getDate()
			const taskEndYear = task.endDate.getFullYear()
			// Account for inclusive display: if task ends on the 1st of next month,
			// it visually belongs to the previous month
			const belongsToMonth =
				(taskEndDay > 1 && taskEndMonth === month && taskEndYear === year) ||
				(taskEndDay === 1 && taskEndMonth === month + 1 && taskEndYear === year)
			if (belongsToMonth && isSponsorTask(task)) {
				count++
			}
		}
	})
	return count
}

const monthGroups = computed(() => {
	const groups = props.timelineData.reduce(
		(groups, date) => {
			const month = date.getMonth()
			const year = date.getFullYear()
			const key = `${year}-${month}`

			const lastGroup = groups[groups.length - 1]
			if (lastGroup?.key === key) {
				lastGroup.width += props.dayWidthPixels
			} else {
				const sponsors = countSponsorsInMonth(year, month)
				const publishDays = countPublishDays(year, month)
				groups.push({
					key,
					label: dayjs(date).format('MMMM YYYY'),
					width: props.dayWidthPixels,
					sponsorInfo: props.tasks ? `[${sponsors} / ${publishDays}]` : '',
				})
			}

			return groups
		},
		[] as Array<{key: string; label: string; width: number; sponsorInfo: string}>,
	)

	return groups
})
</script>

<style scoped lang="scss">
.gantt-timeline {
	background: var(--white);
	border-block-end: 1px solid var(--grey-200);
	position: sticky;
	inset-block-start: 0;
	z-index: 10;
}

.gantt-timeline-months {
	display: flex;

	.timeunit-month {
		background: var(--white);
		font-family: $vikunja-font;
		font-weight: bold;
		border-inline-end: 1px solid var(--grey-200);
		padding: 0.5rem 0;
		text-align: center;
		font-size: 1rem;
		color: var(--grey-800);
	}
}

.gantt-timeline-days {
	display: flex;

	.timeunit {
		.timeunit-wrapper {
			padding: 0.5rem 0;
			font-size: 1rem;
			display: flex;
			flex-direction: column;
			align-items: center;
			inline-size: 100%;
			font-family: $vikunja-font;

			&.today {
				background: var(--primary);
				color: var(--white);
				border-radius: 5px 5px 0 0;
				font-weight: bold;
			}

			&.special-day:not(.today) {
				background: #df4a33;
				color: var(--white);
				border-radius: 5px 5px 0 0;
			}

			.weekday {
				font-size: 0.8rem;
			}
		}
	}
}
</style>
