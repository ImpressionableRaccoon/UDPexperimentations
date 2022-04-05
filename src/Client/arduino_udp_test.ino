#include <MKRGSM.h>

const unsigned int deviceID = 1;
const byte deviceToken[16] = {0xcd, 0x8a, 0x3c, 0x67, 0xc6, 0x78, 0x69, 0x7f, 0x68, 0x07, 0x05, 0x7c, 0x6b, 0x7a, 0x50, 0x1e};

const IPAddress serverIP(95, 189, 110, 110);
const unsigned int serverPort = 4848;

GSMClient client;
GSMLocation location;
GPRS gprs;
GSM gsmAccess;
GSMUDP udp;

void setup() {
  Serial.begin(115200);
  //while (!Serial);

  Serial.println("Board started");
  
  initLed();
  initSensor();
  initConnection();
  initLocation();
}

void loop() {
  updateConnection();
}
