void getMe() {
  Serial.println("getMe");
  byte packetBuffer[19];

  // выбираем метод getMe
  packetBuffer[0] = 0x00;

  // записываем id
  packetBuffer[1] = (byte) ((deviceID >> 8) & 0xFF);
  packetBuffer[2] = (byte) (deviceID & 0xFF);

  // записываем токен
  for (int i = 0; i < 16; i++) packetBuffer[i + 3] = deviceToken[i];

  udp.beginPacket(serverIP, serverPort);
  udp.write(packetBuffer, 19);
  udp.endPacket();
}

void getMeResponse(byte (&data)[1024]) {
  switch (data[1]) {
    case 0x00:
      {
        unsigned int errorCode = (unsigned int)data[2] * 256 + data[3];
        String error = "";
        for (int i = 4; i < 36 && data[i] != 0; i++) error += (char) data[i];
        
        Serial.print("Ошибка ");
        Serial.print(errorCode);
        Serial.print(": ");
        Serial.print(error);
        Serial.println();

        while(1);
      }
      break;
    case 0x01:
      {
        String deviceName = "";
        for (int i = 4; i < 20 && data[i] != 0; i++) deviceName += (char) data[i];

        Serial.println("Авторизация прошла успешно!");
        Serial.println("Название устройства: " + deviceName);
      }
      break;
    default:
      break;
  }
}
