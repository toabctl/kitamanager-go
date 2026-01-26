<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useUiStore } from '@/stores/ui'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const uiStore = useUiStore()

function closeSidebarOnMobile() {
  // Only close on mobile (when sidebar is used as overlay)
  if (window.innerWidth < 768) {
    uiStore.sidebarCollapsed = true
  }
}
</script>

<template>
  <div class="app-layout">
    <!-- Sidebar overlay for mobile -->
    <div
      v-if="!uiStore.sidebarCollapsed"
      class="sidebar-overlay md:hidden"
      @click="uiStore.toggleSidebar"
    ></div>

    <AppSidebar
      :class="{ 'sidebar-open': !uiStore.sidebarCollapsed }"
      @navigate="closeSidebarOnMobile"
    />
    <div class="app-main">
      <AppHeader />
      <main class="app-content">
        <RouterView />
      </main>
    </div>
  </div>
</template>
