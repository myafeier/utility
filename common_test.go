package utility

import (
	//"strconv"
	"testing"
	"time"
)

func TestNumberVerifyCode(t *testing.T) {
	stat1 := NumberVerifyCode("1234")
	stat2 := NumberVerifyCode("123456")
	stat3 := NumberVerifyCode("123")
	stat4 := NumberVerifyCode("1234dx")
	if stat1 != true {
		t.Error("1234 is valid")
	}
	if stat2 != true {
		t.Error("123456 is valid")
	}
	if stat3 == true {
		t.Error("123 is invalid")
	}
	if stat4 == true {
		t.Error("1234dx is invalid")
	}

}

func TestGenerateSign(t *testing.T) {
	test := map[string]string{"a": "1111", "b": "22222", "c": "333333"}

	result := GenerateGeneralSign(test)
	t.Log(result)

}

func TestVerifyApiSign(t *testing.T) {
	clientId := "wx"
	nonceStr := "kdish387dhsdyTGbdhx756"
	//timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	timestamp := "1469426344"
	str := "clientid=" + clientId + "&noncestr=" + nonceStr + "&timestamp=" + timestamp
	str1 := "clientid=wx&noncestr=kdish387dhsdyTGbdhx756&timestamp=1469426344"

	if str1 != str {
		t.Error("string not match!")
	}

	sign := Md5(str)

	result := VerifyApiSign(clientId, nonceStr, timestamp, sign)
	if result == nil {
		t.Log("test success")
	} else {
		t.Error("test failure")
	}

}

func TestPasswordComplexity(t *testing.T) {
	str1 := "1234"
	str2 := "a12345"
	str3 := "Ab123456"
	str4 := "Abasdfasdf"
	num1 := PasswordComplexity(str1)
	num2 := PasswordComplexity(str2)
	num3 := PasswordComplexity(str3)
	num4 := PasswordComplexity(str4)

	if num1 != 0 {
		t.Errorf("num1 shuld be 1 but:%d", num1)
	}
	if num2 != 2 {
		t.Errorf("num1 shuld be 2 but:%d", num2)
	}
	if num3 != 3 {
		t.Errorf("num1 shuld be 3 but:%d", num3)
	}
	if num4 != 2 {
		t.Errorf("num1 shuld be 2 but:%d", num4)
	}
}

func TestGenerateAgeByBirthday(t *testing.T) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		t.Error(err)
	}
	b := time.Date(2007, 12, 1, 0, 0, 0, 0, local)
	t.Log("2007-12-1", GenerateAgeByBirthday(b.Unix()))

	b = time.Date(2018, 12, 1, 0, 0, 0, 0, local)
	t.Log("2018-12-1", GenerateAgeByBirthday(b.Unix()))

	b = time.Date(1976, 1, 1, 0, 0, 0, 0, local)
	t.Log("1976-1-1", GenerateAgeByBirthday(b.Unix()))

	b = time.Date(0, 0, 0, 0, 0, 0, 0, local)

	t.Log("0-0-0:", b.Unix(), GenerateAgeByBirthday(b.Unix()))

}

func TestFormatPhoneNumbuer(t *testing.T) {
	num1 := "18987092886"
	n2 := "0871-23444323"
	n3 := "010-23444234"
	n4 := "34456333"
	n5 := "1777777"

	t.Log(shortTelPattern.MatchString(n5))
	t.Log(FormatPhoneNumber(num1))
	t.Log(FormatPhoneNumber(n2))
	t.Log(FormatPhoneNumber(n3))
	t.Log(FormatPhoneNumber(n4))
	t.Log(FormatPhoneNumber(n5))

}

func TestFormatMoney(t *testing.T) {
	m1 := 105
	m2 := 104
	e1 := "¥ 1.05"
	e2 := "¥ 1.04"

	s1 := FormatMoney(m1)
	if s1 != e1 {
		t.Errorf("Format Money error! Expected:%s,Got:%s", e1, s1)
	}

	s2 := FormatMoney(m2)
	if s2 != e2 {
		t.Errorf("Format Money error! Expected:%s,Got:%s", e2, s2)
	}

}

