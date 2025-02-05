package jenkins

import (
	"context"
	"log"
)

func (p *Jenkins) RollbackAddItem(ctx context.Context, data *JenkinsData) {
	if data.Username == "" {
		log.Println("name is required")
	}

	var jobInfo bool = true
	var credentialInfo bool = true
	_, err := p.DeleteJob(ctx, data)
	if err != nil {
		log.Println(err)
		jobInfo = false
	}
	_, err = p.DeleteCredential(ctx, data)
	if err != nil {
		log.Println(err)
		credentialInfo = false
	}
	if jobInfo && credentialInfo {
		log.Println("rollback item success")
	} else {
		log.Printf("rollback item failed, job info status: %v, credential info status: %v\n", jobInfo, credentialInfo)
	}
}
