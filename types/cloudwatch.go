package types

import (
	"encoding/json"
	"time"
)

// CloudWatchEvent is the event sent by CloudWatch to Lambda functions
type CloudWatchEvent struct {
	Version    string          `json:"version"`     // e.g. 0
	ID         string          `json:"id"`          // e.g. 12345678-1234-1234-1234-123456789012
	DetailType string          `json:"detail-type"` // e.g. AWS API Call via CloudTrail or mobile
	Source     string          `json:"source"`      // aws.partner/salesforce.com/... or aws.ec2 or rudderstack
	AccountID  string          `json:"account"`     // e.g. 123456789012
	Time       time.Time       `json:"time"`        // e.g. 2019-03-01T21:49:13Z
	Region     string          `json:"region"`      // e.g. us-east-1
	Resources  []string        `json:"resources"`   // ARNs of resources
	Detail     json.RawMessage `json:"detail"`      // this is the raw JSON event
}

// SalesforceDetailEvent is the event sent by Salesforce to Lambda functions
type SalesforceDetailEvent struct {
	Payload json.RawMessage `json:"payload"`
}
