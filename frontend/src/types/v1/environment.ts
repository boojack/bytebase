import { environmentNamePrefix } from "@/store";
import { EMPTY_ID, UNKNOWN_ID } from "../const";
import type { EnvironmentSetting_Environment } from "../proto/v1/setting_service";

export const ENVIRONMENT_ALL_NAME = "environments/-";
export const EMPTY_ENVIRONMENT_NAME = `environments/${EMPTY_ID}`;
export const UNKNOWN_ENVIRONMENT_NAME = `environments/${UNKNOWN_ID}`;

export interface Environment extends EnvironmentSetting_Environment {
  order: number;
}

export const unknownEnvironment = (): Environment => {
  return {
    name: UNKNOWN_ENVIRONMENT_NAME,
    id: String(UNKNOWN_ID),
    title: "<<Unknown environment>>",
    order: 0,
    tags: {},
    color: "",
  };
};

export const isValidEnvironmentName = (name: any): name is string => {
  return (
    typeof name === "string" &&
    name.startsWith("environments/") &&
    name !== EMPTY_ENVIRONMENT_NAME &&
    name !== UNKNOWN_ENVIRONMENT_NAME
  );
};

export const formatEnvironmentName = (envId: string): string => {
  return `${environmentNamePrefix}${envId}`;
};
