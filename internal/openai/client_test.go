package openai

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestCompletions(t *testing.T) {
	// go test -v openai.go client_test.go
	msg := "给我介绍一下各个星座。"
	completions, err := Completions(ModelName, msg, 0.7)
	if err != nil {
		log.Fatalf("%s", err)
		return
	}
	fmt.Printf("响应：%s\n", completions)
}

func TestChat(t *testing.T) {
	var msg = []string{"愿与愁歌词", "岁月在默数三四五六 第六天以后\n人们开始存在宇宙 黑夜和白昼\n呼吸第一口气的咽喉 最怕命运小偷\n坏和美好我用血肉 去感受\n问宿命是否 再多久 再持久 再永久\n抵不了不朽\n恋人从挥手 到牵手 到放手 到挥手\n就该足够\n对夜的长吼 我胸口 的伤口 随风陈旧\n你我终会沦为尘埃漂流\n等待花季烟雨稠\n再化降水驻守\n属于你的愿与愁\n时间在倒数你在左右 多想踩碎沙漏\n但能同时在同个宇宙 就不求滞留\n呼吸下一口气的预谋 终究会被没收\n漫天风雪我陪你颤抖 我们别回头\n问宿命是否 再多久 再持久 再永久\n抵不了不朽\n恋人从挥手 到牵手 到放手 到挥手\n就该足够\n对夜的长吼 我胸口 的伤口 随风陈旧\n你我终会沦为尘埃漂流\n等待花季烟雨稠\n再化降水驻守\n属于你的愿与愁\n能爱多久 想多久 是多久 是永久\n爱过就不朽\n那我不走 不分手 不放手 不挥手\n十指紧扣\n分岔路口 我伤口 贪与渴求\n渺小微弱像尘埃漂流\n等待花季烟雨稠\n再化降水驻守\n属于你的愿与愁\n分岔路口 我胸口 的伤口 贪与渴求\n渺小微弱像尘埃漂流\n等待花季烟雨稠\n再化降水驻守\n只为重逢的时候",
		"在上面歌词找灵感，创作一首新歌，包括歌名和歌词。"}
	reply, err := Chat("gpt-3.5-turbo", msg, 0.7)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	fmt.Println("响应: ", reply)
}

func TestEdits(t *testing.T) {
	input := "What day of the wek is it?"
	instruction := "Fix the spelling mistakes"
	reply, err := Edits("text-davinci-edit-001", input, instruction, 0.7)
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}
	fmt.Println("修改完成后文本：", reply)
}

func TestImages(t *testing.T) {
	msg := "总是梦见云层之上飞过子午线\n分不清是黑夜还是白天\n带着装不下的期待匆匆地赶来\n我再想一遍\n想一遍\n我们寻找着在这条路的中间\n我们迷失着在这条路的两端\n每当黄昏阳光把所有都渲染\n你看那金黄多耀眼\n我们奔跑着在这条路的中间\n我们哭泣着在这条路的两端\n每当黄昏阳光把所有都渲染\n我看到夜的黑暗\n晚风吹过金色沙滩海边的晚宴\n那种味道现在还不习惯\n拉斯维加斯往返路上我看见\n这里无人烟\n无人烟\n我们寻找着在这条路的中间\n我们迷失着在这条路的两端\n每当黄昏阳光把所有都渲染\n你看那金黄多耀眼\n我们奔跑着在这条路的中间\n我们哭泣着在这条路的两端\n每当黄昏阳光把所有都渲染\n我看到夜的黑暗\n奔跑着在这条路的中间\n哭泣着在这条路的两端\n每当黄昏阳光把所有都渲染\n我看到夜的黑暗"
	reply, err := ImagesGenerations(msg, "1024x1024", 2)
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	for _, v := range reply {
		fmt.Println(v)
	}
}

func TestImagesEdits(t *testing.T) {
	msg := "请把水印去掉。"
	filePath := "/Users/zhangmai/test/test.png"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	reply, err := ImagesEdits(file, msg, "1024x1024", 2)
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	for _, v := range reply {
		fmt.Println(v)
	}
}
