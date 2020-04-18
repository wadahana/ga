
tags := vp8enc,memu
tags := $(strip $(tags))

ga.zip: clean ga
	@zip -r ga.zip html ga.exe

ga:
	go build -o ./ga.exe -tags "$(tags)" ./main.go 

.PHONY: clean
clean:
	@if [ -f ga ]; then rm ga; fi
	@if [ -f ga.zip ]; then rm ga.zip ; fi
