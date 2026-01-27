<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { useI18n } from 'vue-i18n'
import { apiClient } from '@/api/client'
import type { Group, GroupCreateRequest, GroupUpdateRequest } from '@/api/types'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import GroupForm from './GroupForm.vue'

const route = useRoute()
const toast = useToast()
const confirm = useConfirm()
const { t } = useI18n()

const orgId = ref(Number(route.params.orgId))
const groups = ref<Group[]>([])
const loading = ref(false)

const dialogVisible = ref(false)
const editingGroup = ref<Group | null>(null)

watch(
  () => route.params.orgId,
  (newOrgId) => {
    orgId.value = Number(newOrgId)
    fetchGroups()
  }
)

async function fetchGroups() {
  loading.value = true
  try {
    groups.value = await apiClient.getGroups(orgId.value)
  } catch {
    toast.add({
      severity: 'error',
      summary: t('common.error'),
      detail: t('common.failedToLoad', { resource: t('groups.title') }),
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  editingGroup.value = null
  dialogVisible.value = true
}

function openEditDialog(group: Group) {
  editingGroup.value = group
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
  editingGroup.value = null
}

async function saveGroup(data: GroupCreateRequest | GroupUpdateRequest) {
  try {
    if (editingGroup.value) {
      await apiClient.updateGroup(orgId.value, editingGroup.value.id, data as GroupUpdateRequest)
      toast.add({
        severity: 'success',
        summary: t('common.success'),
        detail: t('groups.updateSuccess'),
        life: 3000
      })
    } else {
      await apiClient.createGroup(orgId.value, data as GroupCreateRequest)
      toast.add({
        severity: 'success',
        summary: t('common.success'),
        detail: t('groups.createSuccess'),
        life: 3000
      })
    }
    closeDialog()
    await fetchGroups()
  } catch {
    toast.add({
      severity: 'error',
      summary: t('common.error'),
      detail: t('common.failedToSave', { resource: t('groups.title') }),
      life: 3000
    })
  }
}

function confirmDelete(group: Group) {
  confirm.require({
    message: t('groups.deleteConfirm'),
    header: t('common.confirmDelete'),
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await apiClient.deleteGroup(orgId.value, group.id)
        toast.add({
          severity: 'success',
          summary: t('common.success'),
          detail: t('groups.deleteSuccess'),
          life: 3000
        })
        await fetchGroups()
      } catch {
        toast.add({
          severity: 'error',
          summary: t('common.error'),
          detail: t('common.failedToDelete', { resource: t('groups.title') }),
          life: 3000
        })
      }
    }
  })
}

onMounted(() => {
  fetchGroups()
})
</script>

<template>
  <div>
    <div class="page-header">
      <h1>{{ t('groups.title') }}</h1>
      <Button :label="t('groups.newGroup')" icon="pi pi-plus" @click="openCreateDialog" />
    </div>

    <div class="card">
      <DataTable
        :value="groups"
        :loading="loading"
        striped-rows
        paginator
        :rows="10"
        :rows-per-page-options="[10, 25, 50]"
      >
        <Column field="id" :header="t('common.id')" sortable style="width: 80px"></Column>
        <Column field="name" :header="t('common.name')" sortable></Column>
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
        <Column :header="t('common.actions')" style="width: 150px">
          <template #body="{ data }">
            <Button
              icon="pi pi-pencil"
              text
              rounded
              :title="t('common.edit')"
              @click="openEditDialog(data)"
            />
            <Button
              icon="pi pi-trash"
              text
              rounded
              severity="danger"
              :title="t('common.delete')"
              @click="confirmDelete(data)"
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <GroupForm
      :visible="dialogVisible"
      :group="editingGroup"
      @close="closeDialog"
      @save="saveGroup"
    />
  </div>
</template>
