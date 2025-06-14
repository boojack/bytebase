// @generated by protoc-gen-es v2.5.2
// @generated from file v1/idp_service.proto (package bytebase.v1, syntax proto3)
/* eslint-disable */

import { enumDesc, fileDesc, messageDesc, serviceDesc, tsEnum } from "@bufbuild/protobuf/codegenv2";
import { file_google_api_annotations } from "../google/api/annotations_pb";
import { file_google_api_client } from "../google/api/client_pb";
import { file_google_api_field_behavior } from "../google/api/field_behavior_pb";
import { file_google_api_resource } from "../google/api/resource_pb";
import { file_google_protobuf_empty, file_google_protobuf_field_mask } from "@bufbuild/protobuf/wkt";
import { file_v1_annotation } from "./annotation_pb";

/**
 * Describes the file v1/idp_service.proto.
 */
export const file_v1_idp_service = /*@__PURE__*/
  fileDesc("ChR2MS9pZHBfc2VydmljZS5wcm90bxILYnl0ZWJhc2UudjEiRQoaR2V0SWRlbnRpdHlQcm92aWRlclJlcXVlc3QSJwoEbmFtZRgBIAEoCUIZ4kEBAvpBEgoQYnl0ZWJhc2UuY29tL0lkUCJFChxMaXN0SWRlbnRpdHlQcm92aWRlcnNSZXF1ZXN0EhEKCXBhZ2Vfc2l6ZRgBIAEoBRISCgpwYWdlX3Rva2VuGAIgASgJInMKHUxpc3RJZGVudGl0eVByb3ZpZGVyc1Jlc3BvbnNlEjkKEmlkZW50aXR5X3Byb3ZpZGVycxgBIAMoCzIdLmJ5dGViYXNlLnYxLklkZW50aXR5UHJvdmlkZXISFwoPbmV4dF9wYWdlX3Rva2VuGAIgASgJIpQBCh1DcmVhdGVJZGVudGl0eVByb3ZpZGVyUmVxdWVzdBI+ChFpZGVudGl0eV9wcm92aWRlchgBIAEoCzIdLmJ5dGViYXNlLnYxLklkZW50aXR5UHJvdmlkZXJCBOJBAQISHAoUaWRlbnRpdHlfcHJvdmlkZXJfaWQYAiABKAkSFQoNdmFsaWRhdGVfb25seRgDIAEoCCKQAQodVXBkYXRlSWRlbnRpdHlQcm92aWRlclJlcXVlc3QSPgoRaWRlbnRpdHlfcHJvdmlkZXIYASABKAsyHS5ieXRlYmFzZS52MS5JZGVudGl0eVByb3ZpZGVyQgTiQQECEi8KC3VwZGF0ZV9tYXNrGAIgASgLMhouZ29vZ2xlLnByb3RvYnVmLkZpZWxkTWFzayJICh1EZWxldGVJZGVudGl0eVByb3ZpZGVyUmVxdWVzdBInCgRuYW1lGAEgASgJQhniQQEC+kESChBieXRlYmFzZS5jb20vSWRQIrMBChtUZXN0SWRlbnRpdHlQcm92aWRlclJlcXVlc3QSOAoRaWRlbnRpdHlfcHJvdmlkZXIYASABKAsyHS5ieXRlYmFzZS52MS5JZGVudGl0eVByb3ZpZGVyEk8KDm9hdXRoMl9jb250ZXh0GAIgASgLMjUuYnl0ZWJhc2UudjEuT0F1dGgySWRlbnRpdHlQcm92aWRlclRlc3RSZXF1ZXN0Q29udGV4dEgAQgkKB2NvbnRleHQiOAooT0F1dGgySWRlbnRpdHlQcm92aWRlclRlc3RSZXF1ZXN0Q29udGV4dBIMCgRjb2RlGAEgASgJIh4KHFRlc3RJZGVudGl0eVByb3ZpZGVyUmVzcG9uc2UizgEKEElkZW50aXR5UHJvdmlkZXISDAoEbmFtZRgBIAEoCRINCgV0aXRsZRgEIAEoCRIOCgZkb21haW4YBSABKAkSLwoEdHlwZRgGIAEoDjIhLmJ5dGViYXNlLnYxLklkZW50aXR5UHJvdmlkZXJUeXBlEjMKBmNvbmZpZxgHIAEoCzIjLmJ5dGViYXNlLnYxLklkZW50aXR5UHJvdmlkZXJDb25maWc6IepBHgoQYnl0ZWJhc2UuY29tL0lkUBIKaWRwcy97aWRwfUoECAIQAyLmAQoWSWRlbnRpdHlQcm92aWRlckNvbmZpZxJCCg1vYXV0aDJfY29uZmlnGAEgASgLMikuYnl0ZWJhc2UudjEuT0F1dGgySWRlbnRpdHlQcm92aWRlckNvbmZpZ0gAEj4KC29pZGNfY29uZmlnGAIgASgLMicuYnl0ZWJhc2UudjEuT0lEQ0lkZW50aXR5UHJvdmlkZXJDb25maWdIABI+CgtsZGFwX2NvbmZpZxgDIAEoCzInLmJ5dGViYXNlLnYxLkxEQVBJZGVudGl0eVByb3ZpZGVyQ29uZmlnSABCCAoGY29uZmlnIpECChxPQXV0aDJJZGVudGl0eVByb3ZpZGVyQ29uZmlnEhAKCGF1dGhfdXJsGAEgASgJEhEKCXRva2VuX3VybBgCIAEoCRIVCg11c2VyX2luZm9fdXJsGAMgASgJEhEKCWNsaWVudF9pZBgEIAEoCRIVCg1jbGllbnRfc2VjcmV0GAUgASgJEg4KBnNjb3BlcxgGIAMoCRIwCg1maWVsZF9tYXBwaW5nGAcgASgLMhkuYnl0ZWJhc2UudjEuRmllbGRNYXBwaW5nEhcKD3NraXBfdGxzX3ZlcmlmeRgIIAEoCBIwCgphdXRoX3N0eWxlGAkgASgOMhwuYnl0ZWJhc2UudjEuT0F1dGgyQXV0aFN0eWxlIoACChpPSURDSWRlbnRpdHlQcm92aWRlckNvbmZpZxIOCgZpc3N1ZXIYASABKAkSEQoJY2xpZW50X2lkGAIgASgJEhUKDWNsaWVudF9zZWNyZXQYAyABKAkSDgoGc2NvcGVzGAQgAygJEjAKDWZpZWxkX21hcHBpbmcYBSABKAsyGS5ieXRlYmFzZS52MS5GaWVsZE1hcHBpbmcSFwoPc2tpcF90bHNfdmVyaWZ5GAYgASgIEjAKCmF1dGhfc3R5bGUYByABKA4yHC5ieXRlYmFzZS52MS5PQXV0aDJBdXRoU3R5bGUSGwoNYXV0aF9lbmRwb2ludBgIIAEoCUIE4kEBAyL3AgoaTERBUElkZW50aXR5UHJvdmlkZXJDb25maWcSDAoEaG9zdBgBIAEoCRIMCgRwb3J0GAIgASgFEhcKD3NraXBfdGxzX3ZlcmlmeRgDIAEoCBIPCgdiaW5kX2RuGAQgASgJEhUKDWJpbmRfcGFzc3dvcmQYBSABKAkSDwoHYmFzZV9kbhgGIAEoCRITCgt1c2VyX2ZpbHRlchgHIAEoCRJTChFzZWN1cml0eV9wcm90b2NvbBgIIAEoDjI4LmJ5dGViYXNlLnYxLkxEQVBJZGVudGl0eVByb3ZpZGVyQ29uZmlnLlNlY3VyaXR5UHJvdG9jb2wSMAoNZmllbGRfbWFwcGluZxgJIAEoCzIZLmJ5dGViYXNlLnYxLkZpZWxkTWFwcGluZyJPChBTZWN1cml0eVByb3RvY29sEiEKHVNFQ1VSSVRZX1BST1RPQ09MX1VOU1BFQ0lGSUVEEAASDQoJU1RBUlRfVExTEAESCQoFTERBUFMQAiJdCgxGaWVsZE1hcHBpbmcSEgoKaWRlbnRpZmllchgBIAEoCRIUCgxkaXNwbGF5X25hbWUYAiABKAkSDQoFcGhvbmUYBCABKAkSDgoGZ3JvdXBzGAUgASgJSgQIAxAEKl4KFElkZW50aXR5UHJvdmlkZXJUeXBlEiYKIklERU5USVRZX1BST1ZJREVSX1RZUEVfVU5TUEVDSUZJRUQQABIKCgZPQVVUSDIQARIICgRPSURDEAISCAoETERBUBADKlIKD09BdXRoMkF1dGhTdHlsZRIhCh1PQVVUSDJfQVVUSF9TVFlMRV9VTlNQRUNJRklFRBAAEg0KCUlOX1BBUkFNUxABEg0KCUlOX0hFQURFUhACMr0IChdJZGVudGl0eVByb3ZpZGVyU2VydmljZRKfAQoTR2V0SWRlbnRpdHlQcm92aWRlchInLmJ5dGViYXNlLnYxLkdldElkZW50aXR5UHJvdmlkZXJSZXF1ZXN0Gh0uYnl0ZWJhc2UudjEuSWRlbnRpdHlQcm92aWRlciJA2kEEbmFtZYrqMBhiYi5pZGVudGl0eVByb3ZpZGVycy5nZXSQ6jABgtPkkwITEhEvdjEve25hbWU9aWRwcy8qfRKHAQoVTGlzdElkZW50aXR5UHJvdmlkZXJzEikuYnl0ZWJhc2UudjEuTGlzdElkZW50aXR5UHJvdmlkZXJzUmVxdWVzdBoqLmJ5dGViYXNlLnYxLkxpc3RJZGVudGl0eVByb3ZpZGVyc1Jlc3BvbnNlIhfaQQCA6jABgtPkkwIKEggvdjEvaWRwcxKyAQoWQ3JlYXRlSWRlbnRpdHlQcm92aWRlchIqLmJ5dGViYXNlLnYxLkNyZWF0ZUlkZW50aXR5UHJvdmlkZXJSZXF1ZXN0Gh0uYnl0ZWJhc2UudjEuSWRlbnRpdHlQcm92aWRlciJN2kEAiuowG2JiLmlkZW50aXR5UHJvdmlkZXJzLmNyZWF0ZZDqMAGY6jABgtPkkwIdOhFpZGVudGl0eV9wcm92aWRlciIIL3YxL2lkcHMS6wEKFlVwZGF0ZUlkZW50aXR5UHJvdmlkZXISKi5ieXRlYmFzZS52MS5VcGRhdGVJZGVudGl0eVByb3ZpZGVyUmVxdWVzdBodLmJ5dGViYXNlLnYxLklkZW50aXR5UHJvdmlkZXIihQHaQR1pZGVudGl0eV9wcm92aWRlcix1cGRhdGVfbWFza4rqMBtiYi5pZGVudGl0eVByb3ZpZGVycy51cGRhdGWQ6jABmOowAYLT5JMCODoRaWRlbnRpdHlfcHJvdmlkZXIyIy92MS97aWRlbnRpdHlfcHJvdmlkZXIubmFtZT1pZHBzLyp9EqUBChZEZWxldGVJZGVudGl0eVByb3ZpZGVyEiouYnl0ZWJhc2UudjEuRGVsZXRlSWRlbnRpdHlQcm92aWRlclJlcXVlc3QaFi5nb29nbGUucHJvdG9idWYuRW1wdHkiR9pBBG5hbWWK6jAbYmIuaWRlbnRpdHlQcm92aWRlcnMuZGVsZXRlkOowAZjqMAGC0+STAhMqES92MS97bmFtZT1pZHBzLyp9EqoBChRUZXN0SWRlbnRpdHlQcm92aWRlchIoLmJ5dGViYXNlLnYxLlRlc3RJZGVudGl0eVByb3ZpZGVyUmVxdWVzdBopLmJ5dGViYXNlLnYxLlRlc3RJZGVudGl0eVByb3ZpZGVyUmVzcG9uc2UiPYrqMBtiYi5pZGVudGl0eVByb3ZpZGVycy51cGRhdGWQ6jABgtPkkwIUOgEqIg8vdjEvaWRwcy8qOnRlc3RCNFoyZ2l0aHViLmNvbS9ieXRlYmFzZS9ieXRlYmFzZS9wcm90by9nZW5lcmF0ZWQtZ28vdjFiBnByb3RvMw", [file_google_api_annotations, file_google_api_client, file_google_api_field_behavior, file_google_api_resource, file_google_protobuf_empty, file_google_protobuf_field_mask, file_v1_annotation]);

