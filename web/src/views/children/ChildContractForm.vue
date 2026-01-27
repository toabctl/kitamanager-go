<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import type { Child, ChildContractCreateRequest } from '@/api/types'
import Dialog from 'primevue/dialog'
import DatePicker from 'primevue/datepicker'
import Button from 'primevue/button'
import Chips from 'primevue/chips'

const props = defineProps<{
  visible: boolean
  child: Child | null
}>()

const emit = defineEmits<{
  close: []
  save: [data: ChildContractCreateRequest]
}>()

const form = ref({
  from: null as Date | null,
  to: null as Date | null,
  attributes: [] as string[]
})

const errors = ref<{
  from?: string
}>({})

const dialogTitle = computed(() =>
  props.child
    ? `New Contract for ${props.child.first_name} ${props.child.last_name}`
    : 'New Contract'
)

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      form.value = {
        from: new Date(),
        to: null,
        attributes: []
      }
      errors.value = {}
    }
  }
)

function validate(): boolean {
  errors.value = {}

  if (!form.value.from) {
    errors.value.from = 'Start date is required'
  }

  return Object.keys(errors.value).length === 0
}

function handleSave() {
  if (validate()) {
    emit('save', {
      from: form.value.from!.toISOString(),
      to: form.value.to ? form.value.to.toISOString() : null,
      attributes: form.value.attributes
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
        <label for="from">Start Date</label>
        <DatePicker
          id="from"
          v-model="form.from"
          date-format="dd.mm.yy"
          :class="{ 'p-invalid': errors.from }"
          placeholder="Contract start date"
          show-icon
        />
        <small v-if="errors.from" class="p-error">{{ errors.from }}</small>
      </div>

      <div class="field">
        <label for="to">End Date (optional)</label>
        <DatePicker
          id="to"
          v-model="form.to"
          date-format="dd.mm.yy"
          placeholder="Contract end date"
          show-icon
        />
      </div>

      <div class="field">
        <label for="attributes">Attributes (care type & extras)</label>
        <Chips
          id="attributes"
          v-model="form.attributes"
          placeholder="e.g. ganztags, ndh, integration_a"
        />
        <small class="text-secondary">
          Press Enter to add each attribute (e.g., ganztags, halbtags, teilzeit, ndh, integration_a)
        </small>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <Button label="Cancel" text @click="$emit('close')" />
        <Button label="Save" @click="handleSave" />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.text-secondary {
  color: var(--text-color-secondary);
}
</style>
