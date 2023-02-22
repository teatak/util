package common

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

const MobileJoiner = ":"
const WechatUnionPrefix = "wechat:"
const Date = "2006-01-02"
const DateTime = "2006-01-02 15:04:05"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Context(dur ...time.Duration) (context.Context, context.CancelFunc) {
	var duration time.Duration = 5 * time.Second
	if len(dur) > 0 {
		duration = dur[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	return ctx, cancel
}

func ResolveRegexString(text string) string {
	replacer := strings.NewReplacer(
		".", "",
		"+", "",
		"*", "",
		"?", "",
		"^", "",
		"$", "",
		"(", "",
		")", "",
		"[", "",
		"]", "",
		"{", "",
		"}", "",
		"|", "",
		"\\", "",
	)
	text = replacer.Replace(text)
	return text
}

func Pinyin(input string) []string {

	arr := []string{}
	a := pinyin.NewArgs()
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	py := strings.Join(pinyin.LazyPinyin(input, a), "")
	if py != "" && py != input {
		arr = append(arr, py)
	}
	a.Style = pinyin.FirstLetter
	py = strings.Join(pinyin.LazyPinyin(input, a), "")
	if py != "" && py != input {
		arr = append(arr, py)
	}

	return arr
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	bytes := h.Sum(nil)
	result := hex.EncodeToString(bytes)
	return result
}

func ShaPass(pass string) string {
	h := sha1.New()
	h.Write([]byte(pass))
	bytes := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(bytes)
}

func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bytes := h.Sum(nil)
	hexString := hex.EncodeToString(bytes)
	return hexString
}

func GenerateUUID() string {
	unix32bits := uint32(time.Now().UTC().Unix())
	buff := make([]byte, 12)
	numRead, err := rand.Read(buff)
	if numRead != len(buff) || err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x-%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:])
}

func MustString(input *string) string {
	output := ""
	if input != nil {
		output = *input
	}
	return output
}

func MaskMobile(mobile string) string {
	if mobile == "" {
		return mobile
	}
	arr := strings.Split(mobile, MobileJoiner)
	code := ""
	number := ""
	if len(arr) > 1 {
		//存在国家码
		code = arr[0]
		number = arr[1]
	} else {
		code = ""
		number = arr[0]
	}
	if len(number) < 10 {
		if len(number) < 6 {
			return number
		}
		mask := strings.Repeat("*", len(number)-5)
		if code == "" {
			return number[:2] + mask + number[len(number)-3:]
		} else {
			return code + MobileJoiner + number[:2] + mask + number[len(number)-3:]
		}
	} else {
		mask := strings.Repeat("*", len(number)-7)
		if code == "" {
			return number[:3] + mask + number[len(number)-4:]
		} else {
			return code + MobileJoiner + number[:3] + mask + number[len(number)-4:]
		}
	}
}

func MaskEmail(email string) string {
	if email == "" {
		return email
	}
	arr := strings.Split(email, "@")
	if len(arr) > 1 {
		account := arr[0]
		domain := arr[1]
		if len(account) > 8 {
			mask := strings.Repeat("*", len(account)-6)
			return account[:4] + mask + account[len(account)-2:] + "@" + domain
		}
		if len(account) > 6 {
			mask := strings.Repeat("*", len(account)-3)
			return account[:2] + mask + account[len(account)-1:] + "@" + domain
		} else {
			mask := strings.Repeat("*", len(account)-1)
			return account[:1] + mask + "@" + domain
		}
	} else {
		return email
	}

}

func MaskString(input string) string {
	l := len(input)
	if l < 6 {
		return input
	}
	if l < 10 {
		mask := strings.Repeat("*", l-5)
		return input[:2] + mask + input[l-3:]
	} else {
		mask := strings.Repeat("*", l-7)
		return input[:3] + mask + input[l-4:]
	}
}

var charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func GenerateUniqueString(length int) string {
	t := time.Now().UnixNano()
	timeUnix := Reverse(ConvertToShortUrl(t))
	length = (length - len(timeUnix))
	if length > 0 {
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return timeUnix + string(b)
	} else {
		return timeUnix
	}
}