func TestLongTelVerify(t *testing.T) {
	s1 := "010-234432"
	s2 := "10-23434433"
	s3 := "2342344"
	s4 := "01023424234"
	//号码区号带零，有破折号
	if LongTelVerify(s1) != true {
		t.Error(s1)
	}
	if LongTelVerify(s2) != false {
		t.Error(s2)
	}
	if LongTelVerify(s3) != false {
		t.Error(s3)
	}
	if LongTelVerify(s4) != false {
		t.Error(s4)
	}
}
func TestLongTelWithoutZeroAndDashVerify(t *testing.T) {
	s1 := "010-234432"
	s2 := "10-23434433"
	s3 := "2342344"
	s4 := "01023424234"
	s5 := "87123424234"
	s6 := "1023424234"
	//号码区号不带零，无破折号
	if LongTelWithoutZeroAndDashVerify(s1) != false {
		t.Error(s1)
	}
	if LongTelWithoutZeroAndDashVerify(s2) != false {
		t.Error(s2)
	}
	if LongTelWithoutZeroAndDashVerify(s3) != false {
		t.Error(s3)
	}
	if LongTelWithoutZeroAndDashVerify(s4) != false {
		t.Error(s4)
	}
	if LongTelWithoutZeroAndDashVerify(s5) != true {
		t.Error(s5)
	}
	if LongTelWithoutZeroAndDashVerify(s6) != true {
		t.Error(s6)
	}
}

func TestLongTelWithOutZeroVerify(t *testing.T) {
	s1 := "010-234432"
	s2 := "10-23434433"
	s3 := "2342344"
	s4 := "01023424234"
	s5 := "87123424234"
	s6 := "1023424234"
	////号码区号不带零，有破折号
	if LongTelWithOutZeroVerify(s1) != false {
		t.Error(s1)
	}
	if LongTelWithOutZeroVerify(s2) != true {
		t.Error(s2)
	}
	if LongTelWithOutZeroVerify(s3) != false {
		t.Error(s3)
	}
	if LongTelWithOutZeroVerify(s4) != false {
		t.Error(s4)
	}
	if LongTelWithOutZeroVerify(s5) != false {
		t.Error(s5)
	}
	if LongTelWithOutZeroVerify(s6) != false {
		t.Error(s6)
	}

}
func TestLongTelWithZeroWithOutDashVerify(t *testing.T) {
	s1 := "010-234432"
	s2 := "10-23434433"
	s3 := "2342344"
	s4 := "01023424234"
	s5 := "87123424234"
	s6 := "1023424234"
	////号码区号带零，无破折号
	if LongTelWithZeroWithOutDashVerify(s1) != false {
		t.Error(s1)
	}
	if LongTelWithZeroWithOutDashVerify(s2) != false {
		t.Error(s2)
	}
	if LongTelWithZeroWithOutDashVerify(s3) != false {
		t.Error(s3)
	}
	if LongTelWithZeroWithOutDashVerify(s4) != true {
		t.Error(s4)
	}
	if LongTelWithZeroWithOutDashVerify(s5) != false {
		t.Error(s5)
	}
	if LongTelWithZeroWithOutDashVerify(s6) != false {
		t.Error(s6)
	}
}

func TestTrimFieldPhoneNumber(t *testing.T) {
	n1 := "01038776788"
	n2 := "053634443876"
	n3 := "18987098663"
	n4 := "017688378833"
	if TrimFieldPhoneNumber(n1) != n1 {
		t.Error(TrimFieldPhoneNumber(n1))
	}
	if TrimFieldPhoneNumber(n2) != n2 {
		t.Error(TrimFieldPhoneNumber(n2))
	}
	if TrimFieldPhoneNumber(n3) != n3 {
		t.Error(TrimFieldPhoneNumber(n3))
	}
	if TrimFieldPhoneNumber(n4) != "17688378833" {
		t.Error(TrimFieldPhoneNumber(n4))
	}
}
