<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Employee, EmployeeContractCreateRequest } from '@/api/types'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import DatePicker from 'primevue/datepicker'
import Button from 'primevue/button'

const { t } = useI18n()

const props = defineProps<{
  visible: boolean
  employee: Employee | null
}>()

const emit = defineEmits<{
  close: []
  save: [data: EmployeeContractCreateRequest]
}>()

const form = ref({
  from: null as Date | null,
  to: null as Date | null,
  position: '',
  weekly_hours: 40,
  salary: 0
})

const errors = ref<{
  from?: string
  position?: string
  weekly_hours?: string
  salary?: string
}>({})

const dialogTitle = computed(() =>
  props.employee
    ? t('contracts.newContractFor', {
        name: `${props.employee.first_name} ${props.employee.last_name}`
      })
    : t('contracts.newContract')
)

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      form.value = {
        from: new Date(),
        to: null,
        position: '',
        weekly_hours: 40,
        salary: 0
      }
      errors.value = {}
    }
  }
)

function validate(): boolean {
  errors.value = {}

  if (!form.value.from) {
    errors.value.from = t('contracts.startDateRequired')
  }

  if (!form.value.position.trim()) {
    errors.value.position = t('employees.positionRequired')
  }

  if (!form.value.weekly_hours || form.value.weekly_hours <= 0) {
    errors.value.weekly_hours = t('employees.weeklyHoursRequired')
  }

  if (!form.value.salary || form.value.salary <= 0) {
    errors.value.salary = t('employees.salaryRequired')
  }

  return Object.keys(errors.value).length === 0
}

function handleSave() {
  if (validate()) {
    emit('save', {
      from: form.value.from!.toISOString(),
      to: form.value.to ? form.value.to.toISOString() : null,
      position: form.value.position,
      weekly_hours: form.value.weekly_hours,
      salary: Math.round(form.value.salary * 100) // Convert to cents
    })
  }
}
</script>

<template>
  <Dialog
    :visible="visible"
    :header="dialogTitle"
    modal
    :closable="true"
    :style="{ width: '500px' }"
    @update:visible="$emit('close')"
  >
    <div class="form-grid">
      <div class="field">
        <label for="from">{{ t('contracts.startDate') }}</label>
        <DatePicker
          id="from"
          v-model="form.from"
          date-format="dd.mm.yy"
          :class="{ 'p-invalid': errors.from }"
          :placeholder="t('contracts.contractStartPlaceholder')"
          show-icon
        />
        <small v-if="errors.from" class="p-error">{{ errors.from }}</small>
      </div>

      <div class="field">
        <label for="to">{{ t('contracts.endDateOptional') }}</label>
        <DatePicker
          id="to"
          v-model="form.to"
          date-format="dd.mm.yy"
          :placeholder="t('contracts.contractEndPlaceholder')"
          show-icon
        />
      </div>

      <div class="field">
        <label for="position">{{ t('employees.position') }}</label>
        <InputText
          id="position"
          v-model="form.position"
          :class="{ 'p-invalid': errors.position }"
          placeholder="e.g., Erzieher"
        />
        <small v-if="errors.position" class="p-error">{{ errors.position }}</small>
      </div>

      <div class="field">
        <label for="weekly_hours">{{ t('employees.weeklyHours') }}</label>
        <InputNumber
          id="weekly_hours"
          v-model="form.weekly_hours"
          :class="{ 'p-invalid': errors.weekly_hours }"
          :min="0"
          :max="60"
          suffix=" h"
        />
        <small v-if="errors.weekly_hours" class="p-error">{{ errors.weekly_hours }}</small>
      </div>

      <div class="field">
        <label for="salary">{{ t('employees.monthlySalary') }}</label>
        <InputNumber
          id="salary"
          v-model="form.salary"
          :class="{ 'p-invalid': errors.salary }"
          mode="currency"
          currency="EUR"
          locale="de-DE"
        />
        <small v-if="errors.salary" class="p-error">{{ errors.salary }}</small>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <Button :label="t('common.cancel')" text @click="$emit('close')" />
        <Button :label="t('common.save')" @click="handleSave" />
      </div>
    </template>
  </Dialog>
</template>
