package api

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const id = "test"
const test_key = "8B0D20CB790418B3CBE3A8B7B0A0A7F114BFFBD2179DF015A7EF086845B15C46"

func TestValidation(t *testing.T) {
	var res map[string]interface{}
	json.Unmarshal([]byte(randomJSON), &res)
	data, _ := json.Marshal(res)

	h := sha256.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)

	//shareKey should be of the format id:key
	key, err := hex.DecodeString(test_key)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		t.Error(err)
		t.FailNow()
	}

	hash := gcm.Seal(nonce, nonce, bs, nil)
	hashStr := base64.StdEncoding.EncodeToString(hash)
	//sanity check?
	hash2, _ := base64.StdEncoding.DecodeString(hashStr)
	if !bytes.Equal(hash, hash2) {
		t.Error("hash is weird?")
		t.FailNow()
	}

	encryptedHash := []byte(id + ":" + hashStr)

	log.Printf("encrypted hash: %v\n", string(encryptedHash))

	s := &Server{
		cfg: Config{
			AESDecryptionKeys: make(map[string][]byte),
		},
	}
	s.cfg.AESDecryptionKeys[id] = key

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	sugar := logger.Sugar()
	sugar.Debugw("logger initiated")

	s.Log = sugar

	err = s.validateSigning([]byte(randomJSON), id+":"+hashStr)
	if err != nil {
		t.Error(err)
	}

}

