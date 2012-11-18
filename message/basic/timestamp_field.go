package basic

import(
    "time"
    "quickfixgo/message"
    )

type UTCTimestampField struct {
  FieldBase
  timeValue time.Time
}

func (f *UTCTimestampField) UTCTimestampValue() time.Time { return f.timeValue}

const (
  utcTimestampFormat = "20060102-15:04:05.000"
  utcTimestampNoMillisFormat = "20060102-15:04:05"
)



//returns utc timestamp field in the format YYYYMMDD-HH:MM:SS.sss 
func NewUTCTimestampField(tag message.Tag, value time.Time) *UTCTimestampField {
  f:=new(UTCTimestampField)
  f.init(tag, value.UTC().Format(utcTimestampFormat))
  f.timeValue=value 

  return f
}

//returns utc timestamp field in the format YYYYMMDD-HH:MM:SS
func NewUTCTimestampFieldNoMillis(tag message.Tag, value time.Time) *UTCTimestampField {
  f:=new(UTCTimestampField)
  f.init(tag, value.UTC().Format(utcTimestampNoMillisFormat))
  f.timeValue=value
  return f
}


//converts a generic field to a utc timestamp field
//check error for convert errors
func ToUTCTimestampField(f message.Field) (*UTCTimestampField, error) {
  //with millisecs
  value, err:=time.Parse(utcTimestampFormat,f.Value())
  if err==nil {
    return NewUTCTimestampField(f.Tag(), value), nil
  }

  //w/o millisecs
  value, err=time.Parse(utcTimestampNoMillisFormat,f.Value())
  if err==nil {
    return NewUTCTimestampFieldNoMillis(f.Tag(), value), nil
  }

  return nil, err
}