/**
 * Describes the message bytebase.v1.GetIdentityProviderRequest.
 * Use `create(GetIdentityProviderRequestSchema)` to create a new message.
 */
export const GetIdentityProviderRequestSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 0);

/**
 * Describes the message bytebase.v1.ListIdentityProvidersRequest.
 * Use `create(ListIdentityProvidersRequestSchema)` to create a new message.
 */
export const ListIdentityProvidersRequestSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 1);

/**
 * Describes the message bytebase.v1.ListIdentityProvidersResponse.
 * Use `create(ListIdentityProvidersResponseSchema)` to create a new message.
 */
export const ListIdentityProvidersResponseSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 2);

/**
 * Describes the message bytebase.v1.CreateIdentityProviderRequest.
 * Use `create(CreateIdentityProviderRequestSchema)` to create a new message.
 */
export const CreateIdentityProviderRequestSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 3);

/**
 * Describes the message bytebase.v1.UpdateIdentityProviderRequest.
 * Use `create(UpdateIdentityProviderRequestSchema)` to create a new message.
 */
export const UpdateIdentityProviderRequestSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 4);

/**
 * Describes the message bytebase.v1.DeleteIdentityProviderRequest.
 * Use `create(DeleteIdentityProviderRequestSchema)` to create a new message.
 */
export const DeleteIdentityProviderRequestSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 5);

/**
 * Describes the message bytebase.v1.TestIdentityProviderRequest.
 * Use `create(TestIdentityProviderRequestSchema)` to create a new message.
 */
export const TestIdentityProviderRequestSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 6);

/**
 * Describes the message bytebase.v1.OAuth2IdentityProviderTestRequestContext.
 * Use `create(OAuth2IdentityProviderTestRequestContextSchema)` to create a new message.
 */
export const OAuth2IdentityProviderTestRequestContextSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 7);

/**
 * Describes the message bytebase.v1.TestIdentityProviderResponse.
 * Use `create(TestIdentityProviderResponseSchema)` to create a new message.
 */
export const TestIdentityProviderResponseSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 8);

/**
 * Describes the message bytebase.v1.IdentityProvider.
 * Use `create(IdentityProviderSchema)` to create a new message.
 */
export const IdentityProviderSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 9);

/**
 * Describes the message bytebase.v1.IdentityProviderConfig.
 * Use `create(IdentityProviderConfigSchema)` to create a new message.
 */
export const IdentityProviderConfigSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 10);

/**
 * Describes the message bytebase.v1.OAuth2IdentityProviderConfig.
 * Use `create(OAuth2IdentityProviderConfigSchema)` to create a new message.
 */
