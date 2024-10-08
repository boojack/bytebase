<template>
  <!-- eslint-disable vue/no-mutating-props -->
  <NInputGroup class="w-full flex items-center overflow-x-hidden">
    <NSelect
      v-model:value="selector.key"
      :options="keyOptions"
      v-bind="commonSelectProps"
      style="width: auto"
    />
    <NSelect
      v-model:value="selector.operator"
      :options="operatorOptions"
      v-bind="commonSelectProps"
      style="width: auto"
    />
    <NSelect
      v-if="selector.operator !== OperatorType.OPERATOR_TYPE_EXISTS"
      v-bind="valueSelectProps"
      @update:value="handleUpdateValues"
    />
    <ErrorTipsButton
      v-if="editable"
      :errors="
        isRequiredEnvironmentSelector
          ? [$t('deployment-config.error.env-in-selector-required')]
          : []
      "
      style="--n-padding: 10px"
      class="shrink-0"
      @click="$emit('remove')"
    >
      <TrashIcon class="w-4 h-4" />
    </ErrorTipsButton>
  </NInputGroup>
</template>

<script lang="ts" setup>
/* eslint-disable vue/no-mutating-props */
import { TrashIcon } from "lucide-vue-next";
import type { SelectOption, SelectProps } from "naive-ui";
import { NInputGroup, NSelect } from "naive-ui";
import type { CSSProperties, PropType } from "vue";
import { computed, watch } from "vue";
import { useI18n } from "vue-i18n";
import type { ComposedDatabase } from "@/types";
import type { LabelSelectorRequirement } from "@/types/proto/v1/project_service";
import { OperatorType } from "@/types/proto/v1/project_service";
import {
  displayDeploymentMatchSelectorKey,
  getAvailableDeploymentConfigMatchSelectorKeyList,
  getLabelValuesFromDatabaseV1List,
  unwrapArray,
  wrapArray,
} from "@/utils";
import { ErrorTipsButton } from "../v2";

const OPERATORS: OperatorType[] = [
  OperatorType.OPERATOR_TYPE_IN,
  OperatorType.OPERATOR_TYPE_NOT_IN,
  OperatorType.OPERATOR_TYPE_EXISTS,
];

const props = defineProps({
  selector: {
    type: Object as PropType<LabelSelectorRequirement>,
    required: true,
  },
  selectors: {
    type: Array as PropType<LabelSelectorRequirement[]>,
    default: () => [],
  },
  index: {
    type: Number,
    default: -1,
  },
  databaseList: {
    type: Array as PropType<ComposedDatabase[]>,
    default: () => [],
  },
  editable: {
    type: Boolean,
    default: false,
  },
});

defineEmits<{
  (event: "remove"): void;
}>();

const { t } = useI18n();

const isRequiredEnvironmentSelector = computed(() => {
  return props.selector.key === "environment" && props.index === 0;
});

const commonSelectProps = computed((): SelectProps => {
  return {
    consistentMenuWidth: false,
    disabled: !props.editable || isRequiredEnvironmentSelector.value,
    showArrow: props.editable && !isRequiredEnvironmentSelector.value,
  };
});

const operators = computed(() => {
  return OPERATORS.filter((op) => {
    if (
      props.selector.key === "environment" &&
      op === OperatorType.OPERATOR_TYPE_NOT_IN
    ) {
      return false;
    }
    return true;
  });
});

const keys = computed(() => {
  return getAvailableDeploymentConfigMatchSelectorKeyList(
    props.databaseList,
    true /* withVirtualLabelKeys */,
    true /* sort */
  );
});
const keyOptions = computed(() => {
  return keys.value.map<SelectOption>((key) => {
    return {
      label: displayDeploymentMatchSelectorKey(key),
      value: key,
    };
  });
});
const operatorOptions = computed(() => {
  return operators.value.map<SelectOption>((op) => {
    return {
      label: operatorToText(op),
      value: op,
    };
  });
});
const allowMultipleValues = computed(() => {
  return props.selector.key !== "environment";
});
const values = computed(() => {
  if (!props.selector.key) return [];
  return getLabelValuesFromDatabaseV1List(
    props.selector.key,
    props.databaseList
  );
});
const valueOptions = computed(() => {
  return values.value.map<SelectOption>((value) => {
    return {
      label: value,
      value,
    };
  });
});

const resetValues = () => {
  props.selector.values = [];
};

const operatorToText = (op: OperatorType) => {
  switch (op) {
    case OperatorType.OPERATOR_TYPE_IN:
      return allowMultipleValues.value ? "in" : "is";
    case OperatorType.OPERATOR_TYPE_NOT_IN:
      return "not in";
    case OperatorType.OPERATOR_TYPE_EXISTS:
      return "exists";
  }
  console.error("[operatorToText] should never reach this line", op);
  return "";
};

const valueSelectProps = computed(() => {
  const multiple = allowMultipleValues.value;
  const { values } = props.selector;
  const selectProps: SelectProps = {
    ...commonSelectProps.value,
    multiple,
    value: multiple ? wrapArray(values) : unwrapArray(values),
    options: valueOptions.value,
    maxTagCount: "responsive",
    disabled: !props.editable,
    showArrow: props.editable,
    placeholder: t("label.placeholder.select-values"),
  };
  const style: CSSProperties = {
    "min-width": props.editable ? (multiple ? "12rem" : "8rem") : undefined,
    width: "auto",
    "overflow-x": "hidden",
  };
  return {
    ...selectProps,
    style,
  };
});

const handleUpdateValues = (values: string | string[]) => {
  if (allowMultipleValues.value) {
    props.selector.values = wrapArray(values);
  } else {
    props.selector.values = [unwrapArray(values)];
  }
};

watch(() => props.selector.key, resetValues);
watch(() => props.selector.operator, resetValues);
</script>
