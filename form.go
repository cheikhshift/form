package form

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/cheikhshift/gos/core"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Handle file upload. Use func Path to get a 
// string path relative to your filesystem.
type File string
// Display textarea
type Paragraph string
// Display Date input. Return string of Date.
// Format YYYY-M-d
type Date string
// Display dropdown.
// Tag properties :
// title : title of field.
// placeholder : Prompt left of field.
// select : comma delimited choices of field.
type Select string
// Display dropdown list with multiple selection support.
// Tag properties :
// title : title of field.
// placeholder : Prompt left of field.
// select : comma delimited choices of field.
type SelectMult Select
// Display radio input.
// Tag properties :
// title : title of field.
// select : comma delimited choices of field.
type Radio Select
// Display Email input.
type Email string
// Display Passworld input
type Password []byte


//select opts,placeholder,title
const (
	MB = 1 << 20
)

var (
	// string : Default input classes of
	// generared HTML input class attribute.
	InputClass  string
	// string : Default submit button attribute class of
	// generated HTML forms.
	ButtonClass string
	// int : Secret form key 
	FormKey     = "a very very very very secret key"
	// int : Maximum upload size
	MaxSize     = 10 //Mb
)

func GetSel(tagstr string) string {
	strs := strings.Split(tagstr, "\",")
	for _, v := range strs {
		if strings.Contains(v, "select") {
			splm := strings.Split(v, ":")
			return strings.Replace(splm[1], "\"", "", -1)
		}
	}

	return ""
}

func GetPl(tagstr string) string {
	strs := strings.Split(tagstr, "\",")
	for _, v := range strs {
		if strings.Contains(v, "placeholder") {
			splm := strings.Split(v, ":")
			return strings.Replace(splm[1], "\"", "", -1)
		}
	}

	return ""
}

func Hash(input string) []byte {
	return sha512.New512_256().Sum([]byte(input))
}

func Path(pathtoformat File) string {
	return fmt.Sprintf("./uploads/%s", pathtoformat)
}

// Parse form from request and set data to specified
// pointer.
func Form(r *http.Request, i interface{}) error {
	//convert form and validate
	var err error
	v := reflect.ValueOf(i).Elem()
	//t := reflect.TypeOf(item)
	if strings.Contains(r.Header.Get("content-type"), "/json") {

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(i)
		if err != nil {
			return err
		}

		_, err = govalidator.ValidateStruct(i)
		if err != nil {
			return err
		}
	} else {

		bso := bson.M{"__elem": "form"}

		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			fieldtype := strings.ToLower(field.Type.String())
			if strings.Contains(fieldtype, "file") {
				file, handler, _ := r.FormFile(field.Name)

				if handler != nil {
					fid := fmt.Sprintf("%s-%s", core.NewLen(15), handler.Filename)
					f, err := os.OpenFile(fmt.Sprintf("./uploads/%s", fid), os.O_WRONLY|os.O_CREATE, 0700)
					if err != nil {

						return err
					}

					io.Copy(f, file)
					f.Close()
					bso[field.Name] = fid
					file.Close()
				}
			} else if strings.Contains(fieldtype, "bool") {
				bso[field.Name] = strings.Contains(r.FormValue(field.Name), "on")
			} else if strings.Contains(fieldtype, "password") {
				bso[field.Name] = Hash(r.FormValue(field.Name))
			} else if strings.Contains(fieldtype, "int") {
				i, _ := strconv.Atoi(r.FormValue(field.Name))
				bso[field.Name] = i
			} else if strings.Contains(fieldtype, "float") {
				f, _ := strconv.ParseFloat(r.FormValue(field.Name), 64)
				bso[field.Name] = f
			} else {
				bso[field.Name] = r.FormValue(field.Name)

			}

		}

		data, _ := json.Marshal(&bso)
		json.Unmarshal([]byte(data), i)

		_, err = govalidator.ValidateStruct(i)
		if err != nil {
			return err
		}
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

// Returns bool : Verifies if request is valid
func ValidateRequest(r *http.Request, usertoken string) bool {
	var token string
	if strings.Contains(r.Header.Get("content-type"), "/json") {
		decoder := json.NewDecoder(r.Body)
		var i bson.M
		err := decoder.Decode(i)
		if err != nil {
			return false
		}

		token = Decrypt([]byte(FormKey), i["formtoken"].(string))

	} else {
		token = Decrypt([]byte(FormKey), r.FormValue("formtoken"))
	}
	return strings.Contains(token, usertoken) && strings.Contains(token, r.URL.Path)
}

// Returns form token.
// Each token is specific to the URL specified.
func GenerateToken(url string, usertoken string) string {
	return Encrypt([]byte(FormKey), fmt.Sprintf("%s@%s@%s", url, time.Now().String(), usertoken))
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
