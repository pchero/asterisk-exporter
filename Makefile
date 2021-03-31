TARGET = asterisk-exporter
VERSION = 0.0.2

TARGET_RELEASE = $(TARGET)-$(VERSION)

build:
	./build.sh $(TARGET_RELEASE)

clean:
	rm $(TARGET)*
