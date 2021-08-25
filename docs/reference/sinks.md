# Sinks

## Console

`console`

Print data to stdout.

### Sample usage of console sink

```yaml
sinks:
 - name: console
```

## Columbus

`columbus`

Upload metadata to a given `type` in [Columbus](https://github.com/odpf/meteor/tree/cb12c3ecf8904cf3f4ce365ca8981ccd132f35d0/docs/reference/github.com/odpf/columbus/README.md). Request will be send via HTTP to given host.

### Sample usage of columbus sink

```yaml
sinks:
 - name: columbus
   config:
     host: https://columbus.com
     type: sample-columbus-type
     mapping:
       ModelFieldName: "new_fieldname"
       Urn: "id"
       Name: "displayName"
```

_**Notes**_

Columbus' Type requires certain fields to be sent, hence why `mapping` config is needed to map value from any of our metadata models to any field name when sending to Columbus.