const randomJSON = `
[
  {
    "_id": "63b8c78d8253c16f91f38b5a",
    "index": 0,
    "guid": "a222d4f2-b082-4ca7-ac15-7cd62edb13c4",
    "isActive": false,
    "balance": "$2,734.38",
    "picture": "http://placehold.it/32x32",
    "age": 22,
    "eyeColor": "brown",
    "name": "Gaines Villarreal",
    "gender": "male",
    "company": "DYMI",
    "email": "gainesvillarreal@dymi.com",
    "phone": "+1 (985) 531-3134",
    "address": "892 Chester Avenue, Boling, Arkansas, 6031",
    "about": "Tempor est ipsum nisi cupidatat commodo cupidatat in officia laboris minim ut. Cupidatat veniam sint amet ea ex adipisicing ut magna excepteur et labore officia amet. Lorem aliquip adipisicing dolor ut minim ipsum Lorem commodo fugiat ea ullamco velit tempor. Labore nisi exercitation mollit cillum qui qui sint. Anim id nulla aute non magna excepteur ullamco tempor aute in quis dolor. Pariatur reprehenderit fugiat dolor ipsum mollit officia exercitation esse dolore. Nisi veniam aliquip sunt ex amet mollit sit Lorem in anim.\r\n",
    "registered": "2015-08-01T09:43:19 +04:00",
    "latitude": -75.546611,
    "longitude": -107.749453,
    "tags": [
      "do",
      "ut",
      "ex",
      "ullamco",
      "esse",
      "consectetur",
      "duis"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Lila Nicholson"
      },
      {
        "id": 1,
        "name": "Susie Jensen"
      },
      {
        "id": 2,
        "name": "Robles Russell"
      }
    ],
    "greeting": "Hello, Gaines Villarreal! You have 3 unread messages.",
    "favoriteFruit": "apple"
  },
  {
    "_id": "63b8c78ddd91a9f174a9c806",
    "index": 1,
    "guid": "c1c7a6cf-515e-4be3-aef2-faee825bfde6",
    "isActive": false,
    "balance": "$1,917.84",
    "picture": "http://placehold.it/32x32",
    "age": 25,
    "eyeColor": "blue",
    "name": "Susanne Scott",
    "gender": "female",
    "company": "ZIGGLES",
    "email": "susannescott@ziggles.com",
    "phone": "+1 (872) 500-3167",
    "address": "156 Prospect Place, Wells, Mississippi, 5725",
    "about": "Dolor sit ut velit tempor duis amet dolore qui irure voluptate. Laboris laboris adipisicing consectetur id et tempor nisi sit eu excepteur. Culpa minim culpa laboris ut aute aliquip nisi consectetur duis id elit. In eu laborum adipisicing exercitation culpa.\r\n",
    "registered": "2018-04-08T07:03:27 +04:00",
    "latitude": 42.034142,
    "longitude": -86.968782,
    "tags": [
      "ipsum",
      "pariatur",
      "sint",
      "cupidatat",
      "aliquip",
      "est",
      "eu"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Josie Rodgers"
      },
      {
        "id": 1,
        "name": "Rae Mcneil"
      },
      {
        "id": 2,
        "name": "Mcfadden Reed"
      }
    ],
    "greeting": "Hello, Susanne Scott! You have 1 unread messages.",
    "favoriteFruit": "banana"
  },
  {
    "_id": "63b8c78dd2229dee15813720",
    "index": 2,
    "guid": "31f6e463-7e42-4193-ad06-ab8f3ab8a847",
    "isActive": true,
    "balance": "$3,300.34",
    "picture": "http://placehold.it/32x32",
    "age": 24,
    "eyeColor": "green",
    "name": "Sanchez Munoz",
    "gender": "male",
    "company": "BOVIS",
    "email": "sanchezmunoz@bovis.com",
    "phone": "+1 (831) 453-3139",
    "address": "439 Claver Place, Clinton, Washington, 4093",
    "about": "Tempor ea dolor dolore sint et sunt ullamco ex ex veniam dolor duis mollit enim. Nostrud tempor fugiat nostrud qui mollit laborum commodo pariatur voluptate exercitation ipsum eu. Deserunt dolor tempor aute laboris ullamco sit nulla ipsum incididunt. Laborum culpa sint officia ea in laboris sint et aliquip amet anim. Do ullamco et pariatur sit reprehenderit consectetur dolore sunt dolor ullamco magna aliquip nulla voluptate. Occaecat elit cillum irure ad proident ad consequat ut ea commodo eiusmod aliqua.\r\n",
    "registered": "2016-10-01T12:02:37 +04:00",
    "latitude": -77.350426,
    "longitude": 170.752199,
    "tags": [
      "esse",
      "laboris",
      "ad",
      "eu",
      "ad",
      "laboris",
      "veniam"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Angelia Ayala"
      },
      {
        "id": 1,
        "name": "Joyce Shepard"
      },
      {
        "id": 2,
        "name": "Henrietta Hogan"
      }
    ],
    "greeting": "Hello, Sanchez Munoz! You have 1 unread messages.",
    "favoriteFruit": "apple"
  },
  {
    "_id": "63b8c78d7bf1a1e520c946c3",
    "index": 3,
    "guid": "4572da83-fe97-41e4-8aff-fd444f044dc2",
    "isActive": false,
    "balance": "$3,395.92",
    "picture": "http://placehold.it/32x32",
    "age": 20,
    "eyeColor": "blue",
    "name": "Ronda Townsend",
    "gender": "female",
    "company": "SAVVY",
    "email": "rondatownsend@savvy.com",
    "phone": "+1 (879) 591-2789",
    "address": "703 Court Square, Vivian, Tennessee, 1990",
    "about": "Officia elit adipisicing elit sit pariatur aliqua cupidatat Lorem pariatur exercitation eu aute. Est deserunt exercitation ex magna aliquip voluptate dolore consequat cupidatat laboris laboris minim ex do. Occaecat anim Lorem nisi laborum nisi amet. Proident officia sit tempor amet nisi cupidatat proident proident voluptate voluptate ipsum ut. Eu nostrud commodo occaecat quis ut eu enim magna et dolore.\r\n",
    "registered": "2018-10-10T05:56:43 +04:00",
    "latitude": -40.264579,
    "longitude": 141.018599,
    "tags": [
      "ullamco",
      "in",
      "non",
      "aliquip",
      "laboris",
      "Lorem",
      "commodo"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Karyn Trujillo"
      },
      {
        "id": 1,
        "name": "Sosa Cooke"
      },
      {
        "id": 2,
        "name": "Hughes Gibson"
      }
    ],
    "greeting": "Hello, Ronda Townsend! You have 2 unread messages.",
    "favoriteFruit": "strawberry"
  },
  {
    "_id": "63b8c78d0a7416268b508491",
    "index": 4,
    "guid": "aab6a6f8-2045-4422-8058-3ca34bd9f68f",
    "isActive": true,
    "balance": "$3,629.54",
    "picture": "http://placehold.it/32x32",
    "age": 32,
    "eyeColor": "brown",
    "name": "Alissa Bridges",
    "gender": "female",
    "company": "EXOSTREAM",
    "email": "alissabridges@exostream.com",
    "phone": "+1 (858) 452-2157",
    "address": "315 Dahlgreen Place, Blairstown, West Virginia, 1941",
    "about": "Do aliqua elit esse anim esse enim anim. Officia cillum pariatur sit tempor minim sunt ut veniam labore nisi veniam. Amet veniam velit culpa duis adipisicing laboris sint nisi id irure cillum. Ea pariatur eiusmod exercitation dolor excepteur magna qui et sint ea. Sit reprehenderit aute aute cupidatat est est laborum tempor est deserunt nisi ipsum fugiat sit.\r\n",
    "registered": "2021-07-11T08:48:07 +04:00",
    "latitude": 11.777498,
    "longitude": -159.848101,
    "tags": [
      "et",
      "dolore",
      "Lorem",
      "sunt",
      "sit",
      "amet",
      "mollit"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Rachel Franks"
      },
      {
        "id": 1,
        "name": "Bond Kelly"
      },
      {
        "id": 2,
        "name": "Shannon Hurst"
      }
    ],
    "greeting": "Hello, Alissa Bridges! You have 3 unread messages.",
    "favoriteFruit": "banana"
  },
  {
    "_id": "63b8c78d1b27e386918f932f",
    "index": 5,
    "guid": "4d9aa863-8bbc-47a3-a18d-65a6c08fd059",
    "isActive": true,
    "balance": "$2,150.26",
    "picture": "http://placehold.it/32x32",
    "age": 28,
    "eyeColor": "blue",
    "name": "Hoffman Romero",
    "gender": "male",
    "company": "DANCITY",
    "email": "hoffmanromero@dancity.com",
    "phone": "+1 (915) 534-3636",
    "address": "332 Rockaway Avenue, Bath, Hawaii, 8729",
    "about": "Reprehenderit est excepteur ullamco proident adipisicing nulla non pariatur non cillum tempor sunt ut. In reprehenderit eu occaecat esse adipisicing elit minim adipisicing exercitation qui esse adipisicing magna. Nisi consequat quis consequat veniam. Irure amet esse sit nisi excepteur culpa id ut consequat culpa ullamco excepteur adipisicing ut.\r\n",
    "registered": "2020-08-08T06:51:09 +04:00",
    "latitude": 75.898447,
    "longitude": 172.542552,
    "tags": [
      "laborum",
      "amet",
      "esse",
      "sunt",
      "aliqua",
      "Lorem",
      "sit"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Lela Roach"
      },
      {
        "id": 1,
        "name": "Rena Hammond"
      },
      {
        "id": 2,
        "name": "Whitaker Buchanan"
      }
    ],
    "greeting": "Hello, Hoffman Romero! You have 1 unread messages.",
    "favoriteFruit": "banana"
  },
  {
    "_id": "63b8c78d514323f7003a96b8",
    "index": 6,
    "guid": "4bdde114-d4c9-479b-8c75-ec097394aff3",
    "isActive": false,
    "balance": "$1,053.36",
    "picture": "http://placehold.it/32x32",
    "age": 23,
    "eyeColor": "blue",
    "name": "Dorsey Hardy",
    "gender": "male",
    "company": "XTH",
    "email": "dorseyhardy@xth.com",
    "phone": "+1 (905) 588-3971",
    "address": "427 Fillmore Place, Dyckesville, Vermont, 8844",
    "about": "Proident aute anim sint do. Proident exercitation nulla ea voluptate dolore occaecat amet ex consequat dolore nulla sint. Est velit nostrud ipsum consequat pariatur dolore consequat fugiat. Proident esse tempor enim nulla ipsum qui elit officia dolore voluptate duis. Labore eu nulla velit consectetur enim irure quis laborum nisi. Aliquip tempor id quis anim aliqua elit non sit reprehenderit exercitation. Dolore veniam nostrud ullamco ipsum dolore enim exercitation.\r\n",
    "registered": "2017-07-10T02:48:59 +04:00",
    "latitude": 82.974639,
    "longitude": 108.059367,
    "tags": [
      "irure",
      "id",
      "quis",
      "occaecat",
      "Lorem",
      "sit",
      "labore"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Castro Browning"
      },
      {
        "id": 1,
        "name": "Maude Bradley"
      },
      {
        "id": 2,
        "name": "Bradford Strickland"
      }
    ],
    "greeting": "Hello, Dorsey Hardy! You have 10 unread messages.",
    "favoriteFruit": "strawberry"
  }
]`