func GenerateString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateNumber(min, max int) int {
	i := rand.Intn(max-min) + min
	return i
}

// 转换url
func ConvertToShortUrl(id int64) string {
	// 1 -- > 1
	// 10-- > a
	// 62-- > Z
	var shortUrl []byte
	count := int64(len(charset))
	for {
		var result byte
		number := id % count
		result = charset[number]
		var tmp []byte
		tmp = append(tmp, result)
		shortUrl = append(tmp, shortUrl...)
		id = (id - number) / count
		if id == 0 {
			break
		}
	}
	return string(shortUrl)
}

func ConvertFromShortUrl(str string) int64 {
	l := len(str)
	count := int64(len(charset))
	var result int64 = 0
	for _, r := range str {
		l--
		i := strings.Index(charset, string(r))
		result += int64(i) * int64(math.Pow(float64(count), float64(l)))
	}
	return result
}

var charsetLower = "0123456789abcdefghjkmnpqrstuvwxyz"

// 转换url
func ConvertToLowerShortUrl(id int64) string {
	var shortUrl []byte
	count := int64(len(charsetLower))
	for {
		var result byte
		number := id % count
		result = charsetLower[number]
		var tmp []byte
		tmp = append(tmp, result)
		shortUrl = append(tmp, shortUrl...)
		id = (id - number) / count
		if id == 0 {
			break
		}
	}
	return string(shortUrl)
}

func ConvertFromLowerShortUrl(str string) int64 {
	l := len(str)
	count := int64(len(charsetLower))
	var result int64 = 0
	for _, r := range str {
		l--
		i := strings.Index(charsetLower, string(r))
		result += int64(i) * int64(math.Pow(float64(count), float64(l)))
	}
	return result
}

func Hmac(str, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(str))
	bytes := h.Sum(nil)
	hexString := hex.EncodeToString(bytes)
	return hexString
}

// 验证Email，必须是字母数字组合 6-20位
func ValidateEmail(e string) bool {
	_, err := mail.ParseAddress(e)
	return err == nil
}

// 验证是否是手机号码
func ValidateMobile(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	var emailRegex = regexp.MustCompile(`^\+d{1,3}\:[1]([3-9])[0-9]{9}$`)
	return emailRegex.MatchString(e)
}

type Type string

const (
	TypeEmail  Type = "EMAIL"
	TypeMobile Type = "MOBILE"
	TypeUnknow Type = "UNKNOW"
)

func IdentityType(e string) Type {
	if ValidateEmail(e) {
		return TypeEmail
	}
	if ValidateMobile(e) {
		return TypeMobile
	}
	return TypeUnknow
}

// 验证用户名，必须是字母数字组合 6-20位
func ValidateName(name string) bool {
	letters := len(name)
	if letters < 6 || letters > 20 {
		return false
	}
	hasSpecial := false
	hasSpace := false
	for _, c := range name {
		switch {
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		case unicode.IsSpace(c):
			hasSpace = true
		default:
			//return false, false, false, false
		}
	}
	fmt.Println(hasSpace, hasSpecial)
	if hasSpace || hasSpecial {
		return false
	}
	return true
}

// 验证密码
func ValidatePassword(password string) bool {
	hasNumber := false
	hasLetter := false
	hasSpecial := false
	hasSpace := false
	letters := len(password)
	if letters < 6 || letters > 20 {
		return false
	}
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		case unicode.IsSpace(c):
			hasSpace = true
		default:
			//return false, false, false, false
		}
	}

	if hasSpace {
		return false
	}

	if hasLetter && hasNumber {
		return true
	}

	if hasSpecial && hasNumber {
		return true
	}

	if hasLetter && hasSpecial {
		return true
	}

	return false
}

func PureMobile(mobile string) string {
	if mobile == "" {
		return mobile
	}
	arr := strings.Split(mobile, MobileJoiner)
	number := ""
	if len(arr) > 1 {
		//存在国家码
		number = arr[1]
	} else {
		number = arr[0]
	}
	return number
}
