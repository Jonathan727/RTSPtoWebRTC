package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/deepch/vdk/av"
)

//Config global
var Config = loadConfig()

//ConfigST struct
type ConfigST struct {
	Server  ServerST            `json:"server"`
	Streams map[string]StreamST `json:"streams"`
}

//ServerST struct
type ServerST struct {
	HTTPPort string `json:"http_port"`
}

//StreamST struct
type StreamST struct {
	URL    string `json:"url"`
	Status bool   `json:"status"`
	Codecs []av.CodecData
	Cl     map[string]viewer
}

type viewer struct {
	c chan av.Packet
}

func loadConfig() *ConfigST {
	var tmp ConfigST
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		log.Fatalln(err)
	}
	for i, v := range tmp.Streams {
		v.Cl = make(map[string]viewer)
		tmp.Streams[i] = v
	}
	return &tmp
}

func (c *ConfigST) cast(uuid string, pck av.Packet) {
	for _, v := range c.Streams[uuid].Cl {
		if len(v.c) < cap(v.c) {
			v.c <- pck
		}
	}
}

func (c *ConfigST) ext(suuid string) bool {
	_, ok := c.Streams[suuid]
	return ok
}

func (c *ConfigST) coAd(suuid string, codecs []av.CodecData) {
	t := c.Streams[suuid]
	t.Codecs = codecs
	c.Streams[suuid] = t
}

func (c *ConfigST) coGe(suuid string) []av.CodecData {
	return c.Streams[suuid].Codecs
}

func (c *ConfigST) clAd(suuid string) (string, chan av.Packet) {
	cuuid := pseudoUUID()
	ch := make(chan av.Packet, 100)
	c.Streams[suuid].Cl[cuuid] = viewer{c: ch}
	return cuuid, ch
}

func (c *ConfigST) list() (string, []string) {
	var res []string
	var fist string
	for k := range c.Streams {
		if fist == "" {
			fist = k
		}
		res = append(res, k)
	}
	return fist, res
}
func (c *ConfigST) clDe(suuid, cuuid string) {
	delete(c.Streams[suuid].Cl, cuuid)
}

func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
