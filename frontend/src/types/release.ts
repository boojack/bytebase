import { getProjectNameReleaseId } from "@/store/modules/v1/common";
import { UNKNOWN_ID } from "./const";
import { Release } from "./proto/v1/release_service";
import type { User } from "./proto/v1/user_service";
import { unknownUser } from "./v1";
import {
  UNKNOWN_PROJECT_NAME,
  unknownProject,
  type ComposedProject,
} from "./v1/project";

export interface ComposedRelease extends Release {
  // Format: projects/{project}
  project: string;
  projectEntity: ComposedProject;
  creatorEntity: User;
}

export const UNKNOWN_RELEASE_NAME = `${UNKNOWN_PROJECT_NAME}/releases/${UNKNOWN_ID}`;

export const unknownRelease = (): ComposedRelease => {
  const projectEntity = unknownProject();
  const release = Release.fromPartial({
    name: `${projectEntity.name}/releases/${UNKNOWN_ID}`,
  });
  return {
    ...release,
    project: projectEntity.name,
    projectEntity,
    creatorEntity: unknownUser(),
  };
};

export const isValidReleaseName = (name: any): name is string => {
  if (typeof name !== "string") return false;
  const [projectName, releaseName] = getProjectNameReleaseId(name);
  return Boolean(
    projectName &&
      projectName !== String(UNKNOWN_ID) &&
      releaseName &&
      releaseName !== String(UNKNOWN_ID)
  );
};
