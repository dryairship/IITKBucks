FROM node:10-alpine AS frontend
WORKDIR /iitkbucks
COPY frontend ./
RUN npm i &&\
    npm run build

FROM golang:alpine AS backend
WORKDIR /iitkbucks
COPY . .
RUN go build ./cmd/iitkbucks

FROM alpine
WORKDIR /iitkbucks
RUN mkdir --p /iitkbucks/build &&\
    mkdir --p /iitkbucks/blocks
COPY --from=backend /iitkbucks/iitkbucks ./
COPY --from=frontend /iitkbucks/build ./build

EXPOSE 8000
ENTRYPOINT ["./iitkbucks"]
