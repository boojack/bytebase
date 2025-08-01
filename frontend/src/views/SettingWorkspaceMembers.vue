<template>
  <div class="w-full mx-auto space-y-4 pb-4">
    <NoPermissionPlaceholder
      v-if="permissionStore.onlyWorkspaceMember"
      class="py-6"
    />
    <NTabs v-else v-model:value="state.selectedTab" type="line" animated>
      <NTabPane name="MEMBERS">
        <template #tab>
          <p class="text-base font-medium leading-7 text-main">
            {{ $t("settings.members.view-by-members") }}
          </p>
        </template>
        <MemberDataTable
          scope="workspace"
          :allow-edit="allowEdit"
          :bindings="memberBindings"
          :selected-bindings="state.selectedMembers"
          :select-disabled="
            (member: MemberBinding) =>
              member.user?.userType === UserType.SYSTEM_BOT
          "
          :on-click-user="onClickUser"
          @update-binding="selectMember"
          @revoke-binding="revokeMember"
          @update-selected-bindings="state.selectedMembers = $event"
        />
      </NTabPane>
      <NTabPane name="ROLES">
        <template #tab>
          <p class="text-base font-medium leading-7 text-main">
            {{ $t("settings.members.view-by-roles") }}
          </p>
        </template>
        <MemberDataTableByRole
          scope="workspace"
          :allow-edit="allowEdit"
          :bindings-by-role="memberBindingsByRole"
          :on-click-user="onClickUser"
          @update-binding="selectMember"
          @revoke-binding="revokeMember"
        />
      </NTabPane>

      <template #suffix>
        <div class="flex items-center space-x-3">
          <SearchBox
            v-model:value="state.searchText"
            :placeholder="$t('settings.members.search-member')"
          />
          <div v-if="allowEdit" class="flex justify-end gap-x-3">
            <NButton
              v-if="state.selectedTab === 'MEMBERS'"
              :disabled="state.selectedMembers.length === 0"
              @click="handleRevokeSelectedMembers"
            >
              {{ $t("settings.members.revoke-access") }}
            </NButton>
            <NButton type="primary" @click="state.showAddMemberPanel = true">
              <template #icon>
                <heroicons-outline:user-add class="w-4 h-4" />
              </template>
              {{ $t("settings.members.grant-access") }}
            </NButton>
          </div>
        </div>
      </template>
    </NTabs>
  </div>

  <EditMemberRoleDrawer
    v-if="state.showAddMemberPanel"
    :member="state.editingMember"
    @close="
      () => {
        state.showAddMemberPanel = false;
        state.editingMember = undefined;
      }
    "
  />
</template>

<script setup lang="ts">
import { computedAsync } from "@vueuse/core";
import { NButton, NTabPane, NTabs, useDialog } from "naive-ui";
import { computed, reactive } from "vue";
import { useI18n } from "vue-i18n";
import EditMemberRoleDrawer from "@/components/Member/EditMemberRoleDrawer.vue";
import MemberDataTable from "@/components/Member/MemberDataTable/index.vue";
import MemberDataTableByRole from "@/components/Member/MemberDataTableByRole.vue";
import type { MemberBinding } from "@/components/Member/types";
import {
  getMemberBindings,
  getMemberBindingsByRole,
} from "@/components/Member/utils";
import NoPermissionPlaceholder from "@/components/misc/NoPermissionPlaceholder.vue";
import { SearchBox } from "@/components/v2";
import {
  pushNotification,
  useCurrentUserV1,
  usePermissionStore,
  useWorkspaceV1Store,
} from "@/store";
import { userBindingPrefix } from "@/types";
import type { User } from "@/types/proto-es/v1/user_service_pb";
import { UserType } from "@/types/proto-es/v1/user_service_pb";
import { hasWorkspacePermissionV2 } from "@/utils";

interface LocalState {
  searchText: string;
  selectedTab: "MEMBERS" | "ROLES";
  // the member should in user:{user} or group:{group} format.
  selectedMembers: string[];
  showAddMemberPanel: boolean;
  editingMember?: MemberBinding;
}

defineProps<{
  onClickUser?: (user: User, event: MouseEvent) => void;
}>();

const { t } = useI18n();
const dialog = useDialog();
const currentUserV1 = useCurrentUserV1();
const workspaceStore = useWorkspaceV1Store();
const permissionStore = usePermissionStore();

const state = reactive<LocalState>({
  searchText: "",
  selectedTab: "MEMBERS",
  selectedMembers: [],
  showAddMemberPanel: false,
});

const allowEdit = computed(() => {
  return hasWorkspacePermissionV2("bb.workspaces.setIamPolicy");
});

const handleRevokeSelectedMembers = () => {
  if (state.selectedMembers.length === 0) {
    return;
  }

  if (
    state.selectedMembers.includes(
      `${userBindingPrefix}${currentUserV1.value.email}`
    )
  ) {
    pushNotification({
      module: "bytebase",
      style: "WARN",
      title: "You cannot revoke yourself",
    });
    return;
  }

  dialog.warning({
    title: t("settings.members.revoke-access"),
    negativeText: t("common.cancel"),
    positiveText: t("common.confirm"),
    onPositiveClick: async () => {
      await workspaceStore.patchIamPolicy(
        state.selectedMembers.map((member) => ({
          member,
          roles: [],
        }))
      );
      pushNotification({
        module: "bytebase",
        style: "INFO",
        title: t("settings.members.revoked"),
      });
      state.selectedMembers = [];
    },
  });
};

const selectMember = (binding: MemberBinding) => {
  state.editingMember = binding;
  state.showAddMemberPanel = true;
};

const revokeMember = async (binding: MemberBinding) => {
  await workspaceStore.patchIamPolicy([
    {
      member: binding.binding,
      roles: [],
    },
  ]);
  pushNotification({
    module: "bytebase",
    style: "INFO",
    title: t("settings.members.revoked"),
  });
};

const memberBindingsByRole = computedAsync(() => {
  return getMemberBindingsByRole({
    policies: [
      {
        level: "WORKSPACE",
        policy: workspaceStore.workspaceIamPolicy,
      },
    ],
    searchText: state.searchText,
    ignoreRoles: new Set([]),
  });
}, new Map());

const memberBindings = computed(() => {
  return getMemberBindings(memberBindingsByRole.value);
});
</script>
