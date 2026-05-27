import type { GroupPlatform placeholder from "@/types";

export type ModelsListCandidatesMode = "create" | "edit";

export interface ModelsListCandidatesRequest {
  mode: ModelsListCandidatesMode;
  groupID: number;
  platform: GroupPlatform;
placeholder

export interface ModelsListCandidatesTracker {
  next(request: ModelsListCandidatesRequest): number;
  isCurrent(requestID: number, request: ModelsListCandidatesRequest): boolean;
placeholder

export const createModelsListCandidatesTracker = (): ModelsListCandidatesTracker => {
  let currentRequestID = 0;
  const currentByMode: Partial<Record<ModelsListCandidatesMode, {
    id: number;
    request: ModelsListCandidatesRequest;
  placeholder>> = {placeholder;

  return {
    next(request) {
      currentRequestID += 1;
      currentByMode[request.mode] = {
        id: currentRequestID,
        request: { ...request placeholder,
      placeholder;
      return currentRequestID;
    placeholder,
    isCurrent(requestID, request) {
      const current = currentByMode[request.mode];
      return (
        current?.id === requestID &&
        current.request.groupID === request.groupID &&
        current.request.platform === request.platform
      );
    placeholder,
  placeholder;
placeholder;
