void sendSensorsData() {
  Serial.println("sendSensorsData");
  
  float temperature = getTemperature();
  float humidity = getHumidity();
  
  byte packetBuffer[27];

  // выбираем метод sendSensorsData
  packetBuffer[0] = 0x01;

  // записываем id
  packetBuffer[1] = (byte) ((deviceID >> 8) & 0xFF);
  packetBuffer[2] = (byte) (deviceID & 0xFF);

  // записываем токен
  for (int i = 0; i < 16; i++) packetBuffer[i + 3] = deviceToken[i];

  // записываем температуру
  auto * p = reinterpret_cast<unsigned char const *>(&temperature);
  for (std::size_t i = 0; i != sizeof(float); ++i) packetBuffer[19 + i] = p[i];

  // записываем влажность
  p = reinterpret_cast<unsigned char const *>(&humidity);
  for (std::size_t i = 0; i != sizeof(float); ++i) packetBuffer[23 + i] = p[i];

  // отправляем данные
  udp.beginPacket(serverIP, serverPort);
  udp.write(packetBuffer, 27);
  udp.endPacket();
}
