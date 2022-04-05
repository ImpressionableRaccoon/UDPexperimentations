void setLEDColorResponse(byte (&data)[1024]) {
  unsigned int id = (unsigned int)data[1] * 256 + data[2];
  if (id != deviceID) return;
  
  byte red = data[3];
  byte green = data[4];
  byte blue = data[5];

  setColor(red, green, blue);
}
