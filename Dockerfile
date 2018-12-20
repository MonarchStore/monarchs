FROM library/golang:1.6 as builder

ARG commit="none"

WORKDIR /go/src/github.com/MonarchStore/monarchs

COPY . .

RUN make dep
RUN make test
RUN make build


#RUN go test -v -race $(go list ./... | grep -v /vendor/)
#RUN go build -v .


FROM library/golang:1.6-alpine
COPY --from=builder /go/src/github.com/MonarchStore/monarchs /go/bin/monarchs

ENV MONARCHS_ADDR "0.0.0.0"
ENV MONARCHS_PORT "6789"
ENV MONARCHS_LOG_LEVEL "info"

LABEL org.label-schema.schema-version "1.0.0"
LABEL org.label-schema.version $commit
LABEL org.label-schema.name "monarchs"
LABEL org.label-schema.description "A hierarchial, NoSQL, in-memory data store with a RESTful API"
LABEL org.label-schema.vcs-url "https://github.com/MonarchStore/monarchs"

USER nobody

ENTRYPOINT ["/go/bin/monarchs"]
