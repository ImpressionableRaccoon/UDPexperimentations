#define REDPIN 5
#define GREENPIN 4
#define BLUEPIN 3

void initLed() {
  pinMode(REDPIN, OUTPUT);
  pinMode(GREENPIN, OUTPUT);
  pinMode(BLUEPIN, OUTPUT);

  Serial.println("LED initialization finished");
}

void setColor(byte red, byte green, byte blue) {
  analogWrite(REDPIN, red);
  analogWrite(GREENPIN, green);
  analogWrite(BLUEPIN, blue);
}
