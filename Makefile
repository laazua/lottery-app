apk := lottery.apk

.PHONY: clear, build

build:
	gogio -target android -appid com.lottery.app -icon lottery.png -o $(apk) .

clear:
	rm -f $(apk)*
