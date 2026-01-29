import { ref, type Ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { useI18n } from 'vue-i18n'
import { getErrorMessage } from '@/api/client'

interface CrudConfig<T, CreateDto, UpdateDto> {
  entityName: string
  entityNameKey?: string // i18n key for entity name (e.g., 'children.title')
  fetchAll: () => Promise<T[]>
  create: (data: CreateDto) => Promise<T>
  update: (id: number, data: UpdateDto) => Promise<T>
  remove: (id: number) => Promise<void>
  getId?: (item: T) => number
}

export function useCrud<T, CreateDto, UpdateDto>(config: CrudConfig<T, CreateDto, UpdateDto>) {
  const toast = useToast()
  const confirm = useConfirm()
  const { t } = useI18n()

  const items: Ref<T[]> = ref([])
  const loading = ref(false)
  const dialogVisible = ref(false)
  const editingItem: Ref<T | null> = ref(null)

  const getId = config.getId || ((item: T) => (item as { id: number }).id)
  const entityLabel = config.entityNameKey ? t(config.entityNameKey) : config.entityName

  async function fetchItems() {
    loading.value = true
    try {
      items.value = await config.fetchAll()
    } catch (error) {
      toast.add({
        severity: 'error',
        summary: t('common.error'),
        detail: getErrorMessage(error, t('common.failedToLoad', { resource: entityLabel })),
        life: 5000
      })
    } finally {
      loading.value = false
    }
  }

  function openCreateDialog() {
    editingItem.value = null
    dialogVisible.value = true
  }

  function openEditDialog(item: T) {
    editingItem.value = item
    dialogVisible.value = true
  }

  function closeDialog() {
    dialogVisible.value = false
    editingItem.value = null
  }

  async function saveItem(data: CreateDto | UpdateDto) {
    try {
      if (editingItem.value) {
        await config.update(getId(editingItem.value), data as UpdateDto)
        toast.add({
          severity: 'success',
          summary: t('common.success'),
          detail: t('common.updateSuccess', { resource: entityLabel }),
          life: 3000
        })
      } else {
        await config.create(data as CreateDto)
        toast.add({
          severity: 'success',
          summary: t('common.success'),
          detail: t('common.createSuccess', { resource: entityLabel }),
          life: 3000
        })
      }
      closeDialog()
      await fetchItems()
    } catch (error) {
      toast.add({
        severity: 'error',
        summary: t('common.error'),
        detail: getErrorMessage(error, t('common.failedToSave', { resource: entityLabel })),
        life: 5000
      })
    }
  }

  function confirmDelete(item: T) {
    confirm.require({
      message: t('common.confirmDeleteMessage', { resource: entityLabel }),
      header: t('common.confirmDelete'),
      icon: 'pi pi-exclamation-triangle',
      acceptClass: 'p-button-danger',
      accept: async () => {
        try {
          await config.remove(getId(item))
          toast.add({
            severity: 'success',
            summary: t('common.success'),
            detail: t('common.deleteSuccess', { resource: entityLabel }),
            life: 3000
          })
          await fetchItems()
        } catch (error) {
          toast.add({
            severity: 'error',
            summary: t('common.error'),
            detail: getErrorMessage(error, t('common.failedToDelete', { resource: entityLabel })),
            life: 5000
          })
        }
      }
    })
  }

  return {
    items,
    loading,
    dialogVisible,
    editingItem,
    fetchItems,
    openCreateDialog,
    openEditDialog,
    closeDialog,
    saveItem,
    confirmDelete
  }
}
