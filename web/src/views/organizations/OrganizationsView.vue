<script setup lang="ts">
import { onMounted } from 'vue'
import { useCrud } from '@/composables/useCrud'
import { useUiStore } from '@/stores/ui'
import { useI18n } from 'vue-i18n'
import { apiClient } from '@/api/client'
import type {
  Organization,
  OrganizationCreateRequest,
  OrganizationUpdateRequest
} from '@/api/types'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import OrganizationForm from './OrganizationForm.vue'

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

onMounted(() => {
  fetchItems()
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
        <Column field="state" :header="t('states.state')" sortable style="width: 150px">
          <template #body="{ data }">
            {{ t(`states.${data.state}`) }}
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
        <Column :header="t('common.actions')" style="width: 120px">
          <template #body="{ data }">
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
  </div>
</template>
