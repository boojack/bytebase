<template>
  <div class="w-full space-y-4">
    <FeatureAttention :feature="PlanFeature.FEATURE_SCHEMA_TEMPLATE" />
    <NTabs v-model:value="state.selectedTab" type="line">
      <NTabPane
        name="FIELD_TEMPLATE"
        :tab="$t('schema-template.field-template.self')"
      >
        <FieldTemplates
          :show-engine-filter="true"
          :readonly="!hasFeature || !allowEdit"
        />
      </NTabPane>
      <NTabPane
        name="TABLE_TEMPLATE"
        :tab="$t('schema-template.table-template.self')"
      >
        <TableTemplates
          :show-engine-filter="true"
          :readonly="!hasFeature || !allowEdit"
        />
      </NTabPane>
      <NTabPane
        name="COLUMN_TYPE_RESTRICTION"
        :tab="$t('schema-template.column-type-restriction.self')"
      >
        <ColumnTypes :readonly="!hasFeature || !allowEdit" />
      </NTabPane>
    </NTabs>
  </div>
</template>

<script lang="ts" setup>
import { NTabs, NTabPane } from "naive-ui";
import { reactive, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { FeatureAttention } from "@/components/FeatureGuard";
import { featureToRef } from "@/store";
import { PlanFeature } from "@/types/proto-es/v1/subscription_service_pb";
import ColumnTypes from "@/views/SchemaTemplate/ColumnTypes.vue";
import FieldTemplates from "@/views/SchemaTemplate/FieldTemplates.vue";
import TableTemplates from "@/views/SchemaTemplate/TableTemplates.vue";

interface LocalState {
  selectedTab: "FIELD_TEMPLATE" | "COLUMN_TYPE_RESTRICTION" | "TABLE_TEMPLATE";
}

defineProps<{
  allowEdit: boolean;
}>();

const route = useRoute();
const router = useRouter();
const state = reactive<LocalState>({
  selectedTab: "FIELD_TEMPLATE",
});

const hasFeature = featureToRef(PlanFeature.FEATURE_SCHEMA_TEMPLATE);

watch(
  () => route.hash,
  (hash) => {
    switch (hash) {
      case "#column-type-restriction":
        state.selectedTab = "COLUMN_TYPE_RESTRICTION";
        break;
      case "#table-template":
        state.selectedTab = "TABLE_TEMPLATE";
        break;
      default:
        state.selectedTab = "FIELD_TEMPLATE";
        break;
    }
  },
  {
    immediate: true,
  }
);

watch(
  () => state.selectedTab,
  (tab) => {
    switch (tab) {
      case "COLUMN_TYPE_RESTRICTION":
        router.push({ hash: "#column-type-restriction" });
        break;
      case "TABLE_TEMPLATE":
        router.push({ hash: "#table-template" });
        break;
      default:
        router.push({ hash: "" });
        break;
    }
  }
);
</script>
