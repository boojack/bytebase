<template>
  <NPerformantEllipsis
    class="w-full cursor-pointer leading-6 hover:text-accent"
  >
    <span @click="$emit('click')">
      {{ table.name }}
    </span>
  </NPerformantEllipsis>
</template>

<script lang="ts" setup>
import { NPerformantEllipsis } from "naive-ui";
import type { TableMetadata } from "@/types/proto-es/v1/database_service_pb";

defineProps<{
  table: TableMetadata;
}>();
defineEmits<{
  (event: "click"): void;
}>();
</script>
