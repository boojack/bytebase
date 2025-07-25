import {
  DatabaseIcon,
  GalleryHorizontalEndIcon,
  HomeIcon,
  LayersIcon,
  LinkIcon,
  SettingsIcon,
  ShieldCheck,
  SquareStackIcon,
  UsersIcon,
  WorkflowIcon,
} from "lucide-vue-next";
import { computed, h } from "vue";
import { useRoute } from "vue-router";
import type { SidebarItem } from "@/components/v2/Sidebar/type";
import { getFlattenRoutes } from "@/components/v2/Sidebar/utils.ts";
import { t } from "@/plugins/i18n";
import workspaceRoutes from "@/router/dashboard/workspace";
import {
  DATABASE_ROUTE_DASHBOARD,
  ENVIRONMENT_V1_ROUTE_DASHBOARD,
  INSTANCE_ROUTE_DASHBOARD,
  PROJECT_V1_ROUTE_DASHBOARD,
  WORKSPACE_ROUTE_AUDIT_LOG,
  WORKSPACE_ROUTE_CUSTOM_APPROVAL,
  WORKSPACE_ROUTE_DATA_CLASSIFICATION,
  WORKSPACE_ROUTE_GLOBAL_MASKING,
  WORKSPACE_ROUTE_IM,
  WORKSPACE_ROUTE_LANDING,
  WORKSPACE_ROUTE_MEMBERS,
  WORKSPACE_ROUTE_RISKS,
  WORKSPACE_ROUTE_ROLES,
  WORKSPACE_ROUTE_SCHEMA_TEMPLATE,
  WORKSPACE_ROUTE_SEMANTIC_TYPES,
  WORKSPACE_ROUTE_SQL_REVIEW,
  WORKSPACE_ROUTE_IDENTITY_PROVIDERS,
  WORKSPACE_ROUTE_USER_PROFILE,
  WORKSPACE_ROUTE_USERS,
} from "@/router/dashboard/workspaceRoutes";
import {
  SETTING_ROUTE_WORKSPACE_ARCHIVE,
  SETTING_ROUTE_WORKSPACE_GENERAL,
  SETTING_ROUTE_WORKSPACE_SUBSCRIPTION,
} from "@/router/dashboard/workspaceSetting";
import { useAppFeature, usePermissionStore } from "@/store";
import { DatabaseChangeMode } from "@/types/proto-es/v1/setting_service_pb";
import { hasWorkspacePermissionV2 } from "@/utils";

export interface DashboardSidebarItem extends SidebarItem {
  navigationId?: string;
  shortcuts?: string[];
  hide?: boolean;
  children?: DashboardSidebarItem[];
}

