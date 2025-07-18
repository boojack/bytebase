<template>
  <NTag
    :class="[selected && 'shadow']"
    :type="tagType"
    round
    :bordered="selected"
    @click="clickable && $emit('click')"
  >
    <template #icon>
      <template v-if="status === PlanCheckRun_Status.RUNNING">
        <TaskSpinner class="h-4 w-4 text-info" />
      </template>
      <template v-else-if="status === PlanCheckRun_Status.DONE">
        <CheckIcon
          v-if="resultStatus === PlanCheckRun_Result_Status.SUCCESS"
          class="text-success"
          :size="16"
        />
        <TriangleAlertIcon
          v-else-if="resultStatus === PlanCheckRun_Result_Status.WARNING"
          :size="16"
        />
        <CircleAlertIcon
          v-else-if="resultStatus === PlanCheckRun_Result_Status.ERROR"
          :size="16"
        />
      </template>
      <template v-else-if="status === PlanCheckRun_Status.FAILED">
        <CircleAlertIcon :size="16" />
      </template>
      <template v-else-if="status === PlanCheckRun_Status.CANCELED">
        <TriangleAlertIcon :size="16" />
      </template>
    </template>
    <span>{{ title }}</span>
  </NTag>
</template>

<script setup lang="ts">
import { maxBy } from "lodash-es";
import { CheckIcon, TriangleAlertIcon, CircleAlertIcon } from "lucide-vue-next";
import { NTag } from "naive-ui";
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import { TaskSpinner } from "@/components/IssueV1/components/common";
import type { PlanCheckRun } from "@/types/proto-es/v1/plan_service_pb";
import {
  PlanCheckRun_Result_Status,
  PlanCheckRun_Status,
  PlanCheckRun_Type,
} from "@/types/proto-es/v1/plan_service_pb";
import { extractPlanCheckRunUID } from "@/utils";
import { planCheckRunResultStatus } from "./common";

const props = defineProps<{
  planCheckRuns: PlanCheckRun[];
  type: PlanCheckRun_Type;
  clickable?: boolean;
  selected?: boolean;
}>();

defineEmits<{
  (event: "click"): void;
}>();

const { t } = useI18n();

const latestPlanCheckRun = computed(() => {
  // Get the latest PlanCheckRun by UID.
  return maxBy(props.planCheckRuns, (check) =>
    Number(extractPlanCheckRunUID(check.name))
  )!;
});

const status = computed(() => {
  return latestPlanCheckRun.value.status;
});

const tagType = computed(() => {
  if (latestPlanCheckRun.value.status === PlanCheckRun_Status.FAILED) {
    return "error";
  }
  if (latestPlanCheckRun.value.status === PlanCheckRun_Status.CANCELED) {
    return "warning";
  }
  if (latestPlanCheckRun.value.status === PlanCheckRun_Status.RUNNING) {
    return "info";
  }
  if (latestPlanCheckRun.value.status !== PlanCheckRun_Status.DONE) {
    // Should not reach here.
    return "default";
  }

  switch (planCheckRunResultStatus(latestPlanCheckRun.value)) {
    case PlanCheckRun_Result_Status.SUCCESS:
      return "default";
    case PlanCheckRun_Result_Status.WARNING:
      return "warning";
    case PlanCheckRun_Result_Status.ERROR:
      return "error";
  }
  // Should not reach here.
  return "default";
});

const resultStatus = computed(() => {
  return planCheckRunResultStatus(latestPlanCheckRun.value);
});

const title = computed(() => {
  const { type } = latestPlanCheckRun.value;
  switch (type) {
    case PlanCheckRun_Type.DATABASE_STATEMENT_FAKE_ADVISE:
      return t("task.check-type.fake");
    case PlanCheckRun_Type.DATABASE_STATEMENT_ADVISE:
      return t("task.check-type.sql-review");
    case PlanCheckRun_Type.DATABASE_CONNECT:
      return t("task.check-type.connection");
    case PlanCheckRun_Type.DATABASE_GHOST_SYNC:
      return t("task.check-type.ghost-sync");
    case PlanCheckRun_Type.DATABASE_STATEMENT_SUMMARY_REPORT:
      return t("task.check-type.summary-report");
    default:
      console.assert(false, `Missing PlanCheckType name of "${type}"`);
      return type;
  }
});
</script>
