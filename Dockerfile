FROM containerize/godep

ADD . /go/src/imgprox

RUN godep restore && go build -o imgprox main.go

ENTRYPOINT ["./imgprox"]