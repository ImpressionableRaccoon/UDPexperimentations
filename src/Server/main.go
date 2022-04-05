package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"html/template"
	"math"
	"net"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type location struct {
	Latitude  float32
	Longitude float32
	Altitude  int32
	Accuracy  int32
}

type device struct {
	Name               string
	Token              []byte
	LastOnline         time.Time
	LastSensorsUpdate  time.Time
	LastLocationUpdate time.Time
	Temperature        float32
	Humidity           float32
	Remote             *net.UDPAddr
	Location           location
}

var devices = map[int]device{
	1: {
		Name:        "Inverter",
		Token:       []byte{205, 138, 60, 103, 198, 120, 105, 127, 104, 7, 5, 124, 107, 122, 80, 30},
		LastOnline:  time.Unix(0, 0),
		Temperature: 0.0,
		Humidity:    0.0,
	},
	2: {
		Name:        "Inverter2",
		Token:       []byte{205, 138, 60, 103, 198, 120, 105, 127, 104, 7, 5, 124, 107, 122, 80, 30},
		LastOnline:  time.Unix(0, 0),
		Temperature: 0.0,
		Humidity:    0.0,
	},
	3: {
		Name:        "Inverter3",
		Token:       []byte{205, 138, 60, 103, 198, 120, 105, 127, 104, 7, 5, 124, 107, 122, 80, 30},
		LastOnline:  time.Unix(0, 0),
		Temperature: 0.0,
		Humidity:    0.0,
	},
}
var UDPServer *net.UDPConn

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/device/", deviceHandler)
	http.HandleFunc("/api/sendColor", sendColorHandler)
	go func() {
		err := http.ListenAndServe(":4849", nil)
		if err != nil {
			panic(err)
		}
	}()

	// Тут мы храним адрес и порт сервера
	serverAddress := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 4848,
	}

	// Создаем UDP-сервер
	var err error
	UDPServer, err = net.ListenUDP("udp", &serverAddress)
	if err != nil {
		panic(err)
	}

	// При выходе из программы выключаем сервер
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(UDPServer)

	fmt.Printf("Сервер успешно запущен на %s\n", UDPServer.LocalAddr().String())

	// Начинаем в бесконечном цикле обрабатывать клиентов, запросы от которых приходят к серверу
	for {
		message := make([]byte, 1024)

		// Читаем сообщение из соединения и записываем его в массив message
		// messageLength - длина сообщения
		// remote - адрес, с которого пришло сообщение
		// messageLength, remote, err := server.ReadFromUDP(message[:])
		_, remote, err := UDPServer.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		//fmt.Printf("Получили сообщение от %s\n", remote)
		//for i := 0; i < messageLength; i++ {
		//	fmt.Printf("%d\t%X\n", i, message[i])
		//}

		switch message[0] {
		case 0x00:
			go getMe(remote, message)
		case 0x01:
			go sendSensorsDataResponse(remote, message)
		case 0x03:
			go sendLocationDataResponse(remote, message)
		}

		//// Преобразуем сообщение в строку и выводим его
		//data := strings.TrimSpace(string(message[:messageLength]))
		//fmt.Printf("Получили: %s от %s\n", data, remote)
	}
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	payload := struct {
		Name    string
		Devices map[int]device
	}{
		"Список устройств",
		devices,
	}

	t, err := template.New("index").ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, payload)
	if err != nil {
		panic(err)
	}
}

func faviconHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "assets/favicon.ico")
}

func deviceHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(path.Base(req.URL.Path))
	if err != nil {
		http.Redirect(w, req, "/", 301)
	}

	d, found := devices[id]

	if !found {
		http.Redirect(w, req, "/", 301)
	}

	payload := struct {
		ID                         int
		Name                       string
		Token                      string
		LastOnline                 time.Time
		LastOnlineDuration         time.Duration
		LastSensorsUpdate          time.Time
		LastSensorsUpdateDuration  time.Duration
		LastLocationUpdate         time.Time
		LastLocationUpdateDuration time.Duration
		Temperature                float32
		Humidity                   float32
		Remote                     *net.UDPAddr
		Location                   location
	}{
		id,
		d.Name,
		strings.ToUpper(hex.EncodeToString(d.Token)),
		d.LastOnline,
		time.Now().Sub(d.LastOnline).Round(time.Second),
		d.LastSensorsUpdate,
		time.Now().Sub(d.LastSensorsUpdate).Round(time.Second),
		d.LastLocationUpdate,
		time.Now().Sub(d.LastLocationUpdate).Round(time.Second),
		d.Temperature,
		d.Humidity,
		d.Remote,
		d.Location,
	}

	t, err := template.New("device").ParseFiles("templates/device.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, payload)
	if err != nil {
		panic(err)
	}
}

func sendColorHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	id, _ := strconv.Atoi(req.FormValue("id"))
	red, _ := strconv.ParseUint(req.FormValue("red"), 10, 8)
	green, _ := strconv.ParseUint(req.FormValue("green"), 10, 8)
	blue, _ := strconv.ParseUint(req.FormValue("blue"), 10, 8)

	_, found := devices[id]

	if !found {
		_, err = fmt.Fprint(w, `"status": "error", "error": "id_not_found"`)
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Printf("New color for id %d\n", id)
	fmt.Printf("Red: %d\n", red)
	fmt.Printf("Green: %d\n", green)
	fmt.Printf("Blue: %d\n", blue)

	setLEDColor(id, uint8(red), uint8(green), uint8(blue))

	_, err = fmt.Fprint(w, `"status": "ok"`)
	if err != nil {
		panic(err)
	}
}

func getMe(addr *net.UDPAddr, payload []byte) {
	id := int(payload[1])*256 + int(payload[2])
	token := payload[3:19]

	fmt.Println("ID:", id)
	fmt.Println("Token: ", token)

	d, found := devices[id]
	if !found || !reflect.DeepEqual(token, d.Token) {
		response := make([]byte, 36)

		response[0] = 0
		response[1] = 0
		response[2] = 1
		response[3] = 145
		for i, value := range "Unauthorized" {
			response[i+4] = (byte)(value)
		}

		sendResponse(addr, response)

		fmt.Println("Device not found!")
		return
	}

	d.LastOnline = time.Now()
	d.Remote = addr
	devices[id] = d

	response := make([]byte, 20)

	response[0] = 0
	response[1] = 1
	response[2] = (byte)((id >> 8) & 0xFF)
	response[3] = (byte)(id & 0xFF)
	for i, value := range d.Name {
		response[i+4] = (byte)(value)
	}

	sendResponse(addr, response)
}

func sendSensorsDataResponse(addr *net.UDPAddr, payload []byte) {
	id := int(payload[1])*256 + int(payload[2])
	token := payload[3:19]

	fmt.Println("ID:", id)
	fmt.Println("Token: ", token)

	d, found := devices[id]
	if found && reflect.DeepEqual(token, d.Token) {
		d.LastOnline = time.Now()
		d.LastSensorsUpdate = time.Now()
		d.Remote = addr

		d.Temperature = math.Float32frombits(binary.LittleEndian.Uint32(payload[19:23]))
		d.Humidity = math.Float32frombits(binary.LittleEndian.Uint32(payload[23:27]))

		devices[id] = d

		fmt.Printf("Получили температуру и влажность от устройства %s, id = %d!\n", d.Name, id)
		fmt.Printf("Температура: %f\n", d.Temperature)
		fmt.Printf("Влажность: %f\n", d.Humidity)

		return
	}
}

func setLEDColor(id int, red, green, blue uint8) {
	response := make([]byte, 6)

	response[0] = 0x02

	response[1] = byte(id / 256)
	response[2] = byte(id % 256)

	response[3] = red
	response[4] = green
	response[5] = blue

	sendResponse(devices[id].Remote, response)
}

func sendLocationDataResponse(addr *net.UDPAddr, payload []byte) {
	id := int(payload[1])*256 + int(payload[2])
	token := payload[3:19]

	fmt.Println("ID:", id)
	fmt.Println("Token: ", token)

	d, found := devices[id]
	if found && reflect.DeepEqual(token, d.Token) {
		d.LastOnline = time.Now()
		d.LastLocationUpdate = time.Now()
		d.Remote = addr

		d.Location.Latitude = math.Float32frombits(binary.LittleEndian.Uint32(payload[19:23]))
		d.Location.Longitude = math.Float32frombits(binary.LittleEndian.Uint32(payload[23:27]))

		var alt [4]byte
		alt[0] = payload[27]
		alt[1] = payload[28]
		alt[2] = payload[29]
		alt[3] = payload[30]
		d.Location.Altitude = *(*int32)(unsafe.Pointer(&alt))

		var acc [4]byte
		acc[0] = payload[31]
		acc[1] = payload[32]
		acc[2] = payload[33]
		acc[3] = payload[34]
		d.Location.Accuracy = *(*int32)(unsafe.Pointer(&acc))

		devices[id] = d

		fmt.Printf("Получили местоположение от устройства %s, id = %d!\n", d.Name, id)
		fmt.Printf("Широта: %f\n", d.Location.Latitude)
		fmt.Printf("Долгота: %f\n", d.Location.Longitude)
		fmt.Printf("Высота: %d\n", d.Location.Altitude)
		fmt.Printf("Точность: %d\n", d.Location.Accuracy)

		return
	}
}

func sendResponse(addr *net.UDPAddr, payload []byte) {
	_, err := UDPServer.WriteToUDP(payload, addr)
	if err != nil {
		fmt.Printf("Не получается отправить запрос: %v\n", err)
	}
}
