<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Chart from 'primevue/chart'
import ChartDataLabels from 'chartjs-plugin-datalabels'
import { Chart as ChartJS } from 'chart.js'
import { apiClient } from '@/api/client'
import type { ChildrenContractCountByMonthResponse } from '@/api/types'

// Register the datalabels plugin
ChartJS.register(ChartDataLabels)

const props = defineProps<{
  orgId: number | null
}>()

const { t } = useI18n()

const loading = ref(false)
const error = ref<string | null>(null)
const statsData = ref<ChildrenContractCountByMonthResponse | null>(null)

const monthLabels = computed(() => [
  t('months.jan'),
  t('months.feb'),
  t('months.mar'),
  t('months.apr'),
  t('months.may'),
  t('months.jun'),
  t('months.jul'),
  t('months.aug'),
  t('months.sep'),
  t('months.oct'),
  t('months.nov'),
  t('months.dec')
])

const chartColors = [
  { border: 'rgb(59, 130, 246)', background: 'rgba(59, 130, 246, 0.7)' }, // blue
  { border: 'rgb(34, 197, 94)', background: 'rgba(34, 197, 94, 0.7)' }, // green
  { border: 'rgb(249, 115, 22)', background: 'rgba(249, 115, 22, 0.7)' }, // orange
  { border: 'rgb(168, 85, 247)', background: 'rgba(168, 85, 247, 0.7)' }, // purple
  { border: 'rgb(236, 72, 153)', background: 'rgba(236, 72, 153, 0.7)' }, // pink
  { border: 'rgb(20, 184, 166)', background: 'rgba(20, 184, 166, 0.7)' }, // teal
  { border: 'rgb(245, 158, 11)', background: 'rgba(245, 158, 11, 0.7)' }, // amber
  { border: 'rgb(239, 68, 68)', background: 'rgba(239, 68, 68, 0.7)' } // red
]

const chartData = computed(() => {
  if (!statsData.value) {
    return { labels: monthLabels.value, datasets: [] }
  }

  const datasets = statsData.value.years.map((yearData, index) => ({
    label: String(yearData.year),
    data: yearData.counts,
    borderColor: chartColors[index % chartColors.length].border,
    backgroundColor: chartColors[index % chartColors.length].background,
    borderWidth: 1
  }))

  return {
    labels: monthLabels.value,
    datasets
  }
})

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  layout: {
    padding: {
      top: 20
    }
  },
  plugins: {
    legend: {
      position: 'top' as const,
      labels: {
        usePointStyle: true,
        padding: 20
      }
    },
    datalabels: {
      anchor: 'end' as const,
      align: 'end' as const,
      color: 'var(--text-color)',
      font: {
        size: 10,
        weight: 'bold' as const
      },
      formatter: (value: number) => (value > 0 ? value : '')
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      grace: '15%',
      ticks: {
        stepSize: 1,
        precision: 0
      }
    }
  }
}))

async function loadStats() {
  if (!props.orgId) {
    statsData.value = null
    return
  }

  loading.value = true
  error.value = null

  try {
    statsData.value = await apiClient.getChildrenContractCountByMonth(props.orgId)
  } catch (err) {
    console.error('Failed to load children monthly stats:', err)
    error.value = t('statistics.chartError')
  } finally {
    loading.value = false
  }
}

watch(
  () => props.orgId,
  () => {
    loadStats()
  }
)

onMounted(() => {
  loadStats()
})
</script>

<template>
  <div class="children-monthly-chart">
    <h3>{{ t('statistics.childrenContractCount') }}</h3>

    <div v-if="loading" class="chart-loading">
      <i class="pi pi-spin pi-spinner"></i>
      {{ t('common.loading') }}
    </div>

    <div v-else-if="error" class="chart-error">
      <i class="pi pi-exclamation-triangle"></i>
      {{ error }}
    </div>

    <div v-else-if="!orgId" class="chart-placeholder">
      {{ t('statistics.selectOrgForStats') }}
    </div>

    <div v-else class="chart-container">
      <Chart type="bar" :data="chartData" :options="chartOptions" style="height: 450px" />
    </div>
  </div>
</template>

<style scoped>
.children-monthly-chart {
  background: var(--surface-card);
  border-radius: 8px;
  padding: 1.5rem;
  margin-top: 2rem;
}

.children-monthly-chart h3 {
  margin: 0 0 1rem;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-color);
}

.chart-container {
  height: 450px;
  min-height: 450px;
}

.chart-loading,
.chart-error,
.chart-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  height: 200px;
  color: var(--text-color-secondary);
}

.chart-error {
  color: var(--red-500);
}
</style>
