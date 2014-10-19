# (WIP) Etcd-Registrator

Etcd-Registrator automatically register/deregisters services for Docker containers based on published ports and metadata from the container environment. It is intended to be used with HaMesh.

## Starting Etcd-Registrator

Registrator assumes the default Docker socket at `file:///var/run/docker.sock` or you can override it with `DOCKER_HOST`. The only  mandatory argument is a registry URI, which specifies and configures the registry backend to use.

	$ registrator <registry-uri>

By default, when registering a service, registrator will assign the service address by attempting to resolve the current hostname. If you would like to force the service address to be a specific address, you can specify the `-ip` argument.

Etcd-Registrator was designed to just be run as a container. You must pass the Docker socket file as a mount to `/tmp/docker.sock`, and it's a good idea to set the hostname to the machine host:

	$ docker run -d \
		-v /var/run/docker.sock:/tmp/docker.sock \
		-h $HOSTNAME progrium/registrator <registry-uri>

### Registry URIs

The registry backend to use is defined by a URI. The scheme is the supported registry name, and an address. Registries based on key-value stores like etcd and Zookeeper (not yet supported) can specify a key path to use to prefix service definitions. Registries may also use query params for other options. See also [Adding support for other service registries](#adding-support-for-other-service-registries).

#### Etcd Key-value Store

	$ registrator etcd:///path/to/services
	$ registrator etcd://192.168.1.100/services

Service definitions are stored as:

	<registry-uri-path>/<service-name>/<service-id>/address = <ip>
	<registry-uri-path>/<service-name>/<service-id>/port = <port>
	<registry-uri-path>/<service-name>/<service-id>/mode = <mode>
	<registry-uri-path>/<service-name>/<service-id>/domains/<domain N> = <domain N> (http mode only)
	
### Public-facing HTTP service with multiple domains

	$ docker run -d --name nginx.0 -P \
		-e "SERVICE_MODE=http" \
		-e "SERVICE_DOMAINS=awesome.website.com,other.website.com" \

Will bind on 0.0.0.0 and port 80 for all given domains.

### Public-facing TCP service

	$ docker run -d --name nginx.1 -P \
		-e "SERVICE_MODE=tcp" \
		-e "SERVICE_EXTERNAL_PORT=8080" \

Will bind on 0.0.0.0 and port 8080, port must be unique.

### Internal TCP service

	$ docker run -d --name nginx.2 -P \
		-e "SERVICE_MODE=tcp" \
		-e "SERVICE_INTERNAL_PORT=8090" \

Will bind on the private network and port 8090, port must be unique.

## License

BSD