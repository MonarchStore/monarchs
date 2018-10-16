Monarchs Helm Chart
===================


Installation
------------


### Use Case 1: Cluster-Only Service

In this case, monarchs will be available only from within the cluster (or via kubectl proxy)

    helm upgrade --install monarchs ./chart/monarchs


### Use Case 2: Expose via Ingress Controller

This example assumes you have the nginx ingress controller configured already.

    helm upgrade --install monarchs \
    	 --set ingress.url=monarchs.mycluster.example.org \
	 ./chart/monarchs


### Use Case 3: Expose via Internal ELB + Set domain with external-dns

This example assumes you have external-dns configured in your cluster, and the Hosted Zone `example.org` is configured
for private VPC access

    helm upgrade --install monarchs \
	 --set external.url=monarchs.example.org \
	 --set elb.type=internal \
	 ./chart/monarchs


### Use Case 4: Expose via External ELB + DNS + SSL on 443

    helm upgrade --install monarchs \
    	 --set elb.type=external \
	 --set elb.ssl_cert="arn:aws:acm:<region>:<accountID>:certificate/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx" \
	 --set external.url=monarchs.example.org \
	 --set service.externalPort=443 \
	 ./chart/monarchs


Values
------

| Value                     | Type   | Default          | Description                                                                |
|---------------------------|--------|------------------|----------------------------------------------------------------------------|
| image.repository          | string | arturom/monarchs |                                                                            |
| image.tag                 | string | latest           |                                                                            |
| image.pullPolicy          | string | Always           |                                                                            |
| service.name              | string | monarchs         | Kubernetes Service Name                                                    |
| service.type              | string | ClusterIP        | Kubernetes Service Type                                                    |
| service.externalPort      | int    | 6789             | The port the service will actually be available on                         |
| service.internalPort      | int    | 6789             | Probably shouldn't change this                                             |
| service.extra_annotations | map    | {}               | Extra annotations for the service                                          |
| ingress.url               | string | ''               | The host for ingress-nginx, if desired                                     |
| ingress.class             | string | nginx            | The Kubernetes Ingress class                                               |
| external.url              | string | ''               | The hostname for external-dns, if desired                                  |
| elb.type                  | string | none             | If ELB is desired, set to 'internal' or 'external'. Overrides service.type |
| elb.ssl_cert              | strng  | none             | The ARN of the AWS SSL/TLS Certificate                                     |
