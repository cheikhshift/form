package form

import (
	"time"
	"github.com/asaskevich/govalidator"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"crypto/sha512"
	"gopkg.in/mgo.v2/bson"
	"github.com/cheikhshift/gos/core"
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"strings"
	"reflect"
	"os"
	"strconv"
)

type File string 
type Paragraph string 
type Date string
type Select string
type SelectMult Select
type Radio Select
type Email string
type Password []byte
//select opts,placeholder,title
const (
        MB = 1 << 20
      
)

var (
	InputClass string
	ButtonClass string
	FormKey  = "a very very very very secret key"
	MaxSize = 10 //Mb
)

func GetSel(tagstr string) string {
	strs := strings.Split(tagstr,"\",")
	for _,v := range strs {
		if strings.Contains(v,"select") {
			splm := strings.Split(v,":")
			return strings.Replace(splm[1], "\"","", -1)
		}
	}

	return ""
}

func GetPl(tagstr string) string {
	strs := strings.Split(tagstr,"\",")
	for _,v := range strs {
		if strings.Contains(v,"placeholder") {
			splm := strings.Split(v,":")
			return strings.Replace(splm[1], "\"","", -1)
		}
	}

	return ""
}

func Hash(input string) []byte {
	return sha512.New512_256().Sum([]byte(input ))
}
	

func Path(pathtoformat File) string {
	return fmt.Sprintf("./uploads/%s", pathtoformat )
}


func Form(r *http.Request, i interface{}) error {
	//convert form and validate 
	var err error
	v := reflect.ValueOf(i).Elem()
			//t := reflect.TypeOf(item)
			bso :=  bson.M{"__elem" : "form"}

			for i := 0; i < v.NumField(); i++ {
				field := v.Type().Field(i)
				fieldtype := strings.ToLower(field.Type.String())				
					if strings.Contains(fieldtype, "file"){			
						file, handler, _ := r.FormFile(field.Name)

						if handler != nil {
		      			fid := fmt.Sprintf("%s-%s", core.NewLen(15),handler.Filename )
				        f, err := os.OpenFile(fmt.Sprintf("./uploads/%s", fid ), os.O_WRONLY|os.O_CREATE, 0700)
				        if err != nil {
				           
				            return err
				        }
				        
				 
				        io.Copy(f, file)
				        f.Close()
				        bso[field.Name] = fid
				        file.Close()
			    	}
				} else if strings.Contains(fieldtype, "bool"){
					bso[field.Name] = strings.Contains( r.FormValue(field.Name), "on")
				} else if strings.Contains(fieldtype, "password"){
					bso[field.Name] = Hash( r.FormValue(field.Name) )
				}  else if strings.Contains(fieldtype, "int")  {
					i, _ := strconv.Atoi(r.FormValue(field.Name))
					bso[field.Name] = i
				} else if strings.Contains(fieldtype, "float") {
					f, _ := strconv.ParseFloat(r.FormValue(field.Name), 64)
					bso[field.Name] = f
				} else  {
					bso[field.Name] = r.FormValue(field.Name)
				
				}
			
			}

	data,_ := json.Marshal(&bso)
	json.Unmarshal([]byte(data), i)

	_, err = govalidator.ValidateStruct(i)
	if err != nil {
			return err
	}

	return nil
	//ValidateStruct
}

func SetKey(key string) {
	FormKey = key
}

func ToBson(data string) bson.M {
	res := bson.M{}
    json.Unmarshal([]byte(data), &res)
    return res
}

func ValidateRequest(r *http.Request, usertoken string) bool {
	token := Decrypt([]byte(FormKey),r.FormValue("formtoken"))
	return strings.Contains(token, usertoken ) && strings.Contains(token, r.URL.Path )
}

func GenerateToken(url string, usertoken string) string {
	return Encrypt([]byte(FormKey), fmt.Sprintf("%s@%s@%s",url, time.Now().String() , usertoken) )
}

func Encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err.Error()
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err.Error()
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

//func GetFile


//func DeleteFile add support to html file filters : file:"mimeType"



// decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err.Error()
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "ciphertext too short"
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}