<template>
  <BBAvatar
    :username="name"
    :email="email"
    :size="size"
    :override-class="overrideClass"
    :override-text-size="overrideTextSize"
  />
</template>

<script lang="ts" setup>
import type { PropType } from "vue";
import { computed } from "vue";
import { BBAvatar } from "@/bbkit";
import type { BBAvatarSizeType } from "@/bbkit/types";
import { extractUserId } from "@/store";
import { UNKNOWN_ID, unknownUser } from "@/types";
import type { User } from "@/types/proto-es/v1/user_service_pb";
import type { VueClass } from "@/utils";

const props = defineProps({
  user: {
    type: Object as PropType<User>,
    default: () => unknownUser(),
  },
  username: {
    type: String,
    default: "?",
  },
  size: {
    type: String as PropType<BBAvatarSizeType>,
    default: "NORMAL",
  },
  overrideClass: {
    type: [String, Object, Array] as PropType<VueClass>,
    default: undefined,
  },
  overrideTextSize: {
    type: String,
    default: undefined,
  },
});

const name = computed((): string => {
  const uid = extractUserId(props.user.name);
  if (uid === String(UNKNOWN_ID)) {
    return props.username;
  }
  return props.user.title;
});
const email = computed((): string => {
  const uid = extractUserId(props.user.name);
  if (uid === String(UNKNOWN_ID)) {
    return "";
  }
  return props.user.email;
});
</script>
