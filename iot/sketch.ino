#include <WiFi.h>
#include <HTTPClient.h>
#include <Wire.h>
#include <Adafruit_BMP085.h>
#include <ArduinoJson.h> 

// ===== WIFI =====
const char *ssid = "Wokwi-GUEST";
const char *password = "";

// ===== API =====
const char *serverName = "https://46mp8r92-8080.asse.devtunnels.ms/sensor";

// ===== PIN ULTRASONIC =====
#define TRIG_PIN 5
#define ECHO_PIN 18

// ===== BMP180 =====
Adafruit_BMP085 bmp;

// ===== VARIABLE =====
float distance;
float temperature;
float pressure;

void setup()
{
  Serial.begin(115200);

  pinMode(TRIG_PIN, OUTPUT);
  pinMode(ECHO_PIN, INPUT);

  // WiFi connect
  WiFi.begin(ssid, password);
  Serial.print("Connecting to WiFi");

  while (WiFi.status() != WL_CONNECTED)
  {
    delay(500);
    Serial.print(".");
  }

  Serial.println("\nConnected!");

  // BMP180 init
  if (!bmp.begin())
  {
    Serial.println("BMP180 not found!");
    while (1)
      ;
  }
}

float readDistance()
{
  digitalWrite(TRIG_PIN, LOW);
  delayMicroseconds(2);
  digitalWrite(TRIG_PIN, HIGH);
  delayMicroseconds(10);
  digitalWrite(TRIG_PIN, LOW);

  long duration = pulseIn(ECHO_PIN, HIGH);
  float dist = duration * 0.034 / 2;
  return dist;
}

void loop()
{
  // ===== READ SENSOR =====
  distance = readDistance();
  temperature = bmp.readTemperature();
  pressure = bmp.readPressure() / 100.0; // hPa

  // ===== SEND TO API =====
  if (WiFi.status() == WL_CONNECTED)
  {
    HTTPClient http;
    http.begin(serverName);
    http.setTimeout(30000);
    http.addHeader("Content-Type", "application/json");

    // JSON body untuk request
    StaticJsonDocument<200> reqDoc;
    reqDoc["distance"] = distance;
    reqDoc["temperature"] = temperature;
    reqDoc["pressure"] = pressure;

    String jsonRequest;
    serializeJson(reqDoc, jsonRequest);

    int httpResponseCode = http.POST(jsonRequest);

    if (httpResponseCode > 0)
    {
      String responsePayload = http.getString();

      // ===== PARSING RESPONSE DARI SERVER =====
      StaticJsonDocument<512> resDoc;
      DeserializationError error = deserializeJson(resDoc, responsePayload);

      if (!error)
      {
        // Mengambil data dari objek "data" sesuai postman
        bool success = resDoc["success"];
        const char *status = resDoc["data"]["status"];
        int statusCode = resDoc["data"]["status_code"];
        float waterLevel = resDoc["data"]["water_level"];

        // ===== DEBUG PRINT HASIL API =====
        Serial.println("\n--- SERVER RESPONSE ---");
        Serial.print("Status: ");
        Serial.println(status);
        Serial.print("Status Code: ");
        Serial.println(statusCode);
        Serial.print("Water Level: ");
        Serial.print(waterLevel);
        Serial.println(" cm");
        Serial.println("-----------------------\n");
      }
      else
      {
        Serial.print("JSON Parsing Error: ");
        Serial.println(error.f_str());
      }
    }
    else
    {
      Serial.print("Error on sending POST: ");
      Serial.println(http.errorToString(httpResponseCode));
    }

    http.end();
  }

  delay(15000); // 15 detik
}