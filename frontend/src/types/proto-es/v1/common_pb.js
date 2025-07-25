// @generated by protoc-gen-es v2.5.2
// @generated from file v1/common.proto (package bytebase.v1, syntax proto3)
/* eslint-disable */

import { enumDesc, fileDesc, messageDesc, tsEnum } from "@bufbuild/protobuf/codegenv2";

/**
 * Describes the file v1/common.proto.
 */
export const file_v1_common = /*@__PURE__*/
  fileDesc("Cg92MS9jb21tb24ucHJvdG8SC2J5dGViYXNlLnYxIigKCFBvc2l0aW9uEgwKBGxpbmUYASABKAUSDgoGY29sdW1uGAIgASgFIiMKBVJhbmdlEg0KBXN0YXJ0GAEgASgFEgsKA2VuZBgCIAEoBSo3CgVTdGF0ZRIVChFTVEFURV9VTlNQRUNJRklFRBAAEgoKBkFDVElWRRABEgsKB0RFTEVURUQQAiqeAwoGRW5naW5lEhYKEkVOR0lORV9VTlNQRUNJRklFRBAAEg4KCkNMSUNLSE9VU0UQARIJCgVNWVNRTBACEgwKCFBPU1RHUkVTEAMSDQoJU05PV0ZMQUtFEAQSCgoGU1FMSVRFEAUSCAoEVElEQhAGEgsKB01PTkdPREIQBxIJCgVSRURJUxAIEgoKBk9SQUNMRRAJEgsKB1NQQU5ORVIQChIJCgVNU1NRTBALEgwKCFJFRFNISUZUEAwSCwoHTUFSSUFEQhANEg0KCU9DRUFOQkFTRRAOEgYKAkRNEA8SDgoKUklTSU5HV0FWRRAQEhQKEE9DRUFOQkFTRV9PUkFDTEUQERINCglTVEFSUk9DS1MQEhIJCgVET1JJUxATEggKBEhJVkUQFBIRCg1FTEFTVElDU0VBUkNIEBUSDAoIQklHUVVFUlkQFhIMCghEWU5BTU9EQhAXEg4KCkRBVEFCUklDS1MQGBIPCgtDT0NLUk9BQ0hEQhAZEgwKCENPU01PU0RCEBoSCQoFVFJJTk8QGxINCglDQVNTQU5EUkEQHCpcCgdWQ1NUeXBlEhgKFFZDU19UWVBFX1VOU1BFQ0lGSUVEEAASCgoGR0lUSFVCEAESCgoGR0lUTEFCEAISDQoJQklUQlVDS0VUEAMSEAoMQVpVUkVfREVWT1BTEAQqTAoMRXhwb3J0Rm9ybWF0EhYKEkZPUk1BVF9VTlNQRUNJRklFRBAAEgcKA0NTVhABEggKBEpTT04QAhIHCgNTUUwQAxIICgRYTFNYEARCNlo0Z2l0aHViLmNvbS9ieXRlYmFzZS9ieXRlYmFzZS9iYWNrZW5kL2dlbmVyYXRlZC1nby92MWIGcHJvdG8z");

/**
 * Describes the message bytebase.v1.Position.
 * Use `create(PositionSchema)` to create a new message.
 */
export const PositionSchema = /*@__PURE__*/
  messageDesc(file_v1_common, 0);

/**
 * Describes the message bytebase.v1.Range.
 * Use `create(RangeSchema)` to create a new message.
 */
export const RangeSchema = /*@__PURE__*/
  messageDesc(file_v1_common, 1);

/**
 * Describes the enum bytebase.v1.State.
 */
export const StateSchema = /*@__PURE__*/
  enumDesc(file_v1_common, 0);

/**
 * @generated from enum bytebase.v1.State
 */
export const State = /*@__PURE__*/
  tsEnum(StateSchema);

/**
 * Describes the enum bytebase.v1.Engine.
 */
export const EngineSchema = /*@__PURE__*/
  enumDesc(file_v1_common, 1);

/**
 * @generated from enum bytebase.v1.Engine
 */
export const Engine = /*@__PURE__*/
  tsEnum(EngineSchema);

/**
 * Describes the enum bytebase.v1.VCSType.
 */
export const VCSTypeSchema = /*@__PURE__*/
  enumDesc(file_v1_common, 2);

/**
 * @generated from enum bytebase.v1.VCSType
 */
export const VCSType = /*@__PURE__*/
  tsEnum(VCSTypeSchema);

/**
 * Describes the enum bytebase.v1.ExportFormat.
 */
export const ExportFormatSchema = /*@__PURE__*/
  enumDesc(file_v1_common, 3);

/**
 * @generated from enum bytebase.v1.ExportFormat
 */
export const ExportFormat = /*@__PURE__*/
  tsEnum(ExportFormatSchema);

