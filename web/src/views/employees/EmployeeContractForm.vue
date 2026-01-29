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
  grade: '',
  step: 1,
  weekly_hours: 39
})

const errors = ref<{
  from?: string
  position?: string
  grade?: string
  step?: string
  weekly_hours?: string
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
        grade: '',
        step: 1,
        weekly_hours: 39
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

  return Object.keys(errors.value).length === 0
}

function handleSave() {
  if (validate()) {
    emit('save', {
      from: form.value.from!.toISOString(),
      to: form.value.to ? form.value.to.toISOString() : null,
      position: form.value.position,
      grade: form.value.grade,
      step: form.value.step,
      weekly_hours: form.value.weekly_hours
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
        <label for="grade">{{ t('employees.grade') }}</label>
        <InputText
          id="grade"
          v-model="form.grade"
          :class="{ 'p-invalid': errors.grade }"
          placeholder="e.g., S8a"
        />
        <small v-if="errors.grade" class="p-error">{{ errors.grade }}</small>
      </div>

      <div class="field">
        <label for="step">{{ t('employees.step') }}</label>
        <InputNumber
          id="step"
          v-model="form.step"
          :class="{ 'p-invalid': errors.step }"
          :min="1"
          :max="6"
        />
        <small v-if="errors.step" class="p-error">{{ errors.step }}</small>
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
    </div>

    <template #footer>
      <div class="dialog-footer">
        <Button :label="t('common.cancel')" text @click="$emit('close')" />
        <Button :label="t('common.save')" @click="handleSave" />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.form-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field label {
  font-weight: 600;
}

.field :deep(input),
.field :deep(.p-inputnumber),
.field :deep(.p-datepicker) {
  width: 100%;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.p-error {
  color: var(--red-500);
}
</style>
