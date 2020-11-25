export GOTRACEBACK=all

all: ../nopanic.raw ../simplepanic.raw ../concurrent.raw ../recovered.raw

.PHONY: simplepanic
simplepanic:
	go build -o $@ -ldflags="-s -w" -trimpath .

../nopanic.raw: simplepanic
	-./$< -mode 0 >$@ 2>&1 || :

../simplepanic.raw: simplepanic
	-./$< -mode 1 >$@ 2>&1 || :

../concurrent.raw: simplepanic
	-./$< -mode 2 >$@ 2>&1 || :

../recovered.raw: simplepanic
	-./$< -mode 3 >$@ 2>&1 || :
