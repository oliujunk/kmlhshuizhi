package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"github.com/go-xorm/xorm"
	"github.com/robfig/cron"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"xorm.io/core"
)

// DataEntity 数据
type DataEntity struct {
	DeviceID int      `json:"deviceId"`
	Entity   []Entity `json:"entity"`
}

// Entity 实体
type Entity struct {
	Datetime string `json:"datetime"`
	EUnit    string `json:"eUnit"`
	EValue   string `json:"eValue"`
	EKey     string `json:"eKey"`
	EName    string `json:"eName"`
	ENum     string `json:"eNum"`
}

type StRsvrR struct {
	STCD string    `xorm:"notnull 'STCD'"`
	TM   time.Time `xorm:"notnull 'TM'"`
	RZ   float64   `xorm:"'RZ'"`
}

var (
	engine *xorm.Engine
	//deviceIDs = [...]int{16072689, 16072690, 16078692}
	deviceIDs = [...]int{16078692}
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	var err error
	//engine, err = xorm.NewEngine("odbc", "driver={SQL Server};server=127.0.0.1;uid=sa;pwd=Xph87510227#;database=HYDDX")
	//engine, err = xorm.NewEngine("odbc", "driver={SQL Server};server=127.0.0.1;uid=sa;pwd=hdhsk@2019;database=HYDDX")
	engine, err = xorm.NewEngine("odbc", "driver={SQL Server};server=127.0.0.1;uid=sa;pwd=lonhwin@2011sql;database=HYDDX")
	if err != nil {
		log.Println(err)
	}

	engine.SetTableMapper(core.SnakeMapper{})

	c := cron.New()
	_ = c.AddFunc("0 */1 * * * *", getData)
	c.Start()

	log.Println("1分钟写入一次数据")

	select {}
}

func getData() {
	cstZone := time.FixedZone("CST", 8*3600)
	for _, deviceID := range deviceIDs {
		log.Printf("deviceID: %d\n", deviceID)
		resp, err := http.Get("http://47.105.215.208:8005/intfa/queryData/" + strconv.Itoa(deviceID))
		if err != nil {
			log.Println(err)
			continue
		}
		result, _ := ioutil.ReadAll(resp.Body)
		var dataEntity DataEntity
		_ = json.Unmarshal(result, &dataEntity)
		//if idIndex == 0 {
		//	for index, entity := range dataEntity.Entity {
		//		var stData StRsvrR
		//		if index == 0 {
		//			stData.STCD = "20000025"
		//			stData.TM = time.Now().In(cstZone)
		//			stData.RZ, _ = strconv.ParseFloat(entity.EValue, 64)
		//		} else if index == 1 {
		//			stData.STCD = "20000024"
		//			stData.TM = time.Now().In(cstZone)
		//			stData.RZ, _ = strconv.ParseFloat(entity.EValue, 64)
		//		} else {
		//			stData.STCD = fmt.Sprintf("200000%02d", index-1)
		//			stData.TM = time.Now().In(cstZone)
		//			stData.RZ, _ = strconv.ParseFloat(entity.EValue, 64)
		//		}
		//		_, err := engine.Insert(&stData)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//	}
		//} else if idIndex == 1 {
		//	for index, entity := range dataEntity.Entity {
		//		var stData StRsvrR
		//		stData.STCD = fmt.Sprintf("200000%02d", index+15)
		//		stData.TM = time.Now().In(cstZone)
		//		stData.RZ, _ = strconv.ParseFloat(entity.EValue, 64)
		//		_, err := engine.Insert(&stData)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//	}
		//} else if idIndex == 2 {
		//	for index, entity := range dataEntity.Entity {
		//		if entity.EValue == "3276.7" {
		//			continue
		//		}
		//		var stData StRsvrR
		//		stData.STCD = fmt.Sprintf("%d", index+53238181)
		//		stData.TM = time.Now().In(cstZone)
		//		stData.RZ, _ = strconv.ParseFloat(entity.EValue, 64)
		//		_, err := engine.Insert(&stData)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//	}
		//}
		for index, entity := range dataEntity.Entity {
			if entity.EValue == "3276.7" {
				continue
			}
			var stData StRsvrR
			stData.STCD = fmt.Sprintf("%d", index+53238181)
			stData.TM = time.Now().In(cstZone)
			stData.RZ, _ = strconv.ParseFloat(entity.EValue, 64)
			_, err := engine.Insert(&stData)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
