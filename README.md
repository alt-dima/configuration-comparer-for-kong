# Configuration comparer for Kong API Gateway
Simple Go application to compare two Kong configurations via Go bindings for Kong's Admin API.

**Written for learing Go** and demonstrate how to use Go bindings for Kong's Admin API (`github.com/kong/go-kong`)

Application compare configuration tree one Kong API Gateway to another ignoring different IDs of the resources.

- Routes (it search routes in the second Kong by all the pathes in the Route of the original Kong)
    - Route configuration (Preserve Host, Strip Path should be the same)
    - Route plugins (that plugins exist and configuration are the same)
    - Target Service
        - Service Configuration (Host, Port, Path should be the same)
        - Service plugins (that plugins exist and configuration are the same)
- Consumers (consumer with the same consumer name should exist)
    - Consumer plugins (that plugins exist and configuration are the same)
    - Consumers credentials comparing **not implemented!**

Plugins configuration comparsion works very primitive (comparing of JSON of the config), it causes false positives in configuration mismatch because of some possible unique fields in the configs (like anonymous consumer ID), thats why I do "hardcoded things" like `delete(pluginRouteClient1.Config, "anonymous")`

## Use

`./configuration-comparer-for-kong https://kong-admin-api.example.com   https://kong-admin-api.example2.com`

Example output:
```
Amount routes in https://kong-admin-api.example.com: 22 
Amount routes in https://kong-admin-api.example2.com: 22 
Service service-kafka-rest-1 plugin config post-function not equals in https://kong-admin-api.example2.com
Service service-kafka-rest-2 plugin config post-function not equals in https://kong-admin-api.example2.com
Service service-kafka-rest-3 plugin config post-function not equals in https://kong-admin-api.example2.com
Service service-kafka-rest-4 plugin config post-function not equals in https://kong-admin-api.example2.com
Service service-kafka-rest-5 plugin config post-function not equals in https://kong-admin-api.example2.com
Amount consumers in https://kong-admin-api.example.com: 20 
Amount consumers in https://kong-admin-api.example2.com: 18 
Consumers credentials comparing not implemented!!! 
Consumer trytryrty@example.com does not exists in https://kong-admin-api.example2.com
Consumer bvnvbnvbn@example.com does not exists in https://kong-admin-api.example2.com
Consumer adfadf@example.com does not exists in https://kong-admin-api.example2.com
```

## Compatibility

`configuration-comparer-for-kong` is compatible the same as `github.com/kong/go-kong` module (currently with with Kong 2.x and 3.x.)

## License

`configuration-comparer-for-kong` is licensed with Apache License Version 2.0.
Please read the LICENSE file for more details.
