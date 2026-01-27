<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { GovernmentFunding, GovernmentFundingCreateRequest } from '@/api/types'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
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

const form = ref({
  name: ''
})

const errors = ref<{ name?: string }>({})

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
          name: props.governmentFunding.name
        }
      } else {
        form.value = {
          name: ''
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
          placeholder="e.g. Berlin"
        />
        <small v-if="errors.name" class="p-error">{{ errors.name }}</small>
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
