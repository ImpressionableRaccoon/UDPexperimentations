void sendLocationData() {
  Serial.println("sendLocationData");
  if (!location.available()) {
    Serial.println("Location not available");
    return;
  }

  float latitude = getLatitude();
  float longitude = getLongitude();
  int altitude = getAltitude();
  int accuracy = getAccuracy();

  Serial.println("Coordinates: ");
  Serial.print(latitude);
  Serial.print(", ");
  Serial.println(longitude);

  Serial.println("Altitude: " + String(altitude));
  Serial.println("Accuracy: " + String(accuracy));
  
  byte packetBuffer[35];

  // выбираем метод sendLocationData
  packetBuffer[0] = 0x03;

  // записываем id
  packetBuffer[1] = (byte) ((deviceID >> 8) & 0xFF);
  packetBuffer[2] = (byte) (deviceID & 0xFF);

  // записываем токен
  for (int i = 0; i < 16; i++) packetBuffer[i + 3] = deviceToken[i];

  // записываем широту
  auto * p = reinterpret_cast<unsigned char const *>(&latitude);
  for (int i = 0; i < 4; i++) packetBuffer[19 + i] = p[i];

  // записываем долготу
  p = reinterpret_cast<unsigned char const *>(&longitude);
  for (int i = 0; i < 4; i++) packetBuffer[23 + i] = p[i];

  // записываем высоту
  p = reinterpret_cast<unsigned char const *>(&altitude);
  for (int i = 0; i < 4; i++) packetBuffer[27 + i] = p[i];

  // записываем точность
  p = reinterpret_cast<unsigned char const *>(&accuracy);
  for (int i = 0; i < 4; i++) packetBuffer[31 + i] = p[i];

  // отправляем данные
  udp.beginPacket(serverIP, serverPort);
  udp.write(packetBuffer, 35);
  udp.endPacket();
}
