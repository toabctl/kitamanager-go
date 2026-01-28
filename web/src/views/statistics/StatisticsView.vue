<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useUiStore } from '@/stores/ui'
import ChildrenMonthlyChart from '@/components/ChildrenMonthlyChart.vue'

const { t } = useI18n()
const uiStore = useUiStore()
</script>

<template>
  <div>
    <div class="page-header">
      <h1>{{ t('nav.statistics') }}</h1>
    </div>

    <div v-if="!uiStore.selectedOrganizationId" class="no-org-message">
      <i class="pi pi-info-circle"></i>
      {{ t('statistics.selectOrgForStats') }}
    </div>

    <div v-else class="charts-container">
      <ChildrenMonthlyChart :org-id="uiStore.selectedOrganizationId" />
    </div>
  </div>
</template>

<style scoped>
.no-org-message {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 3rem;
  background: var(--surface-card);
  border-radius: 8px;
  color: var(--text-color-secondary);
}

.charts-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
</style>
