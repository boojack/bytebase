<template>
  <NRadioGroup
    v-bind="$attrs"
    :value="state.selected"
    class="w-full !grid grid-cols-3 gap-2"
    name="radiogroup"
    @update:value="onSelect"
  >
    <div
      v-for="option in options"
      :key="option.value"
      :class="[
        'col-span-1 h-8 flex flex-row justify-start items-center',
        option.value === -2 ? 'col-span-2' : '',
        option.value === -1 ? 'col-span-3' : '',
      ]"
    >
      <NRadio :value="option.value" :label="option.label" />
      <template v-if="option.value === -2">
        <NInputNumber
          size="small"
          style="width: 5rem"
          :min="1"
          :max="maximumRoleExpiration"
          :placeholder="$t('common.date.days')"
          :disabled="state.selected !== -2"
          :value="customExpirationDays"
          @update:value="handleCustomExpirationDaysChange"
        />
        <span class="ml-2">{{ $t("common.date.days") }}</span>
      </template>
      <template v-else-if="option.value === -1">
        <NDatePicker
          size="small"
          :value="
            state.selected === -1 ? state.expirationTimestampInMS : undefined
          "
          :disabled="state.selected !== -1"
          :actions="null"
          :update-value-on-close="true"
          type="datetime"
          :is-date-disabled="isDateDisabled"
          clearable
          @update:value="(val) => (state.expirationTimestampInMS = val)"
        />
        <span v-if="maximumRoleExpiration" class="ml-3 textinfolabel">
          {{ $t("settings.general.workspace.maximum-role-expiration.self") }}:
          {{ $t("common.date.days", { days: maximumRoleExpiration }) }}
        </span>
      </template>
    </div>
  </NRadioGroup>
</template>

<script lang="ts" setup>
import { useLocalStorage } from "@vueuse/core";
import dayjs from "dayjs";
import { NRadio, NRadioGroup, NDatePicker, NInputNumber } from "naive-ui";
import { computed, reactive, watch, watchEffect } from "vue";
import { useI18n } from "vue-i18n";
import { useSettingV1Store } from "@/store";
import { PresetRoleType } from "@/types";

interface ExpirationOption {
  value: number;
  label: string;
}

interface LocalState {
  selected: number;
  expirationTimestampInMS?: number;
}

const props = defineProps<{
  timestampInMs?: number;
  role: string;
}>();

const emit = defineEmits<{
  (event: "update:timestampInMs", timestampInMS: number | undefined): void;
}>();

const { t } = useI18n();
const settingV1Store = useSettingV1Store();

const state = reactive<LocalState>({
  selected: props.timestampInMs === undefined ? 0 : -1,
});
const customExpirationDays = useLocalStorage(
  "bb.roles.custom-expiration-days",
  1
);

const handleCustomExpirationDaysChange = (val: number | null) => {
  if (!val) {
    return;
  }
  state.expirationTimestampInMS =
    new Date().getTime() + val * 24 * 60 * 60 * 1000;
  customExpirationDays.value = val;
};

const maximumRoleExpiration = computed(() => {
  if (props.role === PresetRoleType.PROJECT_OWNER) {
    return undefined;
  }
  const seconds =
    settingV1Store.workspaceProfileSetting?.maximumRoleExpiration?.seconds?.toNumber();
  if (!seconds) {
    return undefined;
  }
  return Math.floor(seconds / (60 * 60 * 24));
});

watchEffect(() => {
  if (!maximumRoleExpiration.value) {
    return;
  }
  if (maximumRoleExpiration.value < customExpirationDays.value) {
    customExpirationDays.value = maximumRoleExpiration.value;
  }
});

const isDateDisabled = (date: number) => {
  if (date < dayjs().startOf("day").valueOf()) {
    return true;
  }
  if (!maximumRoleExpiration.value) {
    return false;
  }
  return date > dayjs().add(maximumRoleExpiration.value, "days").valueOf();
};

const options = computed((): ExpirationOption[] => {
  let options = [
    {
      value: 1,
      label: t("common.date.days", { days: 1 }),
    },
    {
      value: 3,
      label: t("common.date.days", { days: 3 }),
    },
    {
      value: 30,
      label: t("common.date.days", { days: 30 }),
    },
    {
      value: 90,
      label: t("common.date.days", { days: 90 }),
    },
  ];
  if (maximumRoleExpiration.value) {
    options = options.filter(
      (option) => option.value < maximumRoleExpiration.value!
    );
    options.push({
      value: maximumRoleExpiration.value,
      label: t("common.date.days", { days: maximumRoleExpiration.value }),
    });
  } else {
    options.push({
      value: 0,
      label: t("project.members.never-expires"),
    });
  }
  options.push(
    {
      value: -2,
      label: t("issue.grant-request.custom"),
    },
    {
      value: -1,
      label: t("issue.grant-request.custom"),
    }
  );
  return options;
});

const onSelect = (value: number) => {
  if (value > 0) {
    state.expirationTimestampInMS =
      new Date().getTime() + value * 24 * 60 * 60 * 1000;
  }
  state.selected = value;
};

watch(
  () => props.role,
  () => {
    let value = state.selected;
    if (!options.value.find((o) => o.value === state.selected)) {
      const neverExpire = options.value.find((o) => o.value === 0);
      if (neverExpire) {
        value = neverExpire.value;
      } else {
        value = options.value[0].value;
      }
    }
    onSelect(value);
  },
  { immediate: true }
);

watch(
  () => state.expirationTimestampInMS,
  () => {
    emit("update:timestampInMs", state.expirationTimestampInMS);
  }
);
</script>
