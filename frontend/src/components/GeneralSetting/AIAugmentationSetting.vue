<template>
  <div id="ai" ref="containerRef" class="py-6 lg:flex">
    <div class="text-left lg:w-1/4">
      <div class="flex items-center space-x-2">
        <h1 class="text-2xl font-bold">
          {{ title }}
        </h1>
      </div>
      <span v-if="!allowEdit" class="text-sm text-gray-400">
        {{ $t("settings.general.workspace.only-admin-can-edit") }}
      </span>
    </div>
    <div class="flex-1 lg:px-4">
      <div class="mt-4 lg:mt-0 space-y-4">
        <div>
          <div class="flex items-center gap-x-2">
            <Switch
              v-model:value="state.enabled"
              :text="true"
              :disabled="!allowEdit"
              @update:value="toggleAIEnabled"
            />
            <span class="font-medium">
              {{
                $t(
                  "settings.general.workspace.ai-assistant.enable-ai-assistant"
                )
              }}
            </span>
          </div>
          <div class="mt-1 mb-3 text-sm text-gray-400">
            {{ $t("settings.general.workspace.ai-assistant.description") }}
            <LearnMoreLink
              url="https://docs.bytebase.com/ai-assistant?source=console"
              class="ml-1 text-sm"
            />
          </div>
        </div>
        <template v-if="state.enabled">
          <div>
            <label class="flex items-center gap-x-2 mb-2">
              <span class="font-medium">{{
                $t("settings.general.workspace.ai-assistant.provider.self")
              }}</span>
            </label>
            <NSelect
              style="width: 12rem"
              v-model:value="state.provider"
              :options="providerOptions"
              :disabled="!allowEdit"
              :consistent-menu-width="true"
            />
          </div>
          <div>
            <label class="flex items-center gap-x-2">
              <span class="font-medium">{{
                $t("settings.general.workspace.ai-assistant.api-key.self")
              }}</span>
            </label>
            <div class="mb-3 text-sm text-gray-400">
              <i18n-t
                keypath="settings.general.workspace.ai-assistant.api-key.description"
              >
                <template #viewDoc>
                  <a
                    :href="providerDefault.apiKeyDoc"
                    class="normal-link"
                    target="_blank"
                    >{{
                      $t(
                        "settings.general.workspace.ai-assistant.api-key.find-my-key"
                      )
                    }}
                  </a>
                </template>
              </i18n-t>
            </div>
            <NTooltip placement="top-start" :disabled="allowEdit">
              <template #trigger>
                <BBTextField
                  v-model:value="state.apiKey"
                  :disabled="!allowEdit"
                  :placeholder="
                    $t(
                      'settings.general.workspace.ai-assistant.api-key.placeholder'
                    )
                  "
                />
              </template>
              <span class="text-sm text-gray-400 -translate-y-2">
                {{ $t("settings.general.workspace.only-admin-can-edit") }}
              </span>
            </NTooltip>
          </div>

          <div>
            <label class="flex items-center gap-x-2">
              <span class="font-medium">{{
                $t("settings.general.workspace.ai-assistant.endpoint.self")
              }}</span>
            </label>
            <div class="mb-3 text-sm text-gray-400">
              {{
                $t(
                  "settings.general.workspace.ai-assistant.endpoint.description"
                )
              }}
            </div>
            <NTooltip placement="top-start" :disabled="allowEdit">
              <template #trigger>
                <BBTextField
                  v-model:value="state.endpoint"
                  :required="true"
                  :disabled="!allowEdit"
                  :placeholder="providerDefault.endpoint"
                />
              </template>
              <span class="text-sm text-gray-400 -translate-y-2">
                {{ $t("settings.general.workspace.only-admin-can-edit") }}
              </span>
            </NTooltip>
          </div>

          <div>
            <label class="flex items-center gap-x-2">
              <span class="font-medium">{{
                $t("settings.general.workspace.ai-assistant.model.self")
              }}</span>
            </label>
            <div class="mb-3 text-sm text-gray-400">
              {{
                $t("settings.general.workspace.ai-assistant.model.description")
              }}
            </div>
            <NTooltip placement="top-start" :disabled="allowEdit">
              <template #trigger>
                <BBTextField
                  v-model:value="state.model"
                  :required="true"
                  :disabled="!allowEdit"
                />
              </template>
              <span class="text-sm text-gray-400 -translate-y-2">
                {{ $t("settings.general.workspace.only-admin-can-edit") }}
              </span>
            </NTooltip>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { create } from "@bufbuild/protobuf";