export const OAuth2IdentityProviderConfigSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 11);

/**
 * Describes the message bytebase.v1.OIDCIdentityProviderConfig.
 * Use `create(OIDCIdentityProviderConfigSchema)` to create a new message.
 */
export const OIDCIdentityProviderConfigSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 12);

/**
 * Describes the message bytebase.v1.LDAPIdentityProviderConfig.
 * Use `create(LDAPIdentityProviderConfigSchema)` to create a new message.
 */
export const LDAPIdentityProviderConfigSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 13);

/**
 * Describes the enum bytebase.v1.LDAPIdentityProviderConfig.SecurityProtocol.
 */
export const LDAPIdentityProviderConfig_SecurityProtocolSchema = /*@__PURE__*/
  enumDesc(file_v1_idp_service, 13, 0);

/**
 * @generated from enum bytebase.v1.LDAPIdentityProviderConfig.SecurityProtocol
 */
export const LDAPIdentityProviderConfig_SecurityProtocol = /*@__PURE__*/
  tsEnum(LDAPIdentityProviderConfig_SecurityProtocolSchema);

/**
 * Describes the message bytebase.v1.FieldMapping.
 * Use `create(FieldMappingSchema)` to create a new message.
 */
export const FieldMappingSchema = /*@__PURE__*/
  messageDesc(file_v1_idp_service, 14);

/**
 * Describes the enum bytebase.v1.IdentityProviderType.
 */
export const IdentityProviderTypeSchema = /*@__PURE__*/
  enumDesc(file_v1_idp_service, 0);

/**
 * @generated from enum bytebase.v1.IdentityProviderType
 */
export const IdentityProviderType = /*@__PURE__*/
  tsEnum(IdentityProviderTypeSchema);

/**
 * Describes the enum bytebase.v1.OAuth2AuthStyle.
 */
export const OAuth2AuthStyleSchema = /*@__PURE__*/
  enumDesc(file_v1_idp_service, 1);

/**
 * @generated from enum bytebase.v1.OAuth2AuthStyle
 */
export const OAuth2AuthStyle = /*@__PURE__*/
  tsEnum(OAuth2AuthStyleSchema);

/**
 * @generated from service bytebase.v1.IdentityProviderService
 */
export const IdentityProviderService = /*@__PURE__*/
  serviceDesc(file_v1_idp_service, 0);

