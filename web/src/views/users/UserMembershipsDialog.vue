<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { apiClient } from '@/api/client'
import { useUiStore } from '@/stores/ui'
import type { User, Group, Organization, UserMembership, Role } from '@/api/types'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dropdown from 'primevue/dropdown'
import Tag from 'primevue/tag'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import PickList from 'primevue/picklist'

const props = defineProps<{
  visible: boolean
  user: User | null
}>()

const emit = defineEmits<{
  close: []
  updated: []
}>()

const toast = useToast()
const confirm = useConfirm()
const uiStore = useUiStore()
const loading = ref(false)

// Memberships data
const memberships = ref<UserMembership[]>([])

// Organizations data: [available, assigned]
const orgsData = ref<[Organization[], Organization[]]>([[], []])

// Add to group dialog state
const addGroupDialogVisible = ref(false)
const selectedOrgForAdd = ref<Organization | null>(null)
const availableGroupsForOrg = ref<Group[]>([])
const selectedGroupForAdd = ref<Group | null>(null)
const selectedRoleForAdd = ref<Role>('member')

// Edit role dialog state
const editRoleDialogVisible = ref(false)
const editingMembership = ref<UserMembership | null>(null)
const editingRole = ref<Role>('member')

const roleOptions = [
  { label: 'Admin', value: 'admin' as Role },
  { label: 'Manager', value: 'manager' as Role },
  { label: 'Member', value: 'member' as Role }
]

const dialogTitle = computed(() =>
  props.user ? `Manage Memberships: ${props.user.name}` : 'Manage Memberships'
)

watch(
  () => props.visible,
  async (visible) => {
    if (visible && props.user) {
      await loadData()
    }
  }
)

async function loadData() {
  if (!props.user) return

  loading.value = true
  try {
    const [membershipsResponse, allOrgs, userData] = await Promise.all([
      apiClient.getUserMemberships(props.user.id),
      apiClient.getOrganizations(),
      apiClient.getUser(props.user.id)
    ])

    memberships.value = membershipsResponse.memberships || []

    // Set up organizations picklist
    const assignedOrgIds = new Set((userData.organizations || []).map((o) => o.id))
    orgsData.value = [
      allOrgs.filter((o) => !assignedOrgIds.has(o.id)),
      userData.organizations || []
    ]
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load membership data',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

function getRoleSeverity(role: Role): 'success' | 'info' | 'secondary' {
  switch (role) {
    case 'admin':
      return 'success'
    case 'manager':
      return 'info'
    default:
      return 'secondary'
  }
}

function getRoleLabel(role: Role): string {
  return role.charAt(0).toUpperCase() + role.slice(1)
}

// Add to group functions
function openAddGroupDialog() {
  selectedOrgForAdd.value = null
  availableGroupsForOrg.value = []
  selectedGroupForAdd.value = null
  selectedRoleForAdd.value = 'member'
  addGroupDialogVisible.value = true
}

async function onOrgSelectedForAdd(org: Organization | null) {
  if (!org) {
    availableGroupsForOrg.value = []
    selectedGroupForAdd.value = null
    return
  }

  try {
    const groups = await apiClient.getGroups(org.id)
    // Filter out groups the user is already a member of
    const memberGroupIds = new Set(memberships.value.map((m) => m.group_id))
    availableGroupsForOrg.value = groups.filter((g) => !memberGroupIds.has(g.id))
    selectedGroupForAdd.value = null
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load groups',
      life: 3000
    })
  }
}

async function handleAddToGroup() {
  if (!props.user || !selectedGroupForAdd.value) return

  loading.value = true
  try {
    await apiClient.addUserToGroup(
      props.user.id,
      selectedGroupForAdd.value.id,
      selectedRoleForAdd.value
    )
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'User added to group',
      life: 3000
    })
    addGroupDialogVisible.value = false
    await loadData()
    emit('updated')
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to add user to group',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

