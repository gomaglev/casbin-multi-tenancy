[WIP]

# Multiple tenancy casbin use case sample

A sample from a running system.
It will be migrated to gRPC&gRPC gateway same as
Protomicro in this repo. After migrated, it will be used as one
of the authorization/authentication selection for
dev & unit test.

Casbin basic settings and usage in Gin.
It's not a good idea to use it in every microservice.
Because authorization and authentication code
should be avoid in each microservice service.

In Protomicro, authorization and authentication will be added to
k8s/istio layer when generating source code using cli.

## Main files

### model.conf - how to use `r = sub, dom, obj, act`

### config.toml - only casbin related configurations

### i18n/*toml - why? only if authentication / authorization errors need to be shown in multiple languages
