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
		log.Fatalf("Error: %s", err)
		return
	}
	fmt.Printf("响应：%s\n", completions)
}

func TestChat(t *testing.T) {
	var msg = []string{"写一篇玄幻小说。", "在一个古老而神秘的世界里，有一种被称为灵兽的特殊生物，它们拥有着超凡的能力，可以操纵元素、掌控时间、甚至能够窥探灵魂。这些灵兽被人们视为神圣的存在，只有拥有着强大的修行和天赋才能够驯服它们。故事的主人公是一个名叫李青的年轻人，他出身于一个普通始了他的修行之路。在修行的过程中，李青遇到了很多有趣的人物，包括一位神秘的老修行者、一位美丽的灵兽师、一位冷酷而强大的黑暗法师等等。他们每个人都有着不同的故事和能力，但是他们都在帮助李青成为一名更强大的灵兽师。最终，李青成功驯服了炎凤，并且成为了一们，创造了一个更加美好的世界。故事的结尾，李青回到了他的家乡，他看着那些熟悉的田地和房屋，感受到了自己成长过程中所经历的一切。他知道，自己的梦想已经实现了，而这只是他未来更加伟大的旅程的开始。",
		"继续"}
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
	msg := "一只可爱的猫，但是这只猫的耳朵特别大。"
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
	msg := "请把中间的中文去掉"
	path := "/Users/zhangmai/test/test.png"
	file, err := os.Open(path)
	if err != nil {
		return
	}
	reply, err := ImagesEdits(file, msg, "1024x1024", 2)
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	for _, v := range reply {
		fmt.Println(v)
	}
}
