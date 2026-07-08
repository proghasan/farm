<script setup lang="ts">
import { useAuthStore } from './stores/auth'

const auth = useAuthStore()
if (auth.isAuthenticated) auth.fetchProfile()
</script>

<template>
  <RouterView />
</template>