export const useDashboardSidebar = () => {
  const route = useRoute();
  const permissionStore = usePermissionStore();
  const databaseChangeMode = useAppFeature("bb.feature.database-change-mode");

  const getItemClass = (item: SidebarItem): string[] => {
    const { name: current } = route;
    const isActiveRoute =
      item.name === current?.toString() ||
      current?.toString().startsWith(`${item.name}.`);
    const classes: string[] = [];
    if (isActiveRoute) {
      classes.push("router-link-active", "bg-link-hover");
      return classes;
    }
    if (
      item.name === WORKSPACE_ROUTE_USERS &&
      current?.toString() === WORKSPACE_ROUTE_USER_PROFILE
    ) {
      classes.push("router-link-active", "bg-link-hover");
      return classes;
    }
    return classes;
  };

  const filterSidebarByPermissions = (
    sidebarList: DashboardSidebarItem[]
  ): DashboardSidebarItem[] => {
    return sidebarList
      .filter((item) => {
        const routeConfig = flattenRoutes.value.find(
          (workspaceRoute) => workspaceRoute.name === item.name
        );
        return (routeConfig?.permissions ?? []).every((permission) =>
          hasWorkspacePermissionV2(permission)
        );
      })
      .map((item) => ({
        ...item,
        expand:
          item.expand ||
          (item.children ?? [])
            .reduce((classList, child) => {
              classList.push(...getItemClass(child));
              return classList;
            }, [] as string[])
            .includes("router-link-active"),
        children: filterSidebarByPermissions(item.children ?? []),
      }));
  };

  const flattenRoutes = computed(() => {
    return getFlattenRoutes(workspaceRoutes);
  });

  const dashboardSidebarItemList = computed((): DashboardSidebarItem[] => {
    const sidebarList: DashboardSidebarItem[] = [
      {
        navigationId: "bb.navigation.home",
        title: t("common.home"),
        icon: () => h(HomeIcon),
        name: WORKSPACE_ROUTE_LANDING,
        type: "route",
        shortcuts: ["g", "h"],
      },
      {
        navigationId: "bb.navigation.projects",
        title: t("common.projects"),
        icon: () => h(GalleryHorizontalEndIcon),
        name: PROJECT_V1_ROUTE_DASHBOARD,
        type: "route",
        shortcuts: ["g", "p"],
      },
      {
        navigationId: "bb.navigation.instances",
        title: t("common.instances"),
        icon: () => h(LayersIcon),
        name: INSTANCE_ROUTE_DASHBOARD,
        type: "route",
        shortcuts: ["g", "i"],
      },
      {
        navigationId: "bb.navigation.databases",
        title: t("common.databases"),
        icon: () => h(DatabaseIcon),
        name: DATABASE_ROUTE_DASHBOARD,
        type: "route",
        shortcuts: ["g", "d"],
      },
      {
        navigationId: "bb.navigation.environments",
        title: t("common.environments"),
        icon: () => h(SquareStackIcon),
        name: ENVIRONMENT_V1_ROUTE_DASHBOARD,
        type: "route",
        shortcuts: ["g", "e"],
      },
      {
        type: "divider",
        name: "",
      },
      {
        title: t("settings.sidebar.iam-and-admin"),
        icon: () => h(UsersIcon),
        type: "div",
        children: [
          {
            title: t("settings.sidebar.users-and-groups"),
            name: WORKSPACE_ROUTE_USERS,
            type: "route",
          },
          {
            title: t("settings.sidebar.members"),
            name: WORKSPACE_ROUTE_MEMBERS,
            type: "route",
            hide: permissionStore.onlyWorkspaceMember,
          },
          {
            title: t("settings.sidebar.custom-roles"),
            name: WORKSPACE_ROUTE_ROLES,
            type: "route",
          },
          {
            title: t("settings.sidebar.sso"),
            name: WORKSPACE_ROUTE_IDENTITY_PROVIDERS,
            type: "route",
          },
          {
            title: t("settings.sidebar.audit-log"),
            name: WORKSPACE_ROUTE_AUDIT_LOG,
            type: "route",
          },
        ],
      },
      {
        title: "CI/CD",
        icon: () => h(WorkflowIcon),
        type: "div",
        hide: databaseChangeMode.value === DatabaseChangeMode.EDITOR,
        children: [
          {
            title: t("sql-review.title"),
            name: WORKSPACE_ROUTE_SQL_REVIEW,
            type: "route",
          },
          {
            title: t("custom-approval.risk.risks"),
            name: WORKSPACE_ROUTE_RISKS,
            type: "route",
          },
          {
            title: t("custom-approval.self"),
            name: WORKSPACE_ROUTE_CUSTOM_APPROVAL,
            type: "route",
          },
          {
            title: t("schema-template.self"),
            name: WORKSPACE_ROUTE_SCHEMA_TEMPLATE,
            type: "route",
          },
        ],
      },
      {
        title: t("settings.sidebar.data-access"),
        icon: () => h(ShieldCheck),
        type: "div",
        children: [
          {
            title: t("settings.sensitive-data.semantic-types.self"),
            name: WORKSPACE_ROUTE_SEMANTIC_TYPES,
            type: "route",
          },
          {
            title: t("settings.sidebar.data-classification"),
            name: WORKSPACE_ROUTE_DATA_CLASSIFICATION,
            type: "route",
          },
          {
            title: t("settings.sidebar.global-masking"),
            name: WORKSPACE_ROUTE_GLOBAL_MASKING,
            type: "route",
          },
        ],
      },
      {
        title: t("settings.sidebar.integration"),
        icon: () => h(LinkIcon),
        type: "div",
        hide: databaseChangeMode.value === DatabaseChangeMode.EDITOR,
        children: [
          {
            title: t("settings.sidebar.im-integration"),
            name: WORKSPACE_ROUTE_IM,
            type: "route",
          },
        ],
      },
      {
        title: t("common.settings"),
        icon: () => h(SettingsIcon),
        type: "div",
        hide: !hasWorkspacePermissionV2("bb.settings.get"),
        children: [
          {
            title: t("settings.sidebar.general"),
            name: SETTING_ROUTE_WORKSPACE_GENERAL,
            type: "route",
          },
          {
            title: t("settings.sidebar.subscription"),
            name: SETTING_ROUTE_WORKSPACE_SUBSCRIPTION,
            type: "route",
          },
          {
            title: t("common.archived"),
            name: SETTING_ROUTE_WORKSPACE_ARCHIVE,
            type: "route",
          },
        ],
      },
    ];

    return filterSidebarByPermissions(sidebarList);
  });

  return {
    getItemClass,
    dashboardSidebarItemList,
  };
};
