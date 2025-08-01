<template>
  <NPopover v-if="shouldShow" trigger="hover" placement="bottom">
    <template #trigger>
      <NTag round>
        <span class="text-gray-400 text-sm mr-1">{{
          $t("task.check-type.affected-rows.self")
        }}</span>
        <span class="text-sm font-medium">
          {{ affectedRows }}
        </span>
      </NTag>
    </template>
    <div
      class="flex flex-col justify-start items-start w-72 max-h-128 overflow-auto"
    >
      <p class="text-sm text-gray-400">
        {{ $t("task.check-type.affected-rows.description") }}
      </p>
      <div class="w-full flex flex-col mt-1 gap-y-1">
        <div
          class="w-full flex flex-row justify-between items-center gap-1 group hover:opacity-90 cursor-pointer"
          v-for="task in affectedTasks"
          :key="task.name"
          @click="onClickTask(task)"
        >
          <span class="group-hover:underline truncate">
            {{ databaseForTask(projectOfIssue(issue), task).databaseName }}
          </span>
          <span
            class="shrink-0"
            :class="
              !affectedTaskMap.get(task.name) ? 'opacity-80' : 'font-medium'
            "
          >
            {{ affectedTaskMap.get(task.name) }}
          </span>
        </div>
      </div>
    </div>
  </NPopover>
</template>

<script setup lang="ts">
import { NTag, NPopover } from "naive-ui";
import { computed } from "vue";
import { useIssueContext, projectOfIssue } from "@/components/IssueV1/logic";
import { type PlanCheckRun } from "@/types/proto-es/v1/plan_service_pb";
import type { Task } from "@/types/proto-es/v1/rollout_service_pb";
import { databaseForTask } from "@/utils";
import { flattenTaskV1List } from "@/utils";

const props = defineProps<{
  taskSummaryReportMap: Map<string, PlanCheckRun>;
}>();

const { issue, events } = useIssueContext();

const tasks = computed(() => flattenTaskV1List(issue.value.rolloutEntity));

const affectedTaskMap = computed(() => {
  const tempMap = new Map<string, number>();
  props.taskSummaryReportMap.forEach((planCheckRun, task) => {
    if (
      planCheckRun.results.every(
        (result) =>
          result.report?.case !== "sqlSummaryReport" ||
          result.report.value.affectedRows === undefined
      )
    ) {
      return;
    }

    const count = planCheckRun.results.reduce((acc, result) => {
      if (result.report?.case === "sqlSummaryReport") {
        return acc + (result.report.value.affectedRows || 0);
      }
      return acc;
    }, 0);
    tempMap.set(task, count);
  });
  return new Map(Array.from(tempMap.entries()).sort((a, b) => b[1] - a[1]));
});

const affectedTasks = computed(
  () =>
    Array.from(affectedTaskMap.value.keys())
      .map((task) => tasks.value.find((t) => t.name === task))
      .filter((task) => task) as Task[]
);

const summaryReportResults = computed(() =>
  Array.from(props.taskSummaryReportMap.values()).flatMap(
    (report) => report.results
  )
);

const affectedRows = computed(() =>
  summaryReportResults.value.reduce((acc, result) => {
    if (result.report?.case === "sqlSummaryReport") {
      return acc + (result.report.value.affectedRows || 0);
    }
    return acc;
  }, 0)
);

const shouldShow = computed(() =>
  summaryReportResults.value.some(
    (result) =>
      result.report?.case === "sqlSummaryReport" &&
      result.report.value.affectedRows !== undefined
  )
);

const onClickTask = (task: Task) => {
  events.emit("select-task", { task });
};

defineExpose({ shouldShow });
</script>
