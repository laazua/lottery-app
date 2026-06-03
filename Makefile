apk := lottery.apk

.PHONY: clean, build

build:
	gogio -target android -appid com.lottery.app -icon lottery.png -o $(apk) .

clean:
	rm -f $(apk)*
