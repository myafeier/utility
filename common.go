/**
 * author       :xiafei
 * CreateTime   :16/6/16 下午7:44
 * Description  :
 */

package utility

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	logger "github.com/myafeier/log"
	"net/http"
	"os"
	"sync"

	"github.com/go-macaron/session"

	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"gopkg.in/macaron.v1"
	//"crypto/sha256"
	"bytes"
	"math"
)

func init()  {
	logger.SetPrefix("Utility")
	logger.SetLogLevel(logger.DEBUG)
}


type Common struct {
	sync.Mutex
	SiteSecretKey string
}

var (
	emailPattern                      = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")
	urlPattern                        = regexp.MustCompile(`(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	phonePattern                      = regexp.MustCompile("^1\\d{10}$")
	longTelPattern                    = regexp.MustCompile("^0[1-9]\\d{1,2}-\\d{6,8}$")
	longTelWithOutZeroPattern         = regexp.MustCompile("^[1-9]\\d{1,2}-\\d{6,8}$")
	longTelWithZeroWithOutDashPattern = regexp.MustCompile("^0[1-9]\\d{1,2}\\d{6,8}$")
	longTelWithOutZeroAndDashPattern  = regexp.MustCompile("^[2-9]\\d{9,11}$")
	shortTelPattern                   = regexp.MustCompile("^[2-9]\\d{6,7}$")
	hidPattern                        = regexp.MustCompile("(^(00|k)\\d*$)|(^Z\\w*$)")
	numberVerifyCodePattern           = regexp.MustCompile("^\\d{4,6}$")
	capitalLetterPattern              = regexp.MustCompile("[A-Z]+")
	lowerCaseLetterPattern            = regexp.MustCompile("[a-z]+")
	numberPattern                     = regexp.MustCompile("\\d+")
	InjectionPattern                  = regexp.MustCompile("(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)")
	mutex                             sync.Mutex
	Sha1OffSet                        = "@skdfl.3c"
)

func NewCommon(siteSecretKey string)*Common{
	u:=new(Common)
	if siteSecretKey==""{
		u.SiteSecretKey="sdfk377hYtegxytr42870Ij&8w"
	}
	return u
}

func (self *Common) GetUuid() string {
	self.Lock()
	defer self.Unlock()

	randomFileHandler,err := openRandomDev()
	if err != nil {
		panic(err)
	}
	b := make([]byte, 16)
	randomFileHandler.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (self *Common)ExportCSV(csv, filename string, ctx *macaron.Context) {
	//"\xEF\xBB\xBF"
	aLen := len(csv)
	result := []byte(csv)
	// json rendered fine, write out the result
	ctx.Resp.Header().Set("Content-Disposition", "attachment; filename="+filename+".csv")
	ctx.Resp.Header().Set("Content-Type", "text/csv; charset=UTF-8")
	ctx.Resp.Header().Set("Content-Length", strconv.Itoa(aLen))
	ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
	ctx.Resp.WriteHeader(http.StatusOK)
	ctx.Resp.Write(result)
	return
}

//唯一码生成器
func openRandomDev() (*os.File,error) {
	return os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
}





////获取用户信息
//func GetUserInfoOfContex(ctx *macaron.Context) (result *model.User, err error) {
//	var user interface{}
//	var ok bool
//	if user, ok = ctx.Data["current_user"]; !ok {
//		err = fmt.Errorf("用户未登陆!")
//		return
//	} else {
//		result, ok = user.(*model.User)
//		if !ok {
//			err = fmt.Errorf("用户信息获取失败")
//			return
//		} else {
//			return
//		}
//	}
//}

//获取当前用户的uid
func (self *Common)GetUidOfSession(sess session.Store) (uid int64, err error) {
	suid := sess.Get("uid")
	var ok bool
	if uid, ok = suid.(int64); ok {
		if uid <= 0 {
			err = fmt.Errorf("No user Logined")
		} else {
			return
		}
	} else {
		err = fmt.Errorf("No user Logined")
	}
	return
}

//获取用户的当前角色
func (self *Common)GetUserRidOfSession(sess session.Store) (rid int, err error) {
	srid := sess.Get("current_rid")
	var ok bool
	if rid, ok = srid.(int); ok {
		if rid == 0 {
			err = fmt.Errorf("No rid set")
		} else {
			return
		}
	} else {
		err = fmt.Errorf("No rid set")
	}
	return
}

//获取用户的当前医院
func (self *Common)GetUserHidOfSession(sess session.Store) (hid int64, err error) {
	shid := sess.Get("current_hid")
	var ok bool
	if hid, ok = shid.(int64); ok {
		if hid == 0 {
			err = fmt.Errorf("No hid set")
		} else {
			return
		}
	} else {
		err = fmt.Errorf("No hid set")
	}
	return
}

//获取用户的当前部门
func (self *Common)GetUserDidOfSession(sess session.Store) (did int64, err error) {
	sdid := sess.Get("current_did")
	var ok bool
	if did, ok = sdid.(int64); ok {
		if did == 0 {
			err = fmt.Errorf("No did set")
		} else {
			return
		}
	} else {
		err = fmt.Errorf("No did set")
	}
	return
}

func (self *Common)ParseStrType(str string) (stype string) {
	stype = "unknown"
	if phonePattern.MatchString(str) {
		stype = "phone"
		return
	}
	if emailPattern.MatchString(str) {
		stype = "email"
		return
	}

	if hidPattern.MatchString(str) {
		stype = "hid"
		return
	}
	return
}

func (self *Common)PhoneVefify(phone string) bool {
	return phonePattern.MatchString(phone)
}

func (self *Common)Int64SliceToString(s []int64) []string {
	var r []string
	for _, v := range s {
		r = append(r, strconv.FormatInt(v, 10))
	}
	return r
}
func(self *Common) HidVefify(hid string) bool {
	return hidPattern.MatchString(hid)
}

//号码区号带零，有破折号
func (self *Common)LongTelVerify(phone string) bool {
	return longTelPattern.MatchString(phone)
}

//号码区号不带零，有破折号
func (self *Common)LongTelWithOutZeroVerify(phone string) bool {
	return longTelWithOutZeroPattern.MatchString(phone)
}

//号码区号带零，无破折号
func (self *Common)LongTelWithZeroWithOutDashVerify(phone string) bool {
	return longTelWithZeroWithOutDashPattern.MatchString(phone)
}

//号码区号不带零，无破折号
func (self *Common)LongTelWithoutZeroAndDashVerify(phone string) bool {
	return longTelWithOutZeroAndDashPattern.MatchString(phone)
}
func (self *Common)ShortTelVerify(phone string) bool {
	return shortTelPattern.MatchString(phone)
}

func (self *Common)EmailVerify(email string) bool {
	return emailPattern.MatchString(email)
}

//长度6以上,包含大写字母和数字
func (self *Common)PasswordComplexity(str string) (level int) {
	if len(str) < 6 {
		return
	}
	if capitalLetterPattern.FindString(str) != "" {
		level++
	}
	if numberPattern.FindString(str) != "" {
		level++
	}
	if lowerCaseLetterPattern.FindString(str) != "" {
		level++
	}
	return
}

func (self *Common)NumberVerifyCode(number string) bool {
	return numberVerifyCodePattern.MatchString(number)
}

func (self *Common)GetStringFromContent(key string, ctx *macaron.Context) (result string, err error) {
	var str interface{}
	var ok bool
	if str, ok = ctx.Data[key]; !ok {
		err = fmt.Errorf("no key[%s] in content!", key)
		return
	}
	if result, ok = str.(string); !ok {
		err = fmt.Errorf("key[%s] type is not string in content!,type is %T", key, str)
		return
	}
	return
}

//生成密码
func (self *Common)CreatePasswd(pwd string) string {
	hash := md5.New()
	hash.Write([]byte(pwd + self.SiteSecretKey))
	result := hash.Sum(nil)
	return hex.EncodeToString(result)
}

//生成随机字符串
func (self *Common)GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成随机字符串
func (self *Common)CreateRandomCode(length int) string {
	str := "123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (self *Common)GetFirstPartOfFileName(fn string) string {
	strArray := strings.Split(fn, ".")
	return strArray[0]
}

//Api的验签
func (self *Common)VerifyApiSign(reqParam map[string]string, sign string) (err error) {
	reqLength := len(reqParam)
	paramSlice := make([]string, 0, reqLength)
	for k, _ := range reqParam {
		paramSlice = append(paramSlice, k)
	}
	sort.Strings(paramSlice)
	buf := new(bytes.Buffer)
	for k, v := range paramSlice {
		buf.WriteString(v)
		buf.WriteByte('=')
		buf.WriteString(reqParam[v])
		if k < (reqLength - 1) {
			buf.WriteByte('&')
		}
	}

	cilper := md5.Sum(buf.Bytes())
	trueMd5Value := hex.EncodeToString(cilper[:])

	if trueMd5Value == sign {
		return
	} else {
		err = fmt.Errorf("Error found when verify sign,request sign:%s, the true Sign is:%s,reqeust param:%+v\n", sign, trueMd5Value, reqParam)
		return
	}
	return

}



func (self *Common)GetRandomNumberByTime() string {
	self.Lock()
	defer self.Unlock()

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(99999)
	nowtime := time.Now()
	return fmt.Sprintf("%d%d%d", nowtime.Unix(), nowtime.Nanosecond(), randNum)
}

//生成通用版验签码

func (self *Common)GenerateGeneralSign(paras map[string]string) (sign string) {
	mutex.Lock()
	defer mutex.Unlock()
	var mLen = len(paras)

	sortedKeys := make(sort.StringSlice, 0, mLen)
	for k, _ := range paras {
		sortedKeys = append(sortedKeys, k)
	}
	sortedKeys.Sort()
	var strs string

	for i := 0; i < mLen; i++ {
		strs += sortedKeys[i] + "=" + paras[sortedKeys[i]]
		if i < mLen-1 {
			strs += "&"
		}
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(strings.ToUpper(strs)))
	md5Ciper := md5Ctx.Sum(nil)
	return hex.EncodeToString(md5Ciper)

}

//通过生日计算岁数
func (self *Common)GenerateAgeByBirthday(b int64) int {
	return int(time.Now().Sub(time.Unix(b, 0)).Hours() / float64(365*24))
}

func (self *Common)FormatPhoneNumber(number string) string {
	if phonePattern.MatchString(number) { //手机号码
		return number[:3] + "***" + number[8:]
	}
	if longTelPattern.MatchString(number) { //座机号码
		return number[:4] + "***" + number[8:]
	}
	if shortTelPattern.MatchString(number) { //座机号码
		return number[:1] + "***" + number[5:]
	}
	return number
}

//四舍五入
func (self *Common)Round(num float64) int {
	return int(math.Floor(num + 0.5))
}

//格式化金额（从int至string）
func (self *Common)FormatMoney(amount int) string {
	return fmt.Sprintf("¥%.2f", float32(amount)/100.0)
}

//格式化折扣率（从int至string）
func (self *Common)FormatDiscountRate(rate int) string {
	return fmt.Sprintf("%d%%", rate)
}

//拼接list元素,目前仅支持[]int64,[]int,[]string
func (self *Common)ImplodeInt64(list []int64, sep string) string {

	var temp []string
	for _, v := range list {
		s := strconv.FormatInt(v, 10)
		temp = append(temp, s)
	}
	return strings.Join(temp, sep)

}

func (self *Common)transSimpleTypeToString(s interface{}) (string, error) {
	switch s.(type) {
	case int64:
		return strconv.FormatInt(s.(int64), 10), nil
	case int:
		return strconv.Itoa(s.(int)), nil
	case string:
		return s.(string), nil
	default:
		return "", fmt.Errorf("UnSupported Type")
	}
}

func (self *Common)TrimFieldPhoneNumber(no string) string {
	if len(no) == 12 && strings.HasPrefix(no, "01") {
		return no[1:]
	}
	return no
}


func (self *Common)Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipher := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipher)
}
