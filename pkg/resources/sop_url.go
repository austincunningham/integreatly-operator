package resources

// Add alert strings here
const (
	sopUrlRhoamBase                                           = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/"
	sopUrlPostgresInstanceUnavailable                         = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/postgres_instance_unavailable.asciidoc"
	sopUrlPostgresConnectionFailed                            = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/postgres_connection_failed.asciidoc"
	sopUrlRedisCacheUnavailable                               = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/redis_cache_unavailable.asciidoc"
	sopUrlRedisConnectionFailed                               = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/redis_connection_failed.asciidoc"
	sopUrlPostgresResourceStatusPhasePending                  = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/postgres_resource_status_phase_pending.asciidoc"
	sopUrlPostgresResourceStatusPhaseFailed                   = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/postgres_resource_status_phase_failed.asciidoc"
	sopUrlRedisResourceStatusPhasePending                     = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/redis_resource_status_phase_pending.asciidoc"
	sopUrlRedisResourceStatusPhaseFailed                      = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/redis_resource_status_phase_failed.asciidoc"
	sopUrlPostgresWillFill                                    = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/postgres_storage_alerts.asciidoc"
	sopUrlPostgresFreeableMemoryLow                           = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/postgres_freeable_memory_low.asciidoc"
	sopUrlRedisMemoryUsageHigh                                = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/redis_memory_usage_high.asciidoc"
	sopUrlPostgresCpuUsageHigh                                = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/postgres_cpu_usage_high.asciidoc"
	sopUrlRedisCpuUsageHigh                                   = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/redis_cpu_usage_high.asciidoc"
	sopUrlRedisServiceMaintenanceCritical                     = "No SOP required, notify engineering of required update"
	SopUrlEndpointAvailableAlert                              = "https://github.com/RHCloudServices/integreatly-help/tree/master/sops/2.x/alerts/service_endpoint_down.asciidoc"
	SopUrlAlertsAndTroubleshooting                            = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/alerts_and_troubleshooting.md"
	sopUrlCloudResourceDeletionStatusFailed                   = "https://github.com/RHCloudServices/integreatly-help/tree/master/sops/2.x/alerts/clean_up_cloud_resources_failed_teardown.asciidoc"
	sopUrlSendGridSmtpSecretExists                            = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/2.x/alerts/sendgrid_smtp_secret_not_present.asciidoc"
	SopUrlMarin3rEnvoyApicastProductionContainerDown          = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/Marin3rEnvoyApicastProductionContainerDown.asciidoc"
	SopUrlMarin3rEnvoyApicastStagingContainerDown             = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/Marin3rEnvoyApicastStagingContainerDown.asciidoc"
	SopUrlOperatorInstallDelayed                              = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/OperatorInstallDelayed.asciidoc"
	SopUrlUpgradeExpectedDurationExceeded                     = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/UpgradeExpectedDurationExceeded.asciidoc"
	SopUrlRHMICloudResourceOperatorMetricsServiceEndpointDown = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/RHMICloudResourceOperatorMetricsServiceEndpointDown.asciidoc"
	SopUrlRHMIThreeScaleApicastProductionServiceEndpointDown  = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/RHMIThreeScaleApicastProductionServiceEndpointDown.asciidoc"
	SopUrlRHMIThreeScaleApicastStagingServiceEndpointDown     = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/RHMIThreeScaleApicastStagingServiceEndpointDown.asciidoc"
	SopUrlRHMIThreeScaleBackendListenerServiceEndpointDown    = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/RHMIThreeScaleBackendListenerServiceEndpointDown.asciidoc"
	SopUrlRHMIThreeScaleZyncServiceEndpointDown               = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/RHMIThreeScaleZyncServiceEndpointDown.asciidoc"
	SopUrlRHMIThreeScaleZyncDatabaseServiceEndpointDown       = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/RHMIThreeScaleZyncDatabaseServiceEndpointDown.asciidoc"
	SopUrlThreeScaleBackendWorkerPod                          = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/ThreeScaleBackendWorkerPod.asciidoc"
	SopUrlThreeScaleAdminUIBBT                                = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/ThreeScaleAdminUIBBT.asciidoc"
	SopUrlThreeScaleDeveloperUIBBT                            = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/ThreeScaleDeveloperUIBBT.asciidoc"
	SopUrlThreeScaleSystemAdminUIBBT                          = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/ThreeScaleSystemAdminUIBBT.asciidoc"
	SopUrlPodDistributionIncorrect                            = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/multi-az/pod_distribution.md"
	SopUrlSloRhssoAvailabilityAlert                           = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/SloRhssoAvailabilityAlert.asciidoc"
	SopUrlSloUserSsoAvailabilityAlert                         = "https://github.com/RHCloudServices/integreatly-help/blob/master/sops/rhoam/alerts/SloUserSsoAvailabilityAlert.asciidoc"
)
