FROM scratch

ADD bin/surtr surtr

ENTRYPOINT ["/surtr"]
