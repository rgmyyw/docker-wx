package Algorithm

import (
	"github.com/golang/protobuf/proto"
	"math/rand"
	"wechatdll/Cilent/mm"
)

func wcstf_and_wcste_with_count(count, time_interval int) (mm.Wcstf, mm.Wcste) {
	rs := rand.Intn(60) + 60
	tes := time_interval - rs
	ten := time_interval
	// 组装cswstf
	tfs := (time_interval - rs + rand.Intn(1) + 3) * 1000
	wcstf := mm.Wcstf{}
	wcstf.StartTime = proto.Uint64(uint64(tfs))
	for i := 0; i <= count; i++ {
		tfs = tfs + rand.Intn(500) + 500
		wcstf.EndTime = append(wcstf.EndTime, uint64(tfs))
	}
	tfn := tfs + (rand.Intn(7)+3)*1000 + rand.Intn(1000)
	wcstf.CheckTime = proto.Uint64(uint64(tfn))
	wcstf.Count = proto.Uint32(uint32(count))
	// 组装cswste
	wcste := mm.Wcste{}
	// 不同的版本具有不同的wcste, 这里固定为801
	wcste.Checkid = proto.String("4")
	wcste.StartTime = proto.Uint32(uint32(tes))
	wcste.CheckTime = proto.Uint32(uint32(ten))
	wcste.Count1 = proto.Uint32(uint32(0))
	wcste.Count2 = proto.Uint32(uint32(5124))
	return wcstf, wcste

}
