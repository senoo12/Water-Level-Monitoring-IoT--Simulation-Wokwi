#include <WiFi.h>
#include <HTTPClient.h>
#include <Wire.h>
#include <Adafruit_BMP085.h>

// ===== WIFI =====
const char* ssid = "Wokwi-GUEST";
const char* password = "";

// ===== API =====
const char* serverName = "https://46mp8r92-8080.asse.devtunnels.ms/sensor";

// ===== PIN ULTRASONIC =====
#define TRIG_PIN 5
#define ECHO_PIN 18

// ===== BMP180 =====
Adafruit_BMP085 bmp;

// ===== VARIABLE =====
float distance;
float temperature;
float pressure;

void setup() {
  Serial.begin(115200);

  pinMode(TRIG_PIN, OUTPUT);
  pinMode(ECHO_PIN, INPUT);

  // WiFi connect
  WiFi.begin(ssid, password);
  Serial.print("Connecting to WiFi");

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("\nConnected!");

  // BMP180 init
  if (!bmp.begin()) {
    Serial.println("BMP180 not found!");
    while (1);
  }
}

// ===== ULTRASONIC =====
float readDistance() {
  digitalWrite(TRIG_PIN, LOW);
  delayMicroseconds(2);

  digitalWrite(TRIG_PIN, HIGH);
  delayMicroseconds(10);
  digitalWrite(TRIG_PIN, LOW);

  long duration = pulseIn(ECHO_PIN, HIGH);
  float dist = duration * 0.034 / 2;

  return dist;
}

void loop() {

  // ===== READ SENSOR =====
  distance = readDistance();
  temperature = bmp.readTemperature();
  pressure = bmp.readPressure() / 100.0; // hPa

  // ===== DEBUG PRINT =====
  Serial.println("===== DATA =====");
  Serial.print("Distance: "); Serial.println(distance);
  Serial.print("Temp: "); Serial.println(temperature);
  Serial.print("Pressure: "); Serial.println(pressure);
  Serial.println("================");

  // ===== SEND TO API =====
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;

    http.begin(serverName);
    http.addHeader("Content-Type", "application/json");

    // JSON body (SESUAI API BARU)
    String json = "{";
    json += "\"distance\":" + String(distance, 2) + ",";
    json += "\"temperature\":" + String(temperature, 2) + ",";
    json += "\"pressure\":" + String(pressure, 2);
    json += "}";

    int httpResponseCode = http.POST(json);

    Serial.print("HTTP Response: ");
    Serial.println(httpResponseCode);

    if (httpResponseCode > 0) {
      String response = http.getString();
      Serial.println("Response Body:");
      Serial.println(response);
    } else {
      Serial.print("Error: ");
      Serial.println(http.errorToString(httpResponseCode));
    }

    http.end();
  }

  // ===== DELAY (IMPORTANT untuk ThingSpeak) =====
  delay(15000); // 15 detik (hindari rate limit)
}