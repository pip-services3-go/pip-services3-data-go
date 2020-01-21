# Configuration Guide <br/>

Configuration structure follows the 
[standard configuration](https://github.com/pip-services/pip-services3-container-node/doc/Configuration.md) 
structure. 


### <a name="persistence_memory"></a> Memory

Memory persistence has the following configuration properties:
- options: object - Misc configuration options
  - max_page_size: number - Maximum number of items per page (default: 100)

Example:
```yaml
- descriptor: "pip-services-clusters:persistence:memory:default:1.0"
  options:
    max_page_size: 100
```

For more information on this section read 
[Pip.Services Configuration Guide](https://github.com/pip-services/pip-services3-container-node/doc/Configuration.md#deps)