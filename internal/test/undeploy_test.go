package test

import (
	"fmt"
	"testing"
)

func TestUndeploy(t *testing.T) {
	data := []string{"tokomu", "yeye", "cinta", "hariyadistore", "nobita", "bakulanku", "nono", "salim", "dila", "angga", "slopy", "soundsystem", "priza", "manggojo2shop", "cemilanenak", "manggojo22shop", "dayat", "manggojo3shop", "adam", "stefani", "adang", "naura12", "intan", "momoyo", "aspirasihomman", "aspixxxxxn", "kyky", "sepatujaja", "kopyahkung", "dika", "dadangsilva", "yayan", "dikiiw", "wawan", "dorami"}

	for i := 0; i < len(data); i++ {
		fmt.Println(data[i])
	}
}
