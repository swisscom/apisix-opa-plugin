# apisix-opa-plugin

An Open Policy Agent plugin for [Apache APISIX](https://apisix.apache.org/)


## Remarks

This plugin requires the `OPA_URL` environment variable to be set. As per 
[external-plugin FAQ](https://apisix.apache.org/docs/apisix/external-plugin/#when-managing-by-apisix-the-runner-cant-access-my-environment-variable)
the `nginx_config` configuration must be set to forward this environment variable to the runner,
otherwise it will be hidden.
  
To do so, the APISIX configuration should include a similar snippet of configuration:

```yaml
nginx_config:
  envs:
  - OPA_URL
```

## Plugin Configuration

When setting up the plugin on APISIX, provide a configuration
similar to the following:
```json
{
  "conf": [
    {
      "name": "opa",
      "value": "{\"rule_path\": \"com.swisscom.ini.dna.nb.example/rule1\"}"
    }
  ],
  "disable": false
}
```
