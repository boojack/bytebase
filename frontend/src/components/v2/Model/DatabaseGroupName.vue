<template>
  <component
    :is="link ? 'router-link' : tag"
    v-bind="bindings"
    class="inline-flex items-center gap-x-1"
    :class="link && !plain && 'normal-link'"
  >
    <span>{{ databaseGroup.title }}</span>
  </component>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import { PROJECT_V1_ROUTE_DATABASE_GROUP_DETAIL } from "@/router/dashboard/projectV1";
import { getProjectNameAndDatabaseGroupName } from "@/store/modules/v1/common";
import type { DatabaseGroup } from "@/types/proto-es/v1/database_group_service_pb";

const props = withDefaults(
  defineProps<{
    databaseGroup: DatabaseGroup;
    tag?: string;
    link?: boolean;
    plain?: boolean;
  }>(),
  {
    tag: "span",
    link: true,
    plain: false,
  }
);

const bindings = computed(() => {
  if (!props.databaseGroup) {
    return {};
  }

  const [projectId, databaseGroupName] = getProjectNameAndDatabaseGroupName(
    props.databaseGroup.name
  );
  if (props.link) {
    return {
      to: {
        name: PROJECT_V1_ROUTE_DATABASE_GROUP_DETAIL,
        params: {
          projectId,
          databaseGroupName,
        },
      },
      activeClass: "",
      exactActiveClass: "",
      onClick: (e: MouseEvent) => {
        e.stopPropagation();
      },
    };
  }
  return {};
});
</script>
