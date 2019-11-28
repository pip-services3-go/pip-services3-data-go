# <img src="https://github.com/pip-services/pip-services/raw/master/design/Logo.png" alt="Pip.Services Logo" style="max-width:30%"> <br/> IoC container for Golang

This framework is a part of the [Pip.Services](https://github.com/pip-services/pip-services) project.
It provides an inversion-of-control component container to facilitate the development of composable services and applications.

As all Pip.Services projects this framework is implemented in a variety of different languages: Java, .NET, Python, Node.js, Golang. 

The framework provides a light-weight container that can be embedded inside a service or application, or can be run by itself,
as a system process, for example. Container configuration serves as a recipe for instantiating and configuring components inside the container.  
The default container factory provides generic functionality on-demand, such as logging and performance monitoring.

- [**Conteiner**](https://godoc.org/github.com/pip-services3-go/pip-services3-container-go/container) - Component container and container as a system process
- [**Build**](https://godoc.org/github.com/pip-services3-go/pip-services3-container-go/build) - Container default factory
- [**Config**](https://godoc.org/github.com/pip-services3-go/pip-services3-container-go/config) - Container configuration
- [**Refer**](https://godoc.org/github.com/pip-services3-go/pip-services3-container-go/refer) - Container references

Quick Links:

* [Downloads](https://github.com/pip-services3-go/pip-services3-container-go.git)
* [API Reference](https://godoc.org/github.com/pip-services3-go/pip-services3-container-go/)
* [Building and Testing](https://github.com/pip-services3-go/pip-services3-container-go/blob/master/docs/Development.md)
* [Contributing](https://github.com/pip-services3-go/pip-services3-container-go/blob/master/docs/Development.md#contrib)

## Acknowledgements

The Golang version of Pip.Services is created and maintained by:
- **Volodymyr Tkachenko**
- **Sergey Seroukhov**
- **Mark Zontak**

The documentation is written by:
- **Levichev Dmitry**