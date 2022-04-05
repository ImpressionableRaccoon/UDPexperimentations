void initLocation() {
  location.begin();
}

float getLatitude() {
  return location.latitude();
}

float getLongitude() {
  return location.longitude();
}

int getAltitude() {
  return location.altitude();
}

int getAccuracy() {
  return location.accuracy();
}
