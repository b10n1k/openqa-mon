default: all
all: openqa-mon

## XXX: This is ugly. Switch to go build system
ifeq ($(DESTDIR),)
DESTDIR=/usr/local
endif



openqa-mon: cmd/openqa-mon/openqa-mon.go cmd/openqa-mon/terminal.go cmd/openqa-mon/jobs.go cmd/openqa-mon/config.go
	go build $^

install: openqa-mon
	install openqa-mon ${DESTDIR}/bin/
	install doc/openqa-mon.8 ${DESTDIR}/man/man8/
uninstall:
	rm -f ${DESTDIR}/bin/bin/openqa-mon
	rm -f ${DESTDIR}/man/man8/openqa-mon.8
