<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { GovernmentFunding, GovernmentFundingCreateRequest } from '@/api/types'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Button from 'primevue/button'

const { t } = useI18n()

const props = defineProps<{
  visible: boolean
  governmentFunding: GovernmentFunding | null
}>()

const emit = defineEmits<{
  close: []
  save: [data: GovernmentFundingCreateRequest]
}>()

const stateOptions = [{ value: 'berlin', label: 'Berlin' }]

const form = ref({
  name: '',
  state: ''
})

const errors = ref<{ name?: string; state?: string }>({})

const isEditing = computed(() => !!props.governmentFunding)
const dialogTitle = computed(() =>
  isEditing.value ? t('governmentFundings.edit') : t('governmentFundings.newGovernmentFunding')
)

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      if (props.governmentFunding) {
        form.value = {
          name: props.governmentFunding.name,
          state: props.governmentFunding.state
        }
      } else {
        form.value = {
          name: '',
          state: ''
        }
      }
      errors.value = {}
    }
  }
)

function validate(): boolean {
  errors.value = {}
  if (!form.value.name.trim()) {
    errors.value.name = t('validation.nameRequired')
  }
  if (!form.value.state) {
    errors.value.state = t('validation.required')
  }
  return Object.keys(errors.value).length === 0
}

function handleSave() {
  if (validate()) {
    emit('save', form.value)
  }
}
</script>

<template>
  <Dialog
    :visible="visible"
    :header="dialogTitle"
    modal
    :closable="true"
    :style="{ width: '450px' }"
    @update:visible="$emit('close')"
  >
    <div class="form-grid">
      <div class="field">
        <label for="name">{{ t('common.name') }}</label>
        <InputText
          id="name"
          v-model="form.name"
          :class="{ 'p-invalid': errors.name }"
          placeholder="e.g. Berlin Kita-Förderung"
        />
        <small v-if="errors.name" class="p-error">{{ errors.name }}</small>
      </div>
      <div class="field">
        <label for="state">{{ t('states.state') }}</label>
        <Select
          id="state"
          v-model="form.state"
          :options="stateOptions"
          option-label="label"
          option-value="value"
          :class="{ 'p-invalid': errors.state }"
          :disabled="isEditing"
        />
        <small v-if="errors.state" class="p-error">{{ errors.state }}</small>
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
