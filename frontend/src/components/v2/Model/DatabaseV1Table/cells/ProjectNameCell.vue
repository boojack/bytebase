<template>
  <div class="flex flex-row space-x-2 items-center">
    <ProjectV1Name
      :project="project"
      :link="link"
      tag="div"
      :keyword="keyword"
    />

    <slot name="suffix">
      {{ suffix }}
    </slot>

    <NTooltip v-if="project.workflow === Workflow.VCS">
      <template #trigger>
        <GitIcon class="w-4 h-4 text-control" />
      </template>
      <span v-if="mode === 'ALL_SHORT'" class="w-40">
        {{ $t("alter-schema.vcs-info") }}
      </span>
      <span v-else class="tooltip whitespace-nowrap">
        {{ $t("database.gitops-enabled") }}
      </span>
    </NTooltip>

    <NTooltip v-if="project.state === State.DELETED">
      <template #trigger>
        <heroicons-outline:archive class="w-4 h-4 text-control" />
      </template>
      <span class="whitespace-nowrap">
        {{ $t("common.archived") }}
      </span>
    </NTooltip>
  </div>
</template>

<script setup lang="ts">
import { NTooltip } from "naive-ui";
import GitIcon from "@/components/GitIcon.vue";
import { ProjectV1Name } from "@/components/v2";
import { State } from "@/types/proto/v1/common";
import type { Project } from "@/types/proto/v1/project_service";
import { Workflow } from "@/types/proto/v1/project_service";
import type { Mode } from "../DatabaseV1Table.vue";

withDefaults(
  defineProps<{
    project: Project;
    mode?: Mode;
    link?: boolean;
    keyword?: string;
    suffix?: string;
  }>(),
  {
    mode: "ALL",
    link: false,
    keyword: undefined,
    suffix: "",
  }
);
</script>
