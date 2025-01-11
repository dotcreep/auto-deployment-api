package portainer

import (
	"context"
	"log"
)

func (p *Portainer) RollbackAddStack(ctx context.Context, name string) {
	if name == "" {
		log.Println("name is required")
	}
	var stackInfo bool = true
	var folderInfo bool = true
	_, err := p.DeleteStackByName(ctx, name)
	if err != nil {
		log.Println(err)
		stackInfo = false
	}
	err = RemoveClientDirectory(name)
	if err != nil {
		log.Println(err)
		folderInfo = false
	}
	if stackInfo && folderInfo {
		log.Println("rollback stack success")
	} else {
		log.Printf("rollback stack failed, stack info status: %v, endpoint info status: %v\n", stackInfo, folderInfo)
	}
}
