<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useCrud } from '@/composables/useCrud'
import { useToast } from 'primevue/usetoast'
import { useUiStore } from '@/stores/ui'
import { useI18n } from 'vue-i18n'
import { apiClient } from '@/api/client'
import type {
  Organization,
  OrganizationCreateRequest,
  OrganizationUpdateRequest,
  GovernmentFunding
} from '@/api/types'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Dropdown from 'primevue/dropdown'
import OrganizationForm from './OrganizationForm.vue'

const toast = useToast()
const uiStore = useUiStore()
const { t } = useI18n()

const {
  items: organizations,
  loading,
  dialogVisible,
  editingItem,
  fetchItems,
  openCreateDialog,
  openEditDialog,
  closeDialog,
  saveItem: crudSaveItem,
  confirmDelete
} = useCrud<Organization, OrganizationCreateRequest, OrganizationUpdateRequest>({
  entityName: 'Organization',
  fetchAll: () => apiClient.getOrganizations(),
  create: (data) => apiClient.createOrganization(data),
  update: (id, data) => apiClient.updateOrganization(id, data),
  remove: (id) => apiClient.deleteOrganization(id)
})

// Wrap saveItem to also refresh the sidebar's org list
async function saveItem(data: OrganizationCreateRequest | OrganizationUpdateRequest) {
  await crudSaveItem(data)
  // Refresh sidebar organization dropdown
  await uiStore.fetchOrganizations()
}

// GovernmentFunding assignment
const governmentFundings = ref<GovernmentFunding[]>([])
const governmentFundingDialogVisible = ref(false)
const selectedOrg = ref<Organization | null>(null)
const selectedGovernmentFundingId = ref<number | null>(null)

async function fetchGovernmentFundings() {
  try {
    governmentFundings.value = await apiClient.getGovernmentFundings()
  } catch {
    // GovernmentFundings might not be available to non-superadmins
  }
}

function openGovernmentFundingDialog(org: Organization) {
  selectedOrg.value = org
  selectedGovernmentFundingId.value = org.government_funding_id || null
  governmentFundingDialogVisible.value = true
}

async function saveGovernmentFundingAssignment() {
  if (!selectedOrg.value) return

  try {
    if (selectedGovernmentFundingId.value) {
      await apiClient.assignGovernmentFundingToOrganization(
        selectedOrg.value.id,
        selectedGovernmentFundingId.value
      )
      toast.add({
        severity: 'success',
        summary: t('common.success'),
        detail: t('governmentFundings.updateSuccess'),
        life: 3000
      })
    } else {
      await apiClient.removeGovernmentFundingFromOrganization(selectedOrg.value.id)
      toast.add({
        severity: 'success',
        summary: t('common.success'),
        detail: t('governmentFundings.deleteSuccess'),
        life: 3000
      })
    }
    governmentFundingDialogVisible.value = false
    await fetchItems()
  } catch {
    toast.add({
      severity: 'error',
      summary: t('common.error'),
      detail: t('common.failedToSave', { resource: t('governmentFundings.title') }),
      life: 3000
    })
  }
}

function getGovernmentFundingName(org: Organization): string {
  if (org.government_funding) return org.government_funding.name
  const plan = governmentFundings.value.find((p) => p.id === org.government_funding_id)
  return plan?.name || '-'
}

onMounted(() => {
  fetchItems()
  fetchGovernmentFundings()
})
</script>

<template>
  <div>
    <div class="page-header">
      <h1>{{ t('organizations.title') }}</h1>
      <Button
        :label="t('organizations.newOrganization')"
        icon="pi pi-plus"
        @click="openCreateDialog"
      />
    </div>

    <div class="card">
      <DataTable
        :value="organizations"
        :loading="loading"
        striped-rows
        paginator
        :rows="10"
        :rows-per-page-options="[10, 25, 50]"
      >
        <Column field="id" :header="t('common.id')" sortable style="width: 80px"></Column>
        <Column field="name" :header="t('common.name')" sortable></Column>
        <Column :header="t('governmentFundings.title')" style="width: 150px">
          <template #body="{ data }">
            <span>{{ getGovernmentFundingName(data) }}</span>
          </template>
        </Column>
        <Column field="active" :header="t('common.status')" sortable style="width: 120px">
          <template #body="{ data }">
            <Tag
              :value="data.active ? t('common.active') : t('common.inactive')"
              :severity="data.active ? 'success' : 'danger'"
            />
          </template>
        </Column>
        <Column field="created_at" :header="t('common.created')" sortable style="width: 180px">
          <template #body="{ data }">
            {{ new Date(data.created_at).toLocaleDateString() }}
          </template>
        </Column>
        <Column :header="t('common.actions')" style="width: 180px">
          <template #body="{ data }">
            <Button
              icon="pi pi-money-bill"
              text
              rounded
              :title="t('governmentFundings.title')"
              :aria-label="t('governmentFundings.title')"
              @click="openGovernmentFundingDialog(data)"
            />
            <Button
              icon="pi pi-pencil"
              text
              rounded
              :title="t('common.edit')"
              :aria-label="t('common.edit')"
              @click="openEditDialog(data)"
            />
            <Button
              icon="pi pi-trash"
              text
              rounded
              severity="danger"
              :title="t('common.delete')"
              :aria-label="t('common.delete')"
              @click="confirmDelete(data)"
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <OrganizationForm
      :visible="dialogVisible"
      :organization="editingItem"
      @close="closeDialog"
      @save="saveItem"
    />

    <!-- GovernmentFunding Assignment Dialog -->
    <Dialog
      v-model:visible="governmentFundingDialogVisible"
      :header="t('governmentFundings.title')"
      modal
      :style="{ width: '400px' }"
    >
      <div class="form-grid">
        <div class="field">
          <span class="field-label">{{ t('organizations.title') }}</span>
          <p>{{ selectedOrg?.name }}</p>
        </div>
        <div class="field">
          <label for="government-funding">{{ t('governmentFundings.title') }}</label>
          <Dropdown
            id="government-funding"
            v-model="selectedGovernmentFundingId"
            :options="governmentFundings"
            option-label="name"
            option-value="id"
            :placeholder="t('organizations.selectOrg')"
            :show-clear="true"
            class="w-full"
          />
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" text @click="governmentFundingDialogVisible = false" />
        <Button :label="t('common.save')" @click="saveGovernmentFundingAssignment" />
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.w-full {
  width: 100%;
}
</style>
