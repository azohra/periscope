FROM golang:alpine AS build-env
RUN apk --no-cache add ca-certificates build-base git bzr mercurial gcc
RUN go get -v k8s.io/client-go/...
RUN go get k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1
ADD ./ ./src/github.com/Ascendance/periscope
RUN cd src/github.com/Ascendance/periscope && go install ./...

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build-env /go/bin/periscope /app/periscope
ENTRYPOINT ./periscope
EXPOSE 7000