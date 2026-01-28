<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Organization, OrganizationCreateRequest } from '@/api/types'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Checkbox from 'primevue/checkbox'
import Button from 'primevue/button'
import Select from 'primevue/select'

const { t } = useI18n()

const props = defineProps<{
  visible: boolean
  organization: Organization | null
}>()

const emit = defineEmits<{
  close: []
  save: [data: OrganizationCreateRequest]
}>()

const stateOptions = computed(() => [{ value: 'berlin', label: t('states.berlin') }])

const form = ref({
  name: '',
  active: true,
  state: ''
})

const errors = ref<{ name?: string; state?: string }>({})

const isEditing = computed(() => !!props.organization)
const dialogTitle = computed(() =>
  isEditing.value ? t('organizations.edit') : t('organizations.newOrganization')
)

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      if (props.organization) {
        form.value = {
          name: props.organization.name,
          active: props.organization.active,
          state: props.organization.state
        }
      } else {
        form.value = {
          name: '',
          active: true,
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
          :placeholder="t('organizations.name')"
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
          :placeholder="t('states.selectState')"
          :class="{ 'p-invalid': errors.state }"
          class="w-full"
        />
        <small v-if="errors.state" class="p-error">{{ errors.state }}</small>
      </div>

      <div class="field">
        <div class="flex align-items-center gap-2">
          <Checkbox v-model="form.active" input-id="active" :binary="true" />
          <label for="active">{{ t('common.active') }}</label>
        </div>
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
.flex {
  display: flex;
}

.align-items-center {
  align-items: center;
}

.gap-2 {
  gap: 0.5rem;
}

.w-full {
  width: 100%;
}
</style>
