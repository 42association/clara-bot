#include <FastLED.h>
#include <M5StickCPlus2.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include "../lib/secrets.h"

uint8_t ledColor = 0;
#define KEY_PIN 33
#define DATA_PIN 32
CRGB LED[1];

void setup()
{
	M5.begin();
	WiFi.begin(ssid, password);
	while (WiFi.status() != WL_CONNECTED)
	{
		delay(500);
		M5.Lcd.print(".");
	}
	M5.Lcd.fillScreen(BLACK);
	M5.Lcd.setRotation(3);
	M5.Lcd.setTextSize(4);
	pinMode(KEY_PIN, INPUT_PULLUP);
	FastLED.addLeds<SK6812, DATA_PIN, GRB>(LED, 1);
	LED[0] = CRGB::Blue;
	FastLED.setBrightness(0);
	M5.Speaker.begin();
	M5.Speaker.setVolume(225);
}

void loop()
{
	if (!digitalRead(KEY_PIN))
	{
		M5.Lcd.setCursor(0, 50);
		M5.Lcd.print("  Pressed  ");
		sendMessage();
		FastLED.setBrightness(255);
		FastLED.show();
		while (!digitalRead(KEY_PIN))
			;
	}
	else
	{
		M5.Lcd.setCursor(0, 50);
		M5.Lcd.println(" Released ");
		FastLED.setBrightness(0);
		FastLED.show();
	}
	delay(100);
}