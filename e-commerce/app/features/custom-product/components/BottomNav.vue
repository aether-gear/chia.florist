<script setup lang="ts">
import { TOOL_TABS } from '../constants'

const props = defineProps<{
  design: ReturnType<typeof import('../useCustomDesign').useCustomDesign>
}>()

function handleTabClick(tabId: string) {
  // Toggle: clicking the active tab closes the panel
  if (props.design.activeTab === tabId) {
    props.design.activeTab = null as any
  } else {
    props.design.activeTab = tabId as any
  }
}
</script>

<template>
  <!-- Floating centered pill at the bottom of the canvas area -->
  <nav class="zzz-bottom-bar" role="tablist" aria-label="Tool selection">
    <button
      v-for="tab in TOOL_TABS"
      :key="tab.id"
      class="zzz-tab-btn"
      :class="{ 'tab-active': design.activeTab === tab.id }"
      @click="handleTabClick(tab.id)"
      role="tab"
      :aria-selected="design.activeTab === tab.id"
      :id="'bottom-tab-' + tab.id"
    >
      <!-- ZZZ signature red indicator on top of active tab -->
      <span v-if="design.activeTab === tab.id" class="zzz-indicator"></span>

      <div class="zzz-tab-icon">
        <svg v-if="tab.id === 'text'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M4 7V4h16v3"/><path d="M9 20h6"/><path d="M12 4v16"/>
        </svg>
        <svg v-else-if="tab.id === 'image'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/>
        </svg>
        <svg v-else-if="tab.id === 'brush'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M18.37 2.63 14 7l-1.59-1.59a2 2 0 0 0-2.82 0L8 7l9 9 1.59-1.59a2 2 0 0 0 0-2.82L17 10l4.37-4.37a2.12 2.12 0 1 0-3-3Z"/>
          <path d="M9 8c-2 2.5-2 5-2 5"/>
        </svg>
        <svg v-else-if="tab.id === 'border'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="2" y="2" width="20" height="20" rx="2"/>
        </svg>
        <svg v-else-if="tab.id === 'corner'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round">
          <path d="M3 9V5a2 2 0 012-2h4"/><path d="M3 15v4a2 2 0 002 2h4"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 22C12 22 20 18 20 12C20 6 12 2 12 2C12 2 4 6 4 12C4 18 12 22 12 22Z"/><circle cx="12" cy="12" r="3"/>
        </svg>
      </div>
      <span class="zzz-tab-label">{{ tab.label }}</span>
    </button>
  </nav>
</template>
