const char PINNUMBER[]     = "";
const char GPRS_APN[]      = "internet.beeline.ru";
const char GPRS_LOGIN[]    = "beeline";
const char GPRS_PASSWORD[] = "beeline";

unsigned long tempLastUpdate = 0;
unsigned long tempUpdateInterval = 10000; // 10 seconds

unsigned long locationLastUpdate = 0;
unsigned long locationUpdateInterval = 30000; // 30 seconds (need to be changed to 30 minutes)

void initConnection() {
  bool connected = false;

  Serial.println("Connecting to network");
  
  while (!connected)
    if ((gsmAccess.begin(PINNUMBER) == GSM_READY) &&
        (gprs.attachGPRS(GPRS_APN, GPRS_LOGIN, GPRS_PASSWORD) == GPRS_READY))
      connected = true;
    else {
      Serial.println("Not connected");
      delay(1000);
    }
  Serial.println("Connected");

  udp.begin(2390);

  getMe();
}

void updateConnection() {
  if (tempLastUpdate == 0 || millis() - tempLastUpdate >= tempUpdateInterval) {
    tempLastUpdate = millis();
    sendSensorsData();
  }
  if (locationLastUpdate == 0 || millis() - locationLastUpdate >= locationUpdateInterval) {
    locationLastUpdate = millis();
    sendLocationData();
  }
  receive();
}

void receive() {
  int result = udp.parsePacket();
  if (result) {
    byte packetBuffer[1024];
    udp.read(packetBuffer, 1024);

    switch (packetBuffer[0]) {
      case 0x00:
        getMeResponse(packetBuffer);
        break;
      case 0x02:
        setLEDColorResponse(packetBuffer);
        break;
    }

    Serial.println("Packet received");
    for (int i = 0; i < result; i++) {
      Serial.print(i);
      Serial.print('\t');
      Serial.print((int) packetBuffer[i]);
      Serial.println();
    }
  }
}
