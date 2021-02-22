package xlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobHandler func(conf *ServerConfig, o interface{}) error

type IJob interface {
	GetChanelID() primitive.ObjectID
	GetUserID() primitive.ObjectID
}

type Job struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChannelID  primitive.ObjectID `json:"chanelid" bson:"chanelid"`
	UserID     primitive.ObjectID `json:"userid" bson:"userid"`
	Status     JobStatusType      `json:"status" bson:"status"`
	Message    string             `json:"message,omitempty" bson:"message,omitempty"`
	UpdateDate time.Time          `json:"updatedate" bson:"updatedate"`
	InsertDate time.Time          `json:"insertdate" bson:"insertdate"`
	conf       *ServerConfig      `json:"-" bson:"-"`
}

type JobStatusType int

const (
	jCreated JobStatusType = iota
	jRunning
	jCompleted
	jFailed
)

func (v JobStatusType) String() string {
	var toString = map[JobStatusType]string{
		jCreated:   "created",
		jRunning:   "running",
		jCompleted: "completed",
		jFailed:    "failed",
	}
	return toString[v]
}

// MarshalJSON marshals the enum as a quoted json string
func (v JobStatusType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (v *JobStatusType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	var toID = map[string]JobStatusType{
		"created":   jCreated,
		"running":   jRunning,
		"completed": jCompleted,
		"failed":    jFailed,
	}
	*v = toID[j]
	return nil
}

func (j *Job) setStatusCompleted() error {
	return j.updateJobStatus(jCompleted, "")
}

func (j *Job) setStatusRunning() error {
	return j.updateJobStatus(jRunning, "")
}

func (j *Job) setStatusFailed(message string) error {
	return j.updateJobStatus(jFailed, message)
}

func (j *Job) updateJobStatus(status JobStatusType, message string) error {
	if j.conf == nil {
		return errors.New("job config error")
	}
	j.Status = status
	j.Message = message

	filter := bson.M{"_id": j.ID}
	upd := bson.D{{Key: "$set", Value: j}}

	_, err := j.conf.DB.UpdateOne("jobs", filter, upd)
	if err != nil {
		j.conf.Logger.Fatal(err)
		return err
	}
	return nil

}

func CanExecJob(conf *ServerConfig, ijob IJob) (*Job, error) {
	var job Job
	filter := bson.M{
		"chanelid": ijob.GetChanelID(),
		"$or": []interface{}{
			bson.M{"status": jCreated},
			bson.M{"status": jRunning},
		},
	}

	if err := conf.DB.FindOne("jobs", nil, filter, &job); err == nil {
		if time.Now().UTC().Sub(job.UpdateDate).Minutes() < 30 {
			return nil, err
		}
		job.conf = conf
		job.setStatusFailed("Never finished")
	}
	job.conf = conf
	return &job, nil
}

func CreateJob(conf *ServerConfig, ijob IJob) (*Job, error) {
	j := Job{}
	j.conf = conf
	j.UserID = ijob.GetUserID()
	j.InsertDate = time.Now().UTC()
	j.UpdateDate = j.InsertDate
	j.ChannelID = ijob.GetChanelID()
	j.Status = jCreated
	oid, err := conf.DB.InsertOne("jobs", &j)
	if err != nil {
		return nil, err
	}
	j.ID = oid.(primitive.ObjectID)
	return &j, err
}

func (j *Job) Execute(o interface{}, f JobHandler) error {

	err := j.setStatusRunning()
	if err != nil {
		return err
	}
	err = f(j.conf, o)

	if err != nil {
		_ = j.setStatusFailed(err.Error())
		return err
	}

	_ = j.setStatusCompleted()
	return nil
}