// Edit role functions
function openEditRoleDialog(membership: UserMembership) {
  editingMembership.value = membership
  editingRole.value = membership.role
  editRoleDialogVisible.value = true
}

async function handleUpdateRole() {
  if (!props.user || !editingMembership.value) return

  loading.value = true
  try {
    await apiClient.updateUserGroupRole(
      props.user.id,
      editingMembership.value.group_id,
      editingRole.value
    )
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Role updated',
      life: 3000
    })
    editRoleDialogVisible.value = false
    await loadData()
    emit('updated')
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update role',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

// Remove from group
function confirmRemoveFromGroup(membership: UserMembership) {
  confirm.require({
    message: `Remove this user from group "${membership.group.name}"?`,
    header: 'Confirm Removal',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      if (!props.user) return

      loading.value = true
      try {
        await apiClient.removeUserFromGroup(props.user.id, membership.group_id)
        toast.add({
          severity: 'success',
          summary: 'Success',
          detail: 'User removed from group',
          life: 3000
        })
        await loadData()
        emit('updated')
      } catch {
        toast.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to remove user from group',
          life: 3000
        })
      } finally {
        loading.value = false
      }
    }
  })
}

// Organization assignments save
async function handleSaveOrganizations() {
  if (!props.user) return

  loading.value = true
  try {
    // Get current assignments from API
    const userData = await apiClient.getUser(props.user.id)
    const currentOrgIds = new Set((userData.organizations || []).map((o) => o.id))

    // New assignments from PickList
    const newOrgIds = new Set(orgsData.value[1].map((o) => o.id))

    // Calculate additions and removals
    const orgsToAdd = [...newOrgIds].filter((id) => !currentOrgIds.has(id))
    const orgsToRemove = [...currentOrgIds].filter((id) => !newOrgIds.has(id))

    // Apply changes
    await Promise.all([
      ...orgsToAdd.map((oid) => apiClient.addUserToOrganization(props.user!.id, oid)),
      ...orgsToRemove.map((oid) => apiClient.removeUserFromOrganization(props.user!.id, oid))
    ])

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Organizations updated',
      life: 3000
    })

    // Refresh organizations in the UI store
    await uiStore.fetchOrganizations()
    await loadData()
    emit('updated')
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update organizations',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Dialog
    :visible="visible"
    :header="dialogTitle"
    modal
    :closable="true"
    :style="{ width: '800px' }"
    @update:visible="$emit('close')"
  >
    <TabView>
      <TabPanel value="groups" header="Group Memberships">
        <div class="mb-3 flex justify-between">
          <p>Groups this user belongs to and their roles:</p>
          <Button label="Add to Group" icon="pi pi-plus" size="small" @click="openAddGroupDialog" />
        </div>

        <DataTable
          :value="memberships"
          :loading="loading"
          striped-rows
          :paginator="memberships.length > 5"
          :rows="5"
        >
          <Column header="Organization" style="width: 25%">
            <template #body="{ data }">
              {{ data.group.organization?.name || 'Unknown' }}
            </template>
          </Column>
          <Column field="group.name" header="Group" style="width: 25%"></Column>
          <Column header="Role" style="width: 20%">
            <template #body="{ data }">
              <Tag :value="getRoleLabel(data.role)" :severity="getRoleSeverity(data.role)" />
            </template>
          </Column>
          <Column header="Effective Org Role" style="width: 20%">
            <template #body="{ data }">
              <Tag
                :value="getRoleLabel(data.effective_org_role)"
                :severity="getRoleSeverity(data.effective_org_role)"
              />
            </template>
          </Column>
          <Column header="Actions" style="width: 10%">
            <template #body="{ data }">
              <Button
                icon="pi pi-pencil"
                text
                rounded
                size="small"
                title="Edit Role"
                @click="openEditRoleDialog(data)"
              />
              <Button
                icon="pi pi-trash"
                text
                rounded
                size="small"
                severity="danger"
                title="Remove"
                @click="confirmRemoveFromGroup(data)"
              />
            </template>
          </Column>
          <template #empty>
            <div class="text-center text-muted py-4">User is not a member of any groups</div>
          </template>
        </DataTable>
      </TabPanel>

      <TabPanel value="organizations" header="Organizations">
        <p class="mb-3">Organizations this user has access to:</p>
        <PickList
          v-model="orgsData"
          data-key="id"
          breakpoint="575px"
          :show-source-controls="false"
          :show-target-controls="false"
        >
          <template #sourceheader>Available Organizations</template>
          <template #targetheader>Assigned Organizations</template>
          <template #item="{ item }">
            <span>{{ item.name }}</span>
          </template>
        </PickList>
        <div class="mt-3 flex justify-end">
          <Button label="Save Organizations" :loading="loading" @click="handleSaveOrganizations" />
        </div>
      </TabPanel>
    </TabView>

    <template #footer>
      <div class="dialog-footer">
        <Button label="Close" text @click="$emit('close')" />
      </div>
    </template>
  </Dialog>

  <!-- Add to Group Dialog -->
  <Dialog
    :visible="addGroupDialogVisible"
    header="Add User to Group"
    modal
    :closable="true"
    :style="{ width: '450px' }"
    @update:visible="addGroupDialogVisible = false"
  >
    <div class="form-grid">
      <div class="field">
        <label for="add-org">Organization</label>
        <Dropdown
          id="add-org"
          v-model="selectedOrgForAdd"
          :options="orgsData[1]"
          option-label="name"
          placeholder="Select Organization"
          class="w-full"
          @change="onOrgSelectedForAdd(selectedOrgForAdd)"
        />
      </div>

      <div class="field">
        <label for="add-group">Group</label>
        <Dropdown
          id="add-group"
          v-model="selectedGroupForAdd"
          :options="availableGroupsForOrg"
          option-label="name"
          placeholder="Select Group"
          :disabled="!selectedOrgForAdd"
          class="w-full"
        />
      </div>

      <div class="field">
        <label for="add-role">Role</label>
        <Dropdown
          id="add-role"
          v-model="selectedRoleForAdd"
          :options="roleOptions"
          option-label="label"
          option-value="value"
          placeholder="Select Role"
          class="w-full"
        />
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <Button label="Cancel" text @click="addGroupDialogVisible = false" />
        <Button
          label="Add"
          :loading="loading"
          :disabled="!selectedGroupForAdd"
          @click="handleAddToGroup"
        />
      </div>
    </template>
  </Dialog>

  <!-- Edit Role Dialog -->
  <Dialog
    :visible="editRoleDialogVisible"
    header="Edit Role"
    modal
    :closable="true"
    :style="{ width: '400px' }"
    @update:visible="editRoleDialogVisible = false"
  >
    <div class="form-grid">
      <p v-if="editingMembership">
        Change role for <strong>{{ editingMembership.group.name }}</strong
        >:
      </p>
      <div class="field">
        <label for="edit-role">Role</label>
        <Dropdown
          id="edit-role"
          v-model="editingRole"
          :options="roleOptions"
          option-label="label"
          option-value="value"
          class="w-full"
        />
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <Button label="Cancel" text @click="editRoleDialogVisible = false" />
        <Button label="Save" :loading="loading" @click="handleUpdateRole" />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.mb-3 {
  margin-bottom: 1rem;
}

.mt-3 {
  margin-top: 1rem;
}

.py-4 {
  padding-top: 1.5rem;
  padding-bottom: 1.5rem;
}

.flex {
  display: flex;
}

.justify-between {
  justify-content: space-between;
  align-items: center;
}

.justify-end {
  justify-content: flex-end;
}

.text-center {
  text-align: center;
}

.text-muted {
  color: var(--text-color-secondary);
}

.w-full {
  width: 100%;
}
</style>
