package ssm

import (
	. "constants"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Context struct {
	instanceId  string
	instanceIds []string
	command     string
	commandId   string
	session     *session.Session
}

func (c *Context) GetInstanceIds() []string {
	return c.instanceIds
}

func (c *Context) SetInstanceIds(instanceIds []string) {
	c.instanceIds = instanceIds
}

func (c *Context) GetInstanceId() string {
	return c.instanceId
}

func (c *Context) SetInstanceId(instanceId string) {
	c.instanceId = instanceId
}

func (c *Context) GetCommand() string {
	return c.command
}

func (c *Context) SetCommand(command string) {
	c.command = command
}

func (c *Context) GetCommandId() string {
	return c.commandId
}

func (c *Context) SetCommandId(commandId string) {
	c.commandId = commandId
}

func (c *Context) Init() error {
	err := c.setSession()
	if err != nil {
		return err
	}
	return err
}

func (c *Context) setSession() error {
	var err error
	c.session, err = session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(SSM_REGION),
			Credentials: credentials.NewStaticCredentials(SSM_ACCESS_KEY, SSM_SECRET_KEY, ""),
		},
	})
	return err
}

func (c *Context) RunCommand() (*ssm.SendCommandOutput, error) {
	output, err := ssm.New(c.session).SendCommand(&ssm.SendCommandInput{
		InstanceIds:  aws.StringSlice(c.instanceIds),
		DocumentName: aws.String("AWS-RunShellScript"),
		Parameters: map[string][]*string{
			"commands": {
				aws.String(c.command),
			},
		},
	})

	return output, err
}

func (c *Context) GetResult() (*ssm.GetCommandInvocationOutput, error) {
	output, err := ssm.New(c.session).GetCommandInvocation(&ssm.GetCommandInvocationInput{
		CommandId:  aws.String(c.commandId),
		InstanceId: aws.String(c.instanceId),
	})
	return output, err
}