import { NSelect, NTooltip } from "naive-ui";
import scrollIntoView from "scroll-into-view-if-needed";
import { computed, onMounted, reactive, ref, watch, watchEffect } from "vue";
import { useI18n } from "vue-i18n";
import { BBTextField } from "@/bbkit";
import LearnMoreLink from "@/components/LearnMoreLink.vue";
import { Switch } from "@/components/v2";
import { useSettingV1Store } from "@/store/modules/v1/setting";
import {
  AISettingSchema,
  AISetting_Provider,
  Setting_SettingName,
  ValueSchema as SettingValueSchema,
} from "@/types/proto-es/v1/setting_service_pb";

interface LocalState {
  enabled: boolean;
  apiKey: string;
  endpoint: string;
  model: string;
  provider: AISetting_Provider;
}

const props = defineProps<{
  title: string;
  allowEdit: boolean;
}>();

const settingV1Store = useSettingV1Store();
const containerRef = ref<HTMLDivElement>();
const { t } = useI18n();

const state = reactive<LocalState>({
  enabled: false,
  apiKey: "",
  endpoint: "",
  model: "",
  provider: AISetting_Provider.OPEN_AI,
});

const aiSetting = computed(() => {
  const setting = settingV1Store.getSettingByName(Setting_SettingName.AI);
  if (setting?.value?.value?.case === "aiSetting") {
    return setting.value.value.value;
  }
  return undefined;
});

const getInitialState = (): LocalState => {
  return {
    enabled: aiSetting.value?.enabled ?? false,
    apiKey: "",
    endpoint: aiSetting.value?.endpoint ?? "",
    model: aiSetting.value?.model ?? "",
    provider: aiSetting.value?.provider ?? AISetting_Provider.OPEN_AI,
  };
};

const providerOptions = computed(() =>
  [AISetting_Provider.OPEN_AI, AISetting_Provider.AZURE_OPENAI].map(
    (provider) => {
      let label = "";
      switch (provider) {
        case AISetting_Provider.OPEN_AI:
          label = t("settings.general.workspace.ai-assistant.provider.open_ai");
          break;
        case AISetting_Provider.AZURE_OPENAI:
          label = t(
            "settings.general.workspace.ai-assistant.provider.azure_open_ai"
          );
          break;
      }
      return {
        label,
        value: provider,
      };
    }
  )
);

watchEffect(() => {
  Object.assign(state, getInitialState());
});

const allowSave = computed((): boolean => {
  const initValue = getInitialState();
  const enabledUpdated = state.enabled !== initValue.enabled;
  const openAIKeyUpdated = !!state.apiKey;
  const openAIEndpointUpdated = state.endpoint !== initValue.endpoint;
  const openAIModelUpdated = state.model !== initValue.model;
  return (
    enabledUpdated ||
    openAIKeyUpdated ||
    openAIEndpointUpdated ||
    openAIModelUpdated
  );
});

const providerDefault = computed(() => {
  switch (state.provider) {
    case AISetting_Provider.OPEN_AI:
      return {
        apiKey: "",
        apiKeyDoc: "https://platform.openai.com/account/api-keys",
        endpoint: "https://api.openai.com/v1/chat/completions",
        model: "gpt-3.5-turbo",
      };
    case AISetting_Provider.AZURE_OPENAI:
      return {
        apiKey: "",
        apiKeyDoc: "https://ai.azure.com/",
        endpoint:
          "https://{resource name}.openai.azure.com/openai/deployments/{deployment id}/chat/completions?api-version=2024-06-0",
        model: "gpt-4o",
      };
    default:
      return {
        apiKey: "",
        apiKeyDoc: "",
        endpoint: "",
        model: "",
      };
  }
});

watch(
  () => state.provider,
  () => {
    Object.assign(state, providerDefault.value);
  }
);

const toggleAIEnabled = (on: boolean) => {
  if (!on) {
    return;
  }
  Object.assign(state, {
    endpoint: state.endpoint || providerDefault.value.endpoint,
    model: state.model || providerDefault.value.model,
  });
};

const updateAISetting = async () => {
  await settingV1Store.upsertSetting({
    name: Setting_SettingName.AI,
    value: create(SettingValueSchema, {
      value: {
        case: "aiSetting",
        value: create(AISettingSchema, {
          enabled: state.enabled,
          apiKey: state.apiKey,
          endpoint: state.endpoint,
          model: state.model,
          provider: state.provider,
          version: aiSetting.value?.version ?? "",
        }),
      },
    }),
  });

  Object.assign(state, getInitialState());
};

onMounted(() => {
  if (location.hash === "#ai-assistant") {
    const container = containerRef.value;
    if (!container) return;
    scrollIntoView(container, {
      scrollMode: "if-needed",
    });
  }
});

defineExpose({
  isDirty: allowSave,
  title: props.title,
  update: updateAISetting,
  revert: () => {
    Object.assign(state, getInitialState());
  },
});
</script>
