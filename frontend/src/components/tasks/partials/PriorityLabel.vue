<template>
	<span
		v-if="!done && (showAll || priority >= minimumPriority)"
		:class="{
			'not-so-high': priority > priorities.LOW && priority < priorities.HIGH,
			'high-priority': priority >= priorities.HIGH
		}"
		class="priority-label"
	>
		<span
			v-if="priority >= minimumPriority"
			class="icon"
		>
			<Icon
				v-if="priority === priorities.LOW"
				icon="lightbulb"
			/>
			<Icon
				v-if="priority === priorities.MEDIUM"
				icon="bolt"
			/>
			<Icon
				v-if="priority === priorities.HIGH"
				icon="exclamation-circle"
			/>
			<Icon
				v-if="priority === priorities.URGENT"
				icon="exclamation-triangle"
			/>
			<Icon
				v-if="priority === priorities.DO_NOW"
				icon="bomb"
			/>
		</span>
		<span>
			<template v-if="priority === priorities.UNSET">{{ $t('task.priority.unset') }}</template>
			<template v-if="priority === priorities.LOW">{{ $t('task.priority.low') }}</template>
			<template v-if="priority === priorities.MEDIUM">{{ $t('task.priority.medium') }}</template>
			<template v-if="priority === priorities.HIGH">{{ $t('task.priority.high') }}</template>
			<template v-if="priority === priorities.URGENT">{{ $t('task.priority.urgent') }}</template>
			<template v-if="priority === priorities.DO_NOW">{{ $t('task.priority.doNow') }}</template>
		</span>
	</span>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import {PRIORITIES as priorities} from '@/constants/priorities'
import {useAuthStore} from '@/stores/auth'
	
withDefaults(defineProps<{
	priority: number,
	showAll?: boolean,
	done?: boolean
}>(), {
	priority: priorities.UNSET,
	showAll: false,
	done: false,
})

const authStore = useAuthStore()

const minimumPriority = computed(() => {
	return authStore.settings.frontendSettings.minimumPriority
})
</script>

<style lang="scss" scoped>
.high-priority {
	color: var(--danger);
	width: auto !important; // To override the width set in tasks
}

.not-so-high {
	color: var(--warning);
}

.icon {
	vertical-align: top;
	width: auto !important;
	padding-right: .5rem;
}
</style>